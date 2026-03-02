package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// ============ Object Pools ============

// Buffer Pool - สำหรับ reuse byte buffers
var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// Sized Buffer Pool - สำหรับ buffers ขนาดต่างๆ
var (
	smallBufferPool = sync.Pool{
		New: func() interface{} {
			b := make([]byte, 1024) // 1KB
			return &b
		},
	}
	mediumBufferPool = sync.Pool{
		New: func() interface{} {
			b := make([]byte, 8192) // 8KB
			return &b
		},
	}
	largeBufferPool = sync.Pool{
		New: func() interface{} {
			b := make([]byte, 65536) // 64KB
			return &b
		},
	}
)

// Request Pool - สำหรับ reuse request objects
type Request struct {
	ID        string
	Method    string
	Path      string
	Body      []byte
	Headers   map[string]string
	Timestamp time.Time
}

var requestPool = sync.Pool{
	New: func() interface{} {
		return &Request{
			Headers: make(map[string]string),
		}
	},
}

// Response Pool - สำหรับ reuse response objects
type Response struct {
	StatusCode int
	Body       []byte
	Headers    map[string]string
}

var responsePool = sync.Pool{
	New: func() interface{} {
		return &Response{
			Headers: make(map[string]string),
		}
	},
}

// JSON Encoder Pool
type JSONEncoder struct {
	buffer  *bytes.Buffer
	encoder *json.Encoder
}

var jsonEncoderPool = sync.Pool{
	New: func() interface{} {
		buf := new(bytes.Buffer)
		return &JSONEncoder{
			buffer:  buf,
			encoder: json.NewEncoder(buf),
		}
	},
}

// ============ Pool Statistics ============

type PoolStats struct {
	Gets      int64
	Puts      int64
	News      int64
	HitRate   float64
	mu        sync.Mutex
}

var stats = &PoolStats{}

func (s *PoolStats) RecordGet(isNew bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Gets++
	if isNew {
		s.News++
	}
	if s.Gets > 0 {
		s.HitRate = float64(s.Gets-s.News) / float64(s.Gets) * 100
	}
}

func (s *PoolStats) RecordPut() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Puts++
}

// ============ Ring Buffer Pool ============

type RingBufferPool struct {
	buffers [][]byte
	head    int
	tail    int
	size    int
	mu      sync.Mutex
}

func NewRingBufferPool(count, bufSize int) *RingBufferPool {
	pool := &RingBufferPool{
		buffers: make([][]byte, count),
		size:    count,
	}
	for i := 0; i < count; i++ {
		pool.buffers[i] = make([]byte, bufSize)
	}
	return pool
}

func (p *RingBufferPool) Get() []byte {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.head == p.tail {
		// Pool empty, allocate new
		return make([]byte, cap(p.buffers[0]))
	}

	buf := p.buffers[p.head]
	p.head = (p.head + 1) % p.size
	return buf
}

func (p *RingBufferPool) Put(buf []byte) {
	p.mu.Lock()
	defer p.mu.Unlock()

	next := (p.tail + 1) % p.size
	if next == p.head {
		// Pool full, discard
		return
	}

	p.buffers[p.tail] = buf[:0] // Reset length
	p.tail = next
}

// Global ring buffer pool
var ringPool = NewRingBufferPool(100, 4096)

// ============ Handlers ============

func main() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(logger.New())

	// Routes
	app.Get("/", homeHandler)
	app.Get("/stats", statsHandler)
	app.Get("/memory", memoryHandler)

	// Without Pool (slow)
	app.Get("/without-pool", withoutPoolHandler)

	// With Pool (fast)
	app.Get("/with-pool", withPoolHandler)

	// Benchmark endpoints
	app.Get("/bench/buffer", benchBufferHandler)
	app.Get("/bench/json", benchJSONHandler)
	app.Get("/bench/request", benchRequestHandler)

	log.Println("🚀 Memory Pool Demo running on http://localhost:3000")
	log.Println("")
	log.Println("Try these endpoints:")
	log.Println("  /without-pool  - Creates new objects each time (slow)")
	log.Println("  /with-pool     - Reuses pooled objects (fast)")
	log.Println("  /stats         - Pool statistics")
	log.Println("  /memory        - Memory usage")

	log.Fatal(app.Listen(":3000"))
}

func homeHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Memory Pool Demo",
		"endpoints": fiber.Map{
			"/without-pool": "No pooling (allocates each time)",
			"/with-pool":    "With pooling (reuses objects)",
			"/stats":        "Pool statistics",
			"/memory":       "Memory usage",
			"/bench/*":      "Benchmark endpoints",
		},
	})
}

func statsHandler(c *fiber.Ctx) error {
	stats.mu.Lock()
	defer stats.mu.Unlock()

	return c.JSON(fiber.Map{
		"pool_stats": fiber.Map{
			"gets":     stats.Gets,
			"puts":     stats.Puts,
			"news":     stats.News,
			"hit_rate": fmt.Sprintf("%.2f%%", stats.HitRate),
		},
	})
}

func memoryHandler(c *fiber.Ctx) error {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return c.JSON(fiber.Map{
		"memory": fiber.Map{
			"alloc_mb":       m.Alloc / 1024 / 1024,
			"total_alloc_mb": m.TotalAlloc / 1024 / 1024,
			"sys_mb":         m.Sys / 1024 / 1024,
			"num_gc":         m.NumGC,
			"gc_pause_total": fmt.Sprintf("%v", time.Duration(m.PauseTotalNs)),
		},
		"goroutines": runtime.NumGoroutine(),
	})
}

