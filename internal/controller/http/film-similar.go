package http

import (
	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"github.com/gofiber/fiber/v2"
)

type FilmSimilarController struct {
	service controller.FilmSimilarService
}

func NewFilmSimilarController(service controller.FilmSimilarService) *FilmSimilarController {
	return &FilmSimilarController{
		service: service,
	}
}

func (h *FilmSimilarController) GetAll(c *fiber.Ctx) error {
	filmId := c.Params("id")

	ctx := c.Context()

	similars, err := h.service.GetAll(ctx, filmId)
	if err != nil {
		return httperror.HandleError(c, err)
	}
	return c.JSON(similars)
}
