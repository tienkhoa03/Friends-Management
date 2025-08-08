package service

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrNotPermitted = errors.New("action not permitted")
)
