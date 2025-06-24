package repositories

import "errors"

// Repository errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)
