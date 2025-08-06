package service

import (
	entity "BE_Friends_Management/internal/domain/entity"
	mock "BE_Friends_Management/internal/repository/mock"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserService_GetAllUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	t.Run("success", func(t *testing.T) {
		expectedUsers := []*entity.User{
			{Id: 1, Email: "user1@example.com"},
			{Id: 2, Email: "user2@example.com"},
		}

		mockRepo.EXPECT().GetAllUser().Return(expectedUsers, nil)

		users, err := service.GetAllUser()

		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
	})

	t.Run("empty result", func(t *testing.T) {
		expectedUsers := []*entity.User{}

		mockRepo.EXPECT().GetAllUser().Return(expectedUsers, nil)

		users, err := service.GetAllUser()

		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
	})

	t.Run("repository error", func(t *testing.T) {
		expectedError := errors.New("database connection failed")

		mockRepo.EXPECT().GetAllUser().Return(nil, expectedError)

		users, err := service.GetAllUser()

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, users)
	})
}
func TestUserService_GetUserById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	t.Run("success", func(t *testing.T) {
		userId := int64(1)
		expectedUser := &entity.User{Id: 1, Email: "user1@example.com"}

		mockRepo.EXPECT().GetUserById(userId).Return(expectedUser, nil)

		user, err := service.GetUserById(userId)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("repository error", func(t *testing.T) {
		userId := int64(1)
		expectedError := errors.New("user not found")

		mockRepo.EXPECT().GetUserById(userId).Return(nil, expectedError)

		user, err := service.GetUserById(userId)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, user)
	})
}
func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	t.Run("success", func(t *testing.T) {
		email := "newuser@example.com"
		expectedUser := &entity.User{Id: 1, Email: email}
		inputUser := &entity.User{Email: email}

		mockRepo.EXPECT().CreateUser(inputUser).Return(expectedUser, nil)

		user, err := service.CreateUser(email)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("repository error", func(t *testing.T) {
		email := "duplicate@example.com"
		expectedError := errors.New("email already exists")
		inputUser := &entity.User{Email: email}

		mockRepo.EXPECT().CreateUser(inputUser).Return(nil, expectedError)

		user, err := service.CreateUser(email)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, user)
	})

	t.Run("empty email", func(t *testing.T) {
		email := ""
		expectedError := errors.New("email cannot be empty")
		inputUser := &entity.User{Email: email}

		mockRepo.EXPECT().CreateUser(inputUser).Return(nil, expectedError)

		user, err := service.CreateUser(email)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, user)
	})
}
func TestUserService_DeleteUserById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	t.Run("success", func(t *testing.T) {
		userId := int64(1)

		mockRepo.EXPECT().DeleteUserById(userId).Return(nil)

		err := service.DeleteUserById(userId)

		assert.NoError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		userId := int64(1)
		expectedError := errors.New("user not found")

		mockRepo.EXPECT().DeleteUserById(userId).Return(expectedError)

		err := service.DeleteUserById(userId)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("database error", func(t *testing.T) {
		userId := int64(2)
		expectedError := errors.New("database connection failed")

		mockRepo.EXPECT().DeleteUserById(userId).Return(expectedError)

		err := service.DeleteUserById(userId)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}
func TestUserService_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	t.Run("success", func(t *testing.T) {
		userId := int64(1)
		email := "updated@example.com"
		expectedUser := &entity.User{Id: userId, Email: email}
		inputUser := &entity.User{Id: userId, Email: email}

		mockRepo.EXPECT().UpdateUser(inputUser).Return(expectedUser, nil)

		user, err := service.UpdateUser(userId, email)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("repository error", func(t *testing.T) {
		userId := int64(1)
		email := "updated@example.com"
		expectedError := errors.New("user not found")
		inputUser := &entity.User{Id: userId, Email: email}

		mockRepo.EXPECT().UpdateUser(inputUser).Return(nil, expectedError)

		user, err := service.UpdateUser(userId, email)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, user)
	})

	t.Run("database error", func(t *testing.T) {
		userId := int64(2)
		email := "test@example.com"
		expectedError := errors.New("database connection failed")
		inputUser := &entity.User{Id: userId, Email: email}

		mockRepo.EXPECT().UpdateUser(inputUser).Return(nil, expectedError)

		user, err := service.UpdateUser(userId, email)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, user)
	})

	t.Run("empty email", func(t *testing.T) {
		userId := int64(1)
		email := ""
		expectedError := errors.New("email cannot be empty")
		inputUser := &entity.User{Id: userId, Email: email}

		mockRepo.EXPECT().UpdateUser(inputUser).Return(nil, expectedError)

		user, err := service.UpdateUser(userId, email)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, user)
	})
}




