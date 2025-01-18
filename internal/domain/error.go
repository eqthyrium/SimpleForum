package domain

import "errors"

var (
	// General
	ErrUserNotFound      = errors.New("user is not found")
	ErrInvalidToken      = errors.New("invalid session token")
	ErrInvalidCredential = errors.New("invalid credential")
	ErrInvalidCSRFToken  = errors.New("invalid csrf token")

	// User
	ErrInvalidUser = errors.New("invalid user")

	//Post
	ErrPostNotFound = errors.New("post not found")

	//Category
	ErrCategoryNotFound = errors.New("category not found")
)
