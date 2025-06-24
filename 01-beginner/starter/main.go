package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	// สร้าง Fiber app ใหม่
	app := fiber.New()

	// Route แรก - Hello World
	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello from Go Fiber!",
		})
	})

	// TODO: เพิ่ม middleware logger
	// TODO: เพิ่ม middleware recover

	// TODO: สร้าง struct Todo
	// type Todo struct {
	//     ID    int    `json:"id"`
	//     Title string `json:"title"`
	//     Done  bool   `json:"done"`
	// }

	// TODO: สร้าง slice สำหรับเก็บ todos
	// var todos []Todo
	// var nextID = 1

	// TODO: GET /todos - แสดงรายการ todo ทั้งหมด

	// TODO: GET /todos/:id - แสดง todo เดียว

	// TODO: POST /todos - เพิ่ม todo ใหม่

	// รัน server ที่ port 3000
	app.Listen(":3000")
}
