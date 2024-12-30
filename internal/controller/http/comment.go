package http

import (
	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strings"
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
		return httperror.HandleError(c, err)
	}

	totalPages := (totalRecords + int64(pageSize) - 1) / int64(pageSize)

	return c.JSON(fiber.Map{
		"comments":     comments,
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
		return httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	if len(strings.TrimSpace(*comment.Content)) == 0 {
		return httperror.New(http.StatusBadRequest, "Комментарий не может быть пустым")
	}
	result, err := h.commentService.Create(ctx, comment, userId)
	if err != nil {
		return httperror.HandleError(c, err)
	}
	return c.JSON(fiber.Map{
		"data": result,
	})
}

func (h *CommentController) Update(c *fiber.Ctx) error {
	userId := c.Locals("userClaims").(*dto.Claims).UserID
	role := c.Locals("userClaims").(*dto.Claims).Role

	commentId := c.Params("id")

	ctx := c.Context()
	comment, err := h.commentService.GetOne(ctx, commentId)
	if err != nil {
		return httperror.HandleError(c, err)
	}

	if role != "admin" && comment.User.ID != userId {
		return httperror.New(
			http.StatusForbidden,
			"Недостаточно прав для редактирования комментария",
		)
	}
	var commentDto dto.UpdateCommentDTO
	if err := c.BodyParser(&commentDto); err != nil {
		return httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	result, err := h.commentService.Update(ctx, commentDto, commentId)
	if err != nil {
		return httperror.HandleError(c, err)
	}
	return c.JSON(fiber.Map{
		"data": result,
	})
}

func (h *CommentController) Delete(c *fiber.Ctx) error {
	userId := c.Locals("userClaims").(*dto.Claims).UserID
	role := c.Locals("userClaims").(*dto.Claims).Role
	commentId := c.Params("id")
	ctx := c.Context()
	comment, err := h.commentService.GetOne(ctx, commentId)
	if err != nil {
		return httperror.HandleError(c, err)
	}

	if role != "admin" && comment.User.ID != userId {
		return httperror.New(
			http.StatusForbidden,
			"Недостаточно прав для удаления комментария",
		)
	}
	err = h.commentService.Delete(ctx, commentId)
	if err != nil {
		return httperror.HandleError(c, err)
	}
	return c.JSON(fiber.Map{
		"data": "Комментарий удален",
	})
}
