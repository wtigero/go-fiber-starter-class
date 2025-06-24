package user

import (
	"clean-arch/domain/repositories"
	"clean-arch/usecases/interfaces"
)

// deleteUserUseCase implements DeleteUserUseCase interface
type deleteUserUseCase struct {
	userRepo repositories.UserRepository
}

// NewDeleteUserUseCase creates new instance of delete user use case
func NewDeleteUserUseCase(userRepo repositories.UserRepository) interfaces.DeleteUserUseCase {
	return &deleteUserUseCase{
		userRepo: userRepo,
	}
}

// Execute deletes user by ID
func (uc *deleteUserUseCase) Execute(id int) error {
	if id <= 0 {
		return ErrInvalidUserID
	}

	// Check if user exists
	_, err := uc.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Business rule: Check if user can be deleted
	// (In real app, might check for dependencies, permissions, etc.)

	return uc.userRepo.Delete(id)
}
