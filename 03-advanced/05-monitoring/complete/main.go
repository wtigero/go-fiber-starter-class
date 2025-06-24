package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Metrics structures
type Metrics struct {
	RequestCount  int64                    `json:"request_count"`
	ErrorCount    int64                    `json:"error_count"`
	ResponseTimes []time.Duration          `json:"-"`
	EndpointStats map[string]EndpointMetric `json:"endpoint_stats"`
	mutex         sync.RWMutex
}

type EndpointMetric struct {
	Count     int64         `json:"count"`
	AvgTime   time.Duration `json:"avg_time"`
	TotalTime time.Duration `json:"total_time"`
	ErrorRate float64       `json:"error_rate"`
	Errors    int64         `json:"errors"`
}

type HealthStatus struct {
	Status     string       `json:"status"`
	Timestamp  time.Time    `json:"timestamp"`
	Uptime     string       `json:"uptime"`
	Version    string       `json:"version"`
	System     SystemHealth `json:"system"`
}

type SystemHealth struct {
	MemoryUsage string `json:"memory_usage"`
	Goroutines  int    `json:"goroutines"`
	CPUUsage    string `json:"cpu_usage"`
}

// Sample data
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

// Global metrics instance
var metrics = &Metrics{
	EndpointStats: make(map[string]EndpointMetric),
}

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Record error metrics
			recordError()
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		},
	})

	app.Use(logger.New())

	// Add monitoring middleware
	app.Use(monitoringMiddleware())

	// Basic routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Monitoring API - Complete Version",
			"status":  "ok",
		})
	})

	// Application endpoints
	app.Get("/todos", getTodosHandler)
	app.Post("/todos", createTodoHandler)
	app.Put("/todos/:id", updateTodoHandler)

	// Monitoring endpoints
	app.Get("/health", healthCheckHandler)
	app.Get("/metrics", metricsHandler)
	app.Get("/dashboard", dashboardHandler)

	log.Println("🚀 Monitoring API started on port 3000")
	log.Println("📊 Visit http://localhost:3000/dashboard for monitoring")
	log.Fatal(app.Listen(":3000"))
}

// Application handlers
func getTodosHandler(c *fiber.Ctx) error {
	// Simulate processing time
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

	// Simulate processing time
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
	// Simulate error sometimes
	if c.Params("id") == "999" {
		return fiber.NewError(404, "Todo not found")
	}

	// Simulate processing time
	time.Sleep(75 * time.Millisecond)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Todo updated",
	})
}

// Monitoring middleware
func monitoringMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Record metrics
		duration := time.Since(start)
		endpoint := c.Route().Path
		method := c.Method()
		statusCode := c.Response().StatusCode()

		recordMetrics(method+" "+endpoint, duration, statusCode)

		return err
	}
}

// Record metrics
func recordMetrics(endpoint string, duration time.Duration, statusCode int) {
	metrics.mutex.Lock()
	defer metrics.mutex.Unlock()

	metrics.RequestCount++
	metrics.ResponseTimes = append(metrics.ResponseTimes, duration)

	// Keep only last 1000 response times
	if len(metrics.ResponseTimes) > 1000 {
		metrics.ResponseTimes = metrics.ResponseTimes[len(metrics.ResponseTimes)-1000:]
	}

	if statusCode >= 400 {
		metrics.ErrorCount++
	}

	// Update endpoint stats
	stat := metrics.EndpointStats[endpoint]
	stat.Count++
	stat.TotalTime += duration
	stat.AvgTime = stat.TotalTime / time.Duration(stat.Count)

	if statusCode >= 400 {
		stat.Errors++
	}
	stat.ErrorRate = float64(stat.Errors) / float64(stat.Count) * 100

	metrics.EndpointStats[endpoint] = stat
}

func recordError() {
	metrics.mutex.Lock()
	defer metrics.mutex.Unlock()
	metrics.ErrorCount++
}

// Health check handler
func healthCheckHandler(c *fiber.Ctx) error {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	health := HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Uptime:    time.Since(startTime).String(),
		Version:   "1.0.0",
		System: SystemHealth{
			MemoryUsage: fmt.Sprintf("%.2f MB", float64(m.Alloc)/1024/1024),
			Goroutines:  runtime.NumGoroutine(),
			CPUUsage:    "N/A", // For simple example
		},
	}

	return c.JSON(health)
}

// Metrics handler
func metricsHandler(c *fiber.Ctx) error {
	metrics.mutex.RLock()
	defer metrics.mutex.RUnlock()

	// Calculate average response time
	var avgResponseTime time.Duration
	if len(metrics.ResponseTimes) > 0 {
		var total time.Duration
		for _, t := range metrics.ResponseTimes {
			total += t
		}
		avgResponseTime = total / time.Duration(len(metrics.ResponseTimes))
	}

	// Calculate success rate
	successRate := float64(0)
	if metrics.RequestCount > 0 {
		successRate = float64(metrics.RequestCount-metrics.ErrorCount) / float64(metrics.RequestCount) * 100
	}

	return c.JSON(fiber.Map{
		"success": true,
		"metrics": fiber.Map{
			"requests": fiber.Map{
				"total":             metrics.RequestCount,
				"errors":            metrics.ErrorCount,
				"success":           metrics.RequestCount - metrics.ErrorCount,
				"success_rate":      fmt.Sprintf("%.2f%%", successRate),
				"avg_response_time": avgResponseTime.String(),
			},
			"endpoints": metrics.EndpointStats,
			"system": fiber.Map{
				"uptime":     time.Since(startTime).String(),
				"goroutines": runtime.NumGoroutine(),
			},
		},
	})
}

