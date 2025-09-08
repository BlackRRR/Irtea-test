package domain

import "errors"

var (
	ErrUserTooYoung       = errors.New("user must be at least 18 years old")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidPassword    = errors.New("password must be at least 8 characters long")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
