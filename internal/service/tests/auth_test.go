package tests

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/1Storm3/flibox-api/internal/config"
	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/1Storm3/flibox-api/internal/service"
	"github.com/1Storm3/flibox-api/mocks"
	"github.com/1Storm3/flibox-api/pkg/sys"
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
		config         config.AppConfig
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
			config: config.AppConfig{
				JwtSecretKey: "secret",
				JwtExpiresIn: "1h",
			},
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
			config: config.AppConfig{
				JwtSecretKey: "secret",
				JwtExpiresIn: "1h",
			},
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
			config: config.AppConfig{
				JwtSecretKey: "secret",
				JwtExpiresIn: "1h",
			},
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
			config: config.AppConfig{
				JwtSecretKey: "secret",
				JwtExpiresIn: "1h",
			},
		},
		{
			name:           "Invalid JwtExpiresIn format",
			email:          "test@example.com",
			password:       "password",
			mockUser:       model.User{ID: "1", Role: "user"},
			mockGetErr:     nil,
			mockCheckPass:  true,
			mockToken:      "",
			mockTokenErr:   nil,
			expectedToken:  "",
			expectedError:  sys.NewError(sys.ErrTokenGeneration, "time: invalid duration \"invalid_duration\""),
			expectedStatus: http.StatusInternalServerError,
			config: config.AppConfig{
				JwtSecretKey: "secret",
				JwtExpiresIn: "invalid_duration",
			},
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
			if tt.mockCheckPass && tt.config.JwtExpiresIn != "invalid_duration" {
				mockTokenService.On("GenerateToken", mock.Anything, tt.mockUser.ID, tt.mockUser.Role, mock.Anything).Return(tt.mockToken, tt.mockTokenErr)
			}

			authService := service.NewAuthService(mockUserService, nil, &config.Config{
				App: tt.config,
			},
				mockTokenService,
				nil,
			)

			token, err := authService.Login(context.Background(), model.User{Email: tt.email, Password: tt.password})

			assert.Equal(t, tt.expectedToken, token)
			assert.Equal(t, tt.expectedError, err)

			mockUserService.AssertExpectations(t)
			mockTokenService.AssertExpectations(t)
		})
	}
}

