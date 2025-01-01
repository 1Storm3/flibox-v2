package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/mapper"
	"github.com/1Storm3/flibox-api/pkg/sys"
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
		return sys.HandleError(c, err)
	}
	var sequelsDTO []dto.ResponseFilmDTO
	for _, sequel := range sequels {
		sequelsDTO = append(sequelsDTO, mapper.MapModelFilmToResponseDTO(mapper.MapModelFilmSequelToModelFilm(sequel)))
	}
	return c.JSON(sequelsDTO)
}
