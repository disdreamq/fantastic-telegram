package domain

import "errors"

var (
	// User
	ErrInvalidUserName = errors.New("Invalid user name")
	ErrInvalidEmail    = errors.New("Invalid email")
)
