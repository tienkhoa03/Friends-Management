package repository

import (
	block_relationship "BE_Friends_Management/internal/repository/block_relationship"
	friendship "BE_Friends_Management/internal/repository/friendship"
	subscription "BE_Friends_Management/internal/repository/subscription"
	user "BE_Friends_Management/internal/repository/users"

	"gorm.io/gorm"
)

type Repository struct {
	User              user.UserRepository
	Friendship        friendship.FriendshipRepository
	Subscription      subscription.SubscriptionRepository
	BlockRelationship block_relationship.BlockRelationshipRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:              user.NewUserRepository(db),
		Friendship:        friendship.NewFriendshipRepository(db),
		Subscription:      subscription.NewSubscriptionRepository(db),
		BlockRelationship: block_relationship.NewBlockRelationshipRepository(db),
	}
}
