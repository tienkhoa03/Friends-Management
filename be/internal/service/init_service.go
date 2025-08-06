package service

import (
	repository "BE_Friends_Management/internal/repository"
	auth "BE_Friends_Management/internal/service/auth"
	block_relationship "BE_Friends_Management/internal/service/block_relationship"
	friendship "BE_Friends_Management/internal/service/friendship"
	notification "BE_Friends_Management/internal/service/notification"
	subscription "BE_Friends_Management/internal/service/subscription"
	user "BE_Friends_Management/internal/service/users"
)

type Service struct {
	User              user.UserService
	Friendship        friendship.FriendshipService
	Subscription      subscription.SubscriptionService
	BlockRelationship block_relationship.BlockRelationshipService
	Notification      notification.NotificationService
	Auth              auth.AuthService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User:              user.NewUserService(repos.User),
		Friendship:        friendship.NewFriendshipService(repos.Friendship, repos.User, repos.BlockRelationship),
		Subscription:      subscription.NewSubscriptionService(repos.Subscription, repos.User, repos.Friendship, repos.BlockRelationship),
		BlockRelationship: block_relationship.NewBlockRelationshipService(repos.BlockRelationship, repos.User, repos.Friendship, repos.Subscription),
		Notification:      notification.NewNotificationService(repos.BlockRelationship, repos.User, repos.Friendship, repos.Subscription),
		Auth:              auth.NewAuthService(repos.Auth, repos.User),
	}
}
