package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// โครงสร้างข้อมูล Todo สำหรับ MongoDB
type Todo struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title" validate:"required,min=1"`
	Done      bool               `json:"done" bson:"done"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// Request struct สำหรับสร้าง/แก้ไข Todo
type TodoRequest struct {
	Title string `json:"title" validate:"required,min=1"`
	Done  bool   `json:"done"`
}

var client *mongo.Client
var todoCollection *mongo.Collection
var validate *validator.Validate

func main() {
	// เริ่มต้น validator
	validate = validator.New()

	// เชื่อมต่อ MongoDB
	mongoURI := getEnv("MONGODB_URI", "mongodb://admin:password123@localhost:27017/todoapp?authSource=admin")

	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("❌ ไม่สามารถเชื่อมต่อ MongoDB:", err)
	}

	// ทดสอบการเชื่อมต่อ
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("❌ MongoDB ping failed:", err)
	}

	// กำหนด collection
	todoCollection = client.Database("todoapp").Collection("todos")
	log.Println("✅ เชื่อมต่อ MongoDB สำเร็จ!")

	// สร้าง Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Todo API Intermediate v1.0",
	})

	// Middleware พื้นฐาน
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path} (${latency})\n",
	}))
	app.Use(recover.New())

	// Health check endpoint (ไม่ต้องใช้ API Key)
	app.Get("/health", func(c *fiber.Ctx) error {
		// ทดสอบการเชื่อมต่อ database
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := client.Ping(ctx, nil)
		dbStatus := "connected"
		if err != nil {
			dbStatus = "disconnected"
		}

		return c.JSON(fiber.Map{
			"status":    "OK",
			"database":  dbStatus,
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// กลุ่ม routes ที่ต้องใช้ API Key
	protected := app.Group("/", APIKeyMiddleware())

	// CRUD Routes
	protected.Get("/todos", getTodos)          // ดูรายการทั้งหมด (พร้อม pagination)
	protected.Get("/todos/:id", getTodoByID)   // ดูรายการเดียว
	protected.Post("/todos", createTodo)       // เพิ่มรายการใหม่
	protected.Put("/todos/:id", updateTodo)    // แก้ไขรายการ
	protected.Delete("/todos/:id", deleteTodo) // ลบรายการ

	// รัน server
	port := getEnv("PORT", "3000")
	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📖 API Endpoints:")
	log.Printf("   GET    /health")
	log.Printf("   GET    /todos?page=1&limit=10")
	log.Printf("   GET    /todos/:id")
	log.Printf("   POST   /todos")
	log.Printf("   PUT    /todos/:id")
	log.Printf("   DELETE /todos/:id")
	log.Printf("💡 Headers required: X-API-Key: %s", getEnv("API_SECRET_KEY", "your-secret-key"))

	app.Listen(":" + port)
}

// Middleware สำหรับตรวจสอบ API Key
func APIKeyMiddleware() fiber.Handler {
	expectedKey := getEnv("API_SECRET_KEY", "your-secret-key")

	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"error":   "API Key is required",
				"hint":    "Add header: X-API-Key: your-secret-key",
			})
		}

		if apiKey != expectedKey {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid API Key",
			})
		}

		return c.Next()
	}
}

// GET /todos - ดูรายการ todo ทั้งหมด (พร้อม pagination)
func getTodos(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// รับ query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	// ป้องกันค่าที่ไม่เหมาะสม
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// คำนวณ skip
	skip := (page - 1) * limit

	// สร้าง filter และ options
	filter := bson.M{}
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(skip))
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}}) // เรียงใหม่ -> เก่า

	// ดึงข้อมูล
	cursor, err := todoCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่สามารถดึงข้อมูลได้",
		})
	}
	defer cursor.Close(ctx)

	var todos []Todo
	if err = cursor.All(ctx, &todos); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่สามารถอ่านข้อมูลได้",
		})
	}

	// นับจำนวนทั้งหมด
	total, err := todoCollection.CountDocuments(ctx, filter)
	if err != nil {
		total = 0
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    todos,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GET /todos/:id - ดู todo เดียว
func getTodoByID(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// แปลง string ID เป็น ObjectID
	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "ID ไม่ถูกต้อง",
		})
	}

	var todo Todo
	err = todoCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&todo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"error":   "ไม่พบ todo ที่ระบุ",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "เกิดข้อผิดพลาดในการดึงข้อมูล",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    todo,
	})
}

// POST /todos - เพิ่ม todo ใหม่
func createTodo(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var req TodoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "ข้อมูล JSON ไม่ถูกต้อง",
		})
	}

	// Validate ข้อมูล
	if err := validate.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "ข้อมูลไม่ถูกต้อง",
			"details": err.Error(),
		})
	}

	// สร้าง todo ใหม่
	todo := Todo{
		Title:     req.Title,
		Done:      req.Done,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := todoCollection.InsertOne(ctx, todo)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่สามารถบันทึกข้อมูลได้",
		})
	}

	// เซ็ต ID ที่ได้จากการ insert
	todo.ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "สร้าง todo สำเร็จ",
		"data":    todo,
	})
}

// PUT /todos/:id - แก้ไข todo
func updateTodo(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// แปลง ID
	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "ID ไม่ถูกต้อง",
		})
	}

	var req TodoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "ข้อมูล JSON ไม่ถูกต้อง",
		})
	}

	// Validate ข้อมูล
	if err := validate.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "ข้อมูลไม่ถูกต้อง",
			"details": err.Error(),
		})
	}

	// สร้าง update document
	update := bson.M{
		"$set": bson.M{
			"title":      req.Title,
			"done":       req.Done,
			"updated_at": time.Now(),
		},
	}

	// Update document
	result, err := todoCollection.UpdateOne(ctx, bson.M{"_id": objectId}, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่สามารถอัพเดทข้อมูลได้",
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่พบ todo ที่ระบุ",
		})
	}

	// ดึงข้อมูลที่อัพเดทแล้ว
	var updatedTodo Todo
	err = todoCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&updatedTodo)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่สามารถดึงข้อมูลที่อัพเดทได้",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "อัพเดท todo สำเร็จ",
		"data":    updatedTodo,
	})
}

// DELETE /todos/:id - ลบ todo
func deleteTodo(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// แปลง ID
	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "ID ไม่ถูกต้อง",
		})
	}

	// ลบ document
	result, err := todoCollection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่สามารถลบข้อมูลได้",
		})
	}

	if result.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่พบ todo ที่ระบุ",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "ลบ todo สำเร็จ",
	})
}

// Helper function สำหรับ environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
