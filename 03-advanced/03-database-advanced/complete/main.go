package main

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// Models
type User struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Product struct {
	ID    int     `json:"id" db:"id"`
	Name  string  `json:"name" db:"name"`
	Price float64 `json:"price" db:"price"`
	Stock int     `json:"stock" db:"stock"`
}

type Order struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Total     float64   `json:"total" db:"total"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type OrderItem struct {
	ID        int     `json:"id" db:"id"`
	OrderID   int     `json:"order_id" db:"order_id"`
	ProductID int     `json:"product_id" db:"product_id"`
	Quantity  int     `json:"quantity" db:"quantity"`
	Price     float64 `json:"price" db:"price"`
}

// Repository Interfaces
type UserRepository interface {
	GetAll() ([]User, error)
	GetByID(id int) (*User, error)
	Create(user *User) error
}

type ProductRepository interface {
	GetAll() ([]Product, error)
	GetByID(id int) (*Product, error)
	UpdateStock(id int, newStock int) error
}

type OrderRepository interface {
	Create(order *Order) error
	GetByID(id int) (*Order, error)
	CreateOrderItem(item *OrderItem) error
}

// Repository Implementations
type userRepository struct {
	db *sql.DB
}

type productRepository struct {
	db *sql.DB
}

type orderRepository struct {
	db *sql.DB
}

// User Repository Implementation
func (r *userRepository) GetAll() ([]User, error) {
	rows, err := r.db.Query("SELECT id, email, name, created_at FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepository) GetByID(id int) (*User, error) {
	var user User
	err := r.db.QueryRow("SELECT id, email, name, created_at FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *User) error {
	return r.db.QueryRow(
		"INSERT INTO users (email, name) VALUES ($1, $2) RETURNING id, created_at",
		user.Email, user.Name,
	).Scan(&user.ID, &user.CreatedAt)
}

// Product Repository Implementation
func (r *productRepository) GetAll() ([]Product, error) {
	rows, err := r.db.Query("SELECT id, name, price, stock FROM products ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (r *productRepository) GetByID(id int) (*Product, error) {
	var product Product
	err := r.db.QueryRow("SELECT id, name, price, stock FROM products WHERE id = $1", id).
		Scan(&product.ID, &product.Name, &product.Price, &product.Stock)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) UpdateStock(id int, newStock int) error {
	_, err := r.db.Exec("UPDATE products SET stock = $1 WHERE id = $2", newStock, id)
	return err
}

// Order Repository Implementation
func (r *orderRepository) Create(order *Order) error {
	return r.db.QueryRow(
		"INSERT INTO orders (user_id, total, status) VALUES ($1, $2, $3) RETURNING id, created_at",
		order.UserID, order.Total, order.Status,
	).Scan(&order.ID, &order.CreatedAt)
}

func (r *orderRepository) GetByID(id int) (*Order, error) {
	var order Order
	err := r.db.QueryRow("SELECT id, user_id, total, status, created_at FROM orders WHERE id = $1", id).
		Scan(&order.ID, &order.UserID, &order.Total, &order.Status, &order.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) CreateOrderItem(item *OrderItem) error {
	return r.db.QueryRow(
		"INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4) RETURNING id",
		item.OrderID, item.ProductID, item.Quantity, item.Price,
	).Scan(&item.ID)
}

// Service Layer
type OrderService struct {
	db          *sql.DB
	orderRepo   OrderRepository
	productRepo ProductRepository
}

func NewOrderService(db *sql.DB) *OrderService {
	return &OrderService{
		db:          db,
		orderRepo:   &orderRepository{db: db},
		productRepo: &productRepository{db: db},
	}
}

// Migration System
type Migration struct {
	Version string
	Up      string
	Down    string
}

var migrations = []Migration{
	{
		Version: "001_create_users",
		Up: `CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		)`,
		Down: `DROP TABLE IF EXISTS users`,
	},
	{
		Version: "002_create_products",
		Up: `CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			price DECIMAL(10,2) NOT NULL,
			stock INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP DEFAULT NOW()
		)`,
		Down: `DROP TABLE IF EXISTS products`,
	},
	{
		Version: "003_create_orders",
		Up: `CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			total DECIMAL(10,2) NOT NULL,
			status VARCHAR(50) DEFAULT 'pending',
			created_at TIMESTAMP DEFAULT NOW()
		)`,
		Down: `DROP TABLE IF EXISTS orders`,
	},
	{
		Version: "004_create_order_items",
		Up: `CREATE TABLE IF NOT EXISTS order_items (
			id SERIAL PRIMARY KEY,
			order_id INTEGER REFERENCES orders(id),
			product_id INTEGER REFERENCES products(id),
			quantity INTEGER NOT NULL,
			price DECIMAL(10,2) NOT NULL
		)`,
		Down: `DROP TABLE IF EXISTS order_items`,
	},
}

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

	// เชื่อมต่อ Database
	initDatabase()

	// รัน Migrations
	runMigrations()

	// เตรียมข้อมูลตัวอย่าง
	seedData()

	// สร้าง repositories
	userRepo := &userRepository{db: db}
	productRepo := &productRepository{db: db}
	orderService := NewOrderService(db)

	// Basic routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Database Advanced API - Complete Version",
			"status":  "ok",
		})
	})

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		err := db.Ping()
		if err != nil {
			return c.Status(503).JSON(fiber.Map{
				"status":   "unhealthy",
				"database": "disconnected",
				"error":    err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"status":   "healthy",
			"database": "connected",
			"time":     time.Now(),
		})
	})

	// API routes
	app.Get("/users", getUsersHandler(userRepo))
	app.Post("/users", createUserHandler(userRepo))
	app.Get("/products", getProductsHandler(productRepo))
	app.Post("/orders", createOrderHandler(orderService))
	app.Get("/orders/:id", getOrderHandler(orderService.orderRepo))
	app.Post("/migrate", runMigrationsHandler)

	log.Println("🚀 Database Advanced API started on port 3000")
	log.Println("🗄️ PostgreSQL connected successfully")
	log.Fatal(app.Listen(":3000"))
}

// initDatabase เชื่อมต่อ PostgreSQL
func initDatabase() {
	var err error
	connStr := "user=postgres password=password dbname=ecommerce sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// กำหนด connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("✅ Connected to PostgreSQL")
}

// runMigrations รัน database migrations
func runMigrations() {
	for _, migration := range migrations {
		log.Printf("Running migration: %s", migration.Version)
		_, err := db.Exec(migration.Up)
		if err != nil {
			log.Printf("Migration failed: %v", err)
		} else {
			log.Printf("✅ Migration %s completed", migration.Version)
		}
	}
}

// seedData เพิ่มข้อมูลตัวอย่าง
func seedData() {
	// เพิ่ม users ตัวอย่าง
	db.Exec("INSERT INTO users (email, name) VALUES ($1, $2) ON CONFLICT (email) DO NOTHING",
		"john@example.com", "John Doe")
	db.Exec("INSERT INTO users (email, name) VALUES ($1, $2) ON CONFLICT (email) DO NOTHING",
		"jane@example.com", "Jane Smith")

	// เพิ่ม products ตัวอย่าง
	db.Exec("INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
		"Laptop", 999.99, 10)
	db.Exec("INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
		"Mouse", 29.99, 50)
	db.Exec("INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
		"Keyboard", 79.99, 25)

	log.Println("✅ Sample data seeded")
}

// Handler functions
func getUsersHandler(userRepo UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		users, err := userRepo.GetAll()
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    users,
			"count":   len(users),
		})
	}
}

func createUserHandler(userRepo UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type CreateUserRequest struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		}

		var req CreateUserRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}

		if req.Email == "" || req.Name == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Email and name are required"})
		}

		user := &User{
			Email: req.Email,
			Name:  req.Name,
		}

		if err := userRepo.Create(user); err != nil {
			return err
		}

		return c.Status(201).JSON(fiber.Map{
			"success": true,
			"user":    user,
		})
	}
}

func getProductsHandler(productRepo ProductRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		products, err := productRepo.GetAll()
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    products,
			"count":   len(products),
		})
	}
}

func createOrderHandler(orderService *OrderService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type CreateOrderRequest struct {
			UserID int `json:"user_id"`
			Items  []struct {
				ProductID int `json:"product_id"`
				Quantity  int `json:"quantity"`
			} `json:"items"`
		}

		var req CreateOrderRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}

		// เริ่ม Transaction
		tx, err := orderService.db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback() // Rollback หากมี error

		// คำนวณ total และตรวจสอบ stock
		var total float64
		for _, item := range req.Items {
			// ตรวจสอบ stock และ price
			var stock int
			var price float64
			err := tx.QueryRow("SELECT stock, price FROM products WHERE id = $1 FOR UPDATE",
				item.ProductID).Scan(&stock, &price)
			if err != nil {
				return fiber.NewError(404, "Product not found")
			}

			if stock < item.Quantity {
				return fiber.NewError(400, "Insufficient stock for product")
			}

			// ลด stock
			_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2",
				item.Quantity, item.ProductID)
			if err != nil {
				return err
			}

			total += price * float64(item.Quantity)
		}

		// สร้าง Order
		var orderID int
		var createdAt time.Time
		err = tx.QueryRow("INSERT INTO orders (user_id, total, status) VALUES ($1, $2, $3) RETURNING id, created_at",
			req.UserID, total, "pending").Scan(&orderID, &createdAt)
		if err != nil {
			return err
		}

		// สร้าง Order Items
		for _, item := range req.Items {
			var price float64
			tx.QueryRow("SELECT price FROM products WHERE id = $1", item.ProductID).Scan(&price)

			_, err = tx.Exec("INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)",
				orderID, item.ProductID, item.Quantity, price)
			if err != nil {
				return err
			}
		}

		// Commit Transaction
		if err = tx.Commit(); err != nil {
			return err
		}

		return c.Status(201).JSON(fiber.Map{
			"success":    true,
			"order_id":   orderID,
			"total":      total,
			"status":     "pending",
			"created_at": createdAt,
		})
	}
}

func getOrderHandler(orderRepo OrderRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid order ID"})
		}

		order, err := orderRepo.GetByID(id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Order not found"})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"order":   order,
		})
	}
}

func runMigrationsHandler(c *fiber.Ctx) error {
	runMigrations()
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Migrations completed",
	})
}
