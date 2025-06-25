package service

import (
	"errors"
	"fiber-unit-test/models"
)

type userService struct {
	userRepo models.UserRepository
}

func NewUserService(userRepo models.UserRepository) models.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(user *models.User) (*models.User, error) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// Business logic: check if email already exists
	allUsers, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	for _, existingUser := range allUsers {
		if existingUser.Email == user.Email {
			return nil, errors.New("email already exists")
		}
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUser(id string) (*models.User, error) {
	if id == "" {
		return nil, errors.New("user ID is required")
	}

	return s.userRepo.GetByID(id)
}

func (s *userService) GetAllUsers() ([]*models.User, error) {
	return s.userRepo.GetAll()
}

func (s *userService) UpdateUser(id string, user *models.User) (*models.User, error) {
	if id == "" {
		return nil, errors.New("user ID is required")
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	// Check if user exists
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Business logic: check if email already exists for other users
	allUsers, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	for _, existingUser := range allUsers {
		if existingUser.ID != id && existingUser.Email == user.Email {
			return nil, errors.New("email already exists")
		}
	}

	if err := s.userRepo.Update(id, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(id string) error {
	if id == "" {
		return errors.New("user ID is required")
	}

	return s.userRepo.Delete(id)
}
