package service

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/pkg/sys"
)

type EmailClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type TokenService struct {
}

func NewTokenService() *TokenService {
	return &TokenService{}
}

func (s *TokenService) GenerateToken(jwtKey []byte, userID, role string, duration time.Duration) (string, error) {
	expirationTime := time.Now().Add(duration)
	claims := &dto.Claims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *TokenService) GenerateEmailToken(email string, jwtKey []byte, duration time.Duration) (*string, error) {
	expirationTime := time.Now().Add(duration)
	claims := &EmailClaims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil,
			sys.NewError(sys.ErrUnknown, err.Error())
	}
	return &tokenString, nil
}

func (s *TokenService) ValidateEmailToken(tokenString string, jwtKey []byte) (string, error) {
	claims := &EmailClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(_ *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return "", sys.NewError(sys.ErrInvalidToken, "")
	}
	return claims.Email, nil
}

func (s *TokenService) ParseToken(tokenString string, jwtKey []byte) (*dto.Claims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	claims := &dto.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(_ *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, sys.NewError(sys.ErrInvalidToken, "")
	}
	return claims, nil
}
