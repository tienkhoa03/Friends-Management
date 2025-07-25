package repository

import (
	"BE_Friends_Management/internal/domain/entity"

	"gorm.io/gorm"
)

type PostgreSQLSubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &PostgreSQLSubscriptionRepository{db: db}
}

func (r *PostgreSQLSubscriptionRepository) CreateSubscription(requestorId, targetId int64) error {
	newSubscription := entity.Subscription{RequestorId: requestorId, TargetId: targetId}
	err := r.db.Model(entity.Subscription{}).Create(&newSubscription).Error
	return err
}
