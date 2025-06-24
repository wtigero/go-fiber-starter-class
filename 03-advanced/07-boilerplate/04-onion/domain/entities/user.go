package entities

import (
	"errors"
	valueobjects "onion-arch/domain/value-objects"
	"strings"
	"time"
)

// User represents the core business entity in domain layer
type User struct {
	id        valueobjects.UserID
	name      valueobjects.UserName
	email     valueobjects.Email
	createdAt time.Time
	updatedAt time.Time
}

// Domain errors
var (
	ErrUserNotFound = errors.New("user not found")
)

// NewUser creates a new user entity with validation
func NewUser(name string, email string) (*User, error) {
	userName, err := valueobjects.NewUserName(name)
	if err != nil {
		return nil, err
	}

	userEmail, err := valueobjects.NewEmail(email)
	if err != nil {
		return nil, err
	}

	return &User{
		id:        valueobjects.UserID(0), // Will be set by repository
		name:      userName,
		email:     userEmail,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, nil
}

// NewUserWithID creates user with existing ID (for loading from persistence)
func NewUserWithID(id int, name string, email string, createdAt, updatedAt time.Time) (*User, error) {
	user, err := NewUser(name, email)
	if err != nil {
		return nil, err
	}

	user.id = valueobjects.UserID(id)
	user.createdAt = createdAt
	user.updatedAt = updatedAt

	return user, nil
}

// Getters (immutable access)
func (u *User) ID() int              { return int(u.id) }
func (u *User) Name() string         { return string(u.name) }
func (u *User) Email() string        { return string(u.email) }
func (u *User) CreatedAt() time.Time { return u.createdAt }
func (u *User) UpdatedAt() time.Time { return u.updatedAt }

// SetID sets user ID (used by repository after persistence)
func (u *User) SetID(id int) {
	u.id = valueobjects.UserID(id)
}

// ChangeName changes user name with validation
func (u *User) ChangeName(newName string) error {
	userName, err := valueobjects.NewUserName(newName)
	if err != nil {
		return err
	}

	u.name = userName
	u.updatedAt = time.Now()
	return nil
}

// ChangeEmail changes user email with validation
func (u *User) ChangeEmail(newEmail string) error {
	userEmail, err := valueobjects.NewEmail(newEmail)
	if err != nil {
		return err
	}

	u.email = userEmail
	u.updatedAt = time.Now()
	return nil
}

// Update updates user with new values
func (u *User) Update(name, email string) error {
	if name != "" {
		if err := u.ChangeName(name); err != nil {
			return err
		}
	}

	if email != "" {
		if err := u.ChangeEmail(email); err != nil {
			return err
		}
	}

	return nil
}

// Domain business methods

// IsActive checks if user is active (business rule example)
func (u *User) IsActive() bool {
	// Business rule: User is active if created within last 2 years
	return time.Since(u.createdAt).Hours() < 24*365*2
}

// CanBeDeleted checks if user can be deleted (business rule example)
func (u *User) CanBeDeleted() bool {
	// Business rule: User can be deleted if created more than 30 days ago
	return time.Since(u.createdAt).Hours() > 24*30
}

// GetDisplayName returns formatted display name
func (u *User) GetDisplayName() string {
	name := strings.TrimSpace(string(u.name))
	if name == "" {
		return "Unknown User"
	}
	return name
}
