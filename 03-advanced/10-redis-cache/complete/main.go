package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// ============ Models ============

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Name     string `json:"name"`
}

type Session struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// ============ Global Variables ============

var (
	rdb      *redis.Client
	ctx      = context.Background()
	products = make(map[string]Product) // จำลอง database
	users    = make(map[string]User)    // จำลอง database
)

// Cache Statistics
var cacheStats struct {
	Hits       int64
	Misses     int64
	RateLimited int64
}

// ============ Redis Keys ============

const (
	ProductCachePrefix  = "product:"
	ProductListCacheKey = "products:all"
	SessionPrefix       = "session:"
	RateLimitPrefix     = "ratelimit:"
	CacheTTL            = 5 * time.Minute
	SessionTTL          = 24 * time.Hour
	RateLimitWindow     = 1 * time.Minute
	RateLimitMax        = 100
)

func main() {
	// เชื่อมต่อ Redis
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 10,
	})

	// ทดสอบการเชื่อมต่อ
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatal("❌ ไม่สามารถเชื่อมต่อ Redis:", err)
	}
	log.Println("✅ เชื่อมต่อ Redis สำเร็จ")

	// สร้างข้อมูลตัวอย่าง
	setupSampleData()

	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})

	app.Use(logger.New())
	app.Use(cors.New())

	// Public routes
	app.Get("/", homeHandler)
	app.Get("/stats", statsHandler)
	app.Post("/cache/clear", clearCacheHandler)

	// Auth routes
	app.Post("/register", registerHandler)
	app.Post("/login", loginHandler)

	// Product routes (with caching & rate limiting)
	products := app.Group("/products")
	products.Use(RateLimitMiddleware())
	products.Get("/", CacheMiddleware(ProductListCacheKey, CacheTTL), listProductsHandler)
	products.Get("/:id", getProductHandler)
	products.Post("/", createProductHandler)
	products.Put("/:id", updateProductHandler)
	products.Delete("/:id", deleteProductHandler)

	// Protected routes (require session)
	protected := app.Group("/me")
	protected.Use(SessionMiddleware())
	protected.Get("/profile", profileHandler)
	protected.Post("/logout", logoutHandler)

	log.Println("🚀 Server running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}

// ============ Middleware ============

// RateLimitMiddleware - Sliding window rate limiting with Redis
func RateLimitMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		key := RateLimitPrefix + ip

		// ใช้ Lua script สำหรับ atomic operation
		script := redis.NewScript(`
			local key = KEYS[1]
			local limit = tonumber(ARGV[1])
			local window = tonumber(ARGV[2])
			local now = tonumber(ARGV[3])

			-- ลบ requests เก่า
			redis.call('ZREMRANGEBYSCORE', key, 0, now - window)

			-- นับ requests ปัจจุบัน
			local count = redis.call('ZCARD', key)

			if count < limit then
				-- เพิ่ม request ใหม่
				redis.call('ZADD', key, now, now .. '-' .. math.random())
				redis.call('EXPIRE', key, window / 1000)
				return count + 1
			else
				return -1
			end
		`)

		now := time.Now().UnixMilli()
		result, err := script.Run(ctx, rdb, []string{key}, RateLimitMax, RateLimitWindow.Milliseconds(), now).Int64()

		if err != nil {
			log.Printf("Rate limit error: %v", err)
			return c.Next()
		}

		// Set headers
		c.Set("X-RateLimit-Limit", strconv.Itoa(RateLimitMax))
		c.Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(RateLimitWindow).Unix(), 10))

		if result == -1 {
			atomic.AddInt64(&cacheStats.RateLimited, 1)
			c.Set("X-RateLimit-Remaining", "0")
			return c.Status(429).JSON(fiber.Map{
				"success": false,
				"error":   "Rate limit exceeded. Try again later.",
			})
		}

		c.Set("X-RateLimit-Remaining", strconv.FormatInt(RateLimitMax-result, 10))
		return c.Next()
	}
}

