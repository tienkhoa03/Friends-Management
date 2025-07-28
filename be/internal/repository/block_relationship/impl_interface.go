package repository

import (
	"BE_Friends_Management/internal/domain/entity"

	"gorm.io/gorm"
)

type PostgreSQLBlockRelationshipRepository struct {
	db *gorm.DB
}

func NewBlockRelationshipRepository(db *gorm.DB) BlockRelationshipRepository {
	return &PostgreSQLBlockRelationshipRepository{db: db}
}

func (r *PostgreSQLBlockRelationshipRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *PostgreSQLBlockRelationshipRepository) CreateBlockRelationship(tx *gorm.DB, requestorId, targetId int64) error {
	newBlockRelationship := entity.BlockRelationship{RequestorId: requestorId, TargetId: targetId}
	err := tx.Model(entity.BlockRelationship{}).Create(&newBlockRelationship).Error
	return err
}

func (r *PostgreSQLBlockRelationshipRepository) GetBlockRelationship(requestorId, targetId int64) (*entity.BlockRelationship, error) {
	blockRelationship := entity.BlockRelationship{RequestorId: requestorId, TargetId: targetId}
	err := r.db.Model(entity.BlockRelationship{}).First(&blockRelationship).Error
	if err != nil {
		return nil, err
	}
	return &blockRelationship, nil
}
