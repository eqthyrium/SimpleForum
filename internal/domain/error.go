package domain

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidUser       = errors.New("invalid user")
	ErrInvalidToken      = errors.New("invalid token")
	ErrInvalidCredential = errors.New("invalid credential")
)
