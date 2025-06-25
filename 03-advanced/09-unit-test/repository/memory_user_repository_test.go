package repository

import (
	"fiber-unit-test/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryUserRepository_Create(t *testing.T) {
	repo := NewMemoryUserRepository()

	user := &models.User{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
	}

	err := repo.Create(user)

	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)
	assert.False(t, user.CreatedAt.IsZero())
}

func TestMemoryUserRepository_GetByID(t *testing.T) {
	repo := NewMemoryUserRepository()

	// Create a user first
	user := &models.User{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
	}
	err := repo.Create(user)
	require.NoError(t, err)

	// Test getting the user
	retrievedUser, err := repo.GetByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, retrievedUser.ID)
	assert.Equal(t, user.Name, retrievedUser.Name)
	assert.Equal(t, user.Email, retrievedUser.Email)
	assert.Equal(t, user.Age, retrievedUser.Age)

	// Test getting non-existent user
	_, err = repo.GetByID("non-existent-id")
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}

func TestMemoryUserRepository_GetAll(t *testing.T) {
	repo := NewMemoryUserRepository()

	// Initially empty
	users, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Empty(t, users)

	// Add some users
	user1 := &models.User{Name: "John", Email: "john@example.com", Age: 25}
	user2 := &models.User{Name: "Jane", Email: "jane@example.com", Age: 30}

	err = repo.Create(user1)
	require.NoError(t, err)
	err = repo.Create(user2)
	require.NoError(t, err)

	// Get all users
	users, err = repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, users, 2)

	// Check that both users are present
	userIDs := []string{user1.ID, user2.ID}
	for _, user := range users {
		assert.Contains(t, userIDs, user.ID)
	}
}

func TestMemoryUserRepository_Update(t *testing.T) {
	repo := NewMemoryUserRepository()

	// Create a user first
	user := &models.User{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
	}
	err := repo.Create(user)
	require.NoError(t, err)

	originalCreatedAt := user.CreatedAt

	// Wait a moment to ensure time difference
	time.Sleep(1 * time.Millisecond)

	// Update the user
	updatedUser := &models.User{
		Name:  "John Smith",
		Email: "johnsmith@example.com",
		Age:   26,
	}

	err = repo.Update(user.ID, updatedUser)
	assert.NoError(t, err)

	// Verify the update
	retrievedUser, err := repo.GetByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, retrievedUser.ID)
	assert.Equal(t, "John Smith", retrievedUser.Name)
	assert.Equal(t, "johnsmith@example.com", retrievedUser.Email)
	assert.Equal(t, 26, retrievedUser.Age)
	assert.Equal(t, originalCreatedAt, retrievedUser.CreatedAt) // CreatedAt should not change

	// Test updating non-existent user
	err = repo.Update("non-existent-id", updatedUser)
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}

func TestMemoryUserRepository_Delete(t *testing.T) {
	repo := NewMemoryUserRepository()

	// Create a user first
	user := &models.User{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
	}
	err := repo.Create(user)
	require.NoError(t, err)

	// Delete the user
	err = repo.Delete(user.ID)
	assert.NoError(t, err)

	// Verify the user is deleted
	_, err = repo.GetByID(user.ID)
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())

	// Test deleting non-existent user
	err = repo.Delete("non-existent-id")
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}

func TestMemoryUserRepository_ConcurrencyHandling(t *testing.T) {
	repo := NewMemoryUserRepository()

	// Test concurrent operations
	done := make(chan bool, 10)

	// Create users concurrently
	for i := 0; i < 10; i++ {
		go func(index int) {
			user := &models.User{
				Name:  "User",
				Email: "user@example.com",
				Age:   25,
			}
			err := repo.Create(user)
			assert.NoError(t, err)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all users were created
	users, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, users, 10)
}
