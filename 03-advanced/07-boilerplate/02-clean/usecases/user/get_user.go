package user

import (
	"clean-arch/domain/entities"
	"clean-arch/domain/repositories"
	"clean-arch/usecases/interfaces"
	"errors"
)

var ErrInvalidUserID = errors.New("invalid user ID")

// getUserUseCase implements GetUserUseCase interface
type getUserUseCase struct {
	userRepo repositories.UserRepository
}

// NewGetUserUseCase creates new instance of get user use case
func NewGetUserUseCase(userRepo repositories.UserRepository) interfaces.GetUserUseCase {
	return &getUserUseCase{
		userRepo: userRepo,
	}
}

// Execute retrieves user by ID
func (uc *getUserUseCase) Execute(id int) (*entities.User, error) {
	if id <= 0 {
		return nil, ErrInvalidUserID
	}

	return uc.userRepo.GetByID(id)
}
