package models

import (
	"errors"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	if u.Age < 0 || u.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetAll() ([]*User, error)
	Update(id string, user *User) error
	Delete(id string) error
}

type UserService interface {
	CreateUser(user *User) (*User, error)
	GetUser(id string) (*User, error)
	GetAllUsers() ([]*User, error)
	UpdateUser(id string, user *User) (*User, error)
	DeleteUser(id string) error
}
