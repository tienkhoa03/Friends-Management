package handler

import (
	"BE_Friends_Management/api/handler"
	"BE_Friends_Management/internal/domain/dto"
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

func (m *MockFriendshipService) CreateFriendship(email1, email2 string) error {
	args := m.Called(email1, email2)
	return args.Error(0)
}

func TestFriendshipHandler_CreateFriendship(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		serviceError   error
		expectedStatus int
		setupMock      func(*MockFriendshipService)
	}{
		{
			name: "Success",
			requestBody: dto.CreateFriendshipRequest{
				Friends: []string{"andy@example.com", "john@example.com"},
			},
			serviceError:   nil,
			expectedStatus: http.StatusOK,
			setupMock: func(m *MockFriendshipService) {
				m.On("CreateFriendship", "andy@example.com", "john@example.com").Return(nil)
			},
		},
		{
			name:           "Invalid JSON request",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(m *MockFriendshipService) {},
		},
		{
			name: "Service returns ErrInvalidRequest",
			requestBody: dto.CreateFriendshipRequest{
				Friends: []string{"andy@example.com", "john@example.com"},
			},
			serviceError:   service.ErrInvalidRequest,
			expectedStatus: http.StatusBadRequest,
			setupMock: func(m *MockFriendshipService) {
				m.On("CreateFriendship", "andy@example.com", "john@example.com").Return(service.ErrInvalidRequest)
			},
		},
		{
			name: "Service returns ErrUserNotFound",
			requestBody: dto.CreateFriendshipRequest{
				Friends: []string{"andy@example.com", "john@example.com"},
			},
			serviceError:   service.ErrUserNotFound,
			expectedStatus: http.StatusNotFound,
			setupMock: func(m *MockFriendshipService) {
				m.On("CreateFriendship", "andy@example.com", "john@example.com").Return(service.ErrUserNotFound)
			},
		},
		{
			name: "Service returns ErrAlreadyFriend",
			requestBody: dto.CreateFriendshipRequest{
				Friends: []string{"andy@example.com", "john@example.com"},
			},
			serviceError:   service.ErrAlreadyFriend,
			expectedStatus: http.StatusConflict,
			setupMock: func(m *MockFriendshipService) {
				m.On("CreateFriendship", "andy@example.com", "john@example.com").Return(service.ErrAlreadyFriend)
			},
		},
		{
			name: "Service returns unknown error",
			requestBody: dto.CreateFriendshipRequest{
				Friends: []string{"andy@example.com", "john@example.com"},
			},
			serviceError:   assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			setupMock: func(m *MockFriendshipService) {
				m.On("CreateFriendship", "andy@example.com", "john@example.com").Return(assert.AnError)
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

			handler.CreateFriendship(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}
