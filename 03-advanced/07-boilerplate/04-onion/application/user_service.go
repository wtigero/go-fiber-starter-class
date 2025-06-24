package application

import (
	"errors"
	"onion-arch/domain/entities"
)

var ErrInvalidUserID = errors.New("invalid user ID")

// UserRepository defines contract
type UserRepository interface {
	GetAll() ([]*entities.User, error)
	GetByID(id int) (*entities.User, error)
	Create(user *entities.User) error
	Update(user *entities.User) error
	Delete(id int) error
	EmailExists(email string, excludeID int) (bool, error)
}

// UserService handles business logic
type UserService struct {
	repo UserRepository
}

// NewUserService creates user service
func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAllUsers() ([]*entities.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) GetUserByID(id int) (*entities.User, error) {
	userID, err := value_objects.NewUserID(id)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(userID.Value())
}

func (s *UserService) CreateUser(name, email string) (*entities.User, error) {
	// Validate business rules
	userName, err := value_objects.NewUserName(name)
	if err != nil {
		return nil, err
	}

	userEmail, err := value_objects.NewUserEmail(email)
	if err != nil {
		return nil, err
	}

	// Check email uniqueness
	exists, err := s.repo.EmailExists(userEmail.Value(), 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, value_objects.ErrEmailAlreadyExists
	}

	// Create user
	user := entities.NewUser(userName, userEmail)
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateUser(id int, name, email string) (*entities.User, error) {
	userID, err := value_objects.NewUserID(id)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetByID(userID.Value())
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if name != "" {
		userName, err := value_objects.NewUserName(name)
		if err != nil {
			return nil, err
		}
		user.UpdateName(userName)
	}

	if email != "" {
		userEmail, err := value_objects.NewUserEmail(email)
		if err != nil {
			return nil, err
		}

		// Check email uniqueness
		exists, err := s.repo.EmailExists(userEmail.Value(), userID.Value())
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, value_objects.ErrEmailAlreadyExists
		}

		user.UpdateEmail(userEmail)
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(id int) error {
	userID, err := value_objects.NewUserID(id)
	if err != nil {
		return err
	}

	_, err = s.repo.GetByID(userID.Value())
	if err != nil {
		return err
	}

	return s.repo.Delete(userID.Value())
}
