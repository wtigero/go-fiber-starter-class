package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// โครงสร้างข้อมูล Todo
type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// ตัวแปรสำหรับเก็บข้อมูล todos ใน memory
var todos []Todo
var nextID = 1

func main() {
	// สร้าง Fiber app ใหม่
	app := fiber.New(fiber.Config{})

	// เพิ่ม middleware สำหรับ log การเข้าถึง
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	// เพิ่ม middleware สำหรับจัดการ panic
	app.Use(recover.New())

	// เตรียมข้อมูลตัวอย่าง
	initSampleData()

	// === ROUTES ===

	// Route แรก - Hello World
	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello from Go Fiber Todo API!",
			"version": "1.0.0",
		})
	})

	// GET /todos - แสดงรายการ todo ทั้งหมด
	app.Get("/todos", getTodos)

	// GET /todos/:id - แสดง todo เดียวตาม ID
	app.Get("/todos/:id", getTodoByID)

	// POST /todos - เพิ่ม todo ใหม่
	app.Post("/todos", createTodo)

	// รัน server ที่ port 3000
	app.Listen(":3000")
}

// ฟังก์ชันสำหรับเตรียมข้อมูลตัวอย่าง
func initSampleData() {
	todos = []Todo{
		{ID: 1, Title: "เรียน Go Programming", Done: false},
		{ID: 2, Title: "ทำโปรเจกต์ Todo API", Done: true},
		{ID: 3, Title: "ทบทวน Fiber Framework", Done: false},
	}
	nextID = 4 // ID ถัดไปที่จะใช้
}

// GET /todos - คืนรายการ todo ทั้งหมด
func getTodos(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":    todos,
		"count":   len(todos),
	})
}

// GET /todos/:id - คืน todo เดียวตาม ID
func getTodoByID(c *fiber.Ctx) error {
	// รับ ID จาก URL parameter
	idParam := c.Params("id")

	// แปลง string เป็น int
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "ID ต้องเป็นตัวเลข",
		})
	}

	// ค้นหา todo ที่ตรงกับ ID
	for _, todo := range todos {
		if todo.ID == id {
			return c.JSON(fiber.Map{
				"success": true,
				"data":    todo,
			})
		}
	}

	// หากไม่พบ todo
	return c.Status(404).JSON(fiber.Map{
		"success": false,
		"message": "ไม่พบ todo ที่ระบุ",
	})
}

// POST /todos - เพิ่ม todo ใหม่
func createTodo(c *fiber.Ctx) error {
	// โครงสร้างสำหรับรับข้อมูลจาก client
	type CreateTodoRequest struct {
		Title string `json:"title"`
		Done  bool   `json:"done"`
	}

	var req CreateTodoRequest

	// แปลง JSON จาก request body เป็น struct
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "ข้อมูล JSON ไม่ถูกต้อง",
		})
	}

	// ตรวจสอบว่ามี title หรือไม่
	if req.Title == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "กรุณาระบุ title",
		})
	}

	// สร้าง todo ใหม่
	newTodo := Todo{
		ID:    nextID,
		Title: req.Title,
		Done:  req.Done,
	}

	// เพิ่มเข้าไปใน slice
	todos = append(todos, newTodo)
	nextID++ // เพิ่ม ID สำหรับครั้งถัดไป

	// ส่งผลลัพธ์กลับ
	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "สร้าง todo สำเร็จ",
		"data":    newTodo,
	})
}
