package controllers

import (
	"layered-arch/models"
	"layered-arch/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// UserController handles HTTP requests for users
type UserController struct {
	userService *services.UserService
}

// NewUserController creates new user controller
func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

// GetAllUsers handles GET /users
func (c *UserController) GetAllUsers(ctx *fiber.Ctx) error {
	users := c.userService.GetAllUsers()
	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    users,
		"count":   len(users),
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

	user, err := c.userService.GetUserByID(id)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

// CreateUser handles POST /users
func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	user, err := c.userService.CreateUser(req)
	if err != nil {
		// Determine status code based on error type
		statusCode := 400
		if err.Error() == "email already exists" {
			statusCode = 409
		}

		return ctx.Status(statusCode).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return ctx.Status(201).JSON(fiber.Map{
		"success": true,
		"data":    user,
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

	var req models.UpdateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	user, err := c.userService.UpdateUser(id, req)
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

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    user,
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

	err = c.userService.DeleteUser(id)
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
