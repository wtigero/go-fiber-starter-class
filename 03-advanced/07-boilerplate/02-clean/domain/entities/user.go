package entities

import (
	"errors"
	"strings"
	"time"
)

// User represents the core business entity
type User struct {
	id        int
	name      string
	email     string
	createdAt time.Time
	updatedAt time.Time
}

// User creation and validation errors
var (
	ErrInvalidUserName    = errors.New("invalid user name")
	ErrUserNameTooShort   = errors.New("user name too short (minimum 2 characters)")
	ErrUserNameTooLong    = errors.New("user name too long (maximum 100 characters)")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidEmailFormat = errors.New("invalid email format")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

// NewUser creates a new user with validation
func NewUser(name, email string) (*User, error) {
	user := &User{
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}

	if err := user.SetName(name); err != nil {
		return nil, err
	}

	if err := user.SetEmail(email); err != nil {
		return nil, err
	}

	return user, nil
}

// NewUserWithID creates user with existing ID (for loading from database)
func NewUserWithID(id int, name, email string, createdAt, updatedAt time.Time) (*User, error) {
	user, err := NewUser(name, email)
	if err != nil {
		return nil, err
	}

	user.id = id
	user.createdAt = createdAt
	user.updatedAt = updatedAt

	return user, nil
}

// Getters
func (u *User) ID() int              { return u.id }
func (u *User) Name() string         { return u.name }
func (u *User) Email() string        { return u.email }
func (u *User) CreatedAt() time.Time { return u.createdAt }
func (u *User) UpdatedAt() time.Time { return u.updatedAt }

// SetID sets user ID (used by repository)
func (u *User) SetID(id int) {
	u.id = id
}

// SetName sets and validates user name
func (u *User) SetName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return ErrInvalidUserName
	}
	if len(name) < 2 {
		return ErrUserNameTooShort
	}
	if len(name) > 100 {
		return ErrUserNameTooLong
	}

	u.name = name
	u.updatedAt = time.Now()
	return nil
}

// SetEmail sets and validates user email
func (u *User) SetEmail(email string) error {
	email = strings.ToLower(strings.TrimSpace(email))

	if email == "" {
		return ErrInvalidEmail
	}
	if !isValidEmail(email) {
		return ErrInvalidEmailFormat
	}

	u.email = email
	u.updatedAt = time.Now()
	return nil
}

// Update updates user with new values
func (u *User) Update(name, email string) error {
	if name != "" {
		if err := u.SetName(name); err != nil {
			return err
		}
	}

	if email != "" {
		if err := u.SetEmail(email); err != nil {
			return err
		}
	}

	return nil
}

// isValidEmail performs basic email validation
func isValidEmail(email string) bool {
	if len(email) == 0 {
		return false
	}
	// Basic email validation
	return strings.Contains(email, "@") &&
		strings.Contains(email, ".") &&
		len(strings.Split(email, "@")) == 2
}
