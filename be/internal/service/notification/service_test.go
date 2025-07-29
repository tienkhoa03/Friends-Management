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

func TestNotificationService_GetUpdateRecipients(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock.NewMockUserRepository(ctrl)
	mockBlockRepo := mock.NewMockBlockRelationshipRepository(ctrl)
	mockFriendshipRepo := mock.NewMockFriendshipRepository(ctrl)
	mockSubscriptionRepo := mock.NewMockSubscriptionRepository(ctrl)

	service := NewNotificationService(mockBlockRepo, mockUserRepo, mockFriendshipRepo, mockSubscriptionRepo)

	t.Run("Success", func(t *testing.T) {
		sender := &entity.User{Id: 1, Email: "sender@example.com"}
		friend := &entity.User{Id: 2, Email: "friend@example.com"}
		subscriber := &entity.User{Id: 3, Email: "subscriber@example.com"}
		mentioned := &entity.User{Id: 4, Email: "mentioned@example.com"}

		mockUserRepo.EXPECT().GetUserByEmail("sender@example.com").Return(sender, nil)
		mockBlockRepo.EXPECT().GetBlockRequestorIds(int64(1)).Return([]int64{}, nil)
		mockFriendshipRepo.EXPECT().RetrieveFriendIds(int64(1)).Return([]int64{2}, nil)
		mockSubscriptionRepo.EXPECT().GetAllSubscriberIds(int64(1)).Return([]int64{3}, nil)
		mockUserRepo.EXPECT().GetUsersFromEmails([]string{"mentioned@example.com"}).Return([]*entity.User{mentioned}, nil)
		mockUserRepo.EXPECT().GetUserFromIds(gomock.Any()).Return([]*entity.User{friend, subscriber, mentioned}, nil)

		recipients, err := service.GetUpdateRecipients("sender@example.com", "Hello mentioned@example.com")
		assert.NoError(t, err)
		assert.Len(t, recipients, 3)
	})

	t.Run("SenderNotFound", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserByEmail("sender@example.com").Return(nil, gorm.ErrRecordNotFound)

		recipients, err := service.GetUpdateRecipients("sender@example.com", "Hello world")
		assert.Nil(t, recipients)
		assert.Equal(t, ErrUserNotFound, err)
	})

	t.Run("DatabaseErrorOnGetUser", func(t *testing.T) {
		dbErr := errors.New("database error")
		mockUserRepo.EXPECT().GetUserByEmail("sender@example.com").Return(nil, dbErr)

		recipients, err := service.GetUpdateRecipients("sender@example.com", "Hello world")
		assert.Nil(t, recipients)
		assert.Equal(t, dbErr, err)
	})

	t.Run("ErrorOnGetBlockRequestorIds", func(t *testing.T) {
		sender := &entity.User{Id: 1, Email: "sender@example.com"}
		dbErr := errors.New("database error")

		mockUserRepo.EXPECT().GetUserByEmail("sender@example.com").Return(sender, nil)
		mockBlockRepo.EXPECT().GetBlockRequestorIds(int64(1)).Return(nil, dbErr)

		recipients, err := service.GetUpdateRecipients("sender@example.com", "Hello world")
		assert.Nil(t, recipients)
		assert.Equal(t, dbErr, err)
	})

	t.Run("ErrorOnRetrieveFriendIds", func(t *testing.T) {
		sender := &entity.User{Id: 1, Email: "sender@example.com"}
		dbErr := errors.New("database error")

		mockUserRepo.EXPECT().GetUserByEmail("sender@example.com").Return(sender, nil)
		mockBlockRepo.EXPECT().GetBlockRequestorIds(int64(1)).Return([]int64{}, nil)
		mockFriendshipRepo.EXPECT().RetrieveFriendIds(int64(1)).Return(nil, dbErr)

		recipients, err := service.GetUpdateRecipients("sender@example.com", "Hello world")
		assert.Nil(t, recipients)
		assert.Equal(t, dbErr, err)
	})

	t.Run("ErrorOnGetAllSubscriberIds", func(t *testing.T) {
		sender := &entity.User{Id: 1, Email: "sender@example.com"}
		dbErr := errors.New("database error")

		mockUserRepo.EXPECT().GetUserByEmail("sender@example.com").Return(sender, nil)
		mockBlockRepo.EXPECT().GetBlockRequestorIds(int64(1)).Return([]int64{}, nil)
		mockFriendshipRepo.EXPECT().RetrieveFriendIds(int64(1)).Return([]int64{}, nil)
		mockSubscriptionRepo.EXPECT().GetAllSubscriberIds(int64(1)).Return(nil, dbErr)

		recipients, err := service.GetUpdateRecipients("sender@example.com", "Hello world")
		assert.Nil(t, recipients)
		assert.Equal(t, dbErr, err)
	})

	t.Run("ErrorOnGetUsersFromEmails", func(t *testing.T) {
		sender := &entity.User{Id: 1, Email: "sender@example.com"}
		dbErr := errors.New("database error")

		mockUserRepo.EXPECT().GetUserByEmail("sender@example.com").Return(sender, nil)
		mockBlockRepo.EXPECT().GetBlockRequestorIds(int64(1)).Return([]int64{}, nil)
		mockFriendshipRepo.EXPECT().RetrieveFriendIds(int64(1)).Return([]int64{}, nil)
		mockSubscriptionRepo.EXPECT().GetAllSubscriberIds(int64(1)).Return([]int64{}, nil)
		mockUserRepo.EXPECT().GetUsersFromEmails([]string{"mentioned@example.com"}).Return(nil, dbErr)

		recipients, err := service.GetUpdateRecipients("sender@example.com", "Hello mentioned@example.com")
		assert.Nil(t, recipients)
		assert.Equal(t, dbErr, err)
	})

	t.Run("ErrorOnGetUserFromIds", func(t *testing.T) {
		sender := &entity.User{Id: 1, Email: "sender@example.com"}
		mentioned := &entity.User{Id: 2, Email: "mentioned@example.com"}
		dbErr := errors.New("database error")

		mockUserRepo.EXPECT().GetUserByEmail("sender@example.com").Return(sender, nil)
		mockBlockRepo.EXPECT().GetBlockRequestorIds(int64(1)).Return([]int64{}, nil)
		mockFriendshipRepo.EXPECT().RetrieveFriendIds(int64(1)).Return([]int64{}, nil)
		mockSubscriptionRepo.EXPECT().GetAllSubscriberIds(int64(1)).Return([]int64{}, nil)
		mockUserRepo.EXPECT().GetUsersFromEmails([]string{"mentioned@example.com"}).Return([]*entity.User{mentioned}, nil)
		mockUserRepo.EXPECT().GetUserFromIds(gomock.Any()).Return(nil, dbErr)

		recipients, err := service.GetUpdateRecipients("sender@example.com", "Hello mentioned@example.com")
		assert.Nil(t, recipients)
		assert.Equal(t, dbErr, err)
	})

	t.Run("BlockedUsersExcluded", func(t *testing.T) {
		sender := &entity.User{Id: 1, Email: "sender@example.com"}
		friend := &entity.User{Id: 2, Email: "friend@example.com"}
		// subscriber := &entity.User{Id: 3, Email: "subscriber@example.com"}

		mockUserRepo.EXPECT().GetUserByEmail("sender@example.com").Return(sender, nil)
		mockBlockRepo.EXPECT().GetBlockRequestorIds(int64(1)).Return([]int64{3}, nil) // subscriber is blocked
		mockFriendshipRepo.EXPECT().RetrieveFriendIds(int64(1)).Return([]int64{2}, nil)
		mockSubscriptionRepo.EXPECT().GetAllSubscriberIds(int64(1)).Return([]int64{3}, nil)
		mockUserRepo.EXPECT().GetUsersFromEmails(gomock.Any()).Return([]*entity.User{}, nil)
		mockUserRepo.EXPECT().GetUserFromIds([]int64{2}).Return([]*entity.User{friend}, nil)

		recipients, err := service.GetUpdateRecipients("sender@example.com", "Hello world")
		assert.NoError(t, err)
		assert.Len(t, recipients, 1)
		assert.Equal(t, "friend@example.com", recipients[0].Email)
	})

	t.Run("NoTextEmails", func(t *testing.T) {
		sender := &entity.User{Id: 1, Email: "sender@example.com"}
		friend := &entity.User{Id: 2, Email: "friend@example.com"}

		mockUserRepo.EXPECT().GetUserByEmail("sender@example.com").Return(sender, nil)
		mockBlockRepo.EXPECT().GetBlockRequestorIds(int64(1)).Return([]int64{}, nil)
		mockFriendshipRepo.EXPECT().RetrieveFriendIds(int64(1)).Return([]int64{2}, nil)
		mockSubscriptionRepo.EXPECT().GetAllSubscriberIds(int64(1)).Return([]int64{}, nil)
		mockUserRepo.EXPECT().GetUsersFromEmails(gomock.Any()).Return([]*entity.User{}, nil)
		mockUserRepo.EXPECT().GetUserFromIds([]int64{2}).Return([]*entity.User{friend}, nil)

		recipients, err := service.GetUpdateRecipients("sender@example.com", "Hello world no emails")
		assert.NoError(t, err)
		assert.Len(t, recipients, 1)
	})

	t.Run("EmptyRecipients", func(t *testing.T) {
		sender := &entity.User{Id: 1, Email: "sender@example.com"}

		mockUserRepo.EXPECT().GetUserByEmail("sender@example.com").Return(sender, nil)
		mockBlockRepo.EXPECT().GetBlockRequestorIds(int64(1)).Return([]int64{}, nil)
		mockFriendshipRepo.EXPECT().RetrieveFriendIds(int64(1)).Return([]int64{}, nil)
		mockSubscriptionRepo.EXPECT().GetAllSubscriberIds(int64(1)).Return([]int64{}, nil)
		mockUserRepo.EXPECT().GetUsersFromEmails(gomock.Any()).Return([]*entity.User{}, nil)
		mockUserRepo.EXPECT().GetUserFromIds([]int64{}).Return([]*entity.User{}, nil)

		recipients, err := service.GetUpdateRecipients("sender@example.com", "Hello world")
		assert.NoError(t, err)
		assert.Len(t, recipients, 0)
	})
}
