package domain

import "errors"

var (
	ErrInvalidTitle   = errors.New("Invalid title")
	ErrInvalidContent = errors.New("Invalid content")
	ErrInvalidUserId  = errors.New("Invalid user ID")
	ErrInvalidID      = errors.New("Invalid ID")

	// grpc
	ErrInvalidToken = errors.New("Invalid token")
)
