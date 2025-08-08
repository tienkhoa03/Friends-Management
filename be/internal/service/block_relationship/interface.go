package service

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrAlreadyBlocked = errors.New("requestor has already blocked this target user")
	ErrInvalidRequest = errors.New("two email can not be the same")
	ErrNotSubscribed  = errors.New("can not block if they are friends and have not subscribed")
	ErrNotPermitted   = errors.New("action not permitted")
)

//go:generate mockgen -source=interface.go -destination=../mock/mock_block_service.go

type BlockRelationshipService interface {
	CreateBlockRelationship(authUserId int64, requestorEmail, targetEmail string) error
}
