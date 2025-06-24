package main

import (
	"hexagonal-arch/adapters"
	"hexagonal-arch/core/services"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Create Fiber app
	app := fiber.New()

	// === HEXAGONAL ARCHITECTURE WIRING ===

	// Secondary Adapter (Repository)
	userRepo := adapters.NewMemoryUserRepository()

	// Core Business Logic
	userService := services.NewUserService(userRepo)

	// Primary Adapter (HTTP)
	httpHandler := adapters.NewHttpUserHandler(userService)

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "🔗 Hexagonal Architecture API",
			"pattern": "External → Adapter → Port → Core → Port → Adapter → External",
		})
	})

	// User endpoints
	users := app.Group("/users")
	users.Get("/", httpHandler.GetAllUsers)
	users.Post("/", httpHandler.CreateUser)
	users.Get("/:id", httpHandler.GetUserByID)
	users.Put("/:id", httpHandler.UpdateUser)
	users.Delete("/:id", httpHandler.DeleteUser)

	log.Println("🔗 Hexagonal Architecture API running on :3000")
	log.Fatal(app.Listen(":3000"))
}
