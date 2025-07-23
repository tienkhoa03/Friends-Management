package repository

import (
	friendship "BE_Friends_Management/internal/repository/friendship"
	user "BE_Friends_Management/internal/repository/users"

	"gorm.io/gorm"
)

type Repository struct {
	User       user.UserRepository
	Friendship friendship.FriendshipRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:       user.NewUserRepository(db),
		Friendship: friendship.NewFriendshipRepository(db),
	}
}
