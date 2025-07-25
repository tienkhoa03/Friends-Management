package handler

import (
	"BE_Friends_Management/api/handler"
	"BE_Friends_Management/internal/domain/dto"
	service "BE_Friends_Management/internal/service/subscription"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSubscriptionService struct {
	mock.Mock
}

func (m *MockSubscriptionService) CreateSubscription(requestor, target string) error {
	args := m.Called(requestor, target)
	return args.Error(0)
}

func TestSubscriptionHandler_CreateSubscription(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		requestBody  interface{}
		serviceError error
		expectedCode int
	}{
		{
			name: "Success",
			requestBody: dto.CreateSubscriptionRequest{
				Requestor: "user1@example.com",
				Target:    "user2@example.com",
			},
			serviceError: nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "InvalidJSON",
			requestBody:  "invalid json",
			serviceError: nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "InvalidRequest",
			requestBody: dto.CreateSubscriptionRequest{
				Requestor: "user1@example.com",
				Target:    "user2@example.com",
			},
			serviceError: service.ErrInvalidRequest,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "UserNotFound",
			requestBody: dto.CreateSubscriptionRequest{
				Requestor: "user1@example.com",
				Target:    "user2@example.com",
			},
			serviceError: service.ErrUserNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			name: "AlreadySubscribed",
			requestBody: dto.CreateSubscriptionRequest{
				Requestor: "user1@example.com",
				Target:    "user2@example.com",
			},
			serviceError: service.ErrAlreadySubscribed,
			expectedCode: http.StatusConflict,
		},
		{
			name: "UnknownError",
			requestBody: dto.CreateSubscriptionRequest{
				Requestor: "user1@example.com",
				Target:    "user2@example.com",
			},
			serviceError: errors.New("unknown error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSubscriptionService)
			handler := handler.NewSubscriptionHandler(mockService)

			var jsonData []byte
			if request, ok := tt.requestBody.(dto.CreateSubscriptionRequest); ok {
				mockService.On("CreateSubscription", request.Requestor, request.Target).Return(tt.serviceError)
				jsonData, _ = json.Marshal(request)
			} else {
				jsonData = []byte(tt.requestBody.(string))
			}

			req, _ := http.NewRequest("POST", "/api/subscription", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			handler.CreateSubscription(c)

			assert.Equal(t, tt.expectedCode, w.Code)

			mockService.AssertExpectations(t)
		})
	}
}
