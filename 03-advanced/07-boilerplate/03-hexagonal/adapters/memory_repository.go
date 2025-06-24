package adapters

import (
	"hexagonal-arch/core/domain"
	"hexagonal-arch/core/ports"
	"sync"
)

// memoryUserRepository implements UserRepository port
type memoryUserRepository struct {
	users  []*domain.User
	nextID int
	mutex  sync.RWMutex
}

// NewMemoryUserRepository creates memory repository
func NewMemoryUserRepository() ports.UserRepository {
	return &memoryUserRepository{
		users:  []*domain.User{},
		nextID: 1,
	}
}

func (r *memoryUserRepository) GetAll() ([]*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make([]*domain.User, len(r.users))
	copy(result, r.users)
	return result, nil
}

func (r *memoryUserRepository) GetByID(id int) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func (r *memoryUserRepository) Create(user *domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	user.ID = r.nextID
	r.nextID++
	r.users = append(r.users, user)
	return nil
}

func (r *memoryUserRepository) Update(user *domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, u := range r.users {
		if u.ID == user.ID {
			r.users[i] = user
			return nil
		}
	}
	return domain.ErrUserNotFound
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
	return domain.ErrUserNotFound
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
