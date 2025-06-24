package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO: สร้าง struct Todo สำหรับ MongoDB
// type Todo struct {
//     ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
//     Title     string             `json:"title" bson:"title" validate:"required,min=1"`
//     Done      bool               `json:"done" bson:"done"`
//     CreatedAt time.Time          `json:"created_at" bson:"created_at"`
//     UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
// }

var client *mongo.Client
var todoCollection *mongo.Collection

func main() {
	// เชื่อมต่อ MongoDB
	mongoURI := getEnv("MONGODB_URI", "mongodb://admin:password123@localhost:27017/todoapp?authSource=admin")

	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("ไม่สามารถเชื่อมต่อ MongoDB:", err)
	}

	// ทดสอบการเชื่อมต่อ
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("MongoDB ping failed:", err)
	}

	// กำหนด collection
	todoCollection = client.Database("todoapp").Collection("todos")
	log.Println("เชื่อมต่อ MongoDB สำเร็จ!")

	// สร้าง Fiber app
	app := fiber.New(fiber.Config{})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())

	// TODO: เพิ่ม custom middleware สำหรับ API Key validation
	// protected := app.Group("/", APIKeyMiddleware())

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "OK",
			"database":  "connected",
			"timestamp": "2024-01-01T12:00:00Z",
		})
	})

	// TODO: เพิ่ม routes สำหรับ CRUD operations
	// protected.Get("/todos", getTodos)
	// protected.Get("/todos/:id", getTodoByID)
	// protected.Post("/todos", createTodo)
	// protected.Put("/todos/:id", updateTodo)
	// protected.Delete("/todos/:id", deleteTodo)

	// TODO: สร้างฟังก์ชัน middleware
	// func APIKeyMiddleware() fiber.Handler { ... }

	// TODO: สร้างฟังก์ชัน CRUD handlers
	// func getTodos(c *fiber.Ctx) error { ... }
	// func getTodoByID(c *fiber.Ctx) error { ... }
	// func createTodo(c *fiber.Ctx) error { ... }
	// func updateTodo(c *fiber.Ctx) error { ... }
	// func deleteTodo(c *fiber.Ctx) error { ... }

	// รัน server
	port := getEnv("PORT", "3000")
	log.Printf("Server กำลังรันที่ port %s", port)
	app.Listen(":" + port)
}

// Helper function สำหรับ environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
