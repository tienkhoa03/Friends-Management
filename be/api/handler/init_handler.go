package handler

import (
	service "BE_Friends_Management/internal/service"
)

type Handlers struct {
	User              *UserHandler
	Friendship        *FriendshipHandler
	Subscription      *SubscriptionHandler
	BlockRelationship *BlockRelationshipHandler
}

func NewHandlers(services *service.Service) *Handlers {
	return &Handlers{
		User:              NewUserHandler(services.User),
		Friendship:        NewFriendshipHandler(services.Friendship),
		Subscription:      NewSubscriptionHandler(services.Subscription),
		BlockRelationship: NewBlockRelationshipHandler(services.BlockRelationship),
	}
}
