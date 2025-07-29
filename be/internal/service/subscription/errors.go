package service

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrAlreadySubscribed = errors.New("requestor has already subscribed this target user")
	ErrInvalidRequest    = errors.New("two email can not be the same")
	ErrIsBlocked         = errors.New("requestor has blocked target user")
)
