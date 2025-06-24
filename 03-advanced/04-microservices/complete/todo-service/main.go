package main

import (
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Todo struct
type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// In-memory storage
var todos []Todo
var todoID = 1

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

	// Initialize sample data
	initSampleData()

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "todo-service",
			"status":  "healthy",
			"port":    3002,
			"todos":   len(todos),
		})
	})

	// Basic info endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service":     "Todo Service",
			"version":     "1.0.0",
			"status":      "running",
			"total_todos": len(todos),
		})
	})

	// API endpoints
	app.Post("/todos", createTodoHandler)
	app.Get("/todos/:id", getTodoHandler)
	app.Get("/todos", getAllTodosHandler)
	app.Put("/todos/:id", updateTodoHandler)
	app.Delete("/todos/:id", deleteTodoHandler)

	log.Println("📝 Todo Service started on port 3002")
	log.Fatal(app.Listen(":3002"))
}

func initSampleData() {
	now := time.Now()
	todos = []Todo{
		{
			ID:          1,
			Title:       "เรียน Go Fiber",
			Description: "ศึกษา framework สำหรับสร้าง API",
			Completed:   true,
			UserID:      1,
			CreatedAt:   now.Add(-48 * time.Hour),
			UpdatedAt:   now.Add(-24 * time.Hour),
		},
		{
			ID:          2,
			Title:       "สร้าง Microservices",
			Description: "พัฒนาระบบ microservices ด้วย Go",
			Completed:   false,
			UserID:      1,
			CreatedAt:   now.Add(-24 * time.Hour),
			UpdatedAt:   now.Add(-24 * time.Hour),
		},
		{
			ID:          3,
			Title:       "เขียน Documentation",
			Description: "จัดทำเอกสารสำหรับ API",
			Completed:   false,
			UserID:      2,
			CreatedAt:   now.Add(-12 * time.Hour),
			UpdatedAt:   now.Add(-12 * time.Hour),
		},
		{
			ID:          4,
			Title:       "ทดสอบ API Gateway",
			Description: "ทดสอบการทำงานของ API Gateway",
			Completed:   true,
			UserID:      2,
			CreatedAt:   now.Add(-6 * time.Hour),
			UpdatedAt:   now.Add(-1 * time.Hour),
		},
	}
	todoID = 5
	log.Println("✅ Sample todos loaded")
}

func createTodoHandler(c *fiber.Ctx) error {
	type CreateTodoRequest struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		UserID      int    `json:"user_id"`
	}

	var req CreateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request",
		})
	}

	if req.Title == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Title is required",
		})
	}

	if req.UserID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "User ID is required",
		})
	}

	todo := Todo{
		ID:          todoID,
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
		UserID:      req.UserID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	todos = append(todos, todo)
	todoID++

	log.Printf("📝 Created todo: %s for user %d", todo.Title, todo.UserID)

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"data":    todo,
		"message": "Todo created successfully",
	})
}

func getTodoHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	for _, todo := range todos {
		if strconv.Itoa(todo.ID) == id {
			return c.JSON(fiber.Map{
				"success": true,
				"data":    todo,
			})
		}
	}

	return c.Status(404).JSON(fiber.Map{
		"success": false,
		"error":   "Todo not found",
	})
}

func getAllTodosHandler(c *fiber.Ctx) error {
	// Support filtering by user
	userIDStr := c.Query("user_id")
	completed := c.Query("completed")

	// Support pagination
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Filter todos
	var filteredTodos []Todo
	for _, todo := range todos {
		// Filter by user_id
		if userIDStr != "" {
			userID, err := strconv.Atoi(userIDStr)
			if err == nil && todo.UserID != userID {
				continue
			}
		}

		// Filter by completed status
		if completed != "" {
			if completed == "true" && !todo.Completed {
				continue
			}
			if completed == "false" && todo.Completed {
				continue
			}
		}

		filteredTodos = append(filteredTodos, todo)
	}

	// Apply pagination
	start := (page - 1) * limit
	end := start + limit

	var paginatedTodos []Todo
	if start < len(filteredTodos) {
		if end > len(filteredTodos) {
			end = len(filteredTodos)
		}
		paginatedTodos = filteredTodos[start:end]
	} else {
		paginatedTodos = []Todo{}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    paginatedTodos,
		"pagination": fiber.Map{
			"page":     page,
			"limit":    limit,
			"total":    len(filteredTodos),
			"has_next": end < len(filteredTodos),
			"has_prev": page > 1,
		},
		"count": len(paginatedTodos),
	})
}

func updateTodoHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	type UpdateTodoRequest struct {
		Title       *string `json:"title,omitempty"`
		Description *string `json:"description,omitempty"`
		Completed   *bool   `json:"completed,omitempty"`
	}

	var req UpdateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request",
		})
	}

	for i, todo := range todos {
		if strconv.Itoa(todo.ID) == id {
			// Update fields if provided
			if req.Title != nil {
				todos[i].Title = *req.Title
			}
			if req.Description != nil {
				todos[i].Description = *req.Description
			}
			if req.Completed != nil {
				todos[i].Completed = *req.Completed
			}
			todos[i].UpdatedAt = time.Now()

			log.Printf("📝 Updated todo: %s", todos[i].Title)

			return c.JSON(fiber.Map{
				"success": true,
				"data":    todos[i],
				"message": "Todo updated successfully",
			})
		}
	}

	return c.Status(404).JSON(fiber.Map{
		"success": false,
		"error":   "Todo not found",
	})
}

func deleteTodoHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	for i, todo := range todos {
		if strconv.Itoa(todo.ID) == id {
			// Remove todo from slice
			todos = append(todos[:i], todos[i+1:]...)

			log.Printf("🗑️ Deleted todo: %s", todo.Title)

			return c.JSON(fiber.Map{
				"success": true,
				"message": "Todo deleted successfully",
			})
		}
	}

	return c.Status(404).JSON(fiber.Map{
		"success": false,
		"error":   "Todo not found",
	})
}
