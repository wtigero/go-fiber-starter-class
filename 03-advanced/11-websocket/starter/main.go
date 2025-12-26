package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/contrib/websocket"
)

// TODO: สร้าง Hub สำหรับจัดการ connections
// TODO: สร้าง WebSocket handler
// TODO: สร้าง broadcast function
// TODO: สร้าง room management

func main() {
	app := fiber.New()

	// Static files
	app.Static("/", "./static")

	// WebSocket upgrade middleware
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// TODO: WebSocket endpoint
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// TODO: Handle WebSocket connection
		log.Println("New WebSocket connection")

		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				break
			}
			log.Printf("recv: %s", msg)

			// Echo message back
			if err := c.WriteMessage(mt, msg); err != nil {
				log.Println("write error:", err)
				break
			}
		}
	}))

	log.Fatal(app.Listen(":3000"))
}
