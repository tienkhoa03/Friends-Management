package repository

import (
	"BE_Friends_Management/internal/domain/entity"
)

//go:generate mockgen -source=interface.go -destination=../mock/mock_auth_repository.go

type AuthRepository interface {
	CreateToken(token *entity.UserToken) error
	FindByRefreshToken(refreshToken string) (*entity.UserToken, error)
}
