package valueobjects

import (
	"errors"
	"strings"
)

// UserID represents a unique user identifier
type UserID int

// UserName represents a validated user name
type UserName string

// Email represents a validated email address
type Email string

// Value object errors
var (
	ErrInvalidUserName    = errors.New("invalid user name")
	ErrUserNameTooShort   = errors.New("user name too short (minimum 2 characters)")
	ErrUserNameTooLong    = errors.New("user name too long (maximum 100 characters)")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidEmailFormat = errors.New("invalid email format")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

// NewUserName creates a new validated UserName
func NewUserName(name string) (UserName, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return "", ErrInvalidUserName
	}
	if len(name) < 2 {
		return "", ErrUserNameTooShort
	}
	if len(name) > 100 {
		return "", ErrUserNameTooLong
	}

	return UserName(name), nil
}

// NewEmail creates a new validated Email
func NewEmail(email string) (Email, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	if email == "" {
		return "", ErrInvalidEmail
	}
	if !isValidEmail(email) {
		return "", ErrInvalidEmailFormat
	}

	return Email(email), nil
}

// String returns string representation of UserName
func (un UserName) String() string {
	return string(un)
}

// String returns string representation of Email
func (e Email) String() string {
	return string(e)
}

// Equals compares two UserNames
func (un UserName) Equals(other UserName) bool {
	return un == other
}

// Equals compares two Emails
func (e Email) Equals(other Email) bool {
	return e == other
}

// IsEmpty checks if UserName is empty
func (un UserName) IsEmpty() bool {
	return strings.TrimSpace(string(un)) == ""
}

// IsEmpty checks if Email is empty
func (e Email) IsEmpty() bool {
	return strings.TrimSpace(string(e)) == ""
}

// GetDomain returns email domain part
func (e Email) GetDomain() string {
	parts := strings.Split(string(e), "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// isValidEmail performs basic email validation
func isValidEmail(email string) bool {
	if len(email) == 0 {
		return false
	}
	return strings.Contains(email, "@") &&
		strings.Contains(email, ".") &&
		len(strings.Split(email, "@")) == 2
}
