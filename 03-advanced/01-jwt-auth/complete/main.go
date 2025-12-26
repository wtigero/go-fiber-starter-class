package main

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// User struct สำหรับเก็บข้อมูลผู้ใช้
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"` // ไม่แสดงใน JSON response
	Name     string `json:"name"`
}

// Todo struct สำหรับเก็บข้อมูล todos
type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Done   bool   `json:"done"`
	UserID int    `json:"user_id"`
}

// LoginRequest สำหรับรับข้อมูล login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest สำหรับรับข้อมูล register
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// ตัวแปรเก็บข้อมูล (จำลอง database)
var users []User
var todos []Todo
var userID = 1
var todoID = 1

// JWT Secret Key (ในการใช้งานจริงควรเก็บใน environment variable)
const jwtSecret = "your-secret-key-keep-it-safe"

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

	// เตรียมข้อมูลตัวอย่าง
	setupSampleData()

	// Basic routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "JWT Auth API - Complete Version",
			"status":  "ok",
		})
	})

	// Public routes (ไม่ต้อง login)
	app.Post("/register", registerHandler)
	app.Post("/login", loginHandler)

	// Protected routes (ต้อง login ก่อน)
	protected := app.Group("/", RequireAuth())
	protected.Get("/profile", getProfileHandler)
	protected.Get("/todos", getTodosHandler)

	log.Println("🚀 Server started on port 3000")
	log.Fatal(app.Listen(":3000"))
}

// setupSampleData เตรียมข้อมูลตัวอย่าง
func setupSampleData() {
	// สร้าง user ตัวอย่าง
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	users = append(users, User{
		ID:       1,
		Email:    "demo@example.com",
		Password: string(hashedPassword),
		Name:     "นาย Demo",
	})

	// สร้าง todos ตัวอย่าง
	todos = append(todos, Todo{
		ID:     1,
		Title:  "เรียน JWT Authentication",
		Done:   false,
		UserID: 1,
	})

	userID = 2
	todoID = 2
}

// registerHandler จัดการการสมัครสมาชิก
func registerHandler(c *fiber.Ctx) error {
	var req RegisterRequest

	// รับข้อมูลจาก request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "ข้อมูลไม่ถูกต้อง",
		})
	}

	// ตรวจสอบข้อมูลพื้นฐาน
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "กรุณากรอกข้อมูลให้ครบถ้วน",
		})
	}

	// ตรวจสอบว่า email ซ้ำหรือไม่
	for _, user := range users {
		if user.Email == req.Email {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"error":   "Email นี้ถูกใช้แล้ว",
			})
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "เกิดข้อผิดพลาดในการเข้ารหัส password",
		})
	}

	// สร้าง user ใหม่
	newUser := User{
		ID:       userID,
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
	}

	// บันทึกลง "database"
	users = append(users, newUser)
	userID++

	// ส่ง response (ไม่ส่ง password)
	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "สมัครสมาชิกสำเร็จ",
		"user": fiber.Map{
			"id":    newUser.ID,
			"email": newUser.Email,
			"name":  newUser.Name,
		},
	})
}

// loginHandler จัดการการเข้าสู่ระบบ
func loginHandler(c *fiber.Ctx) error {
	var req LoginRequest

	// รับข้อมูลจาก request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "ข้อมูลไม่ถูกต้อง",
		})
	}

	// หา user จาก email
	var user *User
	for _, u := range users {
		if u.Email == req.Email {
			user = &u
			break
		}
	}

	if user == nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"error":   "Email หรือ Password ไม่ถูกต้อง",
		})
	}

	// เปรียบเทียบ password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"error":   "Email หรือ Password ไม่ถูกต้อง",
		})
	}

	// สร้าง JWT token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token หมดอายุใน 24 ชั่วโมง
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "เกิดข้อผิดพลาดในการสร้าง token",
		})
	}

	// ส่ง response พร้อม token
	return c.JSON(fiber.Map{
		"success": true,
		"message": "เข้าสู่ระบบสำเร็จ",
		"token":   tokenString,
		"user": fiber.Map{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

// RequireAuth middleware สำหรับตรวจสอบ JWT token
func RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ดึง Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"error":   "ต้องใส่ Authorization header",
			})
		}

		// ตรวจสอบว่าขึ้นต้นด้วย "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"error":   "Authorization header ต้องขึ้นต้นด้วย 'Bearer '",
			})
		}

		// ดึง token (ตัดคำว่า "Bearer " ออก)
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// ตรวจสอบ token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"error":   "Token ไม่ถูกต้องหรือหมดอายุ",
			})
		}

		// ดึงข้อมูลจาก token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"error":   "Token format ไม่ถูกต้อง",
			})
		}

		// เก็บข้อมูล user ใน context
		c.Locals("user_id", int(claims["user_id"].(float64)))
		c.Locals("email", claims["email"].(string))
		c.Locals("name", claims["name"].(string))

		return c.Next()
	}
}

// getProfileHandler ส่งข้อมูลโปรไฟล์ของ user ที่ login อยู่
func getProfileHandler(c *fiber.Ctx) error {
	// ดึงข้อมูล user จาก context (มาจาก middleware)
	userIDInterface := c.Locals("user_id")
	email := c.Locals("email").(string)
	name := c.Locals("name").(string)

	userID, _ := userIDInterface.(int)

	return c.JSON(fiber.Map{
		"success": true,
		"user": fiber.Map{
			"id":    userID,
			"email": email,
			"name":  name,
		},
	})
}

// getTodosHandler ส่ง todos ของ user ที่ login อยู่เท่านั้น
func getTodosHandler(c *fiber.Ctx) error {
	// ดึง user_id จาก context
	userIDInterface := c.Locals("user_id")
	userID, _ := userIDInterface.(int)

	// หา todos ของ user นี้เท่านั้น
	var userTodos []Todo
	for _, todo := range todos {
		if todo.UserID == userID {
			userTodos = append(userTodos, todo)
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    userTodos,
		"count":   len(userTodos),
	})
}
