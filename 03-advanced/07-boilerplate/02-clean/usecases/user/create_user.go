package user

import (
	"clean-arch/domain/entities"
	"clean-arch/domain/repositories"
	"clean-arch/usecases/interfaces"
)

// createUserUseCase implements CreateUserUseCase interface
type createUserUseCase struct {
	userRepo repositories.UserRepository
}

// NewCreateUserUseCase creates new instance of create user use case
func NewCreateUserUseCase(userRepo repositories.UserRepository) interfaces.CreateUserUseCase {
	return &createUserUseCase{
		userRepo: userRepo,
	}
}

// Execute creates a new user
func (uc *createUserUseCase) Execute(req interfaces.CreateUserRequest) (*entities.User, error) {
	// Check if email already exists
	exists, err := uc.userRepo.EmailExists(req.Email, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, entities.ErrEmailAlreadyExists
	}

	// Create new user entity (business logic validation happens here)
	user, err := entities.NewUser(req.Name, req.Email)
	if err != nil {
		return nil, err
	}

	// Save to repository
	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
