package service

import "BE_Friends_Management/internal/domain/entity"

//go:generate mockgen -source=interface.go -destination=../mock/mock_user_service.go

type UserService interface {
	GetAllUser() ([]*entity.User, error)
	GetUserById(userId int64) (*entity.User, error)
	DeleteUserById(userId int64) error
	UpdateUser(userId int64, email string, password string) (*entity.User, error)
}
