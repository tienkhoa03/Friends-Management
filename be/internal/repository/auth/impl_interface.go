package repository

import (
	"BE_Friends_Management/internal/domain/entity"

	"gorm.io/gorm"
)

type PostgreSQLAuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &PostgreSQLAuthRepository{db: db}
}

func (r *PostgreSQLAuthRepository) CreateToken(token *entity.UserToken) error {
	result := r.db.Create(token)
	return result.Error
}

func (r *PostgreSQLAuthRepository) FindByRefreshToken(refreshToken string) (*entity.UserToken, error) {
	var userToken = entity.UserToken{}
	err := r.db.Model(&entity.UserToken{}).Where("refresh_token = ?", refreshToken).First(&userToken).Error
	if err != nil {
		return nil, err
	}
	return &userToken, nil
}
