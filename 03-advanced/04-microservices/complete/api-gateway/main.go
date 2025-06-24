package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Service URLs
const (
	UserServiceURL = "http://localhost:3001"
	TodoServiceURL = "http://localhost:3002"
)

// Response structures
type ServiceResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type HealthCheck struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	URL     string `json:"url"`
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

	app.Use(logger.New())
	app.Use(cors.New())

	// API Gateway info
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "API Gateway",
			"version": "1.0.0",
			"status":  "running",
			"services": map[string]string{
				"users": UserServiceURL,
				"todos": TodoServiceURL,
			},
		})
	})

	// Health check for all services
	app.Get("/health", healthCheckHandler)

	// User service routes
	app.Post("/api/users", proxyToUserService)
	app.Get("/api/users/:id", proxyToUserService)
	app.Get("/api/users", proxyToUserService)

	// Todo service routes
	app.Post("/api/todos", proxyToTodoService)
	app.Get("/api/todos/:id", proxyToTodoService)
	app.Get("/api/todos", proxyToTodoService)
	app.Put("/api/todos/:id", proxyToTodoService)
	app.Delete("/api/todos/:id", proxyToTodoService)

	// Aggregate endpoints
	app.Get("/api/dashboard", dashboardHandler)

	log.Println("🌐 API Gateway started on port 3000")
	log.Println("📡 Proxying requests to microservices:")
	log.Println("   👥 Users: " + UserServiceURL)
	log.Println("   📝 Todos: " + TodoServiceURL)
	log.Fatal(app.Listen(":3000"))
}

// Health check handler
func healthCheckHandler(c *fiber.Ctx) error {
	services := []HealthCheck{
		{Service: "api-gateway", Status: "healthy", URL: "localhost:3000"},
	}

	// Check User Service
	userHealth := checkServiceHealth(UserServiceURL + "/health")
	services = append(services, HealthCheck{
		Service: "user-service",
		Status:  userHealth,
		URL:     UserServiceURL,
	})

	// Check Todo Service
	todoHealth := checkServiceHealth(TodoServiceURL + "/health")
	services = append(services, HealthCheck{
		Service: "todo-service",
		Status:  todoHealth,
		URL:     TodoServiceURL,
	})

	// Overall status
	overallStatus := "healthy"
	for _, service := range services {
		if service.Status != "healthy" {
			overallStatus = "degraded"
			break
		}
	}

	return c.JSON(fiber.Map{
		"overall_status": overallStatus,
		"services":       services,
		"timestamp":      time.Now(),
	})
}

// User service proxy
func proxyToUserService(c *fiber.Ctx) error {
	return proxyRequest(c, UserServiceURL, "/users")
}

// Todo service proxy
func proxyToTodoService(c *fiber.Ctx) error {
	return proxyRequest(c, TodoServiceURL, "/todos")
}

// Dashboard aggregator
func dashboardHandler(c *fiber.Ctx) error {
	// Get data from both services concurrently
	usersChan := make(chan ServiceResponse)
	todosChan := make(chan ServiceResponse)

	// Fetch users
	go func() {
		resp, err := makeRequest("GET", UserServiceURL+"/users", nil)
		if err != nil {
			usersChan <- ServiceResponse{Success: false, Error: err.Error()}
			return
		}
		usersChan <- resp
	}()

	// Fetch todos
	go func() {
		resp, err := makeRequest("GET", TodoServiceURL+"/todos", nil)
		if err != nil {
			todosChan <- ServiceResponse{Success: false, Error: err.Error()}
			return
		}
		todosChan <- resp
	}()

	// Wait for both responses
	usersResp := <-usersChan
	todosResp := <-todosChan

	dashboard := fiber.Map{
		"users": usersResp,
		"todos": todosResp,
		"stats": fiber.Map{
			"generated_at": time.Now(),
		},
	}

	// Add some statistics if data is available
	if usersResp.Success && todosResp.Success {
		if usersData, ok := usersResp.Data.(map[string]interface{}); ok {
			if todosData, ok := todosResp.Data.(map[string]interface{}); ok {
				dashboard["stats"] = fiber.Map{
					"total_users":    usersData["count"],
					"total_todos":    todosData["count"],
					"generated_at":   time.Now(),
					"services_up":    2,
					"services_total": 2,
				}
			}
		}
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"dashboard": dashboard,
	})
}

// Generic proxy function
func proxyRequest(c *fiber.Ctx, serviceURL, basePath string) error {
	// Build target URL
	targetURL := serviceURL + basePath

	// Add path parameters
	if c.Params("id") != "" {
		targetURL += "/" + c.Params("id")
	}

	// Add query parameters
	if c.Request().URI().QueryString() != nil {
		targetURL += "?" + string(c.Request().URI().QueryString())
	}

	// Get request body
	var body io.Reader
	if c.Method() != "GET" && c.Method() != "DELETE" {
		body = bytes.NewReader(c.Body())
	}

	// Create request
	req, err := http.NewRequest(c.Method(), targetURL, body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to create request: " + err.Error(),
		})
	}

	// Copy headers
	req.Header.Set("Content-Type", "application/json")

	// Make request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(503).JSON(fiber.Map{
			"success": false,
			"error":   "Service unavailable: " + err.Error(),
		})
	}
	defer resp.Body.Close()

	// Read response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to read response: " + err.Error(),
		})
	}

	// Set response status and headers
	c.Status(resp.StatusCode)
	c.Set("Content-Type", "application/json")

	return c.Send(responseBody)
}

// Helper functions
func checkServiceHealth(url string) string {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "unhealthy"
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return "healthy"
	}
	return "unhealthy"
}

func makeRequest(method, url string, body interface{}) (ServiceResponse, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return ServiceResponse{}, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return ServiceResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return ServiceResponse{}, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ServiceResponse{}, err
	}

	var serviceResp ServiceResponse
	if err := json.Unmarshal(responseBody, &serviceResp); err != nil {
		return ServiceResponse{}, fmt.Errorf("failed to parse response: %w", err)
	}

	return serviceResp, nil
}
