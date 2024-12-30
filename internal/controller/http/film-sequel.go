package http

import (
	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"github.com/gofiber/fiber/v2"
)

type FilmSequelController struct {
	filmSequelService controller.FilmSequelService
}

func NewFilmSequelController(filmSequelService controller.FilmSequelService) *FilmSequelController {
	return &FilmSequelController{
		filmSequelService: filmSequelService,
	}
}

func (h *FilmSequelController) GetAll(c *fiber.Ctx) error {
	filmId := c.Params("id")
	ctx := c.Context()
	sequels, err := h.filmSequelService.GetAll(ctx, filmId)
	if err != nil {
		return httperror.HandleError(c, err)
	}
	return c.JSON(sequels)
}
