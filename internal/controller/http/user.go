package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/mapper"
	"github.com/1Storm3/flibox-api/pkg/sys"
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
		return sys.HandleError(c, err)
	}
	return c.JSON(mapper.MapUserModelToUserResponseDto(user))
}

func (h *UserController) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	ctx := c.Context()

	var userUpdateRequest dto.UpdateUserDTO
	if err := c.BodyParser(&userUpdateRequest); err != nil {
		return sys.NewError(sys.ErrInvalidRequestData, err.Error())
	}

	userUpdateRequest.ID = id

	userUpdate := mapper.MapUpdateUserDTOToUserModel(userUpdateRequest)

	result, err := h.userService.Update(ctx, userUpdate)

	if err != nil {
		return sys.HandleError(c, err)
	}
	updatedUser := mapper.MapModelUserToResponseDTO(result)

	return c.JSON(updatedUser)
}
