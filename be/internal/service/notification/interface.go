package service

import (
	"BE_Friends_Management/internal/domain/entity"
	"errors"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrNotPermitted = errors.New("action not permitted")
)

//go:generate mockgen -source=interface.go -destination=../mock/mock_notification_service.go

type NotificationService interface {
	GetUpdateRecipients(authUserId int64, authUserRole string, senderEmail, text string) ([]*entity.User, error)
}
