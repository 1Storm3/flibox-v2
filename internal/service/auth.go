package service

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/1Storm3/flibox-api/internal/config"
	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/1Storm3/flibox-api/pkg/logger"
	"github.com/1Storm3/flibox-api/pkg/sys"
)

type AuthService struct {
	userService      controller.UserService
	emailService     controller.EmailService
	cfg              *config.Config
	tokenService     controller.TokenService
	takeHTMLTemplate func(appUrl, verificationToken string) (string, error)
}

func NewAuthService(
	userService controller.UserService,
	emailService controller.EmailService,
	cfg *config.Config,
	tokenService controller.TokenService,
	takeHtmlTemplate func(appUrl, verificationToken string) (string, error),
) *AuthService {
	return &AuthService{
		userService:      userService,
		emailService:     emailService,
		cfg:              cfg,
		tokenService:     tokenService,
		takeHTMLTemplate: takeHtmlTemplate,
	}
}

func (s *AuthService) Login(ctx context.Context, req model.User) (string, error) {
	user, err := s.userService.GetOneByEmail(ctx, req.Email)
	if err != nil || !s.userService.CheckPassword(ctx, &user, req.Password) {
		return "", sys.NewError(sys.ErrInvalidCredentials, "")
	}
	jwtKey := []byte(s.cfg.App.JwtSecretKey)
	expiresIn, err := time.ParseDuration(s.cfg.App.JwtExpiresIn)
	if err != nil {
		return "", sys.NewError(sys.ErrTokenGeneration, err.Error())
	}
	tokenString, err := s.tokenService.GenerateToken(jwtKey, user.ID, user.Role, expiresIn)
	if err != nil {
		return "", sys.NewError(sys.ErrTokenGeneration, err.Error())
	}
	return tokenString, nil
}

func (s *AuthService) Register(ctx context.Context, req model.User) (bool, error) {
	existingUser, err := s.userService.GetOneByEmail(ctx, req.Email)
	if err == nil && existingUser.ID != "" {
		return false, sys.NewError(sys.ErrUserAlreadyExists, "")
	}
	existingUser, err = s.userService.GetOneByNickName(ctx, req.NickName)
	if err == nil && existingUser.ID != "" {
		return false, sys.NewError(sys.ErrUserAlreadyExists, "")
	}
	hashedPassword, err := s.userService.HashPassword(ctx, req.Password)
	if err != nil {
		return false, sys.NewError(sys.ErrPasswordHashGeneration, err.Error())
	}

	jwtKey := []byte(s.cfg.App.JwtSecretKey)

	verificationToken, err := s.tokenService.GenerateEmailToken(req.Email, jwtKey, time.Hour*2)
	if err != nil {
		return false, sys.NewError(sys.ErrCreateToken, err.Error())
	}

	newUser := model.User{
		NickName:      req.NickName,
		Name:          req.Name,
		Email:         req.Email,
		Password:      hashedPassword,
		Photo:         req.Photo,
		Role:          "user",
		IsVerified:    false,
		VerifiedToken: verificationToken,
		LastActivity:  time.Now().UTC().Format("2006-01-02 15:04:05"),
		CreatedAt:     time.Now().UTC().Format("2006-01-02 15:04:05"),
		UpdatedAt:     time.Now().UTC().Format("2006-01-02 15:04:05"),
	}

	createdUser, err := s.userService.Create(ctx, newUser)
	if err != nil {
		return false, sys.NewError(sys.ErrCreateUser, err.Error())
	}

	go func() {
		emailBody, err := s.takeHTMLTemplate(s.cfg.App.AppUrl, *verificationToken)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		if err := s.emailService.SendEmail(
			createdUser.Email,
			emailBody,
			"Подтверждение регистрации",
		); err != nil {
			logger.Error("Ошибка при отправке email", zap.Error(err))
			return
		}
	}()

	return true, nil
}

func (s *AuthService) Verify(ctx context.Context, token string) error {
	jwtKey := []byte(s.cfg.App.JwtSecretKey)
	email, err := s.tokenService.ValidateEmailToken(token, jwtKey)
	if err != nil {
		return sys.NewError(sys.ErrInvalidToken, err.Error())
	}
	userNotVerified, err := s.userService.GetOneByEmail(ctx, email)
	if err != nil {
		return sys.NewError(sys.ErrUserNotFound, err.Error())
	}
	user := model.User{
		ID:            userNotVerified.ID,
		IsVerified:    true,
		VerifiedToken: nil,
	}

	if _, err := s.userService.UpdateForVerify(ctx, user); err != nil {
		return sys.NewError(sys.ErrUpdateUser, err.Error())
	}
	return nil
}

func (s *AuthService) Me(ctx context.Context, userId string) (model.User, error) {
	user, err := s.userService.GetOneById(ctx, userId)

	if err != nil {
		return model.User{}, sys.NewError(sys.ErrUserNotFound, err.Error())
	}

	lastActivity := time.Now().Format("2006-01-02 15:04:05")
	_, err = s.userService.Update(ctx, model.User{
		ID:           userId,
		LastActivity: lastActivity,
	})

	if err != nil {
		return model.User{}, sys.NewError(sys.ErrUpdateUser, err.Error())
	}

	return user, nil
}
