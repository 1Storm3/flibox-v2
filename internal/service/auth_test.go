package service

import (
	"context"
	"errors"
	"github.com/1Storm3/flibox-api/internal/config"
	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/1Storm3/flibox-api/mocks"
	"github.com/1Storm3/flibox-api/pkg/sys"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_Login(t *testing.T) {
	tests := []struct {
		name           string
		email          string
		password       string
		mockUser       model.User
		mockGetErr     error
		mockCheckPass  bool
		mockToken      string
		mockTokenErr   error
		expectedToken  string
		expectedError  error
		expectedStatus int
	}{
		{
			name:           "Valid credentials",
			email:          "test@example.com",
			password:       "password",
			mockUser:       model.User{ID: "1", Role: "user"},
			mockGetErr:     nil,
			mockCheckPass:  true,
			mockToken:      "valid_token",
			mockTokenErr:   nil,
			expectedToken:  "valid_token",
			expectedError:  nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid password",
			email:          "test@example.com",
			password:       "wrong_password",
			mockUser:       model.User{ID: "1", Role: "user"},
			mockGetErr:     nil,
			mockCheckPass:  false,
			mockToken:      "",
			mockTokenErr:   nil,
			expectedToken:  "",
			expectedError:  sys.NewError(sys.ErrInvalidCredentials, ""),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "User not found",
			email:          "notfound@example.com",
			password:       "password",
			mockUser:       model.User{},
			mockGetErr:     errors.New("user not found"),
			mockCheckPass:  false,
			mockToken:      "",
			mockTokenErr:   nil,
			expectedToken:  "",
			expectedError:  sys.NewError(sys.ErrInvalidCredentials, ""),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Token generation failure",
			email:          "test@example.com",
			password:       "password",
			mockUser:       model.User{ID: "1", Role: "user"},
			mockGetErr:     nil,
			mockCheckPass:  true,
			mockToken:      "",
			mockTokenErr:   errors.New("token generation failed"),
			expectedToken:  "",
			expectedError:  sys.NewError(sys.ErrTokenGeneration, "token generation failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserService := new(mocks.MockUserService)
			mockTokenService := new(mocks.MockTokenService)
			mockUserService.On("GetOneByEmail", mock.Anything, tt.email).Return(tt.mockUser, tt.mockGetErr)
			if tt.mockGetErr == nil {
				mockUserService.On("CheckPassword", mock.Anything, &tt.mockUser, tt.password).Return(tt.mockCheckPass)
			}
			if tt.mockCheckPass {
				mockTokenService.On("GenerateToken", mock.Anything, tt.mockUser.ID, tt.mockUser.Role, mock.Anything).Return(tt.mockToken, tt.mockTokenErr)
			}

			authService := AuthService{
				userService:  mockUserService,
				tokenService: mockTokenService,
				cfg: &config.Config{
					App: config.AppConfig{
						JwtSecretKey: "secret",
						JwtExpiresIn: "1h",
					},
				},
			}

			token, err := authService.Login(context.Background(), model.User{Email: tt.email, Password: tt.password})

			assert.Equal(t, tt.expectedToken, token)
			assert.Equal(t, tt.expectedError, err)

			mockUserService.AssertExpectations(t)
			mockTokenService.AssertExpectations(t)
		})
	}
}
