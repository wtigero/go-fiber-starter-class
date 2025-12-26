package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof" // pprof endpoints
	"runtime"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// ============ Types ============

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	Tags      []string  `json:"tags"`
}

// ============ Global Variables ============

var (
	// Object pool for buffer reuse
	bufferPool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	// User pool
	userPool = sync.Pool{
		New: func() interface{} {
			return &User{}
		},
	}
)

func main() {
	// Start pprof server on separate port
	go func() {
		log.Println("📊 pprof available at http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatal(err)
		}
	}()

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(logger.New(logger.Config{
		Format: "${time} ${method} ${path} - ${latency}\n",
	}))

	// Routes
	app.Get("/", homeHandler)
	app.Get("/stats", statsHandler)

	// Slow endpoints (for profiling)
	app.Get("/slow/cpu", slowCPUHandler)
	app.Get("/slow/memory", slowMemoryHandler)
	app.Get("/slow/alloc", slowAllocHandler)

	// Optimized endpoints
	app.Get("/fast/cpu", fastCPUHandler)
	app.Get("/fast/memory", fastMemoryHandler)
	app.Get("/fast/alloc", fastAllocHandler)

	// Compare endpoints
	app.Get("/compare/json", compareJSONHandler)

	log.Println("🚀 Server running on http://localhost:3000")
	log.Println("📊 pprof running on http://localhost:6060/debug/pprof/")
	log.Println("")
	log.Println("Try these commands to profile:")
	log.Println("  CPU:    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30")
	log.Println("  Memory: go tool pprof http://localhost:6060/debug/pprof/heap")
	log.Println("  Goroutines: go tool pprof http://localhost:6060/debug/pprof/goroutine")

	log.Fatal(app.Listen(":3000"))
}

// ============ Handlers ============

func homeHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Profiling Demo Server",
		"endpoints": fiber.Map{
			"slow": fiber.Map{
				"/slow/cpu":    "CPU intensive (inefficient)",
				"/slow/memory": "Memory intensive (inefficient)",
				"/slow/alloc":  "High allocations (inefficient)",
			},
			"fast": fiber.Map{
				"/fast/cpu":    "CPU efficient",
				"/fast/memory": "Memory efficient (pooled)",
				"/fast/alloc":  "Low allocations (optimized)",
			},
			"profiling": fiber.Map{
				"pprof": "http://localhost:6060/debug/pprof/",
			},
		},
	})
}

func statsHandler(c *fiber.Ctx) error {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return c.JSON(fiber.Map{
		"memory": fiber.Map{
			"alloc_mb":       m.Alloc / 1024 / 1024,
			"total_alloc_mb": m.TotalAlloc / 1024 / 1024,
			"sys_mb":         m.Sys / 1024 / 1024,
			"num_gc":         m.NumGC,
			"gc_pause_ns":    m.PauseNs[(m.NumGC+255)%256],
		},
		"goroutines": runtime.NumGoroutine(),
		"cpus":       runtime.NumCPU(),
	})
}

// ============ SLOW Handlers (for demonstration) ============

// slowCPUHandler - Inefficient CPU usage
func slowCPUHandler(c *fiber.Ctx) error {
	// BAD: String concatenation in loop
	result := ""
	for i := 0; i < 10000; i++ {
		result += fmt.Sprintf("item-%d,", i) // Creates new string each time!
	}

	return c.JSON(fiber.Map{
		"type":   "slow_cpu",
		"length": len(result),
	})
}

// slowMemoryHandler - Inefficient memory usage
func slowMemoryHandler(c *fiber.Ctx) error {
	// BAD: Creating new slices without pre-allocation
	users := []User{}
	for i := 0; i < 1000; i++ {
		user := User{
			ID:        i,
			Name:      fmt.Sprintf("User %d", i),
			Email:     fmt.Sprintf("user%d@example.com", i),
			CreatedAt: time.Now(),
			Tags:      []string{"tag1", "tag2", "tag3"},
		}
		users = append(users, user) // Slice grows, causing reallocations
	}

	return c.JSON(fiber.Map{
		"type":  "slow_memory",
		"count": len(users),
	})
}

