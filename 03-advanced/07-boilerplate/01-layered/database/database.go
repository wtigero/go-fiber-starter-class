package database

import (
	"layered-arch/models"
	"sync"
)

// Database represents in-memory database
type Database struct {
	users  []models.User
	nextID int
	mutex  sync.RWMutex
}

var (
	instance *Database
	once     sync.Once
)

// GetInstance returns singleton database instance
func GetInstance() *Database {
	once.Do(func() {
		instance = &Database{
			users:  make([]models.User, 0),
			nextID: 1,
		}
	})
	return instance
}

// GetUsers returns all users
func (db *Database) GetUsers() []models.User {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	// Return copy to avoid race conditions
	result := make([]models.User, len(db.users))
	copy(result, db.users)
	return result
}

// GetUserByID returns user by ID
func (db *Database) GetUserByID(id int) (*models.User, bool) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	for _, user := range db.users {
		if user.ID == id {
			return &user, true
		}
	}
	return nil, false
}

// CreateUser creates new user
func (db *Database) CreateUser(user models.User) models.User {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	user.ID = db.nextID
	db.nextID++
	db.users = append(db.users, user)
	return user
}

// UpdateUser updates existing user
func (db *Database) UpdateUser(id int, user models.User) (*models.User, bool) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	for i, u := range db.users {
		if u.ID == id {
			user.ID = id
			user.CreatedAt = u.CreatedAt // Keep original created time
			db.users[i] = user
			return &user, true
		}
	}
	return nil, false
}

// DeleteUser deletes user by ID
func (db *Database) DeleteUser(id int) bool {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	for i, user := range db.users {
		if user.ID == id {
			db.users = append(db.users[:i], db.users[i+1:]...)
			return true
		}
	}
	return false
}
