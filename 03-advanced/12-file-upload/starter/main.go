package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// TODO: สร้าง file upload handler
// TODO: ตรวจสอบประเภทไฟล์
// TODO: จำกัดขนาดไฟล์
// TODO: สร้างชื่อไฟล์ใหม่

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 10MB limit
	})

	app.Use(logger.New())

	// Static files
	app.Static("/files", "./uploads")

	// TODO: เพิ่ม upload routes

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "File Upload Starter",
		})
	})

	log.Fatal(app.Listen(":3000"))
}
