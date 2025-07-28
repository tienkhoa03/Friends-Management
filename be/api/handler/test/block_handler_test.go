package handler

import (
	"BE_Friends_Management/api/handler"
	"BE_Friends_Management/internal/domain/dto"
	service "BE_Friends_Management/internal/service/block_relationship"
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

type MockBlockRelationshipService struct {
	mock.Mock
}

func (m *MockBlockRelationshipService) CreateBlockRelationship(requestor, target string) error {
	args := m.Called(requestor, target)
	return args.Error(0)
}

func TestCreateBlockRelationship(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		serviceError   error
		expectedStatus int
		setupMock      func(*MockBlockRelationshipService)
	}{
		{
			name: "Success",
			requestBody: dto.CreateBlockRequest{
				Requestor: "user1@example.com",
				Target:    "user2@example.com",
			},
			serviceError:   nil,
			expectedStatus: http.StatusOK,
			setupMock: func(m *MockBlockRelationshipService) {
				m.On("CreateBlockRelationship", "user1@example.com", "user2@example.com").Return(nil)
			},
		},
		{
			name:           "Invalid JSON",
			requestBody:    "invalid json",
			serviceError:   nil,
			expectedStatus: http.StatusBadRequest,
			setupMock: func(m *MockBlockRelationshipService) {
				// No service call expected
			},
		},
		{
			name: "Service Error - Invalid Request",
			requestBody: dto.CreateBlockRequest{
				Requestor: "user1@example.com",
				Target:    "user2@example.com",
			},
			serviceError:   service.ErrInvalidRequest,
			expectedStatus: http.StatusBadRequest,
			setupMock: func(m *MockBlockRelationshipService) {
				m.On("CreateBlockRelationship", "user1@example.com", "user2@example.com").Return(service.ErrInvalidRequest)
			},
		},
		{
			name: "Service Error - User Not Found",
			requestBody: dto.CreateBlockRequest{
				Requestor: "user1@example.com",
				Target:    "user2@example.com",
			},
			serviceError:   service.ErrUserNotFound,
			expectedStatus: http.StatusNotFound,
			setupMock: func(m *MockBlockRelationshipService) {
				m.On("CreateBlockRelationship", "user1@example.com", "user2@example.com").Return(service.ErrUserNotFound)
			},
		},
		{
			name: "Service Error - Already Blocked",
			requestBody: dto.CreateBlockRequest{
				Requestor: "user1@example.com",
				Target:    "user2@example.com",
			},
			serviceError:   service.ErrAlreadyBlocked,
			expectedStatus: http.StatusConflict,
			setupMock: func(m *MockBlockRelationshipService) {
				m.On("CreateBlockRelationship", "user1@example.com", "user2@example.com").Return(service.ErrAlreadyBlocked)
			},
		},
		{
			name: "Service Error - Not Subscribed",
			requestBody: dto.CreateBlockRequest{
				Requestor: "user1@example.com",
				Target:    "user2@example.com",
			},
			serviceError:   service.ErrNotSubscribed,
			expectedStatus: http.StatusForbidden,
			setupMock: func(m *MockBlockRelationshipService) {
				m.On("CreateBlockRelationship", "user1@example.com", "user2@example.com").Return(service.ErrNotSubscribed)
			},
		},
		{
			name: "Service Error - Unknown Error",
			requestBody: dto.CreateBlockRequest{
				Requestor: "user1@example.com",
				Target:    "user2@example.com",
			},
			serviceError:   errors.New("unknown error"),
			expectedStatus: http.StatusInternalServerError,
			setupMock: func(m *MockBlockRelationshipService) {
				m.On("CreateBlockRelationship", "user1@example.com", "user2@example.com").Return(errors.New("unknown error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockBlockRelationshipService)
			tt.setupMock(mockService)

			handler := handler.NewBlockRelationshipHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			var jsonBody []byte
			var err error
			if tt.requestBody != "invalid json" {
				jsonBody, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			} else {
				jsonBody = []byte("invalid json")
			}

			c.Request, _ = http.NewRequest("POST", "/api/block", bytes.NewBuffer(jsonBody))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.CreateBlockRelationship(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}
