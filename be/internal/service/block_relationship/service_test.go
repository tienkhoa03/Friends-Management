package service

import (
	"BE_Friends_Management/internal/domain/entity"
	mock "BE_Friends_Management/internal/repository/mock"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestBlockRelationshipService_CreateBlockRelationship(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock.NewMockUserRepository(ctrl)
	mockBlockRepo := mock.NewMockBlockRelationshipRepository(ctrl)
	mockFriendshipRepo := mock.NewMockFriendshipRepository(ctrl)
	mockSubscriptionRepo := mock.NewMockSubscriptionRepository(ctrl)

	db, mockSQL, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	mockBlockRepo.EXPECT().GetDB().Return(gormDB).AnyTimes()

	service := NewBlockRelationshipService(mockBlockRepo, mockUserRepo, mockFriendshipRepo, mockSubscriptionRepo)

	t.Run("Success_NoFriendship_NoSubscription", func(t *testing.T) {
		authUserId := int64(1)
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}
		user2 := &entity.User{Id: 2, Email: "user2@example.com"}

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user2@example.com").Return(user2, nil)
		mockFriendshipRepo.EXPECT().GetFriendship(int64(1), int64(2)).Return(nil, gorm.ErrRecordNotFound)
		mockSubscriptionRepo.EXPECT().GetSubscription(int64(1), int64(2)).Return(nil, gorm.ErrRecordNotFound)

		mockSQL.ExpectBegin()
		mockBlockRepo.EXPECT().CreateBlockRelationship(gomock.Any(), int64(1), int64(2)).Return(nil)
		mockSQL.ExpectCommit()

		err := service.CreateBlockRelationship(authUserId, "user1@example.com", "user2@example.com")
		assert.NoError(t, err)
	})

	t.Run("Success_HasFriendship_HasSubscription", func(t *testing.T) {
		authUserId := int64(1)
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}
		user2 := &entity.User{Id: 2, Email: "user2@example.com"}

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user2@example.com").Return(user2, nil)
		mockFriendshipRepo.EXPECT().GetFriendship(int64(1), int64(2)).Return(&entity.Friendship{}, nil)
		mockSubscriptionRepo.EXPECT().GetSubscription(int64(1), int64(2)).Return(&entity.Subscription{}, nil)

		mockSQL.ExpectBegin()
		mockSubscriptionRepo.EXPECT().DeleteSubscription(gomock.Any(), int64(1), int64(2)).Return(nil)
		mockBlockRepo.EXPECT().CreateBlockRelationship(gomock.Any(), int64(1), int64(2)).Return(nil)
		mockSQL.ExpectCommit()

		err := service.CreateBlockRelationship(authUserId, "user1@example.com", "user2@example.com")
		assert.NoError(t, err)
	})

	t.Run("Success_NoFriendship_HasSubscription", func(t *testing.T) {
		authUserId := int64(1)
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}
		user2 := &entity.User{Id: 2, Email: "user2@example.com"}

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user2@example.com").Return(user2, nil)
		mockFriendshipRepo.EXPECT().GetFriendship(int64(1), int64(2)).Return(nil, gorm.ErrRecordNotFound)
		mockSubscriptionRepo.EXPECT().GetSubscription(int64(1), int64(2)).Return(&entity.Subscription{}, nil)

		mockSQL.ExpectBegin()
		mockSubscriptionRepo.EXPECT().DeleteSubscription(gomock.Any(), int64(1), int64(2)).Return(nil)
		mockBlockRepo.EXPECT().CreateBlockRelationship(gomock.Any(), int64(1), int64(2)).Return(nil)
		mockSQL.ExpectCommit()

		err := service.CreateBlockRelationship(authUserId, "user1@example.com", "user2@example.com")
		assert.NoError(t, err)
	})

	t.Run("RequestorNotFound", func(t *testing.T) {
		authUserId := int64(1)
		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(nil, gorm.ErrRecordNotFound)

		err := service.CreateBlockRelationship(authUserId, "user1@example.com", "user2@example.com")
		assert.Equal(t, ErrUserNotFound, err)
	})

	t.Run("TargetNotFound", func(t *testing.T) {
		authUserId := int64(1)
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user2@example.com").Return(nil, gorm.ErrRecordNotFound)

		err := service.CreateBlockRelationship(authUserId, "user1@example.com", "user2@example.com")
		assert.Equal(t, ErrUserNotFound, err)
	})

	t.Run("SameUser", func(t *testing.T) {
		authUserId := int64(1)
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil).Times(2)

		err := service.CreateBlockRelationship(authUserId, "user1@example.com", "user1@example.com")
		assert.Equal(t, ErrInvalidRequest, err)
	})

	t.Run("AlreadyBlocked", func(t *testing.T) {
		authUserId := int64(1)
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}
		user2 := &entity.User{Id: 2, Email: "user2@example.com"}
		duplicateErr := errors.New("duplicate key constraint")

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user2@example.com").Return(user2, nil)
		mockFriendshipRepo.EXPECT().GetFriendship(int64(1), int64(2)).Return(nil, gorm.ErrRecordNotFound)
		mockSubscriptionRepo.EXPECT().GetSubscription(int64(1), int64(2)).Return(nil, gorm.ErrRecordNotFound)

		mockSQL.ExpectBegin()
		mockBlockRepo.EXPECT().CreateBlockRelationship(gomock.Any(), int64(1), int64(2)).Return(duplicateErr)
		mockSQL.ExpectRollback()

		err := service.CreateBlockRelationship(authUserId, "user1@example.com", "user2@example.com")
		assert.Equal(t, ErrAlreadyBlocked, err)
	})

	t.Run("HasFriendship_NoSubscription_Error", func(t *testing.T) {
		authUserId := int64(1)
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}
		user2 := &entity.User{Id: 2, Email: "user2@example.com"}

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user2@example.com").Return(user2, nil)
		mockFriendshipRepo.EXPECT().GetFriendship(int64(1), int64(2)).Return(&entity.Friendship{}, nil)
		mockSubscriptionRepo.EXPECT().GetSubscription(int64(1), int64(2)).Return(nil, gorm.ErrRecordNotFound)

		mockSQL.ExpectBegin()
		mockSQL.ExpectRollback()

		err := service.CreateBlockRelationship(authUserId, "user1@example.com", "user2@example.com")
		assert.Equal(t, ErrNotSubscribed, err)
	})
	t.Run("HasFriendship_NoSubscription_Error", func(t *testing.T) {
		authUserId := int64(1)
		dbErr := errors.New("database error")

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(nil, dbErr)

		err := service.CreateBlockRelationship(authUserId, "user1@example.com", "user2@example.com")
		assert.Equal(t, dbErr, err)
	})
}
