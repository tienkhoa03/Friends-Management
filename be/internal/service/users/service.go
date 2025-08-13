package service

import (
	"BE_Friends_Management/internal/domain/entity"
	userRepository "BE_Friends_Management/internal/repository/users"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

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
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *userService) DeleteUserById(userId int64) error {
	err := service.repo.DeleteUserById(userId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrUserNotFound
	}
	return err
}

func (service *userService) UpdateUser(userId int64, email string, password string) (*entity.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := entity.User{Id: userId, Email: email, Password: string(hashedPassword)}
	updatedUser, err := service.repo.UpdateUser(&user)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}
