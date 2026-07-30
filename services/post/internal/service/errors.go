package service

import "errors"

var (
	// Post
	ErrPostNotFound       = errors.New("Post not found.")
	ErrLinkedUserNotFound = errors.New("Linked user not found.")
	ErrUpdatePostFailed   = errors.New("Failed to update post.")
	ErrDeletePostFailed   = errors.New("Failed to delete post.")

	// Cache
	ErrCacheUnmarshal = errors.New("Can not unmarshal data from cache.")

	// Unexpected
	ErrUnexpected = errors.New("Unexpected error.")
)
