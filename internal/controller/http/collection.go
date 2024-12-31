package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/mapper"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
)

type CollectionController struct {
	collectionService controller.CollectionService
}

func NewCollectionController(collectionService controller.CollectionService) *CollectionController {
	return &CollectionController{
		collectionService: collectionService,
	}
}

func (h *CollectionController) Update(c *fiber.Ctx) error {
	collectionId := c.Params("id")
	ctx := c.Context()
	var collectionDto dto.UpdateCollectionDTO
	if err := c.BodyParser(&collectionDto); err != nil {
		return httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	collectionDto.ID = collectionId

	collection := mapper.MapUpdateCollectionDTOToCollectionModel(collectionDto)

	result, err := h.collectionService.Update(ctx, collection)

	if err != nil {
		return httperror.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"data": mapper.MapModelCollectionToResponseDTO(result),
	})
}

func (h *CollectionController) Delete(c *fiber.Ctx) error {
	collectionId := c.Params("id")

	ctx := c.Context()

	err := h.collectionService.Delete(ctx, collectionId)

	if err != nil {
		return httperror.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"data": "Коллекция удалена",
	})
}

func (h *CollectionController) GetAllMy(c *fiber.Ctx) error {
	userId := c.Locals("userClaims").(*dto.Claims).UserID
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)
	ctx := c.Context()
	result, totalRecords, err := h.collectionService.GetAllMy(ctx, page, pageSize, userId)
	if err != nil {
		return httperror.HandleError(c, err)
	}
	totalPages := (totalRecords + int64(pageSize) - 1) / int64(pageSize)

	var resultDTO []dto.ResponseCollectionDTO

	for _, result := range result {
		resultDTO = append(resultDTO, mapper.MapModelCollectionToResponseDTO(result))
	}

	return c.JSON(fiber.Map{
		"data":         resultDTO,
		"totalPages":   totalPages,
		"totalRecords": totalRecords,
		"currentPage":  page,
		"pageSize":     pageSize,
	})
}

func (h *CollectionController) GetAll(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)
	ctx := c.Context()
	result, totalRecords, err := h.collectionService.GetAll(ctx, page, pageSize)
	if err != nil {
		return httperror.HandleError(c, err)
	}
	totalPages := (totalRecords + int64(pageSize) - 1) / int64(pageSize)

	var resultDTO []dto.ResponseCollectionDTO

	for _, result := range result {
		resultDTO = append(resultDTO, mapper.MapModelCollectionToResponseDTO(result))
	}

	return c.JSON(fiber.Map{
		"data":         resultDTO,
		"totalPages":   totalPages,
		"totalRecords": totalRecords,
		"currentPage":  page,
		"pageSize":     pageSize,
	})
}

func (h *CollectionController) GetOne(c *fiber.Ctx) error {
	collectionId := c.Params("id")

	ctx := c.Context()

	result, err := h.collectionService.GetOne(ctx, collectionId)

	if err != nil {
		return httperror.HandleError(c, err)
	}
	return c.JSON(fiber.Map{
		"data": mapper.MapModelCollectionToResponseDTO(result),
	})
}

func (h *CollectionController) Create(c *fiber.Ctx) error {
	userId := c.Locals("userClaims").(*dto.Claims).UserID

	ctx := c.Context()
	var collectionDto dto.CreateCollectionDTO
	if err := c.BodyParser(&collectionDto); err != nil {
		return httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	collection := mapper.MapCreateCollectionDTOToCollectionModel(collectionDto)

	result, err := h.collectionService.Create(ctx, collection, userId)

	if err != nil {
		return httperror.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"data": mapper.MapModelCollectionToResponseDTO(result),
	})
}
