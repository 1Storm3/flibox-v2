package http

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/mapper"
	"github.com/1Storm3/flibox-api/pkg/sys"
)

type CommentController struct {
	commentService controller.CommentService
}

func NewCommentController(commentService controller.CommentService) *CommentController {
	return &CommentController{
		commentService: commentService,
	}
}

func (h *CommentController) GetAllByFilmID(c *fiber.Ctx) error {
	filmID := c.Params("filmId")
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)
	ctx := c.Context()
	comments, totalRecords, err := h.commentService.GetAllByFilmId(ctx, filmID, page, pageSize)
	if err != nil {
		return sys.HandleError(c, err)
	}

	totalPages := (totalRecords + int64(pageSize) - 1) / int64(pageSize)

	var resultDTO []dto.CommentResponseDTO
	for _, comment := range comments {
		resultDTO = append(resultDTO, mapper.MapCommentModelToCommentResponseDTO(comment))
	}

	return c.JSON(fiber.Map{
		"comments":     resultDTO,
		"totalPages":   totalPages,
		"totalRecords": totalRecords,
		"currentPage":  page,
		"pageSize":     pageSize,
	})
}

func (h *CommentController) Create(c *fiber.Ctx) error {
	userId := c.Locals("userClaims").(*dto.Claims).UserID

	ctx := c.Context()
	var comment dto.CreateCommentDTO
	if err := c.BodyParser(&comment); err != nil {
		return sys.NewError(sys.ErrInvalidRequestData, err.Error())
	}

	if len(strings.TrimSpace(*comment.Content)) == 0 {
		return sys.NewError(sys.ErrInvalidRequestData, "Комментарий не может быть пустым")
	}

	comment.UserID = userId

	commentModel := mapper.MapCreateCommentDTOToCommentModel(comment)

	result, err := h.commentService.Create(ctx, commentModel)

	if err != nil {
		return sys.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"data": mapper.MapCommentModelToCommentResponseDTO(result),
	})
}

func (h *CommentController) Update(c *fiber.Ctx) error {
	userId := c.Locals("userClaims").(*dto.Claims).UserID
	role := c.Locals("userClaims").(*dto.Claims).Role

	commentId := c.Params("id")

	ctx := c.Context()
	comment, err := h.commentService.GetOne(ctx, commentId)
	if err != nil {
		return sys.HandleError(c, err)
	}

	if role != "admin" && comment.User.ID != userId {
		return sys.NewError(sys.ErrAccessDenied, "")
	}
	var commentDto dto.UpdateCommentDTO
	if err := c.BodyParser(&commentDto); err != nil {
		return sys.NewError(sys.ErrInvalidRequestData, err.Error())
	}
	commentDto.ID = commentId

	commentModel := mapper.MapUpdateCommentDTOToCommentModel(commentDto)

	result, err := h.commentService.Update(ctx, commentModel, commentId)

	if err != nil {
		return sys.HandleError(c, err)
	}
	return c.JSON(fiber.Map{
		"data": mapper.MapCommentModelToCommentResponseDTO(result),
	})
}

func (h *CommentController) Delete(c *fiber.Ctx) error {
	userId := c.Locals("userClaims").(*dto.Claims).UserID
	role := c.Locals("userClaims").(*dto.Claims).Role
	commentId := c.Params("id")
	ctx := c.Context()
	comment, err := h.commentService.GetOne(ctx, commentId)
	if err != nil {
		return sys.HandleError(c, err)
	}

	if role != "admin" && comment.User.ID != userId {
		return sys.NewError(sys.ErrAccessDenied, "")
	}
	err = h.commentService.Delete(ctx, commentId)
	if err != nil {
		return sys.HandleError(c, err)
	}
	return c.JSON(fiber.Map{
		"data": "Комментарий удален",
	})
}
