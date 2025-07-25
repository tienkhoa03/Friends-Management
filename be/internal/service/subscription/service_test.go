package service

import (
	"BE_Friends_Management/internal/domain/entity"
	mock "BE_Friends_Management/internal/repository/mock"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateSubscription(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSubscriptionRepo := mock.NewMockSubscriptionRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)
	service := NewSubscriptionService(mockSubscriptionRepo, mockUserRepo)

	t.Run("successful subscription creation", func(t *testing.T) {
		requestorEmail := "user1@example.com"
		targetEmail := "user2@example.com"
		user1 := &entity.User{Id: 1, Email: requestorEmail}
		user2 := &entity.User{Id: 2, Email: targetEmail}

		mockUserRepo.EXPECT().GetUserByEmail(requestorEmail).Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail(targetEmail).Return(user2, nil)
		mockSubscriptionRepo.EXPECT().CreateSubscription(user1.Id, user2.Id).Return(nil)

		err := service.CreateSubscription(requestorEmail, targetEmail)
		assert.NoError(t, err)
	})

	t.Run("requestor user not found", func(t *testing.T) {
		requestorEmail := "nonexistent@example.com"
		targetEmail := "user2@example.com"

		mockUserRepo.EXPECT().GetUserByEmail(requestorEmail).Return(nil, gorm.ErrRecordNotFound)

		err := service.CreateSubscription(requestorEmail, targetEmail)
		assert.Equal(t, ErrUserNotFound, err)
	})

	t.Run("target user not found", func(t *testing.T) {
		requestorEmail := "user1@example.com"
		targetEmail := "nonexistent@example.com"
		user1 := &entity.User{Id: 1, Email: requestorEmail}

		mockUserRepo.EXPECT().GetUserByEmail(requestorEmail).Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail(targetEmail).Return(nil, gorm.ErrRecordNotFound)

		err := service.CreateSubscription(requestorEmail, targetEmail)
		assert.Equal(t, ErrUserNotFound, err)
	})

	t.Run("same user subscription", func(t *testing.T) {
		email := "user1@example.com"
		user := &entity.User{Id: 1, Email: email}

		mockUserRepo.EXPECT().GetUserByEmail(email).Return(user, nil).Times(2)

		err := service.CreateSubscription(email, email)
		assert.Equal(t, ErrInvalidRequest, err)
	})

	t.Run("already subscribed", func(t *testing.T) {
		requestorEmail := "user1@example.com"
		targetEmail := "user2@example.com"
		user1 := &entity.User{Id: 1, Email: requestorEmail}
		user2 := &entity.User{Id: 2, Email: targetEmail}
		duplicateKeyError := errors.New("duplicate key constraint violation")

		mockUserRepo.EXPECT().GetUserByEmail(requestorEmail).Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail(targetEmail).Return(user2, nil)
		mockSubscriptionRepo.EXPECT().CreateSubscription(user1.Id, user2.Id).Return(duplicateKeyError)

		err := service.CreateSubscription(requestorEmail, targetEmail)
		assert.Equal(t, ErrAlreadySubscribed, err)
	})

	t.Run("database error on requestor lookup", func(t *testing.T) {
		requestorEmail := "user1@example.com"
		targetEmail := "user2@example.com"
		dbError := errors.New("database connection error")

		mockUserRepo.EXPECT().GetUserByEmail(requestorEmail).Return(nil, dbError)

		err := service.CreateSubscription(requestorEmail, targetEmail)
		assert.Equal(t, dbError, err)
	})

	t.Run("database error on target lookup", func(t *testing.T) {
		requestorEmail := "user1@example.com"
		targetEmail := "user2@example.com"
		user1 := &entity.User{Id: 1, Email: requestorEmail}
		dbError := errors.New("database connection error")

		mockUserRepo.EXPECT().GetUserByEmail(requestorEmail).Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail(targetEmail).Return(nil, dbError)

		err := service.CreateSubscription(requestorEmail, targetEmail)
		assert.Equal(t, dbError, err)
	})

	t.Run("database error on subscription creation", func(t *testing.T) {
		requestorEmail := "user1@example.com"
		targetEmail := "user2@example.com"
		user1 := &entity.User{Id: 1, Email: requestorEmail}
		user2 := &entity.User{Id: 2, Email: targetEmail}
		dbError := errors.New("database error")

		mockUserRepo.EXPECT().GetUserByEmail(requestorEmail).Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail(targetEmail).Return(user2, nil)
		mockSubscriptionRepo.EXPECT().CreateSubscription(user1.Id, user2.Id).Return(dbError)

		err := service.CreateSubscription(requestorEmail, targetEmail)
		assert.Equal(t, dbError, err)
	})
}
