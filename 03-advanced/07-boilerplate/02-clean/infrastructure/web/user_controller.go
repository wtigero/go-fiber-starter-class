package web

import (
	"clean-arch/usecases/interfaces"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// UserController handles HTTP requests for users
type UserController struct {
	getAllUseCase  interfaces.GetAllUsersUseCase
	getUserUseCase interfaces.GetUserUseCase
	createUseCase  interfaces.CreateUserUseCase
	updateUseCase  interfaces.UpdateUserUseCase
	deleteUseCase  interfaces.DeleteUserUseCase
}

// NewUserController creates new user controller
func NewUserController(
	getAllUseCase interfaces.GetAllUsersUseCase,
	getUserUseCase interfaces.GetUserUseCase,
	createUseCase interfaces.CreateUserUseCase,
	updateUseCase interfaces.UpdateUserUseCase,
	deleteUseCase interfaces.DeleteUserUseCase,
) *UserController {
	return &UserController{
		getAllUseCase:  getAllUseCase,
		getUserUseCase: getUserUseCase,
		createUseCase:  createUseCase,
		updateUseCase:  updateUseCase,
		deleteUseCase:  deleteUseCase,
	}
}

// GetAllUsers handles GET /users
func (c *UserController) GetAllUsers(ctx *fiber.Ctx) error {
	users, err := c.getAllUseCase.Execute()
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Convert entities to DTOs
	userDTOs := make([]fiber.Map, len(users))
	for i, user := range users {
		userDTOs[i] = fiber.Map{
			"id":         user.ID(),
			"name":       user.Name(),
			"email":      user.Email(),
			"created_at": user.CreatedAt(),
			"updated_at": user.UpdatedAt(),
		}
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    userDTOs,
		"count":   len(userDTOs),
	})
}

// GetUserByID handles GET /users/:id
func (c *UserController) GetUserByID(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid user ID format",
		})
	}

	user, err := c.getUserUseCase.Execute(id)
	if err != nil {
		statusCode := 400
		if err.Error() == "user not found" {
			statusCode = 404
		}
		return ctx.Status(statusCode).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Convert entity to DTO
	userDTO := fiber.Map{
		"id":         user.ID(),
		"name":       user.Name(),
		"email":      user.Email(),
		"created_at": user.CreatedAt(),
		"updated_at": user.UpdatedAt(),
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    userDTO,
	})
}

// CreateUser handles POST /users
func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	var req interfaces.CreateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	user, err := c.createUseCase.Execute(req)
	if err != nil {
		statusCode := 400
		if err.Error() == "email already exists" {
			statusCode = 409
		}
		return ctx.Status(statusCode).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Convert entity to DTO
	userDTO := fiber.Map{
		"id":         user.ID(),
		"name":       user.Name(),
		"email":      user.Email(),
		"created_at": user.CreatedAt(),
		"updated_at": user.UpdatedAt(),
	}

	return ctx.Status(201).JSON(fiber.Map{
		"success": true,
		"data":    userDTO,
		"message": "User created successfully",
	})
}

// UpdateUser handles PUT /users/:id
func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid user ID format",
		})
	}

	var req interfaces.UpdateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	user, err := c.updateUseCase.Execute(id, req)
	if err != nil {
		statusCode := 400
		if err.Error() == "user not found" {
			statusCode = 404
		} else if err.Error() == "email already exists" {
			statusCode = 409
		}
		return ctx.Status(statusCode).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Convert entity to DTO
	userDTO := fiber.Map{
		"id":         user.ID(),
		"name":       user.Name(),
		"email":      user.Email(),
		"created_at": user.CreatedAt(),
		"updated_at": user.UpdatedAt(),
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    userDTO,
		"message": "User updated successfully",
	})
}

// DeleteUser handles DELETE /users/:id
func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid user ID format",
		})
	}

	err = c.deleteUseCase.Execute(id)
	if err != nil {
		statusCode := 400
		if err.Error() == "user not found" {
			statusCode = 404
		}
		return ctx.Status(statusCode).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "User deleted successfully",
	})
}
