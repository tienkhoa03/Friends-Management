package service

import (
	"BE_Friends_Management/internal/domain/entity"
	friendship "BE_Friends_Management/internal/repository/friendship"
	user "BE_Friends_Management/internal/repository/users"
	"errors"
	"strings"

	"gorm.io/gorm"
)

//go:generate mockgen -source=service.go -destination=../mock/mock_friendship_service.go

type FriendshipService interface {
	CreateFriendship(email1, email2 string) error
	RetrieveFriendsList(email string) ([]*entity.User, error)
	CountFriends(friendsList []*entity.User) int64
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

func (service *friendshipService) RetrieveFriendsList(email string) ([]*entity.User, error) {
	user, err := service.userRepo.GetUserByEmail(email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	friendIds, err := service.repo.RetrieveFriendsList(user.Id)
	if err != nil {
		return nil, err
	}
	friends := make([]*entity.User, len(friendIds))
	for i, id := range friendIds {
		friend, err := service.userRepo.GetUserById(id)
		if err != nil {
			return nil, err
		}
		friends[i] = friend
	}
	return friends, nil
}

func (service *friendshipService) CountFriends(friends []*entity.User) int64 {
	var count int64 = 0
	for _, friend := range friends {
		if friend != nil {
			count++
		}
	}
	return count
}
