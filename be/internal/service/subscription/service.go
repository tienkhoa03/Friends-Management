package service

import (
	block_relationship "BE_Friends_Management/internal/repository/block_relationship"
	subscription "BE_Friends_Management/internal/repository/subscription"
	user "BE_Friends_Management/internal/repository/users"
	"errors"
	"strings"

	"gorm.io/gorm"
)

//go:generate mockgen -source=service.go -destination=../mock/mock_subscription_service.go

type SubscriptionService interface {
	CreateSubscription(requestorEmail, targetEmail string) error
}

type subscriptionService struct {
	repo                  subscription.SubscriptionRepository
	userRepo              user.UserRepository
	blockRelationshipRepo block_relationship.BlockRelationshipRepository
}

func NewSubscriptionService(repo subscription.SubscriptionRepository, userRepo user.UserRepository, blockRelationshipRepo block_relationship.BlockRelationshipRepository) SubscriptionService {
	return &subscriptionService{
		repo:                  repo,
		userRepo:              userRepo,
		blockRelationshipRepo: blockRelationshipRepo,
	}
}

func (service *subscriptionService) CreateSubscription(requestorEmail, targetEmail string) error {
	user1, err := service.userRepo.GetUserByEmail(requestorEmail)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}
	user2, err := service.userRepo.GetUserByEmail(targetEmail)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}
	if user1.Id == user2.Id {
		return ErrInvalidRequest
	}
	_, err = service.blockRelationshipRepo.GetBlockRelationship(user1.Id, user2.Id)
	if err == nil {
		return ErrIsBlocked
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	err = service.repo.CreateSubscription(user1.Id, user2.Id)
	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		return ErrAlreadySubscribed
	}
	return err
}