// Dashboard handler
func dashboardHandler(c *fiber.Ctx) error {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>API Monitoring Dashboard</title>
    <style>
        body { 
            font-family: Arial, sans-serif; 
            margin: 0;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 20px;
            border-radius: 10px;
            margin-bottom: 20px;
            text-align: center;
        }
        .stats { 
            display: grid; 
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 20px; 
            margin-bottom: 20px; 
        }
        .stat-card { 
            background: white;
            border: 1px solid #ddd; 
            padding: 20px; 
            border-radius: 10px; 
            text-align: center;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .stat-value { 
            font-size: 2em; 
            font-weight: bold; 
            color: #007bff; 
            margin: 10px 0;
        }
        .stat-label {
            color: #666;
            text-transform: uppercase;
            font-size: 0.9em;
            font-weight: bold;
        }
        .endpoints {
            background: white;
            border-radius: 10px;
            padding: 20px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .endpoint {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 10px 0;
            border-bottom: 1px solid #eee;
        }
        .endpoint:last-child {
            border-bottom: none;
        }
        .error { color: #dc3545; }
        .success { color: #28a745; }
        .warning { color: #ffc107; }
        .refresh {
            position: fixed;
            top: 20px;
            right: 20px;
            background: #007bff;
            color: white;
            padding: 10px 20px;
            border: none;
            border-radius: 5px;
            cursor: pointer;
        }
    </style>
</head>
<body>
    <button class="refresh" onclick="updateDashboard()">🔄 Refresh</button>
    
    <div class="header">
        <h1>📊 API Monitoring Dashboard</h1>
        <p>Real-time monitoring for your Go Fiber API</p>
    </div>
    
    <div class="stats">
        <div class="stat-card">
            <div class="stat-label">Total Requests</div>
            <div class="stat-value" id="total-requests">Loading...</div>
        </div>
        
        <div class="stat-card">
            <div class="stat-label">Success Rate</div>
            <div class="stat-value success" id="success-rate">Loading...</div>
        </div>
        
        <div class="stat-card">
            <div class="stat-label">Avg Response Time</div>
            <div class="stat-value" id="avg-response">Loading...</div>
        </div>
        
        <div class="stat-card">
            <div class="stat-label">Active Goroutines</div>
            <div class="stat-value warning" id="goroutines">Loading...</div>
        </div>
    </div>
    
    <div class="endpoints">
        <h3>📋 Endpoint Statistics</h3>
        <div id="endpoint-stats">Loading...</div>
    </div>
    
    <script>
        function updateDashboard() {
            fetch('/metrics')
                .then(response => response.json())
                .then(data => {
                    const metrics = data.metrics;
                    
                    // Update main stats
                    document.getElementById('total-requests').textContent = metrics.requests.total;
                    document.getElementById('success-rate').textContent = metrics.requests.success_rate;
                    document.getElementById('avg-response').textContent = metrics.requests.avg_response_time;
                    document.getElementById('goroutines').textContent = metrics.system.goroutines;
                    
                    // Update endpoint stats
                    const endpointContainer = document.getElementById('endpoint-stats');
                    endpointContainer.innerHTML = '';
                    
                    Object.keys(metrics.endpoints).forEach(endpoint => {
                        const stats = metrics.endpoints[endpoint];
                        const div = document.createElement('div');
                        div.className = 'endpoint';
                        
                        const errorRateClass = stats.error_rate > 10 ? 'error' : stats.error_rate > 5 ? 'warning' : 'success';
                        
                        div.innerHTML = \`
                            <div><strong>\${endpoint}</strong></div>
                            <div>
                                <span>Calls: \${stats.count}</span> | 
                                <span>Avg: \${stats.avg_time}</span> | 
                                <span class="\${errorRateClass}">Error Rate: \${stats.error_rate.toFixed(2)}%</span>
                            </div>
                        \`;
                        
                        endpointContainer.appendChild(div);
                    });
                })
                .catch(error => {
                    console.error('Error fetching metrics:', error);
                    document.getElementById('total-requests').textContent = 'Error';
                });
        }
        
        // Update every 3 seconds
        updateDashboard();
        setInterval(updateDashboard, 3000);
        
        // Auto-refresh every 30 seconds
        setTimeout(() => {
            location.reload();
        }, 30000);
    </script>
</body>
</html>
    `

	c.Set("Content-Type", "text/html")
	return c.SendString(html)
} 