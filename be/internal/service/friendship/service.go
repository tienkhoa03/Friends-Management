package service

import (
	friendship "BE_Friends_Management/internal/repository/friendship"
	user "BE_Friends_Management/internal/repository/users"
	"errors"
	"strings"

	"gorm.io/gorm"
)

//go:generate mockgen -source=service.go -destination=../mock/mock_friendship_service.go

type FriendshipService interface {
	CreateFriendship(email1, email2 string) error
}

type friendshipService struct {
	repo     friendship.FriendshipRepository
	userRepo user.UserRepository
}

func NewFriendshipService(repo friendship.FriendshipRepository, userRepo user.UserRepository) FriendshipService {
	return &friendshipService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (service *friendshipService) CreateFriendship(email1, email2 string) error {
	user1, err := service.userRepo.GetUserByEmail(email1)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}
	user2, err := service.userRepo.GetUserByEmail(email2)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}
	if user1.Id == user2.Id {
		return ErrInvalidRequest
	}
	if user1.Id < user2.Id {
		err = service.repo.CreateFriendship(user1.Id, user2.Id)
	} else {
		err = service.repo.CreateFriendship(user2.Id, user1.Id)
	}
	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		return ErrAlreadyFriend
	}
	return err
}
