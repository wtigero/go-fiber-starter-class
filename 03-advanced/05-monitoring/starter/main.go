package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// TODO: สร้าง Metrics struct
// type Metrics struct {
//     RequestCount    int64         `json:"request_count"`
//     ErrorCount      int64         `json:"error_count"`
//     ResponseTimes   []time.Duration
//     EndpointStats   map[string]EndpointMetric
//     mutex           sync.RWMutex
// }

// TODO: สร้าง EndpointMetric struct
// type EndpointMetric struct {
//     Count       int64           `json:"count"`
//     AvgTime     time.Duration   `json:"avg_time"`
//     TotalTime   time.Duration   `json:"total_time"`
// }

// TODO: สร้าง HealthStatus struct
// type HealthStatus struct {
//     Status     string                 `json:"status"`
//     Timestamp  time.Time             `json:"timestamp"`
//     Uptime     string                `json:"uptime"`
//     Version    string                `json:"version"`
//     System     SystemHealth          `json:"system"`
// }

// TODO: สร้าง SystemHealth struct
// type SystemHealth struct {
//     MemoryUsage  string `json:"memory_usage"`
//     Goroutines   int    `json:"goroutines"`
//     CPUUsage     string `json:"cpu_usage"`
// }

// ข้อมูลตัวอย่าง
type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var todos = []Todo{
	{ID: 1, Title: "เรียน Monitoring", Done: false},
	{ID: 2, Title: "สร้าง Dashboard", Done: false},
	{ID: 3, Title: "ทดสอบ Metrics", Done: true},
}

var todoID = 4
var startTime = time.Now()

// TODO: สร้างตัวแปร global metrics
// var metrics = &Metrics{
//     EndpointStats: make(map[string]EndpointMetric),
// }

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// TODO: เก็บ error metrics
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		},
	})

	app.Use(logger.New())

	// TODO: เพิ่ม monitoring middleware
	// app.Use(monitoringMiddleware())

	// Basic routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Monitoring API - Ready to track!",
			"status":  "ok",
		})
	})

	// Application endpoints
	app.Get("/todos", getTodosHandler)
	app.Post("/todos", createTodoHandler)
	app.Put("/todos/:id", updateTodoHandler)

	// TODO: สร้าง monitoring endpoints
	// app.Get("/health", healthCheckHandler)
	// app.Get("/metrics", metricsHandler)
	// app.Get("/dashboard", dashboardHandler)

	log.Println("🚀 Monitoring API started on port 3000")
	log.Println("📊 Visit http://localhost:3000/dashboard for monitoring")
	log.Fatal(app.Listen(":3000"))
}

// Application handlers
func getTodosHandler(c *fiber.Ctx) error {
	// จำลองการประมวลผลที่ใช้เวลา
	time.Sleep(50 * time.Millisecond)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    todos,
		"count":   len(todos),
	})
}

func createTodoHandler(c *fiber.Ctx) error {
	type CreateTodoRequest struct {
		Title string `json:"title"`
	}

	var req CreateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.Title == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Title is required"})
	}

	// จำลองการประมวลผลที่ใช้เวลา
	time.Sleep(100 * time.Millisecond)

	todo := Todo{
		ID:    todoID,
		Title: req.Title,
		Done:  false,
	}

	todos = append(todos, todo)
	todoID++

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"todo":    todo,
	})
}

func updateTodoHandler(c *fiber.Ctx) error {
	// จำลองการ error บางครั้ง
	if c.Params("id") == "999" {
		return fiber.NewError(404, "Todo not found")
	}

	// จำลองการประมวลผลที่ใช้เวลา
	time.Sleep(75 * time.Millisecond)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Todo updated",
	})
}

// TODO: สร้าง monitoringMiddleware
// func monitoringMiddleware() fiber.Handler {
//     return func(c *fiber.Ctx) error {
//         start := time.Now()
//
//         // ดำเนินการ request
//         err := c.Next()
//
//         // เก็บ metrics
//         duration := time.Since(start)
//         endpoint := c.Route().Path
//         method := c.Method()
//         statusCode := c.Response().StatusCode()
//
//         recordMetrics(method + " " + endpoint, duration, statusCode)
//
//         return err
//     }
// }