// ❌ WITHOUT POOL - สร้าง object ใหม่ทุกครั้ง
func withoutPoolHandler(c *fiber.Ctx) error {
	// สร้าง buffer ใหม่ทุกครั้ง
	buf := new(bytes.Buffer)

	// สร้าง request object ใหม่
	req := &Request{
		ID:        "req-123",
		Method:    "GET",
		Path:      "/test",
		Body:      make([]byte, 1024),
		Headers:   make(map[string]string),
		Timestamp: time.Now(),
	}

	// สร้าง response object ใหม่
	resp := &Response{
		StatusCode: 200,
		Body:       make([]byte, 1024),
		Headers:    make(map[string]string),
	}

	// ใช้งาน objects
	buf.WriteString("Hello from non-pooled handler")
	req.Headers["Content-Type"] = "application/json"
	resp.Headers["X-Request-ID"] = req.ID

	return c.JSON(fiber.Map{
		"method":  "without_pool",
		"message": buf.String(),
		"note":    "Objects allocated fresh each request",
	})
}

// ✅ WITH POOL - reuse objects จาก pool
func withPoolHandler(c *fiber.Ctx) error {
	// ดึง buffer จาก pool
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)
	stats.RecordGet(false)

	// ดึง request จาก pool
	req := requestPool.Get().(*Request)
	defer func() {
		// Reset ก่อนคืน pool
		req.ID = ""
		req.Method = ""
		req.Path = ""
		req.Body = req.Body[:0]
		for k := range req.Headers {
			delete(req.Headers, k)
		}
		requestPool.Put(req)
		stats.RecordPut()
	}()

	// ดึง response จาก pool
	resp := responsePool.Get().(*Response)
	defer func() {
		resp.StatusCode = 0
		resp.Body = resp.Body[:0]
		for k := range resp.Headers {
			delete(resp.Headers, k)
		}
		responsePool.Put(resp)
	}()

	// ใช้งาน objects
	buf.WriteString("Hello from pooled handler")
	req.ID = "req-123"
	req.Method = "GET"
	req.Path = "/test"
	req.Timestamp = time.Now()
	req.Headers["Content-Type"] = "application/json"

	resp.StatusCode = 200
	resp.Headers["X-Request-ID"] = req.ID

	return c.JSON(fiber.Map{
		"method":  "with_pool",
		"message": buf.String(),
		"note":    "Objects reused from pool",
	})
}

// Benchmark: Buffer Pool
func benchBufferHandler(c *fiber.Ctx) error {
	iterations := 10000

	// Without pool
	start := time.Now()
	for i := 0; i < iterations; i++ {
		buf := new(bytes.Buffer)
		buf.WriteString("test data")
		_ = buf.String()
	}
	withoutPool := time.Since(start)

	// With pool
	start = time.Now()
	for i := 0; i < iterations; i++ {
		buf := bufferPool.Get().(*bytes.Buffer)
		buf.Reset()
		buf.WriteString("test data")
		_ = buf.String()
		bufferPool.Put(buf)
	}
	withPool := time.Since(start)

	return c.JSON(fiber.Map{
		"iterations":   iterations,
		"without_pool": withoutPool.String(),
		"with_pool":    withPool.String(),
		"speedup":      fmt.Sprintf("%.2fx", float64(withoutPool)/float64(withPool)),
	})
}

// Benchmark: JSON Encoder Pool
func benchJSONHandler(c *fiber.Ctx) error {
	data := map[string]interface{}{
		"id":   1,
		"name": "test",
		"tags": []string{"a", "b", "c"},
	}
	iterations := 10000

	// Without pool
	start := time.Now()
	for i := 0; i < iterations; i++ {
		buf := new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.Encode(data)
		_ = buf.Bytes()
	}
	withoutPool := time.Since(start)

	// With pool
	start = time.Now()
	for i := 0; i < iterations; i++ {
		je := jsonEncoderPool.Get().(*JSONEncoder)
		je.buffer.Reset()
		je.encoder.Encode(data)
		_ = je.buffer.Bytes()
		jsonEncoderPool.Put(je)
	}
	withPool := time.Since(start)

	return c.JSON(fiber.Map{
		"iterations":   iterations,
		"without_pool": withoutPool.String(),
		"with_pool":    withPool.String(),
		"speedup":      fmt.Sprintf("%.2fx", float64(withoutPool)/float64(withPool)),
	})
}

// Benchmark: Request Pool
func benchRequestHandler(c *fiber.Ctx) error {
	iterations := 10000

	// Without pool
	start := time.Now()
	for i := 0; i < iterations; i++ {
		req := &Request{
			ID:        "req-123",
			Method:    "GET",
			Path:      "/test",
			Body:      make([]byte, 1024),
			Headers:   make(map[string]string),
			Timestamp: time.Now(),
		}
		req.Headers["Content-Type"] = "application/json"
		_ = req
	}
	withoutPool := time.Since(start)

	// With pool
	start = time.Now()
	for i := 0; i < iterations; i++ {
		req := requestPool.Get().(*Request)
		req.ID = "req-123"
		req.Method = "GET"
		req.Path = "/test"
		req.Timestamp = time.Now()
		req.Headers["Content-Type"] = "application/json"

		// Reset and return
		req.ID = ""
		req.Method = ""
		req.Path = ""
		for k := range req.Headers {
			delete(req.Headers, k)
		}
		requestPool.Put(req)
	}
	withPool := time.Since(start)

	return c.JSON(fiber.Map{
		"iterations":   iterations,
		"without_pool": withoutPool.String(),
		"with_pool":    withPool.String(),
		"speedup":      fmt.Sprintf("%.2fx", float64(withoutPool)/float64(withPool)),
	})
}
