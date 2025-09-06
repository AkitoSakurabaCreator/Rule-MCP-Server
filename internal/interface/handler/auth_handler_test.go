package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(id int) (*domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(username string) (*domain.User, error) {
	args := m.Called(username)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetAll() ([]domain.User, error) {
	args := m.Called()
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) GetActiveUsers() ([]domain.User, error) {
	args := m.Called()
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUsersByRole(role string) ([]domain.User, error) {
	args := m.Called(role)
	return args.Get(0).([]domain.User), args.Error(1)
}

// MockRoleRepository is a mock implementation of RoleRepository
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) GetAll() ([]domain.Role, error) {
	args := m.Called()
	return args.Get(0).([]domain.Role), args.Error(1)
}

func (m *MockRoleRepository) GetByName(name string) (domain.Role, error) {
	args := m.Called(name)
	return args.Get(0).(domain.Role), args.Error(1)
}

func (m *MockRoleRepository) Create(role domain.Role) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *MockRoleRepository) Update(name string, role domain.Role) error {
	args := m.Called(name, role)
	return args.Error(0)
}

func (m *MockRoleRepository) Delete(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    LoginRequest
		mockSetup      func(*MockUserRepository)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Successful login",
			requestBody: LoginRequest{
				Username: "admin",
				Password: "admin123",
			},
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetByUsername", "admin").Return(&domain.User{
					ID:           1,
					Username:     "admin",
					Email:        "admin@example.com",
					Role:         "admin",
					IsActive:     true,
					PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // bcrypt hash for "admin123"
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "Invalid credentials",
			requestBody: LoginRequest{
				Username: "admin",
				Password: "wrongpassword",
			},
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetByUsername", "admin").Return(&domain.User{
					ID:           1,
					Username:     "admin",
					Email:        "admin@example.com",
					Role:         "admin",
					IsActive:     true,
					PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // bcrypt hash for "admin123"
				}, nil)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name: "User not found",
			requestBody: LoginRequest{
				Username: "nonexistent",
				Password: "password",
			},
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetByUsername", "nonexistent").Return((*domain.User)(nil), assert.AnError)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := new(MockUserRepository)
			mockRoleRepo := new(MockRoleRepository)
			tt.mockSetup(mockUserRepo)

			handler := NewAuthHandler("test-secret", mockUserRepo, mockRoleRepo)

			router := gin.New()
			router.POST("/login", handler.Login)

			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError {
				var errorResponse map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.Contains(t, errorResponse, "code")
			} else {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response, "token")
				assert.Contains(t, response, "user")
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockSetup      func(*MockUserRepository)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Successful registration",
			requestBody: map[string]interface{}{
				"username":  "newuser",
				"email":     "newuser@example.com",
				"full_name": "New User",
				"password":  "password123",
			},
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetByUsername", "newuser").Return((*domain.User)(nil), assert.AnError)
				mockRepo.On("GetByEmail", "newuser@example.com").Return((*domain.User)(nil), assert.AnError)
				mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "Username already exists",
			requestBody: map[string]interface{}{
				"username":  "existinguser",
				"email":     "newuser@example.com",
				"full_name": "New User",
				"password":  "password123",
			},
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetByUsername", "existinguser").Return(&domain.User{
					Username: "existinguser",
				}, nil)
			},
			expectedStatus: http.StatusConflict,
			expectedError:  true,
		},
		{
			name: "Email already exists",
			requestBody: map[string]interface{}{
				"username":  "newuser",
				"email":     "existing@example.com",
				"full_name": "New User",
				"password":  "password123",
			},
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetByUsername", "newuser").Return((*domain.User)(nil), assert.AnError)
				mockRepo.On("GetByEmail", "existing@example.com").Return(&domain.User{
					Email: "existing@example.com",
				}, nil)
			},
			expectedStatus: http.StatusConflict,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := new(MockUserRepository)
			mockRoleRepo := new(MockRoleRepository)
			tt.mockSetup(mockUserRepo)

			handler := NewAuthHandler("test-secret", mockUserRepo, mockRoleRepo)

			router := gin.New()
			router.POST("/register", handler.Register)

			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError {
				var errorResponse map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.Contains(t, errorResponse, "code")
			} else {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response, "token")
				assert.Contains(t, response, "user")
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Me(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		authHeader     string
		mockSetup      func(*MockUserRepository)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:       "Valid token",
			authHeader: "Bearer valid-token",
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetByID", 1).Return(&domain.User{
					ID:       1,
					Username: "admin",
					Email:    "admin@example.com",
					Role:     "admin",
					IsActive: true,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "No token",
			authHeader:     "",
			mockSetup:      func(mockRepo *MockUserRepository) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "Invalid token",
			authHeader:     "Bearer invalid-token",
			mockSetup:      func(mockRepo *MockUserRepository) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := new(MockUserRepository)
			mockRoleRepo := new(MockRoleRepository)
			tt.mockSetup(mockUserRepo)

			handler := NewAuthHandler("test-secret", mockUserRepo, mockRoleRepo)

			router := gin.New()
			router.GET("/me", handler.Me)

			req, _ := http.NewRequest("GET", "/me", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError {
				var errorResponse map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.Contains(t, errorResponse, "code")
			} else {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response, "username")
				assert.Contains(t, response, "email")
				assert.Contains(t, response, "role")
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}
