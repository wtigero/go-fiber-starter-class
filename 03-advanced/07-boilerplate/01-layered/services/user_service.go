package services

import (
	"layered-arch/models"
	"layered-arch/repositories"
	"strings"
)

// UserService handles user business logic
type UserService struct {
	userRepo *repositories.UserRepository
}

// NewUserService creates new user service
func NewUserService() *UserService {
	return &UserService{
		userRepo: repositories.NewUserRepository(),
	}
}

// GetAllUsers returns all users
func (s *UserService) GetAllUsers() []models.User {
	return s.userRepo.GetAll()
}

// GetUserByID returns user by ID
func (s *UserService) GetUserByID(id int) (*models.User, error) {
	if id <= 0 {
		return nil, ErrInvalidUserID
	}
	return s.userRepo.GetByID(id)
}

// CreateUser creates new user with business validation
func (s *UserService) CreateUser(req models.CreateUserRequest) (*models.User, error) {
	// Business validation
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Business logic: normalize email
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Name = strings.TrimSpace(req.Name)

	return s.userRepo.Create(req)
}

// UpdateUser updates existing user with business validation
func (s *UserService) UpdateUser(id int, req models.UpdateUserRequest) (*models.User, error) {
	if id <= 0 {
		return nil, ErrInvalidUserID
	}

	// Business validation
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	// Business logic: normalize data
	if req.Email != "" {
		req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	}
	if req.Name != "" {
		req.Name = strings.TrimSpace(req.Name)
	}

	return s.userRepo.Update(id, req)
}

// DeleteUser deletes user by ID
func (s *UserService) DeleteUser(id int) error {
	if id <= 0 {
		return ErrInvalidUserID
	}

	// Business rule: Check if user can be deleted
	// (In real app, might check for dependencies, permissions, etc.)

	return s.userRepo.Delete(id)
}

// validateCreateRequest validates create user request
func (s *UserService) validateCreateRequest(req models.CreateUserRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return ErrInvalidUserName
	}
	if len(req.Name) < 2 {
		return ErrUserNameTooShort
	}
	if len(req.Name) > 100 {
		return ErrUserNameTooLong
	}

	if strings.TrimSpace(req.Email) == "" {
		return ErrInvalidEmail
	}
	if !isValidEmail(req.Email) {
		return ErrInvalidEmailFormat
	}

	return nil
}

// validateUpdateRequest validates update user request
func (s *UserService) validateUpdateRequest(req models.UpdateUserRequest) error {
	if req.Name == "" && req.Email == "" {
		return ErrNoFieldsToUpdate
	}

	if req.Name != "" {
		if strings.TrimSpace(req.Name) == "" {
			return ErrInvalidUserName
		}
		if len(req.Name) < 2 {
			return ErrUserNameTooShort
		}
		if len(req.Name) > 100 {
			return ErrUserNameTooLong
		}
	}

	if req.Email != "" {
		if strings.TrimSpace(req.Email) == "" {
			return ErrInvalidEmail
		}
		if !isValidEmail(req.Email) {
			return ErrInvalidEmailFormat
		}
	}

	return nil
}

// isValidEmail performs basic email validation
func isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if len(email) == 0 {
		return false
	}
	// Basic email validation
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
