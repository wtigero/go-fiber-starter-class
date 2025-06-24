package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// TODO: สร้าง ServiceRegistry สำหรับ service discovery
// type ServiceRegistry struct {
//     services map[string]string
//     mutex    sync.RWMutex
// }

// TODO: สร้าง Health Check Response
// type HealthResponse struct {
//     Gateway  string                 `json:"gateway"`
//     Services map[string]string      `json:"services"`
//     Timestamp time.Time             `json:"timestamp"`
// }

var serviceRegistry = map[string]string{
	"user_service": "http://localhost:3001",
	"todo_service": "http://localhost:3002",
}

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

	// Basic route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message":  "API Gateway - Ready to proxy!",
			"status":   "ok",
			"services": serviceRegistry,
		})
	})

	// TODO: สร้าง health check endpoint
	// app.Get("/health", healthCheckHandler)

	// TODO: สร้าง proxy routes สำหรับ User Service
	// app.Post("/users", proxyToUserService)
	// app.Get("/users/:id", proxyToUserService)

	// TODO: สร้าง proxy routes สำหรับ Todo Service
	// app.Get("/todos", proxyToTodoService)
	// app.Post("/todos", proxyToTodoService)

	log.Println("🚀 API Gateway started on port 3000")
	log.Println("📡 Make sure User Service (3001) and Todo Service (3002) are running")
	log.Fatal(app.Listen(":3000"))
}

// TODO: สร้าง healthCheckHandler
// func healthCheckHandler(c *fiber.Ctx) error {
//     // ตรวจสอบ health ของทุก services
//     services := make(map[string]string)
//
//     for serviceName, serviceURL := range serviceRegistry {
//         if checkServiceHealth(serviceURL) {
//             services[serviceName] = "healthy"
//         } else {
//             services[serviceName] = "unhealthy"
//         }
//     }
//
//     return c.JSON(HealthResponse{
//         Gateway:   "healthy",
//         Services:  services,
//         Timestamp: time.Now(),
//     })
// }

// TODO: สร้าง checkServiceHealth
// func checkServiceHealth(serviceURL string) bool {
//     client := &http.Client{Timeout: 5 * time.Second}
//     resp, err := client.Get(serviceURL + "/health")
//     if err != nil {
//         return false
//     }
//     defer resp.Body.Close()
//     return resp.StatusCode == 200
// }

// TODO: สร้าง proxyToUserService
// func proxyToUserService(c *fiber.Ctx) error {
//     return proxyRequest(c, serviceRegistry["user_service"])
// }

// TODO: สร้าง proxyToTodoService
// func proxyToTodoService(c *fiber.Ctx) error {
//     return proxyRequest(c, serviceRegistry["todo_service"])
// }

// TODO: สร้าง proxyRequest function
// func proxyRequest(c *fiber.Ctx, targetURL string) error {
//     // สร้าง URL เต็ม
//     fullURL := targetURL + c.OriginalURL()
//
//     // สร้าง request ใหม่
//     var body io.Reader
//     if len(c.Body()) > 0 {
//         body = bytes.NewReader(c.Body())
//     }
//
//     req, err := http.NewRequest(c.Method(), fullURL, body)
//     if err != nil {
//         return c.Status(500).JSON(fiber.Map{"error": "Failed to create request"})
//     }
//
//     // Copy headers
//     req.Header.Set("Content-Type", c.Get("Content-Type"))
//
//     // ส่ง request
//     client := &http.Client{Timeout: 30 * time.Second}
//     resp, err := client.Do(req)
//     if err != nil {
//         return c.Status(503).JSON(fiber.Map{
//             "error": "Service unavailable",
//             "service": targetURL,
//         })
//     }
//     defer resp.Body.Close()
//
//     // อ่าน response
//     respBody, err := io.ReadAll(resp.Body)
//     if err != nil {
//         return c.Status(500).JSON(fiber.Map{"error": "Failed to read response"})
//     }
//
//     // ส่ง response กลับ
//     c.Status(resp.StatusCode)
//     c.Set("Content-Type", resp.Header.Get("Content-Type"))
//     return c.Send(respBody)
// }
