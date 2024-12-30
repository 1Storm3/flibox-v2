package middleware

import (
	"github.com/1Storm3/flibox-api/internal/dto"
	"net/http"

	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"github.com/gofiber/fiber/v2"
)

func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userClaims := c.Locals("userClaims").(*dto.Claims)
		userRole := userClaims.Role

		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next()
			}
		}
		return httperror.New(
			http.StatusForbidden,
			"Недостаточно прав для выполнения операции",
		)
	}
}
