package service

import (
	block_relationship "BE_Friends_Management/internal/repository/block_relationship"
	friendship "BE_Friends_Management/internal/repository/friendship"
	subscription "BE_Friends_Management/internal/repository/subscription"
	user "BE_Friends_Management/internal/repository/users"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type subscriptionService struct {
	repo                  subscription.SubscriptionRepository
	userRepo              user.UserRepository
	friendshipRepo        friendship.FriendshipRepository
	blockRelationshipRepo block_relationship.BlockRelationshipRepository
}

func NewSubscriptionService(repo subscription.SubscriptionRepository, userRepo user.UserRepository, friendshipRepo friendship.FriendshipRepository, blockRelationshipRepo block_relationship.BlockRelationshipRepository) SubscriptionService {
	return &subscriptionService{
		repo:                  repo,
		userRepo:              userRepo,
		friendshipRepo:        friendshipRepo,
		blockRelationshipRepo: blockRelationshipRepo,
	}
}

func (service *subscriptionService) CreateSubscription(authUserId int64, requestorEmail, targetEmail string) error {
	requestor, err := service.userRepo.GetUserByEmail(requestorEmail)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}
	if authUserId != requestor.Id {
		return ErrNotPermitted
	}
	target, err := service.userRepo.GetUserByEmail(targetEmail)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}
	if requestor.Id == target.Id {
		return ErrInvalidRequest
	}
	_, err = service.blockRelationshipRepo.GetBlockRelationship(requestor.Id, target.Id)
	if err == nil {
		userId1 := requestor.Id
		userId2 := target.Id
		if userId1 > userId2 {
			userId1, userId2 = userId2, userId1
		}
		_, err := service.friendshipRepo.GetFriendship(userId1, userId2)
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrIsBlocked
		}
		if err != nil {
			return err
		}
		err = service.blockRelationshipRepo.DeleteBlockRelationship(requestor.Id, target.Id)
		if err != nil {
			return err
		}
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	err = service.repo.CreateSubscription(requestor.Id, target.Id)
	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		return ErrAlreadySubscribed
	}
	return err
}
