package middleware

import (
	"github.com/gofiber/fiber/v2"

	"github.com/1Storm3/flibox-api/internal/config"
	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/service"
	"github.com/1Storm3/flibox-api/pkg/sys"
)

func AuthMiddleware(userRepo service.UserRepo, config *config.Config, tokenService controller.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		jwtKey := config.App.JwtSecretKey
		if tokenString == "" {
			return sys.NewError(sys.ErrInvalidToken, "")
		}

		claims, err := tokenService.ParseToken(tokenString, []byte(jwtKey))
		if err != nil {
			return sys.NewError(sys.ErrInvalidToken, "")
		}

		user, err := userRepo.GetOneById(c.Context(), claims.UserID)
		if err != nil {
			return sys.NewError(sys.ErrUserNotFound, err.Error())
		}
		if user.IsBlocked {
			return sys.NewError(sys.ErrUserBlocked, "")
		}
		if !user.IsVerified {
			return sys.NewError(sys.ErrUserNotVerified, "")
		}

		c.Locals("userClaims", claims)
		return c.Next()
	}
}
