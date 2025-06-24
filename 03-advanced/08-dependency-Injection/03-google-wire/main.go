package main

import (
	"log"
	"wire-di/controller"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Create Fiber app
	app := fiber.New()

	// === ⚡ GOOGLE WIRE CODE GENERATION ===

	// ใช้ function ที่ Wire generate ให้
	userController := InitializeUserController()
	// ⚡ Wire สร้างโค้ดแบบนี้ให้ (ใน wire_gen.go):
	// userRepository := repository.NewUserRepository()
	// userService := service.NewUserService(userRepository)
	// userController := controller.NewUserController(userService)

	// Setup routes
	setupRoutes(app, userController)

	log.Println("⚡ Google Wire DI API running on :3000")
	log.Fatal(app.Listen(":3000"))
}

// setupRoutes configures all routes
func setupRoutes(app *fiber.App, userController *controller.UserController) {
	// Health check
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "⚡ Google Wire Code Generation DI API",
			"pattern": "Code Generation - Wire สร้างโค้ดให้",
			"pros":    "Fastest, Compile-time safety, Clean code",
			"cons":    "Complex setup, Learning curve",
		})
	})

	// User endpoints
	users := app.Group("/users")
	users.Get("/", userController.GetAllUsers)
	users.Post("/", userController.CreateUser)
	users.Get("/:id", userController.GetUserByID)
	users.Put("/:id", userController.UpdateUser)
	users.Delete("/:id", userController.DeleteUser)
}
