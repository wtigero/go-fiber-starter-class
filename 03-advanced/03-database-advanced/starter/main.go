package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// TODO: สร้าง Models
// type User struct {
//     ID    int    `json:"id" db:"id"`
//     Email string `json:"email" db:"email"`
//     Name  string `json:"name" db:"name"`
// }

// type Product struct {
//     ID    int     `json:"id" db:"id"`
//     Name  string  `json:"name" db:"name"`
//     Price float64 `json:"price" db:"price"`
//     Stock int     `json:"stock" db:"stock"`
// }

// type Order struct {
//     ID       int     `json:"id" db:"id"`
//     UserID   int     `json:"user_id" db:"user_id"`
//     Total    float64 `json:"total" db:"total"`
//     Status   string  `json:"status" db:"status"`
// }

// TODO: สร้าง Repository Interfaces
// type UserRepository interface {
//     GetAll() ([]User, error)
//     GetByID(id int) (*User, error)
//     Create(user *User) error
// }

// type ProductRepository interface {
//     GetAll() ([]Product, error)
//     GetByID(id int) (*Product, error)
//     UpdateStock(id int, newStock int) error
// }

// type OrderRepository interface {
//     Create(order *Order) error
//     GetByID(id int) (*Order, error)
// }

// TODO: สร้าง Service Layer
// type OrderService struct {
//     db         *sql.DB
//     orderRepo  OrderRepository
//     productRepo ProductRepository
// }

// TODO: สร้าง Migration System
// type Migration struct {
//     Version string
//     Up      string
//     Down    string
// }

var db *sql.DB

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		},
	})

	app.Use(logger.New())

	// TODO: เชื่อมต่อ Database
	// initDatabase()

	// TODO: รัน Migrations
	// runMigrations()

	// Basic routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Database Advanced API - Ready to code!",
			"status":  "ok",
		})
	})

	// TODO: สร้าง API routes
	// app.Get("/users", getUsersHandler)
	// app.Post("/orders", createOrderHandler)
	// app.Get("/orders/:id", getOrderHandler)
	// app.Post("/migrate", runMigrationsHandler)

	log.Println("🚀 Server started on port 3000")
	log.Fatal(app.Listen(":3000"))
}

// TODO: initDatabase เชื่อมต่อ PostgreSQL
// func initDatabase() {
//     var err error
//     connStr := "user=postgres password=password dbname=ecommerce sslmode=disable"
//     db, err = sql.Open("postgres", connStr)
//     if err != nil {
//         log.Fatal("Failed to connect to database:", err)
//     }
//
//     if err = db.Ping(); err != nil {
//         log.Fatal("Failed to ping database:", err)
//     }
//
//     log.Println("✅ Connected to PostgreSQL")
// }

// TODO: runMigrations รัน database migrations
// func runMigrations() {
//     migrations := []Migration{
//         {
//             Version: "001_create_users",
//             Up: `CREATE TABLE IF NOT EXISTS users (
//                 id SERIAL PRIMARY KEY,
//                 email VARCHAR(255) UNIQUE NOT NULL,
//                 name VARCHAR(255) NOT NULL,
//                 created_at TIMESTAMP DEFAULT NOW()
//             )`,
//         },
//         // เพิ่ม migrations อื่นๆ
//     }
//
//     for _, migration := range migrations {
//         log.Printf("Running migration: %s", migration.Version)
//         _, err := db.Exec(migration.Up)
//         if err != nil {
//             log.Printf("Migration failed: %v", err)
//         }
//     }
// }

// TODO: getUsersHandler ดูรายการ users
// func getUsersHandler(c *fiber.Ctx) error {
//     rows, err := db.Query("SELECT id, email, name FROM users")
//     if err != nil {
//         return err
//     }
//     defer rows.Close()
//
//     var users []User
//     for rows.Next() {
//         var user User
//         err := rows.Scan(&user.ID, &user.Email, &user.Name)
//         if err != nil {
//             return err
//         }
//         users = append(users, user)
//     }
//
//     return c.JSON(fiber.Map{
//         "success": true,
//         "data": users,
//     })
// }

// TODO: createOrderHandler สร้างออเดอร์ (ใช้ Transaction)
// func createOrderHandler(c *fiber.Ctx) error {
//     type CreateOrderRequest struct {
//         UserID int `json:"user_id"`
//         Items  []struct {
//             ProductID int `json:"product_id"`
//             Quantity  int `json:"quantity"`
//         } `json:"items"`
//     }
//
//     var req CreateOrderRequest
//     if err := c.BodyParser(&req); err != nil {
//         return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
//     }
//
//     // เริ่ม Transaction
//     tx, err := db.Begin()
//     if err != nil {
//         return err
//     }
//     defer tx.Rollback() // Rollback หากมี error
//
//     // 1. ตรวจสอบและลด Stock
//     var total float64
//     for _, item := range req.Items {
//         // ตรวจสอบ stock
//         var stock int
//         var price float64
//         err := tx.QueryRow("SELECT stock, price FROM products WHERE id = $1", item.ProductID).Scan(&stock, &price)
//         if err != nil {
//             return fiber.NewError(404, "Product not found")
//         }
//
//         if stock < item.Quantity {
//             return fiber.NewError(400, "Insufficient stock")
//         }
//
//         // ลด stock
//         _, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
//         if err != nil {
//             return err
//         }
//
//         total += price * float64(item.Quantity)
//     }
//
//     // 2. สร้าง Order
//     var orderID int
//     err = tx.QueryRow("INSERT INTO orders (user_id, total, status) VALUES ($1, $2, $3) RETURNING id",
//         req.UserID, total, "pending").Scan(&orderID)
//     if err != nil {
//         return err
//     }
//
//     // 3. Commit Transaction
//     if err = tx.Commit(); err != nil {
//         return err
//     }
//
//     return c.Status(201).JSON(fiber.Map{
//         "success": true,
//         "order_id": orderID,
//         "total": total,
//     })
// }
