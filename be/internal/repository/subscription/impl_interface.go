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

func (r *PostgreSQLSubscriptionRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *PostgreSQLSubscriptionRepository) CreateSubscription(requestorId, targetId int64) error {
	newSubscription := entity.Subscription{RequestorId: requestorId, TargetId: targetId}
	err := r.db.Model(&entity.Subscription{}).Create(&newSubscription).Error
	return err
}

func (r *PostgreSQLSubscriptionRepository) DeleteSubscription(tx *gorm.DB, requestorId, targetId int64) error {
	deleteSubscription := entity.Subscription{RequestorId: requestorId, TargetId: targetId}
	err := tx.Model(&entity.Subscription{}).Delete(&deleteSubscription).Error
	return err
}

func (r *PostgreSQLSubscriptionRepository) GetSubscription(requestorId, targetId int64) (*entity.Subscription, error) {
	subscription := entity.Subscription{RequestorId: requestorId, TargetId: targetId}
	err := r.db.Model(&entity.Subscription{}).First(&subscription).Error
	if err != nil {
		return nil, err
	}
	return &subscription, err
}

func (r *PostgreSQLSubscriptionRepository) GetAllSubscriberIds(targetId int64) ([]int64, error) {
	var subscriberIds []int64
	err := r.db.Model(&entity.Subscription{}).Where("target_id = ?", targetId).Pluck("requestor_id", &subscriberIds).Error
	if err != nil {
		return nil, err
	}
	return subscriberIds, nil
}
