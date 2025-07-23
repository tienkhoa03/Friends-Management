package handler

import (
	service "BE_Friends_Management/internal/service"
)

type Handlers struct {
	User       *UserHandler
	Friendship *FriendshipHandler
}

func NewHandlers(services *service.Service) *Handlers {
	return &Handlers{
		User:       NewUserHandler(services.User),
		Friendship: NewFriendshipHandler(services.Friendship),
	}
}
