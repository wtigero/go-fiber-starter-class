package main

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
)

// ============ Types ============

type Message struct {
	Type      string    `json:"type"`
	Content   string    `json:"content,omitempty"`
	From      string    `json:"from,omitempty"`
	Room      string    `json:"room,omitempty"`
	Username  string    `json:"username,omitempty"`
	Users     []string  `json:"users,omitempty"`
	Count     int       `json:"count,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Message   string    `json:"message,omitempty"`
}

type Client struct {
	ID       string
	Conn     *websocket.Conn
	Username string
	Room     string
	Send     chan []byte
}

type Hub struct {
	clients    map[*Client]bool
	rooms      map[string]map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// ============ Global Hub ============

var hub = &Hub{
	clients:    make(map[*Client]bool),
	rooms:      make(map[string]map[*Client]bool),
	broadcast:  make(chan Message, 256),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

func main() {
	// Start hub
	go hub.Run()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		},
	})

	app.Use(logger.New())
	app.Use(cors.New())

	// Serve static files
	app.Static("/", "./static")

	// REST endpoints
	app.Get("/api/rooms", getRoomsHandler)
	app.Get("/api/stats", getStatsHandler)

	// WebSocket upgrade middleware
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket endpoint
	app.Get("/ws", websocket.New(handleWebSocket, websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}))

	log.Println("🚀 WebSocket server running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}

// ============ Hub Methods ============

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			log.Printf("✅ Client connected: %s", client.ID)
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				h.removeFromRoom(client)
				log.Printf("❌ Client disconnected: %s", client.ID)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			h.broadcastToRoom(message)
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) removeFromRoom(client *Client) {
	if client.Room != "" {
		if room, ok := h.rooms[client.Room]; ok {
			delete(room, client)
			// Notify others in room
			h.broadcastToRoom(Message{
				Type:      "user_left",
				Username:  client.Username,
				Room:      client.Room,
				Timestamp: time.Now(),
			})
			// Clean up empty room
			if len(room) == 0 {
				delete(h.rooms, client.Room)
			}
		}
	}
}

func (h *Hub) joinRoom(client *Client, room string) {
	// Leave current room first
	h.removeFromRoom(client)

	// Join new room
	client.Room = room
	if h.rooms[room] == nil {
		h.rooms[room] = make(map[*Client]bool)
	}
	h.rooms[room][client] = true

	// Notify others
	h.broadcastToRoom(Message{
		Type:      "user_joined",
		Username:  client.Username,
		Room:      room,
		Timestamp: time.Now(),
	})

	// Send online users list
	h.sendOnlineUsers(room)
}

func (h *Hub) broadcastToRoom(msg Message) {
	data, _ := json.Marshal(msg)
	room := msg.Room

	if room == "" {
		// Broadcast to all clients
		for client := range h.clients {
			select {
			case client.Send <- data:
			default:
				close(client.Send)
				delete(h.clients, client)
			}
		}
	} else {
		// Broadcast to specific room
		if clients, ok := h.rooms[room]; ok {
			for client := range clients {
				select {
				case client.Send <- data:
				default:
					close(client.Send)
					delete(clients, client)
				}
			}
		}
	}
}

func (h *Hub) sendOnlineUsers(room string) {
	if clients, ok := h.rooms[room]; ok {
		users := make([]string, 0, len(clients))
		for client := range clients {
			users = append(users, client.Username)
		}

		msg := Message{
			Type:  "online_users",
			Room:  room,
			Users: users,
			Count: len(users),
		}
		data, _ := json.Marshal(msg)

		for client := range clients {
			client.Send <- data
		}
	}
}

// ============ WebSocket Handler ============

func handleWebSocket(c *websocket.Conn) {
	client := &Client{
		ID:       uuid.New().String(),
		Conn:     c,
		Username: "Anonymous",
		Send:     make(chan []byte, 256),
	}

	hub.register <- client

	// Send welcome message
	welcome := Message{
		Type:      "welcome",
		Message:   "เชื่อมต่อสำเร็จ! กรุณา join room เพื่อเริ่มแชท",
		Timestamp: time.Now(),
	}
	data, _ := json.Marshal(welcome)
	client.Send <- data

	// Writer goroutine
	go func() {
		defer func() {
			c.Close()
		}()

		for msg := range client.Send {
			if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		}
	}()

	// Reader loop
	defer func() {
		hub.unregister <- client
		c.Close()
	}()

	for {
		_, msgBytes, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Printf("JSON parse error: %v", err)
			continue
		}

		handleMessage(client, msg)
	}
}

func handleMessage(client *Client, msg Message) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	switch msg.Type {
	case "join":
		client.Username = msg.Username
		if client.Username == "" {
			client.Username = "User-" + client.ID[:8]
		}
		hub.joinRoom(client, msg.Room)
		log.Printf("👤 %s joined room: %s", client.Username, msg.Room)

	case "message":
		if client.Room == "" {
			return
		}
		broadcast := Message{
			Type:      "message",
			Content:   msg.Content,
			From:      client.Username,
			Room:      client.Room,
			Timestamp: time.Now(),
		}
		hub.broadcastToRoom(broadcast)
		log.Printf("💬 [%s] %s: %s", client.Room, client.Username, msg.Content)

	case "typing":
		if client.Room == "" {
			return
		}
		broadcast := Message{
			Type:     "typing",
			Username: client.Username,
			Room:     client.Room,
		}
		hub.broadcastToRoom(broadcast)

	case "leave":
		hub.removeFromRoom(client)
		client.Room = ""
		log.Printf("👋 %s left room: %s", client.Username, msg.Room)
	}
}

// ============ REST Handlers ============

func getRoomsHandler(c *fiber.Ctx) error {
	hub.mu.RLock()
	defer hub.mu.RUnlock()

	rooms := make([]fiber.Map, 0)
	for name, clients := range hub.rooms {
		rooms = append(rooms, fiber.Map{
			"name":  name,
			"users": len(clients),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"rooms":   rooms,
	})
}

func getStatsHandler(c *fiber.Ctx) error {
	hub.mu.RLock()
	defer hub.mu.RUnlock()

	return c.JSON(fiber.Map{
		"success":          true,
		"total_clients":    len(hub.clients),
		"total_rooms":      len(hub.rooms),
	})
}
