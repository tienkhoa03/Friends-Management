package repository

//go:generate mockgen -source=interface.go -destination=../mock/mock_subscription_repository.go

type SubscriptionRepository interface {
	CreateSubscription(requestorId, targetId int64) error
}
