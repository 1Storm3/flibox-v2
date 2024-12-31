package http

import (
	"github.com/1Storm3/flibox-api/internal/mapper"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
)

type AuthController struct {
	authService controller.AuthService
}

func NewAuthController(authService controller.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (a *AuthController) Login(c *fiber.Ctx) error {
	var loginData dto.LoginDTO

	ctx := c.Context()

	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message":    "Некорректные данные запроса",
			"statusCode": http.StatusBadRequest,
		})
	}

	user := mapper.MapLoginDTOToUserModel(loginData)

	tokenUser, err := a.authService.Login(ctx, user)

	if err != nil {
		return httperror.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"token": tokenUser,
	})
}

func (a *AuthController) Register(c *fiber.Ctx) error {
	var requestUser dto.RegisterDTO

	ctx := c.Context()
	if err := c.BodyParser(&requestUser); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message":    "Некорректные данные запроса",
			"statusCode": http.StatusBadRequest,
		})
	}

	user := mapper.MapRegisterDTOToUserModel(requestUser)

	result, err := a.authService.Register(ctx, user)
	if err != nil {
		return httperror.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"data": result,
	})
}

func (a *AuthController) Verify(c *fiber.Ctx) error {
	tokenUser := c.Params("token")
	if err := a.authService.Verify(c.Context(), tokenUser); err != nil {
		return httperror.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "Пользователь верифицирован",
	})
}

func (a *AuthController) Me(c *fiber.Ctx) error {
	claims, ok := c.Locals("userClaims").(*dto.Claims)

	ctx := c.Context()

	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message":    "Не удалось получить информацию о пользователе",
			"statusCode": http.StatusUnauthorized,
		})
	}
	result, err := a.authService.Me(ctx, claims.UserID)

	if err != nil {
		return httperror.HandleError(c, err)
	}

	return c.JSON(result)
}
