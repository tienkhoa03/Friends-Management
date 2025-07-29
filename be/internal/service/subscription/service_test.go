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

func TestSubscriptionService_CreateSubscription(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock.NewMockUserRepository(ctrl)
	mockSubscriptionRepo := mock.NewMockSubscriptionRepository(ctrl)
	mockFriendshipRepo := mock.NewMockFriendshipRepository(ctrl)
	mockBlockRepo := mock.NewMockBlockRelationshipRepository(ctrl)

	service := NewSubscriptionService(mockSubscriptionRepo, mockUserRepo, mockFriendshipRepo, mockBlockRepo)

	t.Run("Success", func(t *testing.T) {
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}
		user2 := &entity.User{Id: 2, Email: "user2@example.com"}

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user2@example.com").Return(user2, nil)
		mockBlockRepo.EXPECT().GetBlockRelationship(int64(1), int64(2)).Return(nil, gorm.ErrRecordNotFound)
		mockSubscriptionRepo.EXPECT().CreateSubscription(int64(1), int64(2)).Return(nil)

		err := service.CreateSubscription("user1@example.com", "user2@example.com")
		assert.NoError(t, err)
	})

	t.Run("RequestorNotFound", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(nil, gorm.ErrRecordNotFound)

		err := service.CreateSubscription("user1@example.com", "user2@example.com")
		assert.Equal(t, ErrUserNotFound, err)
	})

	t.Run("TargetNotFound", func(t *testing.T) {
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user2@example.com").Return(nil, gorm.ErrRecordNotFound)

		err := service.CreateSubscription("user1@example.com", "user2@example.com")
		assert.Equal(t, ErrUserNotFound, err)
	})

	t.Run("SameUser", func(t *testing.T) {
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)

		err := service.CreateSubscription("user1@example.com", "user1@example.com")
		assert.Equal(t, ErrInvalidRequest, err)
	})

	t.Run("UserIsBlocked", func(t *testing.T) {
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}
		user2 := &entity.User{Id: 2, Email: "user2@example.com"}

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user2@example.com").Return(user2, nil)
		mockBlockRepo.EXPECT().GetBlockRelationship(int64(1), int64(2)).Return(&entity.BlockRelationship{}, nil)
		mockFriendshipRepo.EXPECT().GetFriendship(int64(1), int64(2)).Return(nil, gorm.ErrRecordNotFound)
		err := service.CreateSubscription("user1@example.com", "user2@example.com")
		assert.Equal(t, ErrIsBlocked, err)
	})

	t.Run("SuccessRemoveBlock", func(t *testing.T) {
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}
		user2 := &entity.User{Id: 2, Email: "user2@example.com"}

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user2@example.com").Return(user2, nil)
		mockBlockRepo.EXPECT().GetBlockRelationship(int64(1), int64(2)).Return(&entity.BlockRelationship{}, nil)
		mockFriendshipRepo.EXPECT().GetFriendship(int64(1), int64(2)).Return(&entity.Friendship{}, nil)
		mockBlockRepo.EXPECT().DeleteBlockRelationship(int64(1), int64(2)).Return(nil)
		mockSubscriptionRepo.EXPECT().CreateSubscription(int64(1), int64(2)).Return(nil)
		err := service.CreateSubscription("user1@example.com", "user2@example.com")
		assert.NoError(t, err)
	})

	t.Run("AlreadySubscribed", func(t *testing.T) {
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}
		user2 := &entity.User{Id: 2, Email: "user2@example.com"}
		duplicateErr := errors.New("duplicate key constraint")

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user2@example.com").Return(user2, nil)
		mockBlockRepo.EXPECT().GetBlockRelationship(int64(1), int64(2)).Return(nil, gorm.ErrRecordNotFound)
		mockSubscriptionRepo.EXPECT().CreateSubscription(int64(1), int64(2)).Return(duplicateErr)

		err := service.CreateSubscription("user1@example.com", "user2@example.com")
		assert.Equal(t, ErrAlreadySubscribed, err)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		dbErr := errors.New("database error")
		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(nil, dbErr)

		err := service.CreateSubscription("user1@example.com", "user2@example.com")
		assert.Equal(t, dbErr, err)
	})
}
