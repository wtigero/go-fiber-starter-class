package services

import "errors"

// Service layer errors
var (
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrInvalidUserName    = errors.New("invalid user name")
	ErrUserNameTooShort   = errors.New("user name too short (minimum 2 characters)")
	ErrUserNameTooLong    = errors.New("user name too long (maximum 100 characters)")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidEmailFormat = errors.New("invalid email format")
	ErrNoFieldsToUpdate   = errors.New("no fields to update")
)
