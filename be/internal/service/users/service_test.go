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
		password := "123"
		role := "user"

		expectedUser := &entity.User{Id: userId, Email: email, Role: role}

		mockRepo.EXPECT().UpdateUser(gomock.Any()).Return(expectedUser, nil)

		user, err := service.UpdateUser(userId, email, password)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser.Id, user.Id)
		assert.Equal(t, expectedUser.Email, user.Email)
		assert.Equal(t, expectedUser.Role, user.Role)
	})

	t.Run("repository error", func(t *testing.T) {
		userId := int64(1)
		email := "updated@example.com"
		password := "123"
		expectedError := errors.New("user not found")
		mockRepo.EXPECT().UpdateUser(gomock.Any()).Return(nil, expectedError)

		user, err := service.UpdateUser(userId, email, password)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, user)
	})

	t.Run("database error", func(t *testing.T) {
		userId := int64(2)
		email := "test@example.com"
		password := "123"
		expectedError := errors.New("database connection failed")

		mockRepo.EXPECT().UpdateUser(gomock.Any()).Return(nil, expectedError)

		user, err := service.UpdateUser(userId, email, password)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, user)
	})

	t.Run("empty email", func(t *testing.T) {
		userId := int64(1)
		email := ""
		password := "123"
		expectedError := errors.New("email cannot be empty")

		mockRepo.EXPECT().UpdateUser(gomock.Any()).Return(nil, expectedError)

		user, err := service.UpdateUser(userId, email, password)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, user)
	})
}
