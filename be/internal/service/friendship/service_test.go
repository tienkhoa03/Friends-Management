package service

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	entity "BE_Friends_Management/internal/domain/entity"
	mock "BE_Friends_Management/internal/repository/mock"
)

func TestCreateFriendship(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFriendshipRepo := mock.NewMockFriendshipRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)
	service := NewFriendshipService(mockFriendshipRepo, mockUserRepo)

	t.Run("Success - user1.Id < user2.Id", func(t *testing.T) {
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}
		user2 := &entity.User{Id: 2, Email: "user2@example.com"}

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user2@example.com").Return(user2, nil)
		mockFriendshipRepo.EXPECT().CreateFriendship(int64(1), int64(2)).Return(nil)

		err := service.CreateFriendship("user1@example.com", "user2@example.com")
		assert.NoError(t, err)
	})

	t.Run("Success - user1.Id > user2.Id", func(t *testing.T) {
		user1 := &entity.User{Id: 2, Email: "user1@example.com"}
		user2 := &entity.User{Id: 1, Email: "user2@example.com"}

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user2@example.com").Return(user2, nil)
		mockFriendshipRepo.EXPECT().CreateFriendship(int64(1), int64(2)).Return(nil)

		err := service.CreateFriendship("user1@example.com", "user2@example.com")
		assert.NoError(t, err)
	})

	t.Run("Error - user1 not found", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(nil, gorm.ErrRecordNotFound)

		err := service.CreateFriendship("user1@example.com", "user2@example.com")
		assert.Equal(t, ErrUserNotFound, err)
	})

	t.Run("Error - user2 not found", func(t *testing.T) {
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user2@example.com").Return(nil, gorm.ErrRecordNotFound)

		err := service.CreateFriendship("user1@example.com", "user2@example.com")
		assert.Equal(t, ErrUserNotFound, err)
	})

	t.Run("Error - same user", func(t *testing.T) {
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}
		user2 := &entity.User{Id: 1, Email: "user1@example.com"}

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user2, nil)

		err := service.CreateFriendship("user1@example.com", "user1@example.com")
		assert.Equal(t, ErrInvalidRequest, err)
	})

	t.Run("Error - already friends", func(t *testing.T) {
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}
		user2 := &entity.User{Id: 2, Email: "user2@example.com"}
		duplicateErr := errors.New("duplicate key value violates unique constraint")

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user2@example.com").Return(user2, nil)
		mockFriendshipRepo.EXPECT().CreateFriendship(int64(1), int64(2)).Return(duplicateErr)

		err := service.CreateFriendship("user1@example.com", "user2@example.com")
		assert.Equal(t, ErrAlreadyFriend, err)
	})

	t.Run("Error - user repo error", func(t *testing.T) {
		repoErr := errors.New("database error")
		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(nil, repoErr)

		err := service.CreateFriendship("user1@example.com", "user2@example.com")
		assert.Equal(t, repoErr, err)
	})

	t.Run("Error - friendship repo error", func(t *testing.T) {
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}
		user2 := &entity.User{Id: 2, Email: "user2@example.com"}
		repoErr := errors.New("database error")

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user2@example.com").Return(user2, nil)
		mockFriendshipRepo.EXPECT().CreateFriendship(int64(1), int64(2)).Return(repoErr)

		err := service.CreateFriendship("user1@example.com", "user2@example.com")
		assert.Equal(t, repoErr, err)
	})
}
