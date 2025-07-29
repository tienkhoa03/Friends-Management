package service

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrAlreadyFriend  = errors.New("users are already friends")
	ErrInvalidRequest = errors.New("two email can not be the same")
	ErrIsBlocked      = errors.New("one user has blocked another")
)
