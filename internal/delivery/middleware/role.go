package middleware

import (
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/pkg/sys"
	"github.com/gofiber/fiber/v2"
)

func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userClaims, ok := c.Locals("userClaims").(*dto.Claims)
		if !ok || userClaims == nil {
			return sys.NewError(sys.ErrAccessDenied, "Недостаточно прав")
		}

		userRole := userClaims.Role

		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next()
			}
		}
		return sys.NewError(sys.ErrAccessDenied, "Недостаточно прав")
	}
}
