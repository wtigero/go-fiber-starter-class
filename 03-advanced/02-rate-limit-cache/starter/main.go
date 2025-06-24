package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// TODO: สร้าง RateLimiter struct
// type RateLimiter struct {
//     requests map[string][]time.Time  // IP -> list of request times
//     mutex    sync.RWMutex
// }

// TODO: สร้าง Cache struct
// type Cache struct {
//     data   map[string]CacheItem
//     mutex  sync.RWMutex
// }

// TODO: สร้าง CacheItem struct
// type CacheItem struct {
//     Value     interface{}
//     ExpiresAt time.Time
// }

// TODO: สร้าง Stats struct
// type Stats struct {
//     TotalRequests int `json:"total_requests"`
//     CacheHits     int `json:"cache_hits"`
//     CacheMisses   int `json:"cache_misses"`
//     RateLimited   int `json:"rate_limited"`
// }

// ข้อมูลตัวอย่าง
type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var todos = []Todo{
	{ID: 1, Title: "เรียน Rate Limiting", Done: false},
	{ID: 2, Title: "เรียน Caching", Done: false},
	{ID: 3, Title: "ทำแบบฝึกหัด", Done: true},
}

// TODO: สร้างตัวแปร global instances
// var rateLimiter = &RateLimiter{
//     requests: make(map[string][]time.Time),
// }
// var cache = &Cache{
//     data: make(map[string]CacheItem),
// }
// var stats = &Stats{}

func main() {
	app := fiber.New()

	app.Use(logger.New())

	// Basic endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Rate Limiting & Cache API - Ready to code!",
			"status":  "ok",
		})
	})

	// TODO: เพิ่ม middleware สำหรับ rate limiting และ caching
	// app.Use("/todos", rateLimitMiddleware())
	// app.Use("/todos", cacheMiddleware())

	// Endpoints
	app.Get("/todos", getTodosHandler)
	app.Get("/stats", getStatsHandler)
	app.Get("/cache/clear", clearCacheHandler)

	log.Println("🚀 Server started on port 3000")
	log.Fatal(app.Listen(":3000"))
}

// getTodosHandler ส่งรายการ todos (จะใช้เวลาจำลองการประมวลผล)
func getTodosHandler(c *fiber.Ctx) error {
	// จำลองการประมวลผลที่ใช้เวลา
	time.Sleep(100 * time.Millisecond)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    todos,
		"count":   len(todos),
		"cached":  false, // TODO: เปลี่ยนเป็น true ถ้ามาจาก cache
	})
}

// getStatsHandler ส่งสถิติการใช้งาน
func getStatsHandler(c *fiber.Ctx) error {
	// TODO: ส่งข้อมูลสถิติจริง
	return c.JSON(fiber.Map{
		"total_requests": 0,
		"cache_hits":     0,
		"cache_misses":   0,
		"rate_limited":   0,
	})
}

// clearCacheHandler ล้าง cache
func clearCacheHandler(c *fiber.Ctx) error {
	// TODO: ล้าง cache จริง
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Cache cleared",
	})
}

// TODO: สร้าง rateLimitMiddleware
// func rateLimitMiddleware() fiber.Handler {
//     return func(c *fiber.Ctx) error {
//         ip := c.IP()
//
//         // ตรวจสอบ rate limit
//         if !rateLimiter.Allow(ip) {
//             stats.RateLimited++
//             return c.Status(429).JSON(fiber.Map{
//                 "error": "Rate limit exceeded. Try again later.",
//                 "retry_after": "60s",
//             })
//         }
//
//         stats.TotalRequests++
//         return c.Next()
//     }
// }

// TODO: สร้าง cacheMiddleware
// func cacheMiddleware() fiber.Handler {
//     return func(c *fiber.Ctx) error {
//         cacheKey := c.Route().Path + c.OriginalURL()
//
//         // ตรวจสอบใน cache ก่อน
//         if cachedData := cache.Get(cacheKey); cachedData != nil {
//             stats.CacheHits++
//             return c.JSON(cachedData)
//         }
//
//         stats.CacheMisses++
//         return c.Next() // ดำเนินการต่อแล้วเก็บ cache
//     }
// }
