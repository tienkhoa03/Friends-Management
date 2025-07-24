package repository

import (
	"BE_Friends_Management/internal/domain/entity"

	"gorm.io/gorm"
)

type PostgreSQLFriendshipRepository struct {
	db *gorm.DB
}

func NewFriendshipRepository(db *gorm.DB) FriendshipRepository {
	return &PostgreSQLFriendshipRepository{db: db}
}

func (r *PostgreSQLFriendshipRepository) CreateFriendship(userId1, userId2 int64) error {
	newFriendship := entity.Friendship{UserId1: userId1, UserId2: userId2}
	err := r.db.Model(entity.Friendship{}).Create(&newFriendship).Error
	return err
}

func (r *PostgreSQLFriendshipRepository) RetrieveFriendsList(userId int64) ([]int64, error) {
	var friends []int64
	err := r.db.Model(entity.Friendship{}).
		Where("user_id1 = ? OR user_id2 = ?", userId, userId).
		Select("CASE WHEN user_id1 = ? THEN user_id2 ELSE user_id1 END", userId).
		Scan(&friends).Error
	if err != nil {
		return nil, err
	}
	return friends, nil
}
