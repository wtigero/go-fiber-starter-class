package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fiber-unit-test/models"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(user *models.User) (*models.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUser(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetAllUsers() ([]*models.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(id string, user *models.User) (*models.User, error) {
	args := m.Called(id, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupTestApp(userHandler *UserHandler) *fiber.App {
	app := fiber.New()

	api := app.Group("/api/v1")
	users := api.Group("/users")

	users.Post("/", userHandler.CreateUser)
	users.Get("/", userHandler.GetAllUsers)
	users.Get("/:id", userHandler.GetUser)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)

	return app
}

func TestUserHandler_CreateUser(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockService := &MockUserService{}
		handler := NewUserHandler(mockService)
		app := setupTestApp(handler)

		user := &models.User{
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   25,
		}

		createdUser := &models.User{
			ID:    "user-id",
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   25,
		}

		mockService.On("CreateUser", mock.AnythingOfType("*models.User")).Return(createdUser, nil)

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/api/v1/users/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		var response models.User
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, createdUser.ID, response.ID)
		assert.Equal(t, createdUser.Name, response.Name)

		mockService.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockService := &MockUserService{}
		handler := NewUserHandler(mockService)
		app := setupTestApp(handler)

		req, _ := http.NewRequest("POST", "/api/v1/users/", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]string
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Invalid request body", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := &MockUserService{}
		handler := NewUserHandler(mockService)
		app := setupTestApp(handler)

		user := &models.User{
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   25,
		}

		mockService.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil, errors.New("email already exists"))

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/api/v1/users/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]string
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "email already exists", response["error"])

		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_GetUser(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		mockService := &MockUserService{}
		handler := NewUserHandler(mockService)
		app := setupTestApp(handler)

		expectedUser := &models.User{
			ID:    "user-id",
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   25,
		}

		mockService.On("GetUser", "user-id").Return(expectedUser, nil)

		req, _ := http.NewRequest("GET", "/api/v1/users/user-id", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response models.User
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, expectedUser.ID, response.ID)
		assert.Equal(t, expectedUser.Name, response.Name)

		mockService.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockService := &MockUserService{}
		handler := NewUserHandler(mockService)
		app := setupTestApp(handler)

		mockService.On("GetUser", "non-existent-id").Return(nil, errors.New("user not found"))

		req, _ := http.NewRequest("GET", "/api/v1/users/non-existent-id", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]string
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "user not found", response["error"])

		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_GetAllUsers(t *testing.T) {
	t.Run("successful get all", func(t *testing.T) {
		mockService := &MockUserService{}
		handler := NewUserHandler(mockService)
		app := setupTestApp(handler)

		expectedUsers := []*models.User{
			{ID: "1", Name: "John", Email: "john@example.com", Age: 25},
			{ID: "2", Name: "Jane", Email: "jane@example.com", Age: 30},
		}

		mockService.On("GetAllUsers").Return(expectedUsers, nil)

		req, _ := http.NewRequest("GET", "/api/v1/users/", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response []*models.User
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Len(t, response, 2)
		assert.Equal(t, expectedUsers[0].ID, response[0].ID)

		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := &MockUserService{}
		handler := NewUserHandler(mockService)
		app := setupTestApp(handler)

		mockService.On("GetAllUsers").Return(nil, errors.New("database error"))

		req, _ := http.NewRequest("GET", "/api/v1/users/", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var response map[string]string
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "database error", response["error"])

		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_UpdateUser(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		mockService := &MockUserService{}
		handler := NewUserHandler(mockService)
		app := setupTestApp(handler)

		user := &models.User{
			Name:  "John Smith",
			Email: "johnsmith@example.com",
			Age:   26,
		}

		updatedUser := &models.User{
			ID:    "user-id",
			Name:  "John Smith",
			Email: "johnsmith@example.com",
			Age:   26,
		}

		mockService.On("UpdateUser", "user-id", mock.AnythingOfType("*models.User")).Return(updatedUser, nil)

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("PUT", "/api/v1/users/user-id", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response models.User
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, updatedUser.ID, response.ID)
		assert.Equal(t, updatedUser.Name, response.Name)

		mockService.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockService := &MockUserService{}
		handler := NewUserHandler(mockService)
		app := setupTestApp(handler)

		req, _ := http.NewRequest("PUT", "/api/v1/users/user-id", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]string
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Invalid request body", response["error"])

		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_DeleteUser(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		mockService := &MockUserService{}
		handler := NewUserHandler(mockService)
		app := setupTestApp(handler)

		mockService.On("DeleteUser", "user-id").Return(nil)

		req, _ := http.NewRequest("DELETE", "/api/v1/users/user-id", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)

		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := &MockUserService{}
		handler := NewUserHandler(mockService)
		app := setupTestApp(handler)

		mockService.On("DeleteUser", "user-id").Return(errors.New("user not found"))

		req, _ := http.NewRequest("DELETE", "/api/v1/users/user-id", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]string
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "user not found", response["error"])

		mockService.AssertExpectations(t)
	})
}
