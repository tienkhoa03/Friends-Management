package service

import "errors"

var (
	ErrUserNotFound      = errors.New("User not found")
	ErrAlreadySubscribed = errors.New("Requestor has already subscribed this target")
	ErrInvalidRequest    = errors.New("Two email can not be the same")
)
