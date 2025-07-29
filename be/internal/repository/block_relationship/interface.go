package repository

import (
	"BE_Friends_Management/internal/domain/entity"

	"gorm.io/gorm"
)

//go:generate mockgen -source=interface.go -destination=../mock/mock_block_repository.go

type BlockRelationshipRepository interface {
	GetDB() *gorm.DB
	CreateBlockRelationship(tx *gorm.DB, requestorId, targetId int64) error
	GetBlockRelationship(requestorId, targetId int64) (*entity.BlockRelationship, error)
	GetBlockRequestorIds(targetId int64) ([]int64, error)
	DeleteBlockRelationship(requestorId, targetId int64) error
}
