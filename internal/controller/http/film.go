package http

import (
	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"github.com/gofiber/fiber/v2"
	"strings"
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"films":        films,
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

	return c.JSON(film)
}
