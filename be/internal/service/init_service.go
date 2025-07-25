package service

import (
	repository "BE_Friends_Management/internal/repository"
	friendship "BE_Friends_Management/internal/service/friendship"
	subscription "BE_Friends_Management/internal/service/subscription"
	user "BE_Friends_Management/internal/service/users"
)

type Service struct {
	User         user.UserService
	Friendship   friendship.FriendshipService
	Subscription subscription.SubscriptionService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User:         user.NewUserService(repos.User),
		Friendship:   friendship.NewFriendshipService(repos.Friendship, repos.User),
		Subscription: subscription.NewSubscriptionService(repos.Subscription, repos.User),
	}
}
