package http

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/mapper"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
)

type FilmController struct {
	filmService controller.FilmService
}

func NewFilmController(filmService controller.FilmService) *FilmController {
	return &FilmController{
		filmService: filmService,
	}
}

func (h *FilmController) Search(c *fiber.Ctx) error {
	match := c.Query("match")
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)
	genresStr := c.Query("genres")
	ctx := c.Context()
	var genres []string
	if genresStr != "" {
		genres = strings.Split(genresStr, ",")
	}
	films, totalRecords, err := h.filmService.Search(ctx, match, genres, page, pageSize)

	if err != nil {
		return httperror.HandleError(c, err)
	}
	totalPages := (totalRecords + int64(pageSize) - 1) / int64(pageSize)

	var filmsDTO []dto.FilmSearchResponseDTO
	for _, film := range films {
		filmsDTO = append(filmsDTO, mapper.MapModelFilmToResponseSearchDTO(film))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"films":        filmsDTO,
		"totalPages":   totalPages,
		"totalRecords": totalRecords,
		"currentPage":  page,
		"pageSize":     pageSize,
	})
}

func (h *FilmController) GetOneByID(c *fiber.Ctx) error {
	filmId := c.Params("id")
	ctx := c.Context()
	film, err := h.filmService.GetOne(ctx, filmId)

	if err != nil {
		return httperror.HandleError(c, err)
	}

	return c.JSON(mapper.MapModelFilmToResponseDTO(film))
}
