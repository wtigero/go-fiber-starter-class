package main

import (
	"layered-arch/controllers"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())

	// Initialize controller
	userController := controllers.NewUserController()

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "🏢 Layered Architecture Example",
			"pattern": "Presentation → Business → Data Access",
			"layers": fiber.Map{
				"controllers":  "HTTP handlers, input validation",
				"services":     "Business logic, business validation",
				"repositories": "Data operations, database queries",
				"models":       "Data structures, entities",
			},
		})
	})

	// User routes
	users := app.Group("/users")
	users.Get("/", userController.GetAllUsers)
	users.Post("/", userController.CreateUser)
	users.Get("/:id", userController.GetUserByID)
	users.Put("/:id", userController.UpdateUser)
	users.Delete("/:id", userController.DeleteUser)

	// Start server
	log.Println("🏢 Layered Architecture API running on :3000")
	log.Println("📋 Layers: Controllers → Services → Repositories → Database")
	log.Fatal(app.Listen(":3000"))
}
