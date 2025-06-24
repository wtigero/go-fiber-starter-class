package main

import (
	"constructor-di/controller"
	"constructor-di/repository"
	"constructor-di/service"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Create Fiber app
	app := fiber.New()

	// === 🔌 MANUAL CONSTRUCTOR INJECTION ===

	// 1. สร้าง Repository (ชั้นล่างสุด - ไม่มี dependencies)
	userRepo := repository.NewUserRepository()

	// 2. สร้าง Service (ฉีด userRepo เข้าไป)
	userService := service.NewUserService(userRepo)

	// 3. สร้าง Controller (ฉีด userService เข้าไป)
	userController := controller.NewUserController(userService)

	// Setup routes
	setupRoutes(app, userController)

	log.Println("🔨 Constructor Injection API running on :3000")
	log.Fatal(app.Listen(":3000"))
}

// setupRoutes configures all routes
func setupRoutes(app *fiber.App, userController *controller.UserController) {
	// Health check
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "🔨 Constructor Injection API",
			"pattern": "Manual DI - สร้าง dependencies เอง",
			"pros":    "Simple, Fast, Safe",
			"cons":    "Verbose when many dependencies",
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
