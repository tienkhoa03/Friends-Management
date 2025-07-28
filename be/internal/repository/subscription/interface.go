package repository

import (
	"BE_Friends_Management/internal/domain/entity"

	"gorm.io/gorm"
)

//go:generate mockgen -source=interface.go -destination=../mock/mock_subscription_repository.go

type SubscriptionRepository interface {
	GetDB() *gorm.DB
	CreateSubscription(requestorId, targetId int64) error
	DeleteSubscription(tx *gorm.DB, requestorId, targetId int64) error
	GetSubscription(requestorId, targetId int64) (*entity.Subscription, error)
}
