package service

import (
	repository "BE_Friends_Management/internal/repository"
	friendship "BE_Friends_Management/internal/service/friendship"
	user "BE_Friends_Management/internal/service/users"
)

type Service struct {
	User       *user.UserService
	Friendship *friendship.FriendshipService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User:       user.NewUserService(repos.User),
		Friendship: friendship.NewFriendshipService(repos.Friendship, repos.User),
	}
}
