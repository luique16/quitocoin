package errorpkg

import "errors"

var (
	ErrEmailExists      = errors.New("email already exists")
	ErrUserNotFound     = errors.New("user not found")
	ErrNameRequired     = errors.New("name is required")
	ErrEmailRequired    = errors.New("email is required")
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrPasswordRequired = errors.New("password is required")
	ErrWeakPassword     = errors.New("password must be at least 8 characters with uppercase, lowercase, digit, and special character")
)
