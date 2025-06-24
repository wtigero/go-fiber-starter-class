package main

import (
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// RateLimiter จัดการการจำกัดอัตราการเรียก API
type RateLimiter struct {
	requests map[string][]time.Time // IP -> list of request times
	mutex    sync.RWMutex
	limit    int           // จำนวนครั้งที่อนุญาต
	window   time.Duration // ช่วงเวลา
}

// NewRateLimiter สร้าง RateLimiter ใหม่
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// Allow ตรวจสอบว่า IP นี้สามารถเรียก API ได้หรือไม่
func (r *RateLimiter) Allow(ip string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	now := time.Now()

	// ดึงรายการ requests ของ IP นี้
	requests, exists := r.requests[ip]
	if !exists {
		r.requests[ip] = []time.Time{now}
		return true
	}

	// กรองเฉพาะ requests ในช่วง window
	var validRequests []time.Time
	for _, reqTime := range requests {
		if now.Sub(reqTime) < r.window {
			validRequests = append(validRequests, reqTime)
		}
	}

	// ตรวจสอบว่าเกิน limit หรือไม่
	if len(validRequests) >= r.limit {
		r.requests[ip] = validRequests
		return false
	}

	// เพิ่ม request ใหม่
	validRequests = append(validRequests, now)
	r.requests[ip] = validRequests
	return true
}

// GetRemainingRequests ส่งจำนวน requests ที่เหลือ
func (r *RateLimiter) GetRemainingRequests(ip string) int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	requests, exists := r.requests[ip]
	if !exists {
		return r.limit
	}

	now := time.Now()
	validCount := 0
	for _, reqTime := range requests {
		if now.Sub(reqTime) < r.window {
			validCount++
		}
	}

	remaining := r.limit - validCount
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Cache จัดการ cache ใน memory
type Cache struct {
	data  map[string]CacheItem
	mutex sync.RWMutex
}

// CacheItem รายการใน cache
type CacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}

// NewCache สร้าง Cache ใหม่
func NewCache() *Cache {
	cache := &Cache{
		data: make(map[string]CacheItem),
	}

	// เริ่ม cleanup goroutine
	go cache.cleanup()
	return cache
}

// Get ดึงข้อมูลจาก cache
func (c *Cache) Get(key string) interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return nil
	}

	if time.Now().After(item.ExpiresAt) {
		return nil // หมดอายุแล้ว
	}

	return item.Value
}

// Set เก็บข้อมูลลง cache
func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(duration),
	}
}

// Clear ล้าง cache ทั้งหมด
func (c *Cache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data = make(map[string]CacheItem)
}

// cleanup ลบข้อมูลที่หมดอายุ
func (c *Cache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()
		for key, item := range c.data {
			if now.After(item.ExpiresAt) {
				delete(c.data, key)
			}
		}
		c.mutex.Unlock()
	}
}

// Stats เก็บสถิติการใช้งาน
type Stats struct {
	TotalRequests int `json:"total_requests"`
	CacheHits     int `json:"cache_hits"`
	CacheMisses   int `json:"cache_misses"`
	RateLimited   int `json:"rate_limited"`
	mutex         sync.RWMutex
}

// IncrementTotal เพิ่มจำนวน total requests
func (s *Stats) IncrementTotal() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.TotalRequests++
}

// IncrementCacheHits เพิ่มจำนวน cache hits
func (s *Stats) IncrementCacheHits() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.CacheHits++
}

// IncrementCacheMisses เพิ่มจำนวน cache misses
func (s *Stats) IncrementCacheMisses() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.CacheMisses++
}

// IncrementRateLimited เพิ่มจำนวน rate limited
func (s *Stats) IncrementRateLimited() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.RateLimited++
}

// GetStats ส่งสถิติปัจจุบัน
func (s *Stats) GetStats() Stats {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return Stats{
		TotalRequests: s.TotalRequests,
		CacheHits:     s.CacheHits,
		CacheMisses:   s.CacheMisses,
		RateLimited:   s.RateLimited,
	}
}

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
	{ID: 4, Title: "ทดสอบ Cache", Done: false},
	{ID: 5, Title: "ทดสอบ Rate Limit", Done: false},
}

// Global instances
var rateLimiter = NewRateLimiter(10, 1*time.Minute) // 10 requests ต่อนาที
var cache = NewCache()
var stats = &Stats{}

