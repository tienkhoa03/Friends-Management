package service

import (
	"BE_Friends_Management/internal/domain/entity"
	user "BE_Friends_Management/internal/repository/users"
)

//go:generate mockgen -source=service.go -destination=../mock/mock_user_service.go

type UserService interface {
	GetAllUser() ([]*entity.User, error)
	GetUserById(userId int64) (*entity.User, error)
	CreateUser(email string) (*entity.User, error)
	DeleteUserById(userId int64) error
	UpdateUser(userId int64, email string) (*entity.User, error)
}

type userService struct {
	repo user.UserRepository
}

func NewUserService(repo user.UserRepository) UserService {
	return &userService{repo: repo}
}

func (service *userService) GetAllUser() ([]*entity.User, error) {
	users, err := service.repo.GetAllUser()
	return users, err
}

func (service *userService) GetUserById(userId int64) (*entity.User, error) {
	user, err := service.repo.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *userService) CreateUser(email string) (*entity.User, error) {
	user := entity.User{Email: email}
	newUser, err := service.repo.CreateUser(&user)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func (service *userService) DeleteUserById(userId int64) error {
	err := service.repo.DeleteUserById(userId)
	return err
}

func (service *userService) UpdateUser(userId int64, email string) (*entity.User, error) {
	user := entity.User{Id: userId, Email: email}
	updatedUser, err := service.repo.UpdateUser(&user)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}
