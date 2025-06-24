package repositories

import "clean-arch/domain/entities"

// UserRepository defines the contract for user data operations
// This is an interface (dependency inversion)
type UserRepository interface {
	// GetAll returns all users
	GetAll() ([]*entities.User, error)

	// GetByID returns user by ID
	GetByID(id int) (*entities.User, error)

	// GetByEmail returns user by email
	GetByEmail(email string) (*entities.User, error)

	// Create creates new user and returns it with assigned ID
	Create(user *entities.User) error

	// Update updates existing user
	Update(user *entities.User) error

	// Delete deletes user by ID
	Delete(id int) error

	// EmailExists checks if email already exists (excluding user with given ID)
	EmailExists(email string, excludeID int) (bool, error)
}
