package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/mapper"
	"github.com/1Storm3/flibox-api/pkg/sys"
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
		return sys.HandleError(c, err)
	}
	var similarsDTO []dto.ResponseFilmDTO
	for _, similar := range similars {
		similarsDTO = append(similarsDTO, mapper.MapModelFilmToResponseDTO(mapper.MapModelFilmSimilarToModelFilm(similar)))
	}
	return c.JSON(similarsDTO)
}
