package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// TODO: สร้าง User struct
// type User struct {
//     ID    int    `json:"id"`
//     Name  string `json:"name"`
//     Email string `json:"email"`
// }

// TODO: สร้างตัวแปรเก็บข้อมูล users
// var users []User
// var userID = 1

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

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "user-service",
			"status":  "healthy",
			"port":    3001,
		})
	})

	// Basic info endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "User Service",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	// TODO: สร้าง API endpoints
	// app.Post("/users", createUserHandler)
	// app.Get("/users/:id", getUserHandler)
	// app.Get("/users", getAllUsersHandler)

	log.Println("👥 User Service started on port 3001")
	log.Fatal(app.Listen(":3001"))
}

// TODO: สร้าง createUserHandler
// func createUserHandler(c *fiber.Ctx) error {
//     type CreateUserRequest struct {
//         Name  string `json:"name"`
//         Email string `json:"email"`
//     }
//
//     var req CreateUserRequest
//     if err := c.BodyParser(&req); err != nil {
//         return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
//     }
//
//     if req.Name == "" || req.Email == "" {
//         return c.Status(400).JSON(fiber.Map{"error": "Name and email are required"})
//     }
//
//     user := User{
//         ID:    userID,
//         Name:  req.Name,
//         Email: req.Email,
//     }
//
//     users = append(users, user)
//     userID++
//
//     return c.Status(201).JSON(fiber.Map{
//         "success": true,
//         "user":    user,
//     })
// }

// TODO: สร้าง getUserHandler
// func getUserHandler(c *fiber.Ctx) error {
//     id := c.Params("id")
//
//     for _, user := range users {
//         if strconv.Itoa(user.ID) == id {
//             return c.JSON(fiber.Map{
//                 "success": true,
//                 "user":    user,
//             })
//         }
//     }
//
//     return c.Status(404).JSON(fiber.Map{
//         "success": false,
//         "error":   "User not found",
//     })
// }

// TODO: สร้าง getAllUsersHandler
// func getAllUsersHandler(c *fiber.Ctx) error {
//     return c.JSON(fiber.Map{
//         "success": true,
//         "users":   users,
//         "count":   len(users),
//     })
// }
