package middleware

import (
	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/service"
	"net/http"

	"github.com/1Storm3/flibox-api/internal/config"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(userRepo service.UserRepo, config *config.Config, tokenService controller.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		jwtKey := config.App.JwtSecretKey
		if tokenString == "" {
			return httperror.New(
				http.StatusUnauthorized,
				"Отсутствует токен",
			)
		}

		claims, err := tokenService.ParseToken(tokenString, []byte(jwtKey))
		if err != nil {
			return httperror.New(
				http.StatusUnauthorized,
				"Недействительный токен")
		}

		user, err := userRepo.GetOneById(c.Context(), claims.UserID)
		if err != nil {
			return httperror.New(
				http.StatusUnauthorized,
				"Ошибка получения информации о пользователе",
			)
		}
		if user.IsBlocked {
			return httperror.New(
				http.StatusForbidden,
				"Пользователь заблокирован")
		}
		if !user.IsVerified {
			return httperror.New(
				http.StatusForbidden,
				"Пользователь не верифицирован")
		}

		c.Locals("userClaims", claims)
		return c.Next()
	}
}
