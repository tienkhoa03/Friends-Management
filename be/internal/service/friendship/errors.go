package service

import "errors"

var (
	ErrUserNotFound   = errors.New("User not found")
	ErrAlreadyFriend  = errors.New("Users are already friends")
	ErrInvalidRequest = errors.New("Two email can not be the same")
	ErrIsBlocked      = errors.New("One user has blocked another")
)
