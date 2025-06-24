package database

import (
	"clean-arch/domain/entities"
	"clean-arch/domain/repositories"
	"sync"
)

// memoryUserRepository implements UserRepository interface
type memoryUserRepository struct {
	users  []*entities.User
	nextID int
	mutex  sync.RWMutex
}

// NewMemoryUserRepository creates new in-memory user repository
func NewMemoryUserRepository() repositories.UserRepository {
	return &memoryUserRepository{
		users:  make([]*entities.User, 0),
		nextID: 1,
	}
}

// GetAll returns all users
func (r *memoryUserRepository) GetAll() ([]*entities.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Return copy to avoid race conditions
	result := make([]*entities.User, len(r.users))
	copy(result, r.users)
	return result, nil
}

// GetByID returns user by ID
func (r *memoryUserRepository) GetByID(id int) (*entities.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.ID() == id {
			return user, nil
		}
	}
	return nil, entities.ErrUserNotFound
}

// GetByEmail returns user by email
func (r *memoryUserRepository) GetByEmail(email string) (*entities.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Email() == email {
			return user, nil
		}
	}
	return nil, entities.ErrUserNotFound
}

// Create creates new user and assigns ID
func (r *memoryUserRepository) Create(user *entities.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	user.SetID(r.nextID)
	r.nextID++
	r.users = append(r.users, user)
	return nil
}

// Update updates existing user
func (r *memoryUserRepository) Update(user *entities.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, u := range r.users {
		if u.ID() == user.ID() {
			r.users[i] = user
			return nil
		}
	}
	return entities.ErrUserNotFound
}

// Delete deletes user by ID
func (r *memoryUserRepository) Delete(id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, user := range r.users {
		if user.ID() == id {
			r.users = append(r.users[:i], r.users[i+1:]...)
			return nil
		}
	}
	return entities.ErrUserNotFound
}

// EmailExists checks if email already exists (excluding user with given ID)
func (r *memoryUserRepository) EmailExists(email string, excludeID int) (bool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Email() == email && user.ID() != excludeID {
			return true, nil
		}
	}
	return false, nil
}
