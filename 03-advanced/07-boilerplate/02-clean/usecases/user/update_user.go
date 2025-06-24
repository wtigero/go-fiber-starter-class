package user

import (
	"clean-arch/domain/entities"
	"clean-arch/domain/repositories"
	"clean-arch/usecases/interfaces"
	"errors"
)

var ErrNoFieldsToUpdate = errors.New("no fields to update")

// updateUserUseCase implements UpdateUserUseCase interface
type updateUserUseCase struct {
	userRepo repositories.UserRepository
}

// NewUpdateUserUseCase creates new instance of update user use case
func NewUpdateUserUseCase(userRepo repositories.UserRepository) interfaces.UpdateUserUseCase {
	return &updateUserUseCase{
		userRepo: userRepo,
	}
}

// Execute updates an existing user
func (uc *updateUserUseCase) Execute(id int, req interfaces.UpdateUserRequest) (*entities.User, error) {
	if id <= 0 {
		return nil, ErrInvalidUserID
	}

	if req.Name == "" && req.Email == "" {
		return nil, ErrNoFieldsToUpdate
	}

	// Get existing user
	user, err := uc.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check email uniqueness if email is being updated
	if req.Email != "" && req.Email != user.Email() {
		exists, err := uc.userRepo.EmailExists(req.Email, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, entities.ErrEmailAlreadyExists
		}
	}

	// Update user (business logic validation happens in entity)
	if err := user.Update(req.Name, req.Email); err != nil {
		return nil, err
	}

	// Save to repository
	if err := uc.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
