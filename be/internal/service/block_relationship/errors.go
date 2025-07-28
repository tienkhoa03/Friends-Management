package service

import "errors"

var (
	ErrUserNotFound   = errors.New("User not found")
	ErrAlreadyBlocked = errors.New("Requestor has already blocked this target user")
	ErrInvalidRequest = errors.New("Two email can not be the same")
	ErrNotSubscribed  = errors.New("Can not block if they are friends and have not subscribed")
)
