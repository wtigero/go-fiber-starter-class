package ports

import "hexagonal-arch/core/domain"

// PRIMARY PORTS (driving) - what external world can do with our app

// UserService defines use cases
type UserService interface {
	GetAllUsers() ([]*domain.User, error)
	GetUserByID(id int) (*domain.User, error)
	CreateUser(name, email string) (*domain.User, error)
	UpdateUser(id int, name, email string) (*domain.User, error)
	DeleteUser(id int) error
}

// SECONDARY PORTS (driven) - what our app needs from external world

// UserRepository defines storage interface
type UserRepository interface {
	GetAll() ([]*domain.User, error)
	GetByID(id int) (*domain.User, error)
	Create(user *domain.User) error
	Update(user *domain.User) error
	Delete(id int) error
	EmailExists(email string, excludeID int) (bool, error)
}
