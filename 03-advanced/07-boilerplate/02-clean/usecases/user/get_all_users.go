package user

import (
	"clean-arch/domain/entities"
	"clean-arch/domain/repositories"
	"clean-arch/usecases/interfaces"
)

// getAllUsersUseCase implements GetAllUsersUseCase interface
type getAllUsersUseCase struct {
	userRepo repositories.UserRepository
}

// NewGetAllUsersUseCase creates new instance of get all users use case
func NewGetAllUsersUseCase(userRepo repositories.UserRepository) interfaces.GetAllUsersUseCase {
	return &getAllUsersUseCase{
		userRepo: userRepo,
	}
}

// Execute retrieves all users
func (uc *getAllUsersUseCase) Execute() ([]*entities.User, error) {
	return uc.userRepo.GetAll()
}
