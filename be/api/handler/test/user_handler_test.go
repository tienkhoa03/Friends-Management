package handler

import (
	"BE_Friends_Management/api/handler"
	"BE_Friends_Management/constant"
	"BE_Friends_Management/internal/domain/entity"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetAllUser() ([]*entity.User, error) {
	args := m.Called()
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockUserService) GetUserById(id int64) (*entity.User, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserService) CreateUser(email string) (*entity.User, error) {
	args := m.Called(email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserService) DeleteUserById(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) UpdateUser(id int64, email string) (*entity.User, error) {
	args := m.Called(id, email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func TestUserHandler_GetAllUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMock      func(*MockUserService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Success with users",
			setupMock: func(m *MockUserService) {
				users := []*entity.User{
					{Id: 1, Email: "user1@example.com", CreatedAt: time.Date(2025, 8, 1, 10, 0, 0, 0, time.UTC)},
					{Id: 2, Email: "user2@example.com", CreatedAt: time.Date(2025, 8, 2, 10, 0, 0, 0, time.UTC)},
				}
				m.On("GetAllUser").Return(users, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, constant.Success.GetResponseMessage(), response["message"])
				assert.NotNil(t, response["data"])

				data := response["data"].([]interface{})
				assert.Len(t, data, 2)
			},
		},
		{
			name: "Success with empty users",
			setupMock: func(m *MockUserService) {
				users := []*entity.User{}
				m.On("GetAllUser").Return(users, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, constant.Success.GetResponseMessage(), response["message"])
				assert.NotNil(t, response["data"])

				data := response["data"].([]interface{})
				assert.Len(t, data, 0)
			},
		},
		{
			name:           "Service returns database error",
			expectedStatus: http.StatusInternalServerError,
			setupMock: func(m *MockUserService) {
				m.On("GetAllUser").Return(nil, errors.New("database connection failed"))
			},
		},
		{
			name:           "Service returns generic error",
			expectedStatus: http.StatusInternalServerError,
			setupMock: func(m *MockUserService) {
				m.On("GetAllUser").Return(nil, errors.New("unexpected error"))
			},
		},
		{
			name: "Service returns nil users with no error",
			setupMock: func(m *MockUserService) {
				m.On("GetAllUser").Return([]*entity.User{}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				data := response["data"].([]interface{})
				assert.Len(t, data, 0)
				assert.Equal(t, constant.Success.GetResponseMessage(), response["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			userHandler := handler.NewUserHandler(mockService)

			tt.setupMock(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/api/users", nil)

			userHandler.GetAllUser(c)
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
			mockService.AssertExpectations(t)
		})
	}
}
func TestUserHandler_GetUserById(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		setupMock      func(*MockUserService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "Success - valid user ID",
			userID: "1",
			setupMock: func(m *MockUserService) {
				user := &entity.User{
					Id:        1,
					Email:     "user1@example.com",
					CreatedAt: time.Date(2025, 8, 1, 10, 0, 0, 0, time.UTC),
				}
				m.On("GetUserById", int64(1)).Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, constant.Success.GetResponseMessage(), response["message"])
				assert.NotNil(t, response["data"])

				data := response["data"].(map[string]interface{})
				assert.Equal(t, float64(1), data["id"])
				assert.Equal(t, "user1@example.com", data["email"])
			},
		},
		{
			name:           "Invalid user ID - non-numeric",
			userID:         "abc",
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid user ID - empty string",
			userID:         "",
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid user ID - special characters",
			userID:         "1@#",
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid user ID - decimal number",
			userID:         "1.5",
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Service returns database error",
			userID: "1",
			setupMock: func(m *MockUserService) {
				m.On("GetUserById", int64(1)).Return((*entity.User)(nil), errors.New("database connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "Service returns user not found error",
			userID: "999",
			setupMock: func(m *MockUserService) {
				m.On("GetUserById", int64(999)).Return((*entity.User)(nil), errors.New("user not found"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "Success - large user ID",
			userID: "9223372036854775807",
			setupMock: func(m *MockUserService) {
				user := &entity.User{
					Id:        9223372036854775807,
					Email:     "user@example.com",
					CreatedAt: time.Date(2025, 8, 1, 10, 0, 0, 0, time.UTC),
				}
				m.On("GetUserById", int64(9223372036854775807)).Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, constant.Success.GetResponseMessage(), response["message"])
				assert.NotNil(t, response["data"])
			},
		},
		{
			name:           "Invalid user ID - overflow",
			userID:         "9223372036854775808",
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Success - negative user ID",
			userID: "-1",
			setupMock: func(m *MockUserService) {
				user := &entity.User{
					Id:        -1,
					Email:     "user@example.com",
					CreatedAt: time.Date(2025, 8, 1, 10, 0, 0, 0, time.UTC),
				}
				m.On("GetUserById", int64(-1)).Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, constant.Success.GetResponseMessage(), response["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			userHandler := handler.NewUserHandler(mockService)

			tt.setupMock(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/api/users/"+tt.userID, nil)
			c.Params = gin.Params{
				{Key: "id", Value: tt.userID},
			}

			userHandler.GetUserById(c)
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
			mockService.AssertExpectations(t)
		})
	}
}
func TestUserHandler_CreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		email          string
		setupMock      func(*MockUserService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:  "Success - valid email",
			email: "user@example.com",
			setupMock: func(m *MockUserService) {
				user := &entity.User{
					Id:        1,
					Email:     "user@example.com",
					CreatedAt: time.Date(2025, 8, 1, 10, 0, 0, 0, time.UTC),
				}
				m.On("CreateUser", "user@example.com").Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, constant.Success.GetResponseMessage(), response["message"])
				assert.NotNil(t, response["data"])

				data := response["data"].(map[string]interface{})
				assert.Equal(t, float64(1), data["id"])
				assert.Equal(t, "user@example.com", data["email"])
			},
		},
		{
			name:  "Success - empty email",
			email: "",
			setupMock: func(m *MockUserService) {
				user := &entity.User{
					Id:        2,
					Email:     "",
					CreatedAt: time.Date(2025, 8, 1, 10, 0, 0, 0, time.UTC),
				}
				m.On("CreateUser", "").Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, constant.Success.GetResponseMessage(), response["message"])
				assert.NotNil(t, response["data"])
			},
		},
		{
			name:  "Success - email with special characters",
			email: "user+test@example.com",
			setupMock: func(m *MockUserService) {
				user := &entity.User{
					Id:        3,
					Email:     "user+test@example.com",
					CreatedAt: time.Date(2025, 8, 1, 10, 0, 0, 0, time.UTC),
				}
				m.On("CreateUser", "user+test@example.com").Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, constant.Success.GetResponseMessage(), response["message"])
				assert.NotNil(t, response["data"])

				data := response["data"].(map[string]interface{})
				assert.Equal(t, "user+test@example.com", data["email"])
			},
		},
		{
			name:  "Service returns duplicate email error",
			email: "duplicate@example.com",
			setupMock: func(m *MockUserService) {
				m.On("CreateUser", "duplicate@example.com").Return((*entity.User)(nil), errors.New("email already exists"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:  "Service returns database error",
			email: "user@example.com",
			setupMock: func(m *MockUserService) {
				m.On("CreateUser", "user@example.com").Return((*entity.User)(nil), errors.New("database connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:  "Service returns validation error",
			email: "invalid-email",
			setupMock: func(m *MockUserService) {
				m.On("CreateUser", "invalid-email").Return((*entity.User)(nil), errors.New("invalid email format"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:  "Success - very long email",
			email: "verylongusernamethatexceedsnormallimits@verylongdomainnamethatisunusuallylongbutvalid.com",
			setupMock: func(m *MockUserService) {
				user := &entity.User{
					Id:        4,
					Email:     "verylongusernamethatexceedsnormallimits@verylongdomainnamethatisunusuallylongbutvalid.com",
					CreatedAt: time.Date(2025, 8, 1, 10, 0, 0, 0, time.UTC),
				}
				m.On("CreateUser", "verylongusernamethatexceedsnormallimits@verylongdomainnamethatisunusuallylongbutvalid.com").Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, constant.Success.GetResponseMessage(), response["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			userHandler := handler.NewUserHandler(mockService)

			tt.setupMock(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/api/users", nil)
			c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			c.Request.PostForm = map[string][]string{
				"email": {tt.email},
			}

			userHandler.CreateUser(c)
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
			mockService.AssertExpectations(t)
		})
	}
}
func TestUserHandler_DeleteUserById(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		setupMock      func(*MockUserService)
		expectedStatus int
	}{
		{
			name:   "Success - valid user ID",
			userID: "1",
			setupMock: func(m *MockUserService) {
				m.On("DeleteUserById", int64(1)).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid user ID - non-numeric",
			userID:         "abc",
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid user ID - empty string",
			userID:         "",
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid user ID - special characters",
			userID:         "1@#",
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid user ID - decimal number",
			userID:         "1.5",
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Service returns database error",
			userID: "1",
			setupMock: func(m *MockUserService) {
				m.On("DeleteUserById", int64(1)).Return(errors.New("database connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "Service returns user not found error",
			userID: "999",
			setupMock: func(m *MockUserService) {
				m.On("DeleteUserById", int64(999)).Return(errors.New("user not found"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "Success - large user ID",
			userID: "9223372036854775807",
			setupMock: func(m *MockUserService) {
				m.On("DeleteUserById", int64(9223372036854775807)).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid user ID - overflow",
			userID:         "9223372036854775808",
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Success - negative user ID",
			userID: "-1",
			setupMock: func(m *MockUserService) {
				m.On("DeleteUserById", int64(-1)).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Service returns constraint violation error",
			userID: "1",
			setupMock: func(m *MockUserService) {
				m.On("DeleteUserById", int64(1)).Return(errors.New("foreign key constraint violation"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			userHandler := handler.NewUserHandler(mockService)

			tt.setupMock(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("DELETE", "/api/users/"+tt.userID, nil)
			c.Params = gin.Params{
				{Key: "id", Value: tt.userID},
			}

			userHandler.DeleteUserById(c)
			assert.Equal(t, tt.expectedStatus, w.Code)

			mockService.AssertExpectations(t)
		})
	}
}
func TestUserHandler_UpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		email          string
		setupMock      func(*MockUserService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "Success - valid user ID and email",
			userID: "1",
			email:  "updated@example.com",
			setupMock: func(m *MockUserService) {
				user := &entity.User{
					Id:        1,
					Email:     "updated@example.com",
					CreatedAt: time.Date(2025, 8, 1, 10, 0, 0, 0, time.UTC),
				}
				m.On("UpdateUser", int64(1), "updated@example.com").Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, constant.Success.GetResponseMessage(), response["message"])
				assert.NotNil(t, response["data"])

				data := response["data"].(map[string]interface{})
				assert.Equal(t, float64(1), data["id"])
				assert.Equal(t, "updated@example.com", data["email"])
			},
		},
		{
			name:           "Invalid user ID - non-numeric",
			userID:         "abc",
			email:          "user@example.com",
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid user ID - empty string",
			userID:         "",
			email:          "user@example.com",
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid user ID - special characters",
			userID:         "1@#",
			email:          "user@example.com",
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid user ID - decimal number",
			userID:         "1.5",
			email:          "user@example.com",
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Success - empty email",
			userID: "1",
			email:  "",
			setupMock: func(m *MockUserService) {
				user := &entity.User{
					Id:        1,
					Email:     "",
					CreatedAt: time.Date(2025, 8, 1, 10, 0, 0, 0, time.UTC),
				}
				m.On("UpdateUser", int64(1), "").Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, constant.Success.GetResponseMessage(), response["message"])
				assert.NotNil(t, response["data"])
			},
		},
		{
			name:   "Success - email with special characters",
			userID: "1",
			email:  "user+test@example.com",
			setupMock: func(m *MockUserService) {
				user := &entity.User{
					Id:        1,
					Email:     "user+test@example.com",
					CreatedAt: time.Date(2025, 8, 1, 10, 0, 0, 0, time.UTC),
				}
				m.On("UpdateUser", int64(1), "user+test@example.com").Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, constant.Success.GetResponseMessage(), response["message"])
				assert.NotNil(t, response["data"])

				data := response["data"].(map[string]interface{})
				assert.Equal(t, "user+test@example.com", data["email"])
			},
		},
		{
			name:   "Service returns user not found error",
			userID: "999",
			email:  "user@example.com",
			setupMock: func(m *MockUserService) {
				m.On("UpdateUser", int64(999), "user@example.com").Return((*entity.User)(nil), errors.New("user not found"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "Service returns database error",
			userID: "1",
			email:  "user@example.com",
			setupMock: func(m *MockUserService) {
				m.On("UpdateUser", int64(1), "user@example.com").Return((*entity.User)(nil), errors.New("database connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "Service returns duplicate email error",
			userID: "1",
			email:  "duplicate@example.com",
			setupMock: func(m *MockUserService) {
				m.On("UpdateUser", int64(1), "duplicate@example.com").Return((*entity.User)(nil), errors.New("email already exists"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "Success - large user ID",
			userID: "9223372036854775807",
			email:  "user@example.com",
			setupMock: func(m *MockUserService) {
				user := &entity.User{
					Id:        9223372036854775807,
					Email:     "user@example.com",
					CreatedAt: time.Date(2025, 8, 1, 10, 0, 0, 0, time.UTC),
				}
				m.On("UpdateUser", int64(9223372036854775807), "user@example.com").Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, constant.Success.GetResponseMessage(), response["message"])
				assert.NotNil(t, response["data"])
			},
		},
		{
			name:           "Invalid user ID - overflow",
			userID:         "9223372036854775808",
			email:          "user@example.com",
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Success - negative user ID",
			userID: "-1",
			email:  "user@example.com",
			setupMock: func(m *MockUserService) {
				user := &entity.User{
					Id:        -1,
					Email:     "user@example.com",
					CreatedAt: time.Date(2025, 8, 1, 10, 0, 0, 0, time.UTC),
				}
				m.On("UpdateUser", int64(-1), "user@example.com").Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, constant.Success.GetResponseMessage(), response["message"])
			},
		},
		{
			name:   "Service returns validation error",
			userID: "1",
			email:  "invalid-email",
			setupMock: func(m *MockUserService) {
				m.On("UpdateUser", int64(1), "invalid-email").Return((*entity.User)(nil), errors.New("invalid email format"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "Success - very long email",
			userID: "1",
			email:  "verylongusernamethatexceedsnormallimits@verylongdomainnamethatisunusuallylongbutvalid.com",
			setupMock: func(m *MockUserService) {
				user := &entity.User{
					Id:        1,
					Email:     "verylongusernamethatexceedsnormallimits@verylongdomainnamethatisunusuallylongbutvalid.com",
					CreatedAt: time.Date(2025, 8, 1, 10, 0, 0, 0, time.UTC),
				}
				m.On("UpdateUser", int64(1), "verylongusernamethatexceedsnormallimits@verylongdomainnamethatisunusuallylongbutvalid.com").Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, constant.Success.GetResponseMessage(), response["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			userHandler := handler.NewUserHandler(mockService)

			tt.setupMock(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("PUT", "/api/users/"+tt.userID, nil)
			c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			c.Request.PostForm = map[string][]string{
				"email": {tt.email},
			}
			c.Params = gin.Params{
				{Key: "id", Value: tt.userID},
			}

			userHandler.UpdateUser(c)
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
			mockService.AssertExpectations(t)
		})
	}
}
