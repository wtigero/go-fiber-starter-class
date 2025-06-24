package services

import (
	"errors"
	"hexagonal-arch/core/domain"
	"hexagonal-arch/core/ports"
)

var ErrInvalidUserID = errors.New("invalid user ID")

// userService implements UserService port
type userService struct {
	userRepo ports.UserRepository
}

// NewUserService creates new user service
func NewUserService(userRepo ports.UserRepository) ports.UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetAllUsers() ([]*domain.User, error) {
	return s.userRepo.GetAll()
}

func (s *userService) GetUserByID(id int) (*domain.User, error) {
	if id <= 0 {
		return nil, ErrInvalidUserID
	}
	return s.userRepo.GetByID(id)
}

func (s *userService) CreateUser(name, email string) (*domain.User, error) {
	// Check email uniqueness
	exists, err := s.userRepo.EmailExists(email, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrEmailAlreadyExists
	}

	// Create user
	user, err := domain.NewUser(name, email)
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUser(id int, name, email string) (*domain.User, error) {
	if id <= 0 {
		return nil, ErrInvalidUserID
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check email uniqueness if changing
	if email != "" && email != user.Email {
		exists, err := s.userRepo.EmailExists(email, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, domain.ErrEmailAlreadyExists
		}
	}

	if err := user.Update(name, email); err != nil {
		return nil, err
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(id int) error {
	if id <= 0 {
		return ErrInvalidUserID
	}

	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	return s.userRepo.Delete(id)
}
