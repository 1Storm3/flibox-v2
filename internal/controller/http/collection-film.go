package http

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
)

type CollectionFilmController struct {
	collectionFilmService controller.CollectionFilmService
}

func NewCollectionFilmController(collectionFilmService controller.CollectionFilmService) *CollectionFilmController {
	return &CollectionFilmController{
		collectionFilmService: collectionFilmService,
	}
}

func (h *CollectionFilmController) Add(c *fiber.Ctx) error {
	var filmRequest dto.CreateCollectionFilmDTO
	collectionId := c.Params("id")
	if err := c.BodyParser(&filmRequest); err != nil {
		return httperror.New(
			http.StatusBadRequest,
			err.Error(),
		)
	}
	ctx := c.Context()
	err := h.collectionFilmService.Add(ctx, collectionId, strconv.Itoa(filmRequest.FilmID))
	if err != nil {
		return httperror.HandleError(c, err)
	}
	return c.JSON(fiber.Map{
		"data": "Фильм добавлен в подборку",
	})
}

func (h *CollectionFilmController) Delete(c *fiber.Ctx) error {
	var filmRequest dto.DeleteCollectionFilmDTO
	collectionId := c.Params("id")
	if err := c.BodyParser(&filmRequest); err != nil {
		return httperror.New(
			http.StatusBadRequest,
			err.Error(),
		)
	}
	ctx := c.Context()
	err := h.collectionFilmService.Delete(ctx, collectionId, strconv.Itoa(filmRequest.FilmID))
	if err != nil {
		return httperror.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"data": "Фильмы удалены из подборки",
	})
}

func (h *CollectionFilmController) GetFilmsByCollectionId(c *fiber.Ctx) error {
	collectionID := c.Params("id")
	if collectionID == "" {
		return httperror.New(
			http.StatusBadRequest,
			"Неверный формат ID коллекции",
		)
	}

	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 20)

	ctx := c.Context()
	films, totalRecords, err := h.collectionFilmService.GetFilmsByCollectionId(ctx, collectionID, page, pageSize)
	if err != nil {
		return httperror.HandleError(c, err)
	}

	totalPages := (totalRecords + int64(pageSize) - 1) / int64(pageSize)

	return c.JSON(fiber.Map{
		"data":         films,
		"totalPages":   totalPages,
		"totalRecords": totalRecords,
		"currentPage":  page,
		"pageSize":     pageSize,
	})
}
