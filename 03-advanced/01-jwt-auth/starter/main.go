package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// TODO: สร้าง User struct
// type User struct {
//     ID       int    `json:"id"`
//     Email    string `json:"email"`
//     Password string `json:"-"`
//     Name     string `json:"name"`
// }

// TODO: สร้าง Todo struct
// type Todo struct {
//     ID     int    `json:"id"`
//     Title  string `json:"title"`
//     Done   bool   `json:"done"`
//     UserID int    `json:"user_id"`
// }

// TODO: สร้างตัวแปรเก็บข้อมูล (จำลอง database)
// var users []User
// var todos []Todo
// var userID = 1
// var todoID = 1

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// สำหรับทดสอบว่า server ทำงาน
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "JWT Auth API - Ready to code!",
			"status":  "ok",
		})
	})

	// TODO: สร้าง API routes
	// app.Post("/register", registerHandler)
	// app.Post("/login", loginHandler)

	// TODO: สร้าง protected routes (ต้อง login ก่อน)
	// protected := app.Group("/", RequireAuth())
	// protected.Get("/profile", getProfileHandler)
	// protected.Get("/todos", getTodosHandler)

	log.Fatal(app.Listen(":3000"))
}

// TODO: สร้าง registerHandler
// func registerHandler(c *fiber.Ctx) error {
//     // 1. รับข้อมูลจาก request body
//     // 2. ตรวจสอบว่า email ซ้ำหรือไม่
//     // 3. hash password
//     // 4. บันทึกข้อมูล user
//     // 5. ส่ง response กลับ
// }

// TODO: สร้าง loginHandler
// func loginHandler(c *fiber.Ctx) error {
//     // 1. รับ email, password จาก request
//     // 2. หา user จาก email
//     // 3. เปรียบเทียบ password
//     // 4. สร้าง JWT token
//     // 5. ส่ง token กลับ
// }

// TODO: สร้าง RequireAuth middleware
// func RequireAuth() fiber.Handler {
//     return func(c *fiber.Ctx) error {
//         // 1. ดึง Authorization header
//         // 2. ตรวจสอบ token
//         // 3. เก็บ user info ไว้ใน context
//         // 4. เรียก c.Next()
//     }
// }

// TODO: สร้าง getProfileHandler
// func getProfileHandler(c *fiber.Ctx) error {
//     // ส่งข้อมูล user ที่ login อยู่
// }

// TODO: สร้าง getTodosHandler
// func getTodosHandler(c *fiber.Ctx) error {
//     // ส่ง todos ของ user ที่ login อยู่เท่านั้น
// }
