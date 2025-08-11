package handler

import (
	"BE_Friends_Management/api/handler"
	"BE_Friends_Management/internal/domain/dto"
	"BE_Friends_Management/internal/domain/entity"
	service "BE_Friends_Management/internal/service/friendship"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFriendshipService struct {
	mock.Mock
}

func (m *MockFriendshipService) CreateFriendship(authUserId int64, email1, email2 string) error {
	args := m.Called(authUserId, email1, email2)
	return args.Error(0)
}
func (m *MockFriendshipService) RetrieveFriendsList(authUserId int64, authUserRole string, email string) ([]*entity.User, error) {
	args := m.Called(authUserId, authUserRole, email)
	return args.Get(0).([]*entity.User), args.Error(1)
}
func (m *MockFriendshipService) RetrieveCommonFriends(authUserId int64, authUserRole string, email1, email2 string) ([]*entity.User, error) {
	args := m.Called(authUserId, authUserRole, email1, email2)
	return args.Get(0).([]*entity.User), args.Error(1)
}
func (m *MockFriendshipService) CountFriends(users []*entity.User) int64 {
	args := m.Called(users)
	return int64(args.Int(0))
}

func TestFriendshipHandler_CreateFriendship(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		authUserId     int64
		requestBody    interface{}
		serviceError   error
		expectedStatus int
		setupMock      func(*MockFriendshipService)
	}{
		{
			name:       "Success",
			authUserId: 1,
			requestBody: dto.CreateFriendshipRequest{
				Friends: []string{"andy@example.com", "john@example.com"},
			},
			serviceError:   nil,
			expectedStatus: http.StatusOK,
			setupMock: func(m *MockFriendshipService) {
				m.On("CreateFriendship", int64(1), "andy@example.com", "john@example.com").Return(nil)
			},
		},
		{
			name:           "Invalid JSON request",
			authUserId:     1,
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(m *MockFriendshipService) {},
		},
		{
			name:       "Service returns ErrInvalidRequest",
			authUserId: 1,
			requestBody: dto.CreateFriendshipRequest{
				Friends: []string{"andy@example.com", "john@example.com"},
			},
			serviceError:   service.ErrInvalidRequest,
			expectedStatus: http.StatusBadRequest,
			setupMock: func(m *MockFriendshipService) {
				m.On("CreateFriendship", int64(1), "andy@example.com", "john@example.com").Return(service.ErrInvalidRequest)
			},
		},
		{
			name:       "Service returns ErrUserNotFound",
			authUserId: 1,
			requestBody: dto.CreateFriendshipRequest{
				Friends: []string{"andy@example.com", "john@example.com"},
			},
			serviceError:   service.ErrUserNotFound,
			expectedStatus: http.StatusNotFound,
			setupMock: func(m *MockFriendshipService) {
				m.On("CreateFriendship", int64(1), "andy@example.com", "john@example.com").Return(service.ErrUserNotFound)
			},
		},
		{
			name:       "Service returns ErrAlreadyFriend",
			authUserId: 1,
			requestBody: dto.CreateFriendshipRequest{
				Friends: []string{"andy@example.com", "john@example.com"},
			},
			serviceError:   service.ErrAlreadyFriend,
			expectedStatus: http.StatusConflict,
			setupMock: func(m *MockFriendshipService) {
				m.On("CreateFriendship", int64(1), "andy@example.com", "john@example.com").Return(service.ErrAlreadyFriend)
			},
		},
		{
			name:       "Service returns unknown error",
			authUserId: 1,
			requestBody: dto.CreateFriendshipRequest{
				Friends: []string{"andy@example.com", "john@example.com"},
			},
			serviceError:   assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			setupMock: func(m *MockFriendshipService) {
				m.On("CreateFriendship", int64(1), "andy@example.com", "john@example.com").Return(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockFriendshipService)
			tt.setupMock(mockService)

			handler := handler.NewFriendshipHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			var reqBody []byte
			if str, ok := tt.requestBody.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, _ = json.Marshal(tt.requestBody)
			}

			c.Request = httptest.NewRequest(http.MethodPost, "/api/friendship", bytes.NewBuffer(reqBody))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Set("authUserId", 1)
			handler.CreateFriendship(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}
func TestFriendshipHandler_RetrieveFriendsList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		emailParam     string
		authUserId     int64
		authUserRole   string
		serviceUsers   []*entity.User
		serviceError   error
		expectedStatus int
		setupMock      func(*MockFriendshipService)
	}{
		{
			name:         "Success",
			authUserId:   1,
			authUserRole: "user",
			emailParam:   "andy@example.com",
			serviceUsers: []*entity.User{
				{Email: "john@example.com"},
				{Email: "jane@example.com"},
			},
			serviceError:   nil,
			expectedStatus: http.StatusOK,
			setupMock: func(m *MockFriendshipService) {
				users := []*entity.User{
					{Email: "john@example.com"},
					{Email: "jane@example.com"},
				}
				m.On("RetrieveFriendsList", int64(1), "user", "andy@example.com").Return(users, nil)
				m.On("CountFriends", users).Return(2)
			},
		},
		{
			name:           "Missing email parameter",
			authUserId:     1,
			authUserRole:   "user",
			emailParam:     "",
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(m *MockFriendshipService) {},
		},
		{
			name:           "Service returns ErrUserNotFound",
			authUserId:     1,
			authUserRole:   "user",
			emailParam:     "nonexistent@example.com",
			serviceUsers:   []*entity.User{},
			serviceError:   service.ErrUserNotFound,
			expectedStatus: http.StatusNotFound,
			setupMock: func(m *MockFriendshipService) {
				m.On("RetrieveFriendsList", int64(1), "user", "nonexistent@example.com").Return([]*entity.User{}, service.ErrUserNotFound)
			},
		},
		{
			name:           "Service returns unknown error",
			authUserId:     1,
			authUserRole:   "user",
			emailParam:     "andy@example.com",
			serviceUsers:   []*entity.User{},
			serviceError:   assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			setupMock: func(m *MockFriendshipService) {
				m.On("RetrieveFriendsList", int64(1), "user", "andy@example.com").Return([]*entity.User{}, assert.AnError)
			},
		},
		{
			name:           "Success with empty friends list",
			emailParam:     "lonely@example.com",
			serviceUsers:   []*entity.User{},
			serviceError:   nil,
			expectedStatus: http.StatusOK,
			setupMock: func(m *MockFriendshipService) {
				users := []*entity.User{}
				m.On("RetrieveFriendsList", int64(1), "user", "lonely@example.com").Return(users, nil)
				m.On("CountFriends", users).Return(0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockFriendshipService)
			tt.setupMock(mockService)

			handler := handler.NewFriendshipHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req := httptest.NewRequest(http.MethodGet, "/api/friendship/friends", nil)
			if tt.emailParam != "" {
				q := req.URL.Query()
				q.Add("email", tt.emailParam)
				req.URL.RawQuery = q.Encode()
			}
			c.Request = req
			c.Set("authUserId", 1)
			c.Set("authUserRole", "user")
			handler.RetrieveFriendsList(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestFriendshipHandler_RetrieveCommonFriends(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		authUserId     int64
		authUserRole   string
		email1Param    string
		email2Param    string
		serviceUsers   []*entity.User
		serviceError   error
		expectedStatus int
		setupMock      func(*MockFriendshipService)
	}{
		{
			name:         "Success",
			authUserId:   1,
			authUserRole: "user",
			email1Param:  "andy@example.com",
			email2Param:  "john@example.com",
			serviceUsers: []*entity.User{
				{Email: "jane@example.com"},
				{Email: "bob@example.com"},
			},
			serviceError:   nil,
			expectedStatus: http.StatusOK,
			setupMock: func(m *MockFriendshipService) {
				users := []*entity.User{
					{Email: "jane@example.com"},
					{Email: "bob@example.com"},
				}
				m.On("RetrieveCommonFriends", int64(1), "user", "andy@example.com", "john@example.com").Return(users, nil)
				m.On("CountFriends", users).Return(2)
			},
		},
		{
			name:           "Missing email1 parameter",
			authUserId:     1,
			authUserRole:   "user",
			email1Param:    "",
			email2Param:    "john@example.com",
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(m *MockFriendshipService) {},
		},
		{
			name:           "Missing email2 parameter",
			authUserId:     1,
			authUserRole:   "user",
			email1Param:    "andy@example.com",
			email2Param:    "",
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(m *MockFriendshipService) {},
		},
		{
			name:           "Missing both email parameters",
			authUserId:     1,
			authUserRole:   "user",
			email1Param:    "",
			email2Param:    "",
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(m *MockFriendshipService) {},
		},
		{
			name:           "Service returns ErrInvalidRequest",
			authUserId:     1,
			authUserRole:   "user",
			email1Param:    "andy@example.com",
			email2Param:    "john@example.com",
			serviceUsers:   []*entity.User{},
			serviceError:   service.ErrInvalidRequest,
			expectedStatus: http.StatusBadRequest,
			setupMock: func(m *MockFriendshipService) {
				m.On("RetrieveCommonFriends", int64(1), "user", "andy@example.com", "john@example.com").Return([]*entity.User{}, service.ErrInvalidRequest)
			},
		},
		{
			name:           "Service returns ErrUserNotFound",
			authUserId:     1,
			authUserRole:   "user",
			email1Param:    "nonexistent@example.com",
			email2Param:    "john@example.com",
			serviceUsers:   []*entity.User{},
			serviceError:   service.ErrUserNotFound,
			expectedStatus: http.StatusNotFound,
			setupMock: func(m *MockFriendshipService) {
				m.On("RetrieveCommonFriends", int64(1), "user", "nonexistent@example.com", "john@example.com").Return([]*entity.User{}, service.ErrUserNotFound)
			},
		},
		{
			name:           "Service returns unknown error",
			authUserId:     1,
			authUserRole:   "user",
			email1Param:    "andy@example.com",
			email2Param:    "john@example.com",
			serviceUsers:   []*entity.User{},
			serviceError:   assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			setupMock: func(m *MockFriendshipService) {
				m.On("RetrieveCommonFriends", int64(1), "user", "andy@example.com", "john@example.com").Return([]*entity.User{}, assert.AnError)
			},
		},
		{
			name:           "Success with no common friends",
			authUserId:     1,
			authUserRole:   "user",
			email1Param:    "andy@example.com",
			email2Param:    "john@example.com",
			serviceUsers:   []*entity.User{},
			serviceError:   nil,
			expectedStatus: http.StatusOK,
			setupMock: func(m *MockFriendshipService) {
				users := []*entity.User{}
				m.On("RetrieveCommonFriends", int64(1), "user", "andy@example.com", "john@example.com").Return(users, nil)
				m.On("CountFriends", users).Return(0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockFriendshipService)
			tt.setupMock(mockService)

			handler := handler.NewFriendshipHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req := httptest.NewRequest(http.MethodGet, "/api/friendship/common-friends", nil)
			q := req.URL.Query()
			if tt.email1Param != "" {
				q.Add("email1", tt.email1Param)
			}
			if tt.email2Param != "" {
				q.Add("email2", tt.email2Param)
			}
			req.URL.RawQuery = q.Encode()
			c.Request = req
			c.Set("authUserId", 1)
			c.Set("authUserRole", "user")
			handler.RetrieveCommonFriends(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}
