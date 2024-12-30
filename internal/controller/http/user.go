package http

import (
	"errors"
	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/mapper"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"net/http"
)

type UserController struct {
	userService controller.UserService
}

func NewUserController(userService controller.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (h *UserController) GetOneByNickName(c *fiber.Ctx) error {
	nickName := c.Params("nickName")

	ctx := c.Context()

	user, err := h.userService.GetOneByNickName(ctx, nickName)
	if err != nil {
		return httperror.HandleError(c, err)
	}
	return c.JSON(user)
}

func (h *UserController) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	ctx := c.Context()

	var userUpdateRequest dto.UpdateUserDTO
	if err := c.BodyParser(&userUpdateRequest); err != nil {
		return httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	userUpdateRequest.ID = id

	result, err := h.userService.Update(ctx, userUpdateRequest)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httperror.New(
				http.StatusNotFound,
				"Пользователь не найден",
			)
		}
		return err
	}
	updatedUser := mapper.MapModelUserToResponseDTO(result)

	return c.JSON(updatedUser)
}