// TODO: สร้าง recordMetrics
// func recordMetrics(endpoint string, duration time.Duration, statusCode int) {
//     metrics.mutex.Lock()
//     defer metrics.mutex.Unlock()
//
//     metrics.RequestCount++
//     metrics.ResponseTimes = append(metrics.ResponseTimes, duration)
//
//     if statusCode >= 400 {
//         metrics.ErrorCount++
//     }
//
//     // อัปเดต endpoint stats
//     stat := metrics.EndpointStats[endpoint]
//     stat.Count++
//     stat.TotalTime += duration
//     stat.AvgTime = stat.TotalTime / time.Duration(stat.Count)
//     metrics.EndpointStats[endpoint] = stat
// }

// TODO: สร้าง healthCheckHandler
// func healthCheckHandler(c *fiber.Ctx) error {
//     var m runtime.MemStats
//     runtime.ReadMemStats(&m)
//
//     health := HealthStatus{
//         Status:    "healthy",
//         Timestamp: time.Now(),
//         Uptime:    time.Since(startTime).String(),
//         Version:   "1.0.0",
//         System: SystemHealth{
//             MemoryUsage: fmt.Sprintf("%.2f MB", float64(m.Alloc)/1024/1024),
//             Goroutines:  runtime.NumGoroutine(),
//             CPUUsage:    "N/A", // สำหรับตัวอย่างง่าย ๆ
//         },
//     }
//
//     return c.JSON(health)
// }

// TODO: สร้าง metricsHandler
// func metricsHandler(c *fiber.Ctx) error {
//     metrics.mutex.RLock()
//     defer metrics.mutex.RUnlock()
//
//     // คำนวณ average response time
//     var avgResponseTime time.Duration
//     if len(metrics.ResponseTimes) > 0 {
//         var total time.Duration
//         for _, t := range metrics.ResponseTimes {
//             total += t
//         }
//         avgResponseTime = total / time.Duration(len(metrics.ResponseTimes))
//     }
//
//     return c.JSON(fiber.Map{
//         "success": true,
//         "metrics": fiber.Map{
//             "requests": fiber.Map{
//                 "total":             metrics.RequestCount,
//                 "errors":            metrics.ErrorCount,
//                 "success":           metrics.RequestCount - metrics.ErrorCount,
//                 "avg_response_time": avgResponseTime.String(),
//             },
//             "endpoints": metrics.EndpointStats,
//         },
//     })
// }

// TODO: สร้าง dashboardHandler
// func dashboardHandler(c *fiber.Ctx) error {
//     html := `
//     <!DOCTYPE html>
//     <html>
//     <head>
//         <title>API Monitoring Dashboard</title>
//         <style>
//             body { font-family: Arial, sans-serif; margin: 20px; }
//             .stats { display: flex; gap: 20px; margin-bottom: 20px; }
//             .stat-card {
//                 border: 1px solid #ddd;
//                 padding: 20px;
//                 border-radius: 8px;
//                 flex: 1;
//                 text-align: center;
//             }
//             .stat-value { font-size: 2em; font-weight: bold; color: #007bff; }
//             h1 { color: #333; }
//         </style>
//     </head>
//     <body>
//         <h1>📊 API Monitoring Dashboard</h1>
//
//         <div class="stats">
//             <div class="stat-card">
//                 <h3>Total Requests</h3>
//                 <div class="stat-value" id="total-requests">Loading...</div>
//             </div>
//
//             <div class="stat-card">
//                 <h3>Error Rate</h3>
//                 <div class="stat-value" id="error-rate">Loading...</div>
//             </div>
//
//             <div class="stat-card">
//                 <h3>Avg Response Time</h3>
//                 <div class="stat-value" id="avg-response">Loading...</div>
//             </div>
//         </div>
//
//         <script>
//             function updateDashboard() {
//                 fetch('/metrics')
//                     .then(response => response.json())
//                     .then(data => {
//                         const metrics = data.metrics;
//                         document.getElementById('total-requests').textContent = metrics.requests.total;
//
//                         const errorRate = ((metrics.requests.errors / metrics.requests.total) * 100).toFixed(2);
//                         document.getElementById('error-rate').textContent = errorRate + '%';
//
//                         document.getElementById('avg-response').textContent = metrics.requests.avg_response_time;
//                     })
//                     .catch(error => console.error('Error:', error));
//             }
//
//             // อัปเดตทุก 5 วินาที
//             updateDashboard();
//             setInterval(updateDashboard, 5000);
//         </script>
//     </body>
//     </html>
//     `
//
//     c.Set("Content-Type", "text/html")
//     return c.SendString(html)
// }
