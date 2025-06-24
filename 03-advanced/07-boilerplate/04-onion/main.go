package main

import (
	"log"
	"onion-arch/application"
	"onion-arch/infrastructure"
	"onion-arch/presentation"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Create Fiber app
	app := fiber.New()

	// === ONION ARCHITECTURE DEPENDENCY INJECTION ===

	// Infrastructure Layer (outermost)
	userRepo := infrastructure.NewMemoryUserRepository()

	// Application Layer (business logic)
	userService := application.NewUserService(userRepo)

	// Presentation Layer (adapters)
	userController := presentation.NewUserController(userService)

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "🧅 Onion Architecture API",
			"pattern": "Infrastructure → Application → Domain ← Application ← Infrastructure",
		})
	})

	// User endpoints
	users := app.Group("/users")
	users.Get("/", userController.GetAllUsers)
	users.Post("/", userController.CreateUser)
	users.Get("/:id", userController.GetUserByID)
	users.Put("/:id", userController.UpdateUser)
	users.Delete("/:id", userController.DeleteUser)

	log.Println("🧅 Onion Architecture API running on :3000")
	log.Fatal(app.Listen(":3000"))
}
