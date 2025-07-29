package handler

import (
	"BE_Friends_Management/api/handler"
	"BE_Friends_Management/internal/domain/dto"
	"BE_Friends_Management/internal/domain/entity"
	service "BE_Friends_Management/internal/service/notification"
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

type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) GetUpdateRecipients(sender, text string) ([]*entity.User, error) {
	args := m.Called(sender, text)
	return args.Get(0).([]*entity.User), args.Error(1)
}

func TestGetUpdateRecipients(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		request      interface{}
		mockSetup    func(*MockNotificationService)
		expectedCode int
		expectPanic  bool
	}{
		{
			name: "Success",
			request: dto.GetUpdateRecipientsRequest{
				Sender: "sender@example.com",
				Text:   "Hello @user1@example.com",
			},
			mockSetup: func(m *MockNotificationService) {
				users := []*entity.User{
					{Email: "user1@example.com"},
					{Email: "user2@example.com"},
				}
				m.On("GetUpdateRecipients", "sender@example.com", "Hello @user1@example.com").Return(users, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid Request",
			request:      "invalid json",
			mockSetup:    func(m *MockNotificationService) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "User Not Found",
			request: dto.GetUpdateRecipientsRequest{
				Sender: "sender@example.com",
				Text:   "Hello",
			},
			mockSetup: func(m *MockNotificationService) {
				m.On("GetUpdateRecipients", "sender@example.com", "Hello").Return([]*entity.User{}, service.ErrUserNotFound)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "Unknown Error",
			request: dto.GetUpdateRecipientsRequest{
				Sender: "sender@example.com",
				Text:   "Hello",
			},
			mockSetup: func(m *MockNotificationService) {
				m.On("GetUpdateRecipients", "sender@example.com", "Hello").Return([]*entity.User{}, errors.New("database error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockNotificationService)
			tt.mockSetup(mockService)
			handler := handler.NewNotificationHandler(mockService)

			var body []byte
			if tt.name == "Invalid Request" {
				body = []byte("invalid json")
			} else {
				body, _ = json.Marshal(tt.request)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/update-recipients", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			handler.GetUpdateRecipients(c)
			assert.Equal(t, tt.expectedCode, w.Code)

			mockService.AssertExpectations(t)
		})
	}
}
