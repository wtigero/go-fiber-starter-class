package infrastructure

import (
	"onion-arch/application"
	"onion-arch/domain/entities"
	"sync"
)

// MemoryUserRepository implements application.UserRepository
type MemoryUserRepository struct {
	users  []*entities.User
	nextID int
	mutex  sync.RWMutex
}

// NewMemoryUserRepository creates memory repository
func NewMemoryUserRepository() application.UserRepository {
	return &MemoryUserRepository{
		users:  []*entities.User{},
		nextID: 1,
	}
}

func (r *MemoryUserRepository) GetAll() ([]*entities.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make([]*entities.User, len(r.users))
	copy(result, r.users)
	return result, nil
}

func (r *MemoryUserRepository) GetByID(id int) (*entities.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.ID() == id {
			return user, nil
		}
	}
	return nil, value_objects.ErrUserNotFound
}

func (r *MemoryUserRepository) Create(user *entities.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	user.SetID(r.nextID)
	r.nextID++
	r.users = append(r.users, user)
	return nil
}

func (r *MemoryUserRepository) Update(user *entities.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, u := range r.users {
		if u.ID() == user.ID() {
			r.users[i] = user
			return nil
		}
	}
	return value_objects.ErrUserNotFound
}

func (r *MemoryUserRepository) Delete(id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, user := range r.users {
		if user.ID() == id {
			r.users = append(r.users[:i], r.users[i+1:]...)
			return nil
		}
	}
	return value_objects.ErrUserNotFound
}

func (r *MemoryUserRepository) EmailExists(email string, excludeID int) (bool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Email() == email && user.ID() != excludeID {
			return true, nil
		}
	}
	return false, nil
}
