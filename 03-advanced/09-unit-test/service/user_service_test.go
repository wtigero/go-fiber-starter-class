package service

import (
	"errors"
	"fiber-unit-test/mocks"
	"fiber-unit-test/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserService_CreateUser(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository()
		service := NewUserService(mockRepo)

		user := &models.User{
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   25,
		}

		// Mock GetAll to return empty list (no existing users)
		mockRepo.On("GetAll").Return([]*models.User{}, nil)
		mockRepo.On("Create", user).Return(nil)

		result, err := service.CreateUser(user)

		assert.NoError(t, err)
		assert.Equal(t, user, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository()
		service := NewUserService(mockRepo)

		user := &models.User{
			Name:  "", // Invalid: empty name
			Email: "john@example.com",
			Age:   25,
		}

		result, err := service.CreateUser(user)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "name is required", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("email already exists", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository()
		service := NewUserService(mockRepo)

		user := &models.User{
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   25,
		}

		existingUser := &models.User{
			ID:    "existing-id",
			Name:  "Jane Doe",
			Email: "john@example.com", // Same email
			Age:   30,
		}

		mockRepo.On("GetAll").Return([]*models.User{existingUser}, nil)

		result, err := service.CreateUser(user)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "email already exists", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error on GetAll", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository()
		service := NewUserService(mockRepo)

		user := &models.User{
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   25,
		}

		mockRepo.On("GetAll").Return(nil, errors.New("database error"))

		result, err := service.CreateUser(user)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "database error", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error on Create", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository()
		service := NewUserService(mockRepo)

		user := &models.User{
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   25,
		}

		mockRepo.On("GetAll").Return([]*models.User{}, nil)
		mockRepo.On("Create", user).Return(errors.New("database error"))

		result, err := service.CreateUser(user)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "database error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetUser(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository()
		service := NewUserService(mockRepo)

		expectedUser := &models.User{
			ID:    "user-id",
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   25,
		}

		mockRepo.On("GetByID", "user-id").Return(expectedUser, nil)

		result, err := service.GetUser("user-id")

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty ID", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository()
		service := NewUserService(mockRepo)

		result, err := service.GetUser("")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "user ID is required", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository()
		service := NewUserService(mockRepo)

		mockRepo.On("GetByID", "non-existent-id").Return(nil, errors.New("user not found"))

		result, err := service.GetUser("non-existent-id")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "user not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetAllUsers(t *testing.T) {
	t.Run("successful get all", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository()
		service := NewUserService(mockRepo)

		expectedUsers := []*models.User{
			{ID: "1", Name: "John", Email: "john@example.com", Age: 25},
			{ID: "2", Name: "Jane", Email: "jane@example.com", Age: 30},
		}

		mockRepo.On("GetAll").Return(expectedUsers, nil)

		result, err := service.GetAllUsers()

		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository()
		service := NewUserService(mockRepo)

		mockRepo.On("GetAll").Return(nil, errors.New("database error"))

		result, err := service.GetAllUsers()

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "database error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository()
		service := NewUserService(mockRepo)

		existingUser := &models.User{
			ID:    "user-id",
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   25,
		}

		updatedUser := &models.User{
			Name:  "John Smith",
			Email: "johnsmith@example.com",
			Age:   26,
		}

		mockRepo.On("GetByID", "user-id").Return(existingUser, nil)
		mockRepo.On("GetAll").Return([]*models.User{existingUser}, nil)
		mockRepo.On("Update", "user-id", updatedUser).Return(nil)

		result, err := service.UpdateUser("user-id", updatedUser)

		assert.NoError(t, err)
		assert.Equal(t, updatedUser, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty ID", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository()
		service := NewUserService(mockRepo)

		user := &models.User{Name: "John", Email: "john@example.com", Age: 25}

		result, err := service.UpdateUser("", user)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "user ID is required", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository()
		service := NewUserService(mockRepo)

		user := &models.User{
			Name:  "", // Invalid
			Email: "john@example.com",
			Age:   25,
		}

		result, err := service.UpdateUser("user-id", user)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "name is required", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository()
		service := NewUserService(mockRepo)

		user := &models.User{Name: "John", Email: "john@example.com", Age: 25}

		mockRepo.On("GetByID", "non-existent-id").Return(nil, errors.New("user not found"))

		result, err := service.UpdateUser("non-existent-id", user)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "user not found", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("email already exists for another user", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository()
		service := NewUserService(mockRepo)

		existingUser := &models.User{
			ID:    "user-id",
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   25,
		}

		anotherUser := &models.User{
			ID:    "another-id",
			Name:  "Jane Doe",
			Email: "jane@example.com",
			Age:   30,
		}

		updatedUser := &models.User{
			Name:  "John Smith",
			Email: "jane@example.com", // Same as another user
			Age:   26,
		}

		mockRepo.On("GetByID", "user-id").Return(existingUser, nil)
		mockRepo.On("GetAll").Return([]*models.User{existingUser, anotherUser}, nil)

		result, err := service.UpdateUser("user-id", updatedUser)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "email already exists", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository()
		service := NewUserService(mockRepo)

		mockRepo.On("Delete", "user-id").Return(nil)

		err := service.DeleteUser("user-id")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty ID", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository()
		service := NewUserService(mockRepo)

		err := service.DeleteUser("")

		assert.Error(t, err)
		assert.Equal(t, "user ID is required", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository()
		service := NewUserService(mockRepo)

		mockRepo.On("Delete", "user-id").Return(errors.New("database error"))

		err := service.DeleteUser("user-id")

		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
