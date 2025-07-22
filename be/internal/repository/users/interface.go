package repository

import "BE_Friends_Management/internal/domain/entity"

type UserRepository interface {
	CreateUser(user *entity.User) (*entity.User, error)
	GetAllUser() []*entity.User
	GetUserById(userId int64) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	UpdateUser(user *entity.User) (*entity.User, error)
	DeleteUserById(userId int64) error
}
