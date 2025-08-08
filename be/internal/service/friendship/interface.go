package service

import (
	"BE_Friends_Management/internal/domain/entity"
	"errors"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrAlreadyFriend  = errors.New("users are already friends")
	ErrInvalidRequest = errors.New("two email can not be the same")
	ErrIsBlocked      = errors.New("one user has blocked another")
	ErrNotPermitted   = errors.New("action not permitted")
)

//go:generate mockgen -source=service.go -destination=../mock/mock_friendship_service.go

type FriendshipService interface {
	CreateFriendship(authUserId int64, email1, email2 string) error
	RetrieveFriendsList(authUserId int64, email string) ([]*entity.User, error)
	RetrieveCommonFriends(authUserId int64, email1, email2 string) ([]*entity.User, error)
	CountFriends(friendsList []*entity.User) int64
}
