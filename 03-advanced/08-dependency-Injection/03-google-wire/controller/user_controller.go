package controller

import (
	"strconv"
	"wire-di/service"

	"github.com/gofiber/fiber/v2"
)

// UserController handles HTTP requests for users
type UserController struct {
	userService service.UserService // ⚡ Wire จะ generate การ inject
}

// NewUserController creates a new user controller
// ⚡ Wire: Wire จะ generate การส่ง UserService มาให้
func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService, // Wire จะ inject dependency ให้
	}
}

func (ctrl *UserController) GetAllUsers(c *fiber.Ctx) error {
	users, err := ctrl.userService.GetAllUsers()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": users})
}

func (ctrl *UserController) GetUserByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	user, err := ctrl.userService.GetUserByID(id)
	if err != nil {
		statusCode := 404
		if err.Error() != "user not found" {
			statusCode = 400
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": user})
}

func (ctrl *UserController) CreateUser(c *fiber.Ctx) error {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	user, err := ctrl.userService.CreateUser(req.Name, req.Email)
	if err != nil {
		statusCode := 400
		if err.Error() == "email already exists" {
			statusCode = 409
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"data": user})
}

func (ctrl *UserController) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	user, err := ctrl.userService.UpdateUser(id, req.Name, req.Email)
	if err != nil {
		statusCode := 400
		if err.Error() == "user not found" {
			statusCode = 404
		} else if err.Error() == "email already exists" {
			statusCode = 409
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": user})
}

func (ctrl *UserController) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	err = ctrl.userService.DeleteUser(id)
	if err != nil {
		statusCode := 400
		if err.Error() == "user not found" {
			statusCode = 404
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User deleted"})
}