func main() {
	app := fiber.New()

	app.Use(logger.New())

	// Basic endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Rate Limiting & Cache API - Complete Version",
			"status":  "ok",
			"config": fiber.Map{
				"rate_limit": "10 requests/minute",
				"cache_ttl":  "30 seconds",
			},
		})
	})

	// Apply middleware สำหรับ /todos endpoint
	app.Use("/todos", rateLimitMiddleware())
	app.Use("/todos", cacheMiddleware())

	// Endpoints
	app.Get("/todos", getTodosHandler)
	app.Get("/stats", getStatsHandler)
	app.Get("/cache/clear", clearCacheHandler)

	log.Println("🚀 Server started on port 3000")
	log.Println("📊 Try: curl http://localhost:3000/todos")
	log.Println("🔥 Spam: for i in {1..15}; do curl http://localhost:3000/todos; echo; done")
	log.Fatal(app.Listen(":3000"))
}

// rateLimitMiddleware ตรวจสอบ rate limit
func rateLimitMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()

		// ตรวจสอบ rate limit
		if !rateLimiter.Allow(ip) {
			stats.IncrementRateLimited()

			remaining := rateLimiter.GetRemainingRequests(ip)
			return c.Status(429).JSON(fiber.Map{
				"error":       "Rate limit exceeded. Try again later.",
				"retry_after": "60s",
				"remaining":   remaining,
				"limit":       10,
				"window":      "1 minute",
			})
		}

		stats.IncrementTotal()

		// เพิ่ม headers ให้ client ทราบสถานะ
		remaining := rateLimiter.GetRemainingRequests(ip)
		c.Set("X-RateLimit-Limit", "10")
		c.Set("X-RateLimit-Remaining", string(rune(remaining)))
		c.Set("X-RateLimit-Reset", string(rune(time.Now().Add(1*time.Minute).Unix())))

		return c.Next()
	}
}

// cacheMiddleware ตรวจสอบและจัดการ cache
func cacheMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// สร้าง cache key จาก path และ query
		cacheKey := c.Route().Path + c.OriginalURL()

		// ตรวจสอบใน cache ก่อน
		if cachedData := cache.Get(cacheKey); cachedData != nil {
			stats.IncrementCacheHits()

			// เพิ่ม header บอกว่ามาจาก cache
			c.Set("X-Cache", "HIT")
			c.Set("X-Cache-TTL", "30s")

			return c.JSON(cachedData)
		}

		stats.IncrementCacheMisses()
		c.Set("X-Cache", "MISS")

		// เก็บ response สำหรับ cache ไว้
		c.Locals("cache_key", cacheKey)

		return c.Next()
	}
}

// getTodosHandler ส่งรายการ todos (จะใช้เวลาจำลองการประมวลผล)
func getTodosHandler(c *fiber.Ctx) error {
	// จำลองการประมวลผลที่ใช้เวลา (เหมือน database query)
	time.Sleep(100 * time.Millisecond)

	response := fiber.Map{
		"success":   true,
		"data":      todos,
		"count":     len(todos),
		"cached":    false,
		"timestamp": time.Now().Format("15:04:05"),
	}

	// เก็บลง cache สำหรับครั้งต่อไป
	if cacheKey := c.Locals("cache_key"); cacheKey != nil {
		cache.Set(cacheKey.(string), response, 30*time.Second)
	}

	return c.JSON(response)
}

// getStatsHandler ส่งสถิติการใช้งาน
func getStatsHandler(c *fiber.Ctx) error {
	currentStats := stats.GetStats()

	// เพิ่มข้อมูลเพิ่มเติม
	return c.JSON(fiber.Map{
		"success": true,
		"stats":   currentStats,
		"cache_info": fiber.Map{
			"ttl": "30 seconds",
			"hit_rate": func() float64 {
				total := currentStats.CacheHits + currentStats.CacheMisses
				if total == 0 {
					return 0
				}
				return float64(currentStats.CacheHits) / float64(total) * 100
			}(),
		},
		"rate_limit_info": fiber.Map{
			"limit":  10,
			"window": "1 minute",
		},
	})
}

// clearCacheHandler ล้าง cache
func clearCacheHandler(c *fiber.Ctx) error {
	cache.Clear()

	return c.JSON(fiber.Map{
		"success":   true,
		"message":   "Cache cleared successfully",
		"timestamp": time.Now().Format("15:04:05"),
	})
}
