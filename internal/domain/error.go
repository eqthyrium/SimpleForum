package domain

import "errors"

var (
	// General
	ErrUserNotFound      = errors.New("user is not found")
	ErrInvalidToken      = errors.New("invalid session token")
	ErrInvalidCredential = errors.New("invalid credential")
	ErrInvalidCSRFToken  = errors.New("invalid csrf token")
	ErrNotValidContent   = errors.New("content is not valid")

	// User
	ErrInvalidUser = errors.New("invalid user")

	//Post
	ErrPostNotFound = errors.New("post not found")

	//Category
	ErrNoCategories = errors.New("there is no categories")

	//Commentary
	ErrCommentaryNotFound = errors.New("commentary not found")
	//Request to be moderator
	ErrRepeatedRequest = errors.New("repeated request")
	// Moderation
	ErrInvalidOperation = errors.New("invalid operation")
)
