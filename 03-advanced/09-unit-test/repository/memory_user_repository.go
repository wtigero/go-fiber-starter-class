package repository

import (
	"errors"
	"fiber-unit-test/models"
	"sync"
	"time"

	"github.com/google/uuid"
)

type MemoryUserRepository struct {
	users map[string]*models.User
	mutex sync.RWMutex
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: make(map[string]*models.User),
	}
}

func (r *MemoryUserRepository) Create(user *models.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	user.CreatedAt = time.Now()

	r.users[user.ID] = user
	return nil
}

func (r *MemoryUserRepository) GetByID(id string) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *MemoryUserRepository) GetAll() ([]*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	users := make([]*models.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

func (r *MemoryUserRepository) Update(id string, user *models.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	existingUser, exists := r.users[id]
	if !exists {
		return errors.New("user not found")
	}

	user.ID = id
	user.CreatedAt = existingUser.CreatedAt
	r.users[id] = user
	return nil
}

func (r *MemoryUserRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[id]; !exists {
		return errors.New("user not found")
	}

	delete(r.users, id)
	return nil
}
