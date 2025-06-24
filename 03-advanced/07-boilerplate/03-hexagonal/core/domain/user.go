package domain

import (
	"errors"
	"strings"
	"time"
)

// User represents core business entity
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Domain errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidInput       = errors.New("invalid input")
)

// NewUser creates new user with validation
func NewUser(name, email string) (*User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(strings.ToLower(email))

	if name == "" || email == "" {
		return nil, ErrInvalidInput
	}

	return &User{
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// Update updates user fields
func (u *User) Update(name, email string) error {
	if name != "" {
		u.Name = strings.TrimSpace(name)
		u.UpdatedAt = time.Now()
	}
	if email != "" {
		u.Email = strings.TrimSpace(strings.ToLower(email))
		u.UpdatedAt = time.Now()
	}
	return nil
}
