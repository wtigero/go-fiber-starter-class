package repository

import (
	"dig-di/models"
	"errors"
	"sync"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

// UserRepository defines user data operations
type UserRepository interface {
	GetAll() ([]*models.User, error)
	GetByID(id int) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id int) error
	EmailExists(email string, excludeID int) (bool, error)
}

// memoryUserRepository implements UserRepository
type memoryUserRepository struct {
	users  []*models.User
	nextID int
	mutex  sync.RWMutex
}

// NewUserRepository creates a new user repository
// 🏗️ Dig: เพื่อให้ Dig เรียกใช้ได้ ต้อง return interface
func NewUserRepository() UserRepository {
	return &memoryUserRepository{
		users:  []*models.User{},
		nextID: 1,
	}
}

func (r *memoryUserRepository) GetAll() ([]*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make([]*models.User, len(r.users))
	copy(result, r.users)
	return result, nil
}

func (r *memoryUserRepository) GetByID(id int) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

func (r *memoryUserRepository) Create(user *models.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	user.ID = r.nextID
	r.nextID++
	r.users = append(r.users, user)
	return nil
}

func (r *memoryUserRepository) Update(user *models.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, u := range r.users {
		if u.ID == user.ID {
			r.users[i] = user
			return nil
		}
	}
	return ErrUserNotFound
}

func (r *memoryUserRepository) Delete(id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, user := range r.users {
		if user.ID == id {
			r.users = append(r.users[:i], r.users[i+1:]...)
			return nil
		}
	}
	return ErrUserNotFound
}

func (r *memoryUserRepository) EmailExists(email string, excludeID int) (bool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Email == email && user.ID != excludeID {
			return true, nil
		}
	}
	return false, nil
}
