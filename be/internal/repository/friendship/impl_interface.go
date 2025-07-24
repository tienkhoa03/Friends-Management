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
	result := r.db.Model(entity.Friendship{}).Create(&newFriendship)
	return result.Error
}
