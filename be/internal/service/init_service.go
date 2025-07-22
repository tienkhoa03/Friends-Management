package service

import (
	repository "BE_Friends_Management/internal/repository"
	user "BE_Friends_Management/internal/service/users"
)

type Service struct {
	User *user.UserService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{User: user.NewUserService(repos.User)}
}
