package repository

import (
	"BE_Friends_Management/internal/domain/entity"

	"gorm.io/gorm"
)

//go:generate mockgen -source=interface.go -destination=../mock/mock_friendship_repository.go

type FriendshipRepository interface {
	GetDB() *gorm.DB
	CreateFriendship(userId1, userId2 int64) error
	RetrieveFriendIds(userId int64) ([]int64, error)
	GetFriendship(userId1, userId2 int64) (*entity.Friendship, error)
}
