package interfaces

import "clean-arch/domain/entities"

// GetAllUsersUseCase defines contract for getting all users
type GetAllUsersUseCase interface {
	Execute() ([]*entities.User, error)
}

// GetUserUseCase defines contract for getting user by ID
type GetUserUseCase interface {
	Execute(id int) (*entities.User, error)
}

// CreateUserUseCase defines contract for creating user
type CreateUserUseCase interface {
	Execute(req CreateUserRequest) (*entities.User, error)
}

// UpdateUserUseCase defines contract for updating user
type UpdateUserUseCase interface {
	Execute(id int, req UpdateUserRequest) (*entities.User, error)
}

// DeleteUserUseCase defines contract for deleting user
type DeleteUserUseCase interface {
	Execute(id int) error
}

// Request/Response DTOs
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}
