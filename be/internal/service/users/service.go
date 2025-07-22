package service

import (
	"BE_Friends_Management/internal/domain/entity"
	user "BE_Friends_Management/internal/repository/users"
)

type UserService struct {
	repo user.UserRepository
}

func NewUserService(repo user.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (service *UserService) GetAllUser() []*entity.User {
	users := service.repo.GetAllUser()
	return users
}

func (service *UserService) GetUserById(userId int64) (*entity.User, error) {
	user, err := service.repo.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *UserService) CreateUser(email string) (*entity.User, error) {
	user := entity.User{Email: email}
	newUser, err := service.repo.CreateUser(&user)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func (service *UserService) DeleteUserById(userId int64) error {
	err := service.repo.DeleteUserById(userId)
	return err
}

func (service *UserService) UpdateUser(userId int64, email string) (*entity.User, error) {
	user := entity.User{Id: userId, Email: email}
	updatedUser, err := service.repo.UpdateUser(&user)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}