// slowAllocHandler - High allocation rate
func slowAllocHandler(c *fiber.Ctx) error {
	var results []string
	for i := 0; i < 1000; i++ {
		// BAD: Creating new buffer each time
		buf := new(bytes.Buffer)
		buf.WriteString(fmt.Sprintf("data-%d", i))
		results = append(results, buf.String())
	}

	return c.JSON(fiber.Map{
		"type":  "slow_alloc",
		"count": len(results),
	})
}

// ============ FAST Handlers (optimized) ============

// fastCPUHandler - Efficient CPU usage
func fastCPUHandler(c *fiber.Ctx) error {
	// GOOD: Use strings.Builder with pre-allocation
	var builder bytes.Buffer
	builder.Grow(10000 * 15) // Pre-allocate approximate size

	for i := 0; i < 10000; i++ {
		fmt.Fprintf(&builder, "item-%d,", i)
	}

	return c.JSON(fiber.Map{
		"type":   "fast_cpu",
		"length": builder.Len(),
	})
}

// fastMemoryHandler - Efficient memory usage
func fastMemoryHandler(c *fiber.Ctx) error {
	// GOOD: Pre-allocate slice with known capacity
	users := make([]User, 0, 1000)

	for i := 0; i < 1000; i++ {
		// GOOD: Reuse from pool
		user := userPool.Get().(*User)
		user.ID = i
		user.Name = fmt.Sprintf("User %d", i)
		user.Email = fmt.Sprintf("user%d@example.com", i)
		user.CreatedAt = time.Now()
		user.Tags = []string{"tag1", "tag2", "tag3"}

		users = append(users, *user)

		// Reset and return to pool
		user.ID = 0
		user.Name = ""
		user.Email = ""
		user.Tags = nil
		userPool.Put(user)
	}

	return c.JSON(fiber.Map{
		"type":  "fast_memory",
		"count": len(users),
	})
}

// fastAllocHandler - Low allocation rate
func fastAllocHandler(c *fiber.Ctx) error {
	results := make([]string, 0, 1000)

	for i := 0; i < 1000; i++ {
		// GOOD: Reuse buffer from pool
		buf := bufferPool.Get().(*bytes.Buffer)
		buf.Reset()
		fmt.Fprintf(buf, "data-%d", i)
		results = append(results, buf.String())
		bufferPool.Put(buf)
	}

	return c.JSON(fiber.Map{
		"type":  "fast_alloc",
		"count": len(results),
	})
}

// ============ Comparison Handler ============

func compareJSONHandler(c *fiber.Ctx) error {
	// Generate test data
	data := generateTestData(100)

	// Measure slow method
	start := time.Now()
	for i := 0; i < 100; i++ {
		json.Marshal(data) // Not reusing encoder
	}
	slowDuration := time.Since(start)

	// Measure fast method (with buffer pool)
	start = time.Now()
	for i := 0; i < 100; i++ {
		buf := bufferPool.Get().(*bytes.Buffer)
		buf.Reset()
		enc := json.NewEncoder(buf)
		enc.Encode(data)
		bufferPool.Put(buf)
	}
	fastDuration := time.Since(start)

	return c.JSON(fiber.Map{
		"slow_method_ms": slowDuration.Milliseconds(),
		"fast_method_ms": fastDuration.Milliseconds(),
		"improvement":    fmt.Sprintf("%.1fx faster", float64(slowDuration)/float64(fastDuration)),
	})
}

func generateTestData(n int) []User {
	users := make([]User, n)
	for i := 0; i < n; i++ {
		users[i] = User{
			ID:        i,
			Name:      fmt.Sprintf("User %d", i),
			Email:     fmt.Sprintf("user%d@example.com", i),
			CreatedAt: time.Now(),
			Tags:      []string{"tag1", "tag2", fmt.Sprintf("tag%d", rand.Intn(100))},
		}
	}
	return users
}