func TestAuthService_Me(t *testing.T) {
	tests := []struct {
		name           string
		userId         string
		mockUser       model.User
		mockGetErr     error
		expectedUser   model.User
		expectedError  error
		mockUpdateErr  error
		expectedStatus int
	}{
		{
			name:           "Success",
			userId:         "1",
			mockUser:       model.User{ID: "1", LastActivity: time.Now().Format("2006-01-02 15:04:05")},
			mockGetErr:     nil,
			expectedUser:   model.User{ID: "1", LastActivity: time.Now().Format("2006-01-02 15:04:05")},
			expectedError:  nil,
			mockUpdateErr:  nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "User not found",
			userId:         "1",
			mockUser:       model.User{},
			mockGetErr:     errors.New("user not found"),
			expectedUser:   model.User{},
			mockUpdateErr:  nil,
			expectedError:  sys.NewError(sys.ErrUserNotFound, "user not found"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Update user failure",
			userId:         "1",
			mockUser:       model.User{ID: "1", LastActivity: time.Now().Format("2006-01-02 15:04:05")},
			mockGetErr:     nil,
			expectedUser:   model.User{},
			expectedError:  sys.NewError(sys.ErrUpdateUser, "update user failed"),
			mockUpdateErr:  errors.New("update user failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserService := new(mocks.MockUserService)
			mockUserService.On("GetOneById", mock.Anything, tt.userId).Return(tt.mockUser, tt.mockGetErr)
			if tt.mockGetErr == nil {
				mockUserService.On("Update", mock.Anything, tt.mockUser).Return(tt.mockUser, tt.mockUpdateErr)
			}

			authService := service.NewAuthService(mockUserService, nil, nil, nil, nil)

			user, err := authService.Me(context.Background(), tt.userId)

			assert.Equal(t, tt.expectedUser, user)
			assert.Equal(t, tt.expectedError, err)

			mockUserService.AssertExpectations(t)
		})
	}
}

func TestAuthService_Verify(t *testing.T) {
	tests := []struct {
		name            string
		token           string
		email           string
		mockValidateErr error
		mockGetErr      error
		mockUpdateErr   error
		expectedError   error
		expectedStatus  int
		config          config.AppConfig
	}{
		{
			name:            "Success",
			token:           "token",
			email:           "email",
			mockValidateErr: nil,
			mockUpdateErr:   nil,
			mockGetErr:      nil,
			expectedError:   nil,
			expectedStatus:  http.StatusOK,
			config: config.AppConfig{
				JwtSecretKey: "secret",
			},
		},
		{
			name:            "Invalid token",
			token:           "token",
			email:           "email",
			mockGetErr:      nil,
			mockUpdateErr:   nil,
			mockValidateErr: errors.New("invalid token"),
			expectedError:   sys.NewError(sys.ErrInvalidToken, "invalid token"),
			expectedStatus:  http.StatusUnauthorized,
			config: config.AppConfig{
				JwtSecretKey: "secret",
			},
		},
		{
			name:            "getting user failed",
			token:           "token",
			email:           "email",
			mockGetErr:      errors.New("getting user failed"),
			mockValidateErr: nil,
			mockUpdateErr:   nil,
			expectedError:   sys.NewError(sys.ErrUserNotFound, "getting user failed"),
			expectedStatus:  http.StatusInternalServerError,
			config: config.AppConfig{
				JwtSecretKey: "secret",
			},
		},
		{
			name:            "update user failed",
			token:           "token",
			email:           "email",
			mockGetErr:      nil,
			mockValidateErr: nil,
			mockUpdateErr:   errors.New("update user failed"),
			expectedError:   sys.NewError(sys.ErrUpdateUser, "update user failed"),
			expectedStatus:  http.StatusInternalServerError,
			config: config.AppConfig{
				JwtSecretKey: "secret",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserService := new(mocks.MockUserService)
			mockTokenService := new(mocks.MockTokenService)
			mockTokenService.On("ValidateEmailToken", tt.token, mock.Anything).Return(tt.email, tt.mockValidateErr)
			if tt.mockValidateErr == nil {
				mockUserService.On("GetOneByEmail", mock.Anything, tt.email).Return(model.User{
					IsVerified: false,
				}, tt.mockGetErr)
			}
			if tt.mockValidateErr == nil && tt.mockGetErr == nil {
				mockUserService.On("UpdateForVerify", mock.Anything, model.User{
					IsVerified: true,
				}).Return(model.User{
					IsVerified: true,
				}, tt.mockUpdateErr)
			}

			authService := service.NewAuthService(mockUserService, nil, &config.Config{
				App: tt.config,
			}, mockTokenService, nil)

			err := authService.Verify(context.Background(), tt.token)

			assert.Equal(t, tt.expectedError, err)

			mockUserService.AssertExpectations(t)
		})
	}
}

func TestAuthService_Register(t *testing.T) {
	type mockBehavior func(userService *mocks.MockUserService, tokenService *mocks.MockTokenService, emailService *mocks.MockEmailService)
	token := "token"
	var wg sync.WaitGroup
	tests := []struct {
		name          string
		input         model.User
		mockBehavior  mockBehavior
		expectedError error
		expectedOK    bool
		mockTakeHTML  func(appUrl, verificationToken string) (string, error)
	}{
		{
			name: "User with email already exists",
			input: model.User{
				Email:    "existing@example.com",
				NickName: "nickname",
				Password: "password123",
			},
			mockBehavior: func(userService *mocks.MockUserService, tokenService *mocks.MockTokenService, emailService *mocks.MockEmailService) {
				userService.On("GetOneByEmail", mock.Anything, "existing@example.com").Return(model.User{ID: "123"}, nil)
			},
			mockTakeHTML:  nil,
			expectedError: sys.NewError(sys.ErrUserAlreadyExists, ""),
			expectedOK:    false,
		},
		{
			name: "User with nickname already exists",
			input: model.User{
				Email:    "new@example.com",
				NickName: "existingNickname",
				Password: "password123",
			},
			mockBehavior: func(userService *mocks.MockUserService, tokenService *mocks.MockTokenService, emailService *mocks.MockEmailService) {
				userService.On("GetOneByEmail", mock.Anything, "new@example.com").Return(model.User{}, errors.New("not found"))
				userService.On("GetOneByNickName", mock.Anything, "existingNickname").Return(model.User{ID: "123"}, nil)
			},
			mockTakeHTML:  nil,
			expectedError: sys.NewError(sys.ErrUserAlreadyExists, ""),
			expectedOK:    false,
		},
		{
			name: "Error hashing password",
			input: model.User{
				Email:    "new@example.com",
				NickName: "newNickname",
				Password: "password123",
			},
			mockBehavior: func(userService *mocks.MockUserService, tokenService *mocks.MockTokenService, emailService *mocks.MockEmailService) {
				userService.On("GetOneByEmail", mock.Anything, "new@example.com").Return(model.User{}, errors.New("not found"))
				userService.On("GetOneByNickName", mock.Anything, "newNickname").Return(model.User{}, errors.New("not found"))
				userService.On("HashPassword", mock.Anything, "password123").Return("", errors.New("hash error"))
			},
			mockTakeHTML:  nil,
			expectedError: sys.NewError(sys.ErrPasswordHashGeneration, "hash error"),
			expectedOK:    false,
		},
		{
			name: "Error generating verification token",
			input: model.User{
				Email:    "new@example.com",
				NickName: "newNickname",
				Password: "password123",
			},
			mockBehavior: func(userService *mocks.MockUserService, tokenService *mocks.MockTokenService, emailService *mocks.MockEmailService) {
				userService.On("GetOneByEmail", mock.Anything, "new@example.com").Return(model.User{}, errors.New("not found"))
				userService.On("GetOneByNickName", mock.Anything, "newNickname").Return(model.User{}, errors.New("not found"))
				userService.On("HashPassword", mock.Anything, "password123").Return("hashedPassword", nil)
				tokenService.On("GenerateEmailToken", "new@example.com", mock.Anything, mock.Anything).Return(nil, errors.New("token error"))
			},
			mockTakeHTML:  nil,
			expectedError: sys.NewError(sys.ErrCreateToken, "token error"),
			expectedOK:    false,
		},
		{
			name: "Create user error",
			input: model.User{
				Email:    "new@example.com",
				NickName: "newNickname",
				Password: "password123",
			},
			mockBehavior: func(userService *mocks.MockUserService, tokenService *mocks.MockTokenService, emailService *mocks.MockEmailService) {
				userService.On("GetOneByEmail", mock.Anything, "new@example.com").Return(model.User{}, errors.New("not found"))
				userService.On("GetOneByNickName", mock.Anything, "newNickname").Return(model.User{}, errors.New("not found"))
				userService.On("HashPassword", mock.Anything, "password123").Return("hashedPassword", nil)
				tokenService.On("GenerateEmailToken", "new@example.com", mock.Anything, mock.Anything).Return(&token, nil)
				userService.On("Create", mock.Anything, mock.Anything).Return(model.User{ID: "123"}, errors.New("create error"))
			},
			mockTakeHTML:  nil,
			expectedError: sys.NewError(sys.ErrCreateUser, "create error"),
			expectedOK:    false,
		},
		{
			name: "Successful registration with error on take html template",
			input: model.User{
				Email:    "new@example.com",
				NickName: "newNickname",
				Password: "password123",
			},
			mockBehavior: func(userService *mocks.MockUserService, tokenService *mocks.MockTokenService, emailService *mocks.MockEmailService) {
				userService.On("GetOneByEmail", mock.Anything, "new@example.com").Return(model.User{}, errors.New("not found"))
				userService.On("GetOneByNickName", mock.Anything, "newNickname").Return(model.User{}, errors.New("not found"))
				userService.On("HashPassword", mock.Anything, "password123").Return("hashedPassword", nil)
				tokenService.On("GenerateEmailToken", "new@example.com", mock.Anything, mock.Anything).Return(&token, nil)
				userService.On("Create", mock.Anything, mock.Anything).Return(model.User{ID: "123"}, nil)
			},
			mockTakeHTML: func(appUrl, verificationToken string) (string, error) {
				return "", errors.New("take html error")
			},
			expectedError: nil,
			expectedOK:    true,
		},
		{
			name: "Successful registration with error on sending email",
			input: model.User{
				Email:    "new@example.com",
				NickName: "newNickname",
				Password: "password123",
			},
			mockBehavior: func(userService *mocks.MockUserService, tokenService *mocks.MockTokenService, emailService *mocks.MockEmailService) {
				userService.On("GetOneByEmail", mock.Anything, "new@example.com").Return(model.User{}, errors.New("not found"))
				userService.On("GetOneByNickName", mock.Anything, "newNickname").Return(model.User{}, errors.New("not found"))
				userService.On("HashPassword", mock.Anything, "password123").Return("hashedPassword", nil)
				tokenService.On("GenerateEmailToken", "new@example.com", mock.Anything, mock.Anything).Return(&token, nil)
				userService.On("Create", mock.Anything, mock.Anything).Return(model.User{ID: "123"}, nil)
				emailService.On("SendEmail", "", mock.Anything, "Подтверждение регистрации").Run(func(args mock.Arguments) {
					defer wg.Done()
					fmt.Println("Горутина завершила выполнение")
				}).Return(errors.New("send email error"))
			},
			mockTakeHTML: func(appUrl, verificationToken string) (string, error) {
				return "<html></html>", nil
			},
			expectedError: nil,
			expectedOK:    true,
		},
		{
			name: "Successful registration",
			input: model.User{
				Email:    "new@example.com",
				NickName: "newNickname",
				Password: "password123",
			},
			mockBehavior: func(userService *mocks.MockUserService, tokenService *mocks.MockTokenService, emailService *mocks.MockEmailService) {
				userService.On("GetOneByEmail", mock.Anything, "new@example.com").Return(model.User{}, errors.New("not found"))
				userService.On("GetOneByNickName", mock.Anything, "newNickname").Return(model.User{}, errors.New("not found"))
				userService.On("HashPassword", mock.Anything, "password123").Return("hashedPassword", nil)
				tokenService.On("GenerateEmailToken", "new@example.com", mock.Anything, mock.Anything).Return(&token, nil)
				userService.On("Create", mock.Anything, mock.Anything).Return(model.User{ID: "123"}, nil)
				emailService.On("SendEmail", "", mock.Anything, "Подтверждение регистрации").Run(func(args mock.Arguments) {
					defer wg.Done()
					fmt.Println("Горутина завершила выполнение")
				}).Return(nil)
			},
			mockTakeHTML: func(appUrl, verificationToken string) (string, error) {
				return "<html></html>", nil
			},
			expectedError: nil,
			expectedOK:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserService := new(mocks.MockUserService)
			mockTokenService := new(mocks.MockTokenService)
			mockEmailService := new(mocks.MockEmailService)

			tt.mockBehavior(mockUserService, mockTokenService, mockEmailService)

			authService := service.NewAuthService(
				mockUserService,
				mockEmailService,
				&config.Config{
					App: config.AppConfig{
						JwtSecretKey: "secret",
						AppUrl:       "http://example.com",
					},
				}, mockTokenService,
				tt.mockTakeHTML,
			)

			if tt.name == "Successful registration with error on sending email" || tt.name == "Successful registration" {
				wg.Add(1)
				ok, err := authService.Register(context.Background(), tt.input)
				wg.Wait()
				assert.Equal(t, tt.expectedOK, ok)
				assert.Equal(t, tt.expectedError, err)
			} else {
				ok, err := authService.Register(context.Background(), tt.input)
				assert.Equal(t, tt.expectedOK, ok)
				assert.Equal(t, tt.expectedError, err)
			}

			mockUserService.AssertExpectations(t)
			mockTokenService.AssertExpectations(t)
			mockEmailService.AssertExpectations(t)
		})
	}
}
