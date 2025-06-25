package main

import (
	"bytes"
	"encoding/json"
	"fiber-unit-test/models"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserAPI_Integration(t *testing.T) {
	app := SetupApp()

	t.Run("complete user lifecycle", func(t *testing.T) {
		// Create a user
		user := models.User{
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   25,
		}

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/api/v1/users/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		var createdUser models.User
		err = json.NewDecoder(resp.Body).Decode(&createdUser)
		require.NoError(t, err)
		assert.NotEmpty(t, createdUser.ID)
		assert.Equal(t, user.Name, createdUser.Name)
		assert.Equal(t, user.Email, createdUser.Email)
		assert.Equal(t, user.Age, createdUser.Age)

		userID := createdUser.ID

		// Get the created user
		req, _ = http.NewRequest("GET", "/api/v1/users/"+userID, nil)
		resp, err = app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var retrievedUser models.User
		err = json.NewDecoder(resp.Body).Decode(&retrievedUser)
		require.NoError(t, err)
		assert.Equal(t, createdUser.ID, retrievedUser.ID)
		assert.Equal(t, createdUser.Name, retrievedUser.Name)

		// Update the user
		updatedUser := models.User{
			Name:  "John Smith",
			Email: "johnsmith@example.com",
			Age:   26,
		}

		body, _ = json.Marshal(updatedUser)
		req, _ = http.NewRequest("PUT", "/api/v1/users/"+userID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err = app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var responseUser models.User
		err = json.NewDecoder(resp.Body).Decode(&responseUser)
		require.NoError(t, err)
		assert.Equal(t, userID, responseUser.ID)
		assert.Equal(t, updatedUser.Name, responseUser.Name)
		assert.Equal(t, updatedUser.Email, responseUser.Email)

		// Get all users
		req, _ = http.NewRequest("GET", "/api/v1/users/", nil)
		resp, err = app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var users []*models.User
		err = json.NewDecoder(resp.Body).Decode(&users)
		require.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, userID, users[0].ID)

		// Delete the user
		req, _ = http.NewRequest("DELETE", "/api/v1/users/"+userID, nil)
		resp, err = app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)

		// Verify user is deleted
		req, _ = http.NewRequest("GET", "/api/v1/users/"+userID, nil)
		resp, err = app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("duplicate email validation", func(t *testing.T) {
		// Create first user
		user1 := models.User{
			Name:  "John Doe",
			Email: "duplicate@example.com",
			Age:   25,
		}

		body, _ := json.Marshal(user1)
		req, _ := http.NewRequest("POST", "/api/v1/users/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		// Try to create another user with same email
		user2 := models.User{
			Name:  "Jane Doe",
			Email: "duplicate@example.com", // Same email
			Age:   30,
		}

		body, _ = json.Marshal(user2)
		req, _ = http.NewRequest("POST", "/api/v1/users/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err = app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var errorResponse map[string]string
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "email already exists", errorResponse["error"])
	})

	t.Run("validation errors", func(t *testing.T) {
		tests := []struct {
			name          string
			user          models.User
			expectedError string
		}{
			{
				name:          "empty name",
				user:          models.User{Name: "", Email: "test@example.com", Age: 25},
				expectedError: "name is required",
			},
			{
				name:          "empty email",
				user:          models.User{Name: "Test", Email: "", Age: 25},
				expectedError: "email is required",
			},
			{
				name:          "negative age",
				user:          models.User{Name: "Test", Email: "test@example.com", Age: -1},
				expectedError: "age must be between 0 and 150",
			},
			{
				name:          "age too high",
				user:          models.User{Name: "Test", Email: "test@example.com", Age: 151},
				expectedError: "age must be between 0 and 150",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				body, _ := json.Marshal(tt.user)
				req, _ := http.NewRequest("POST", "/api/v1/users/", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")

				resp, err := app.Test(req)
				require.NoError(t, err)
				assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

				var errorResponse map[string]string
				err = json.NewDecoder(resp.Body).Decode(&errorResponse)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedError, errorResponse["error"])
			})
		}
	})
}
