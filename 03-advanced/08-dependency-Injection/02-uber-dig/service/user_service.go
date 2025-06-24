package service

import (
	"dig-di/models"
	"dig-di/repository"
	"errors"
	"strings"
)

var ErrInvalidUserID = errors.New("invalid user ID")

// UserService defines user business logic operations
type UserService interface {
	GetAllUsers() ([]*models.User, error)
	GetUserByID(id int) (*models.User, error)
	CreateUser(name, email string) (*models.User, error)
	UpdateUser(id int, name, email string) (*models.User, error)
	DeleteUser(id int) error
}

// userService implements UserService
type userService struct {
	repo repository.UserRepository // 🏗️ Dig จะ inject ให้อัตโนมัติ
}

// NewUserService creates a new user service
// 🏗️ Dig: Dig จะส่ง UserRepository มาให้อัตโนมัติ
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo, // Dig inject dependency ให้
	}
}

func (s *userService) GetAllUsers() ([]*models.User, error) {
	return s.repo.GetAll()
}

func (s *userService) GetUserByID(id int) (*models.User, error) {
	if id <= 0 {
		return nil, ErrInvalidUserID
	}
	return s.repo.GetByID(id)
}

func (s *userService) CreateUser(name, email string) (*models.User, error) {
	// Validate input
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(strings.ToLower(email))

	if name == "" || email == "" {
		return nil, errors.New("name and email are required")
	}

	// Check email uniqueness
	exists, err := s.repo.EmailExists(email, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, repository.ErrEmailAlreadyExists
	}

	// Create user
	user := models.NewUser(name, email)
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUser(id int, name, email string) (*models.User, error) {
	if id <= 0 {
		return nil, ErrInvalidUserID
	}

	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if name != "" {
		user.Name = strings.TrimSpace(name)
	}

	if email != "" {
		email = strings.TrimSpace(strings.ToLower(email))

		// Check email uniqueness
		exists, err := s.repo.EmailExists(email, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, repository.ErrEmailAlreadyExists
		}

		user.Email = email
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(id int) error {
	if id <= 0 {
		return ErrInvalidUserID
	}

	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}
