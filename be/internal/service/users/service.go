package service

import (
	"BE_Friends_Management/internal/domain/entity"
	userRepository "BE_Friends_Management/internal/repository/users"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=service.go -destination=../mock/mock_user_service.go

type UserService interface {
	GetAllUser() ([]*entity.User, error)
	GetUserById(userId int64) (*entity.User, error)
	DeleteUserById(userId int64) error
	UpdateUser(userId int64, email string, password string) (*entity.User, error)
}

type userService struct {
	repo userRepository.UserRepository
}

func NewUserService(repo userRepository.UserRepository) UserService {
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

func (service *userService) DeleteUserById(userId int64) error {
	err := service.repo.DeleteUserById(userId)
	return err
}

func (service *userService) UpdateUser(userId int64, email string, password string) (*entity.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := entity.User{Id: userId, Email: email, Password: string(hashedPassword)}
	updatedUser, err := service.repo.UpdateUser(&user)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}
