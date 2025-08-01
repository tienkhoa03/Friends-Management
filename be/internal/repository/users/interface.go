package repository

import (
	"BE_Friends_Management/internal/domain/entity"

	"gorm.io/gorm"
)

//go:generate mockgen -source=interface.go -destination=../mock/mock_user_repository.go

type UserRepository interface {
	GetDB() *gorm.DB
	CreateUser(user *entity.User) (*entity.User, error)
	GetAllUser() ([]*entity.User, error)
	GetUserById(userId int64) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetUsersFromIds(userIds []int64) ([]*entity.User, error)
	GetUsersFromEmails(emails []string) ([]*entity.User, error)
	UpdateUser(user *entity.User) (*entity.User, error)
	DeleteUserById(userId int64) error
}
