package repositories

import (
	"layered-arch/database"
	"layered-arch/models"
	"time"
)

// UserRepository handles user data operations
type UserRepository struct {
	db *database.Database
}

// NewUserRepository creates new user repository
func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.GetInstance(),
	}
}

// GetAll returns all users
func (r *UserRepository) GetAll() []models.User {
	return r.db.GetUsers()
}

// GetByID returns user by ID
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	user, found := r.db.GetUserByID(id)
	if !found {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// Create creates new user
func (r *UserRepository) Create(req models.CreateUserRequest) (*models.User, error) {
	// Check if email already exists
	users := r.db.GetUsers()
	for _, user := range users {
		if user.Email == req.Email {
			return nil, ErrEmailAlreadyExists
		}
	}

	// Create user
	user := models.User{
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdUser := r.db.CreateUser(user)
	return &createdUser, nil
}

// Update updates existing user
func (r *UserRepository) Update(id int, req models.UpdateUserRequest) (*models.User, error) {
	// Get existing user
	existingUser, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check email uniqueness if email is being updated
	if req.Email != "" && req.Email != existingUser.Email {
		users := r.db.GetUsers()
		for _, user := range users {
			if user.Email == req.Email && user.ID != id {
				return nil, ErrEmailAlreadyExists
			}
		}
	}

	// Update fields
	updatedUser := *existingUser
	if req.Name != "" {
		updatedUser.Name = req.Name
	}
	if req.Email != "" {
		updatedUser.Email = req.Email
	}
	updatedUser.UpdatedAt = time.Now()

	// Save to database
	result, found := r.db.UpdateUser(id, updatedUser)
	if !found {
		return nil, ErrUserNotFound
	}

	return result, nil
}

// Delete deletes user by ID
func (r *UserRepository) Delete(id int) error {
	deleted := r.db.DeleteUser(id)
	if !deleted {
		return ErrUserNotFound
	}
	return nil
}
