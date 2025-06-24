package main

import (
	"clean-arch/infrastructure/database"
	"clean-arch/infrastructure/web"
	"clean-arch/usecases/user"
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

	// Dependency Injection (Composition Root)
	// Infrastructure Layer
	userRepo := database.NewMemoryUserRepository()

	// Use Cases Layer
	getAllUsersUseCase := user.NewGetAllUsersUseCase(userRepo)
	getUserUseCase := user.NewGetUserUseCase(userRepo)
	createUserUseCase := user.NewCreateUserUseCase(userRepo)
	updateUserUseCase := user.NewUpdateUserUseCase(userRepo)
	deleteUserUseCase := user.NewDeleteUserUseCase(userRepo)

	// Interface Adapters Layer
	userController := web.NewUserController(
		getAllUsersUseCase,
		getUserUseCase,
		createUserUseCase,
		updateUserUseCase,
		deleteUserUseCase,
	)

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "🧹 Clean Architecture Example",
			"pattern": "Entities ← Use Cases ← Interface Adapters ← Frameworks",
			"layers": fiber.Map{
				"entities":           "Core business objects, business rules",
				"use_cases":          "Application business rules, orchestration",
				"interface_adapters": "Controllers, presenters, data converters",
				"frameworks":         "HTTP server, database, external tools",
			},
			"benefits": []string{
				"Framework independent",
				"Database independent",
				"Highly testable",
				"UI independent",
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
	log.Println("🧹 Clean Architecture API running on :3000")
	log.Println("🎯 Dependencies point INWARD: Frameworks → Adapters → Use Cases → Entities")
	log.Fatal(app.Listen(":3000"))
}
