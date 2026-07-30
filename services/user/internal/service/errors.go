package service

import "errors"

var (
	// Registration
	ErrEmptyPassword           = errors.New("Empty password not allowed.")
	ErrInvalidPasswordLength   = errors.New("Invalid length of user password.")
	ErrUserAlreadyExists       = errors.New("User with this email already exists.")
	ErrCanNotCalculatePassHash = errors.New("Can not calculate password hash.")

	// Auth
	ErrWrongPassword = errors.New("Wrong password.")
	ErrCanNotLogin   = errors.New("Can not login.")

	// User
	ErrUserNotFound     = errors.New("User not found.")
	ErrUserTimeOut      = errors.New("Time out error while processing user.")
	ErrMethodNotAllowed = errors.New("Method not allowed.")

	// Unexpected
	ErrUnexpected = errors.New("Unexpected error.")
)
