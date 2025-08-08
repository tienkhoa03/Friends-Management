package service

import (
	blockRelationship "BE_Friends_Management/internal/repository/block_relationship"
	friendship "BE_Friends_Management/internal/repository/friendship"
	subscription "BE_Friends_Management/internal/repository/subscription"
	user "BE_Friends_Management/internal/repository/users"
	"errors"
	"strings"

	"gorm.io/gorm"
)

//go:generate mockgen -source=service.go -destination=../mock/mock_block_service.go

type BlockRelationshipService interface {
	CreateBlockRelationship(authUserId int64, requestorEmail, targetEmail string) error
}

type blockRelationshipService struct {
	repo             blockRelationship.BlockRelationshipRepository
	userRepo         user.UserRepository
	friendshipRepo   friendship.FriendshipRepository
	subscriptionRepo subscription.SubscriptionRepository
}

func NewBlockRelationshipService(repo blockRelationship.BlockRelationshipRepository, userRepo user.UserRepository, friendshipRepo friendship.FriendshipRepository, subscriptionRepo subscription.SubscriptionRepository) BlockRelationshipService {
	return &blockRelationshipService{
		repo:             repo,
		userRepo:         userRepo,
		friendshipRepo:   friendshipRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

func (service *blockRelationshipService) CreateBlockRelationship(authUserId int64, requestorEmail, targetEmail string) error {
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
	userId1 := requestor.Id
	userId2 := target.Id
	if userId1 > userId2 {
		userId1, userId2 = userId2, userId1
	}
	_, errFriendship := service.friendshipRepo.GetFriendship(userId1, userId2)
	if (errFriendship != nil) && !errors.Is(errFriendship, gorm.ErrRecordNotFound) {
		return errFriendship
	}
	_, errSubscription := service.subscriptionRepo.GetSubscription(requestor.Id, target.Id)
	if (errSubscription != nil) && !errors.Is(errSubscription, gorm.ErrRecordNotFound) {
		return errSubscription
	}
	db := service.repo.GetDB()
	err = db.Transaction(func(tx *gorm.DB) error {
		if errFriendship == nil && errSubscription != nil {
			return ErrNotSubscribed
		}
		if errSubscription == nil {
			err := service.subscriptionRepo.DeleteSubscription(tx, requestor.Id, target.Id)
			if err != nil {
				return err
			}
		}
		err := service.repo.CreateBlockRelationship(tx, requestor.Id, target.Id)
		if err != nil && strings.Contains(err.Error(), "duplicate key") {
			return ErrAlreadyBlocked
		}
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
