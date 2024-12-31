package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"github.com/1Storm3/flibox-api/pkg/logger"
)

type HistoryFilmsController struct {
	historyFilmService controller.HistoryFilmsService
	recommendService   controller.RecommendService
}

func NewHistoryFilmsController(
	historyFilmService controller.HistoryFilmsService,
	recommendService controller.RecommendService,
) *HistoryFilmsController {
	return &HistoryFilmsController{
		historyFilmService: historyFilmService,
		recommendService:   recommendService,
	}
}

func (h *HistoryFilmsController) Add(c *fiber.Ctx) error {
	userID := c.Locals("userClaims").(*dto.Claims).UserID
	filmID := c.Params("Id")
	ctx := c.Context()
	err := h.historyFilmService.Add(ctx, filmID, userID)
	if err != nil {
		return httperror.HandleError(c, err)
	}
	go func() {
		err := h.recommendService.CreateRecommendations(dto.RecommendationsParams{
			UserID: userID,
		})
		if err != nil {
			logger.Info("Произошла ошибка при создании рекомендаций")
			logger.Error(err.Error())
		}
	}()

	return c.JSON(fiber.Map{
		"data": "Фильм добавлен в историю просмотра",
	})
}