// CacheMiddleware - Cache response in Redis
func CacheMiddleware(cacheKey string, ttl time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// พยายามอ่านจาก cache
		cached, err := rdb.Get(ctx, cacheKey).Result()
		if err == nil {
			atomic.AddInt64(&cacheStats.Hits, 1)
			c.Set("X-Cache", "HIT")
			c.Set("Content-Type", "application/json")
			return c.SendString(cached)
		}

		atomic.AddInt64(&cacheStats.Misses, 1)
		c.Set("X-Cache", "MISS")

		// ดำเนินการต่อและ cache response
		if err := c.Next(); err != nil {
			return err
		}

		// Cache response body
		body := c.Response().Body()
		if len(body) > 0 && c.Response().StatusCode() == 200 {
			rdb.Set(ctx, cacheKey, body, ttl)
		}

		return nil
	}
}

// SessionMiddleware - Validate session from Redis
func SessionMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := c.Cookies("session_id")
		if sessionID == "" {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"error":   "กรุณาเข้าสู่ระบบ",
			})
		}

		// ดึง session จาก Redis
		sessionData, err := rdb.Get(ctx, SessionPrefix+sessionID).Result()
		if err == redis.Nil {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"error":   "Session หมดอายุ กรุณาเข้าสู่ระบบใหม่",
			})
		}
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"error":   "เกิดข้อผิดพลาด",
			})
		}

		var session Session
		if err := json.Unmarshal([]byte(sessionData), &session); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"error":   "Session ไม่ถูกต้อง",
			})
		}

		// เก็บข้อมูล session ใน context
		c.Locals("session", session)
		c.Locals("user_id", session.UserID)

		return c.Next()
	}
}

// ============ Handlers ============

func homeHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Redis Cache & Session API",
		"version": "1.0.0",
		"endpoints": fiber.Map{
			"products":    "GET/POST /products",
			"auth":        "POST /register, /login",
			"protected":   "GET /me/profile (requires session)",
			"stats":       "GET /stats",
			"clear_cache": "POST /cache/clear",
		},
	})
}

func statsHandler(c *fiber.Ctx) error {
	info, _ := rdb.Info(ctx, "stats").Result()
	dbSize, _ := rdb.DBSize(ctx).Result()

	return c.JSON(fiber.Map{
		"cache": fiber.Map{
			"hits":         atomic.LoadInt64(&cacheStats.Hits),
			"misses":       atomic.LoadInt64(&cacheStats.Misses),
			"hit_rate":     calculateHitRate(),
			"rate_limited": atomic.LoadInt64(&cacheStats.RateLimited),
		},
		"redis": fiber.Map{
			"connected": true,
			"db_size":   dbSize,
			"info":      info,
		},
	})
}

func clearCacheHandler(c *fiber.Ctx) error {
	// ลบเฉพาะ cache keys (ไม่ลบ sessions)
	keys, _ := rdb.Keys(ctx, "product*").Result()
	if len(keys) > 0 {
		rdb.Del(ctx, keys...)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("ลบ cache %d keys", len(keys)),
	})
}

// ============ Product Handlers ============

func listProductsHandler(c *fiber.Ctx) error {
	// แปลง map เป็น slice
	productList := make([]Product, 0, len(products))
	for _, p := range products {
		productList = append(productList, p)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    productList,
		"count":   len(productList),
	})
}

func getProductHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	cacheKey := ProductCachePrefix + id

	// ลองอ่านจาก cache
	cached, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		atomic.AddInt64(&cacheStats.Hits, 1)
		c.Set("X-Cache", "HIT")

		var product Product
		json.Unmarshal([]byte(cached), &product)
		return c.JSON(fiber.Map{
			"success": true,
			"data":    product,
		})
	}

	atomic.AddInt64(&cacheStats.Misses, 1)
	c.Set("X-Cache", "MISS")

	// อ่านจาก "database"
	product, exists := products[id]
	if !exists {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่พบสินค้า",
		})
	}

	// Cache ผลลัพธ์
	productJSON, _ := json.Marshal(product)
	rdb.Set(ctx, cacheKey, productJSON, CacheTTL)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    product,
	})
}

func createProductHandler(c *fiber.Ctx) error {
	var product Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "ข้อมูลไม่ถูกต้อง",
		})
	}

	product.ID = uuid.New().String()
	product.CreatedAt = time.Now()
	products[product.ID] = product

	// Invalidate cache
	rdb.Del(ctx, ProductListCacheKey)

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"data":    product,
	})
}

func updateProductHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	if _, exists := products[id]; !exists {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่พบสินค้า",
		})
	}

	var product Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "ข้อมูลไม่ถูกต้อง",
		})
	}

	product.ID = id
	products[id] = product

	// Invalidate cache
	rdb.Del(ctx, ProductCachePrefix+id, ProductListCacheKey)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    product,
	})
}

func deleteProductHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	if _, exists := products[id]; !exists {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่พบสินค้า",
		})
	}

	delete(products, id)

	// Invalidate cache
	rdb.Del(ctx, ProductCachePrefix+id, ProductListCacheKey)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "ลบสินค้าสำเร็จ",
	})
}

// ============ Auth Handlers ============

func registerHandler(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "ข้อมูลไม่ถูกต้อง",
		})
	}

	// ตรวจสอบ email ซ้ำ
	for _, u := range users {
		if u.Email == req.Email {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"error":   "Email นี้ถูกใช้แล้ว",
			})
		}
	}

	user := User{
		ID:       uuid.New().String(),
		Email:    req.Email,
		Password: req.Password, // ในจริงต้อง hash
		Name:     req.Name,
	}
	users[user.ID] = user

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "สมัครสมาชิกสำเร็จ",
		"user": fiber.Map{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

func loginHandler(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "ข้อมูลไม่ถูกต้อง",
		})
	}

	// หา user
	var user *User
	for _, u := range users {
		if u.Email == req.Email && u.Password == req.Password {
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

	// สร้าง session
	sessionID := uuid.New().String()
	session := Session{
		UserID:    user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(SessionTTL),
	}

	sessionJSON, _ := json.Marshal(session)
	rdb.Set(ctx, SessionPrefix+sessionID, sessionJSON, SessionTTL)

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  session.ExpiresAt,
		HTTPOnly: true,
		Secure:   false, // true ใน production
		SameSite: "Lax",
	})

	return c.JSON(fiber.Map{
		"success": true,
		"message": "เข้าสู่ระบบสำเร็จ",
		"user": fiber.Map{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

func profileHandler(c *fiber.Ctx) error {
	session := c.Locals("session").(Session)

	return c.JSON(fiber.Map{
		"success": true,
		"user": fiber.Map{
			"id":    session.UserID,
			"email": session.Email,
			"name":  session.Name,
		},
		"session": fiber.Map{
			"created_at": session.CreatedAt,
			"expires_at": session.ExpiresAt,
		},
	})
}

func logoutHandler(c *fiber.Ctx) error {
	sessionID := c.Cookies("session_id")
	if sessionID != "" {
		rdb.Del(ctx, SessionPrefix+sessionID)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"message": "ออกจากระบบสำเร็จ",
	})
}

// ============ Helpers ============

func setupSampleData() {
	// สร้าง products ตัวอย่าง
	sampleProducts := []Product{
		{ID: "1", Name: "iPhone 15 Pro", Price: 48900, Description: "สมาร์ทโฟนรุ่นใหม่ล่าสุด", CreatedAt: time.Now()},
		{ID: "2", Name: "MacBook Pro M3", Price: 89900, Description: "แล็ปท็อปสำหรับมืออาชีพ", CreatedAt: time.Now()},
		{ID: "3", Name: "AirPods Pro 2", Price: 8990, Description: "หูฟังไร้สายระดับพรีเมียม", CreatedAt: time.Now()},
	}
	for _, p := range sampleProducts {
		products[p.ID] = p
	}

	// สร้าง user ตัวอย่าง
	users["1"] = User{ID: "1", Email: "demo@example.com", Password: "123456", Name: "Demo User"}

	log.Printf("📦 โหลดข้อมูลตัวอย่าง: %d products, %d users", len(products), len(users))
}

func errorHandler(c *fiber.Ctx, err error) error {
	return c.Status(500).JSON(fiber.Map{
		"success": false,
		"error":   err.Error(),
	})
}

func calculateHitRate() string {
	hits := atomic.LoadInt64(&cacheStats.Hits)
	misses := atomic.LoadInt64(&cacheStats.Misses)
	total := hits + misses
	if total == 0 {
		return "0%"
	}
	return fmt.Sprintf("%.1f%%", float64(hits)/float64(total)*100)
}
