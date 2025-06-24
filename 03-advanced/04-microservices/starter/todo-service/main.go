package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// TODO: สร้าง Todo struct
// type Todo struct {
//     ID     int    `json:"id"`
//     Title  string `json:"title"`
//     Done   bool   `json:"done"`
//     UserID int    `json:"user_id"`
// }

// TODO: สร้างตัวแปรเก็บข้อมูล todos
// var todos []Todo
// var todoID = 1

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
			"service": "todo-service",
			"status":  "healthy",
			"port":    3002,
		})
	})

	// Basic info endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "Todo Service",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	// TODO: สร้าง API endpoints
	// app.Get("/todos", getAllTodosHandler)
	// app.Post("/todos", createTodoHandler)
	// app.Get("/todos/user/:user_id", getTodosByUserHandler)

	log.Println("📝 Todo Service started on port 3002")
	log.Fatal(app.Listen(":3002"))
}

// TODO: สร้าง createTodoHandler
// func createTodoHandler(c *fiber.Ctx) error {
//     type CreateTodoRequest struct {
//         Title  string `json:"title"`
//         UserID int    `json:"user_id"`
//     }
//
//     var req CreateTodoRequest
//     if err := c.BodyParser(&req); err != nil {
//         return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
//     }
//
//     if req.Title == "" || req.UserID == 0 {
//         return c.Status(400).JSON(fiber.Map{"error": "Title and user_id are required"})
//     }
//
//     todo := Todo{
//         ID:     todoID,
//         Title:  req.Title,
//         Done:   false,
//         UserID: req.UserID,
//     }
//
//     todos = append(todos, todo)
//     todoID++
//
//     return c.Status(201).JSON(fiber.Map{
//         "success": true,
//         "todo":    todo,
//     })
// }

// TODO: สร้าง getAllTodosHandler
// func getAllTodosHandler(c *fiber.Ctx) error {
//     return c.JSON(fiber.Map{
//         "success": true,
//         "todos":   todos,
//         "count":   len(todos),
//     })
// }

// TODO: สร้าง getTodosByUserHandler
// func getTodosByUserHandler(c *fiber.Ctx) error {
//     userID := c.Params("user_id")
//
//     var userTodos []Todo
//     for _, todo := range todos {
//         if strconv.Itoa(todo.UserID) == userID {
//             userTodos = append(userTodos, todo)
//         }
//     }
//
//     return c.JSON(fiber.Map{
//         "success": true,
//         "todos":   userTodos,
//         "count":   len(userTodos),
//         "user_id": userID,
//     })
// }
