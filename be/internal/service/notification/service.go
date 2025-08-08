package service

import (
	"BE_Friends_Management/internal/domain/entity"
	notificationRepository "BE_Friends_Management/internal/repository/block_relationship"
	friendshipRepository "BE_Friends_Management/internal/repository/friendship"
	subscriptionRepository "BE_Friends_Management/internal/repository/subscription"
	userRepository "BE_Friends_Management/internal/repository/users"
	utils "BE_Friends_Management/pkg/utils"
	"errors"

	"gorm.io/gorm"
)

//go:generate mockgen -source=service.go -destination=../mock/mock_notification_service.go

type NotificationService interface {
	GetUpdateRecipients(authUserId int64, authUserRole string, senderEmail, text string) ([]*entity.User, error)
}

type notificationService struct {
	blockRepo        notificationRepository.BlockRelationshipRepository
	userRepo         userRepository.UserRepository
	friendshipRepo   friendshipRepository.FriendshipRepository
	subscriptionRepo subscriptionRepository.SubscriptionRepository
}

func NewNotificationService(blockRepo notificationRepository.BlockRelationshipRepository, userRepo userRepository.UserRepository, friendshipRepo friendshipRepository.FriendshipRepository, subscriptionRepo subscriptionRepository.SubscriptionRepository) NotificationService {
	return &notificationService{
		blockRepo:        blockRepo,
		userRepo:         userRepo,
		friendshipRepo:   friendshipRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

func (service *notificationService) GetUpdateRecipients(authUserId int64, authUserRole string, senderEmail, text string) ([]*entity.User, error) {
	sender, err := service.userRepo.GetUserByEmail(senderEmail)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	senderId := sender.Id
	if authUserRole == "user" && authUserId != senderId {
		return nil, ErrNotPermitted
	}

	blockRequestorIds, err := service.blockRepo.GetBlockRequestorIds(senderId)
	if err != nil {
		return nil, err
	}

	friendIds, err := service.friendshipRepo.RetrieveFriendIds(senderId)
	if err != nil {
		return nil, err
	}

	subscriberIds, err := service.subscriptionRepo.GetAllSubscriberIds(senderId)
	if err != nil {
		return nil, err
	}

	mentionedEmails := utils.ExtractEmails(text)
	mentionedUsers, err := service.userRepo.GetUsersFromEmails(mentionedEmails)
	if err != nil {
		return nil, err
	}
	var mentionedIds []int64
	for _, mentionedUser := range mentionedUsers {
		mentionedIds = append(mentionedIds, mentionedUser.Id)
	}

	var recipientIdsSet = make(map[int64]bool)
	for _, friendId := range friendIds {
		recipientIdsSet[friendId] = true
	}
	for _, subscriberId := range subscriberIds {
		recipientIdsSet[subscriberId] = true
	}
	for _, mentionedId := range mentionedIds {
		recipientIdsSet[mentionedId] = true
	}
	for _, blockRequestorId := range blockRequestorIds {
		recipientIdsSet[blockRequestorId] = false
	}

	recipientIds := make([]int64, 0, len(recipientIdsSet))
	for recipientId := range recipientIdsSet {
		if recipientIdsSet[recipientId] {
			recipientIds = append(recipientIds, recipientId)
		}
	}
	recipients, err := service.userRepo.GetUsersFromIds(recipientIds)
	if err != nil {
		return nil, err
	}
	return recipients, nil
}
