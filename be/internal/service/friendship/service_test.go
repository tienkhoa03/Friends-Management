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
func TestRetrieveFriendsList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFriendshipRepo := mock.NewMockFriendshipRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)
	service := NewFriendshipService(mockFriendshipRepo, mockUserRepo)

	t.Run("Success - retrieve friends list", func(t *testing.T) {
		user := &entity.User{Id: 1, Email: "user@example.com"}
		friend1 := &entity.User{Id: 2, Email: "friend1@example.com"}
		friend2 := &entity.User{Id: 3, Email: "friend2@example.com"}
		friendIds := []int64{2, 3}

		mockUserRepo.EXPECT().GetUserByEmail("user@example.com").Return(user, nil)
		mockFriendshipRepo.EXPECT().RetrieveFriendsList(int64(1)).Return(friendIds, nil)
		mockUserRepo.EXPECT().GetUserById(int64(2)).Return(friend1, nil)
		mockUserRepo.EXPECT().GetUserById(int64(3)).Return(friend2, nil)

		friends, err := service.RetrieveFriendsList("user@example.com")
		assert.NoError(t, err)
		assert.Len(t, friends, 2)
		assert.Equal(t, friend1, friends[0])
		assert.Equal(t, friend2, friends[1])
	})

	t.Run("Success - empty friends list", func(t *testing.T) {
		user := &entity.User{Id: 1, Email: "user@example.com"}
		friendIds := []int64{}

		mockUserRepo.EXPECT().GetUserByEmail("user@example.com").Return(user, nil)
		mockFriendshipRepo.EXPECT().RetrieveFriendsList(int64(1)).Return(friendIds, nil)

		friends, err := service.RetrieveFriendsList("user@example.com")
		assert.NoError(t, err)
		assert.Len(t, friends, 0)
	})

	t.Run("Error - user not found", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserByEmail("user@example.com").Return(nil, gorm.ErrRecordNotFound)

		friends, err := service.RetrieveFriendsList("user@example.com")
		assert.Nil(t, friends)
		assert.Equal(t, ErrUserNotFound, err)
	})

	t.Run("Error - user repo error", func(t *testing.T) {
		repoErr := errors.New("database error")
		mockUserRepo.EXPECT().GetUserByEmail("user@example.com").Return(nil, repoErr)

		friends, err := service.RetrieveFriendsList("user@example.com")
		assert.Nil(t, friends)
		assert.Equal(t, repoErr, err)
	})

	t.Run("Error - friendship repo error", func(t *testing.T) {
		user := &entity.User{Id: 1, Email: "user@example.com"}
		repoErr := errors.New("friendship database error")

		mockUserRepo.EXPECT().GetUserByEmail("user@example.com").Return(user, nil)
		mockFriendshipRepo.EXPECT().RetrieveFriendsList(int64(1)).Return(nil, repoErr)

		friends, err := service.RetrieveFriendsList("user@example.com")
		assert.Nil(t, friends)
		assert.Equal(t, repoErr, err)
	})

	t.Run("Error - friend lookup error", func(t *testing.T) {
		user := &entity.User{Id: 1, Email: "user@example.com"}
		friendIds := []int64{2, 3}
		repoErr := errors.New("friend lookup error")

		mockUserRepo.EXPECT().GetUserByEmail("user@example.com").Return(user, nil)
		mockFriendshipRepo.EXPECT().RetrieveFriendsList(int64(1)).Return(friendIds, nil)
		mockUserRepo.EXPECT().GetUserById(int64(2)).Return(nil, repoErr)

		friends, err := service.RetrieveFriendsList("user@example.com")
		assert.Nil(t, friends)
		assert.Equal(t, repoErr, err)
	})
}
func TestRetrieveCommonFriends(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFriendshipRepo := mock.NewMockFriendshipRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)
	service := NewFriendshipService(mockFriendshipRepo, mockUserRepo)

	t.Run("Success - common friends found", func(t *testing.T) {
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}
		user2 := &entity.User{Id: 2, Email: "user2@example.com"}
		commonFriend1 := &entity.User{Id: 3, Email: "friend1@example.com"}
		commonFriend2 := &entity.User{Id: 4, Email: "friend2@example.com"}
		
		friendIdsOfUser1 := []int64{3, 4, 5}
		friendIdsOfUser2 := []int64{3, 4, 6}

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user2@example.com").Return(user2, nil)
		mockFriendshipRepo.EXPECT().RetrieveFriendsList(int64(1)).Return(friendIdsOfUser1, nil)
		mockFriendshipRepo.EXPECT().RetrieveFriendsList(int64(2)).Return(friendIdsOfUser2, nil)
		mockUserRepo.EXPECT().GetUserById(int64(3)).Return(commonFriend1, nil)
		mockUserRepo.EXPECT().GetUserById(int64(4)).Return(commonFriend2, nil)

		commonFriends, err := service.RetrieveCommonFriends("user1@example.com", "user2@example.com")
		assert.NoError(t, err)
		assert.Len(t, commonFriends, 2)
		assert.Equal(t, commonFriend1, commonFriends[0])
		assert.Equal(t, commonFriend2, commonFriends[1])
	})

	t.Run("Success - no common friends", func(t *testing.T) {
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}
		user2 := &entity.User{Id: 2, Email: "user2@example.com"}
		
		friendIdsOfUser1 := []int64{3, 4}
		friendIdsOfUser2 := []int64{5, 6}

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user2@example.com").Return(user2, nil)
		mockFriendshipRepo.EXPECT().RetrieveFriendsList(int64(1)).Return(friendIdsOfUser1, nil)
		mockFriendshipRepo.EXPECT().RetrieveFriendsList(int64(2)).Return(friendIdsOfUser2, nil)

		commonFriends, err := service.RetrieveCommonFriends("user1@example.com", "user2@example.com")
		assert.NoError(t, err)
		assert.Len(t, commonFriends, 0)
	})

	t.Run("Success - empty friends lists", func(t *testing.T) {
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}
		user2 := &entity.User{Id: 2, Email: "user2@example.com"}
		
		friendIdsOfUser1 := []int64{}
		friendIdsOfUser2 := []int64{}

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user2@example.com").Return(user2, nil)
		mockFriendshipRepo.EXPECT().RetrieveFriendsList(int64(1)).Return(friendIdsOfUser1, nil)
		mockFriendshipRepo.EXPECT().RetrieveFriendsList(int64(2)).Return(friendIdsOfUser2, nil)

		commonFriends, err := service.RetrieveCommonFriends("user1@example.com", "user2@example.com")
		assert.NoError(t, err)
		assert.Len(t, commonFriends, 0)
	})

	t.Run("Error - user1 not found", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(nil, gorm.ErrRecordNotFound)

		commonFriends, err := service.RetrieveCommonFriends("user1@example.com", "user2@example.com")
		assert.Nil(t, commonFriends)
		assert.Equal(t, ErrUserNotFound, err)
	})

	t.Run("Error - user2 not found", func(t *testing.T) {
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user2@example.com").Return(nil, gorm.ErrRecordNotFound)

		commonFriends, err := service.RetrieveCommonFriends("user1@example.com", "user2@example.com")
		assert.Nil(t, commonFriends)
		assert.Equal(t, ErrUserNotFound, err)
	})

	t.Run("Error - same user", func(t *testing.T) {
		user := &entity.User{Id: 1, Email: "user@example.com"}

		mockUserRepo.EXPECT().GetUserByEmail("user@example.com").Return(user, nil).Times(2)

		commonFriends, err := service.RetrieveCommonFriends("user@example.com", "user@example.com")
		assert.Nil(t, commonFriends)
		assert.Equal(t, ErrInvalidRequest, err)
	})

	t.Run("Error - user1 repo error", func(t *testing.T) {
		repoErr := errors.New("database error")
		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(nil, repoErr)

		commonFriends, err := service.RetrieveCommonFriends("user1@example.com", "user2@example.com")
		assert.Nil(t, commonFriends)
		assert.Equal(t, repoErr, err)
	})

	t.Run("Error - user2 repo error", func(t *testing.T) {
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}
		repoErr := errors.New("database error")

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user2@example.com").Return(nil, repoErr)

		commonFriends, err := service.RetrieveCommonFriends("user1@example.com", "user2@example.com")
		assert.Nil(t, commonFriends)
		assert.Equal(t, repoErr, err)
	})

	t.Run("Error - friendship repo error for user1", func(t *testing.T) {
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}
		user2 := &entity.User{Id: 2, Email: "user2@example.com"}
		repoErr := errors.New("friendship database error")

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user2@example.com").Return(user2, nil)
		mockFriendshipRepo.EXPECT().RetrieveFriendsList(int64(1)).Return(nil, repoErr)

		commonFriends, err := service.RetrieveCommonFriends("user1@example.com", "user2@example.com")
		assert.Nil(t, commonFriends)
		assert.Equal(t, repoErr, err)
	})

	t.Run("Error - friendship repo error for user2", func(t *testing.T) {
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}
		user2 := &entity.User{Id: 2, Email: "user2@example.com"}
		friendIdsOfUser1 := []int64{3, 4}
		repoErr := errors.New("friendship database error")

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user2@example.com").Return(user2, nil)
		mockFriendshipRepo.EXPECT().RetrieveFriendsList(int64(1)).Return(friendIdsOfUser1, nil)
		mockFriendshipRepo.EXPECT().RetrieveFriendsList(int64(2)).Return(nil, repoErr)

		commonFriends, err := service.RetrieveCommonFriends("user1@example.com", "user2@example.com")
		assert.Nil(t, commonFriends)
		assert.Equal(t, repoErr, err)
	})

	t.Run("Error - common friend lookup error", func(t *testing.T) {
		user1 := &entity.User{Id: 1, Email: "user1@example.com"}
		user2 := &entity.User{Id: 2, Email: "user2@example.com"}
		
		friendIdsOfUser1 := []int64{3, 4}
		friendIdsOfUser2 := []int64{3, 5}
		repoErr := errors.New("friend lookup error")

		mockUserRepo.EXPECT().GetUserByEmail("user1@example.com").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByEmail("user2@example.com").Return(user2, nil)
		mockFriendshipRepo.EXPECT().RetrieveFriendsList(int64(1)).Return(friendIdsOfUser1, nil)
		mockFriendshipRepo.EXPECT().RetrieveFriendsList(int64(2)).Return(friendIdsOfUser2, nil)
		mockUserRepo.EXPECT().GetUserById(int64(3)).Return(nil, repoErr)

		commonFriends, err := service.RetrieveCommonFriends("user1@example.com", "user2@example.com")
		assert.Nil(t, commonFriends)
		assert.Equal(t, repoErr, err)
	})
}

func TestCountFriends(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFriendshipRepo := mock.NewMockFriendshipRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)
	service := NewFriendshipService(mockFriendshipRepo, mockUserRepo)

	t.Run("Success - count non-nil friends", func(t *testing.T) {
		friends := []*entity.User{
			{Id: 1, Email: "friend1@example.com"},
			{Id: 2, Email: "friend2@example.com"},
			{Id: 3, Email: "friend3@example.com"},
		}

		count := service.CountFriends(friends)
		assert.Equal(t, int64(3), count)
	})

	t.Run("Success - count with nil friends", func(t *testing.T) {
		friends := []*entity.User{
			{Id: 1, Email: "friend1@example.com"},
			nil,
			{Id: 3, Email: "friend3@example.com"},
			nil,
		}

		count := service.CountFriends(friends)
		assert.Equal(t, int64(2), count)
	})

	t.Run("Success - empty friends list", func(t *testing.T) {
		friends := []*entity.User{}

		count := service.CountFriends(friends)
		assert.Equal(t, int64(0), count)
	})

	t.Run("Success - all nil friends", func(t *testing.T) {
		friends := []*entity.User{nil, nil, nil}

		count := service.CountFriends(friends)
		assert.Equal(t, int64(0), count)
	})
}


