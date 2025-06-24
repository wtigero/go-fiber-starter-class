package main

import (
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// User struct
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// In-memory storage
var users []User
var userID = 1

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		},
	})

	app.Use(logger.New())

	// Initialize sample data
	initSampleData()

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "user-service",
			"status":  "healthy",
			"port":    3001,
			"users":   len(users),
		})
	})

	// Basic info endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service":     "User Service",
			"version":     "1.0.0",
			"status":      "running",
			"total_users": len(users),
		})
	})

	// API endpoints
	app.Post("/users", createUserHandler)
	app.Get("/users/:id", getUserHandler)
	app.Get("/users", getAllUsersHandler)

	log.Println("👥 User Service started on port 3001")
	log.Fatal(app.Listen(":3001"))
}

func initSampleData() {
	users = []User{
		{
			ID:        1,
			Name:      "John Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:        2,
			Name:      "Jane Smith",
			Email:     "jane@example.com",
			CreatedAt: time.Now().Add(-12 * time.Hour),
		},
		{
			ID:        3,
			Name:      "Bob Johnson",
			Email:     "bob@example.com",
			CreatedAt: time.Now().Add(-6 * time.Hour),
		},
	}
	userID = 4
	log.Println("✅ Sample users loaded")
}

func createUserHandler(c *fiber.Ctx) error {
	type CreateUserRequest struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request",
		})
	}

	if req.Name == "" || req.Email == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Name and email are required",
		})
	}

	// Check if email already exists
	for _, existingUser := range users {
		if existingUser.Email == req.Email {
			return c.Status(409).JSON(fiber.Map{
				"success": false,
				"error":   "Email already exists",
			})
		}
	}

	user := User{
		ID:        userID,
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
	}

	users = append(users, user)
	userID++

	log.Printf("👤 Created user: %s (%s)", user.Name, user.Email)

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"data":    user,
		"message": "User created successfully",
	})
}

func getUserHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	for _, user := range users {
		if strconv.Itoa(user.ID) == id {
			return c.JSON(fiber.Map{
				"success": true,
				"data":    user,
			})
		}
	}

	return c.Status(404).JSON(fiber.Map{
		"success": false,
		"error":   "User not found",
	})
}

func getAllUsersHandler(c *fiber.Ctx) error {
	// Support pagination
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	start := (page - 1) * limit
	end := start + limit

	var paginatedUsers []User
	if start < len(users) {
		if end > len(users) {
			end = len(users)
		}
		paginatedUsers = users[start:end]
	} else {
		paginatedUsers = []User{}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    paginatedUsers,
		"pagination": fiber.Map{
			"page":     page,
			"limit":    limit,
			"total":    len(users),
			"has_next": end < len(users),
			"has_prev": page > 1,
		},
		"count": len(paginatedUsers),
	})
}
