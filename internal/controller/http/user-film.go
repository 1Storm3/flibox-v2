package http

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/shared/helper"
	"github.com/1Storm3/flibox-api/pkg/kafka"
	"github.com/1Storm3/flibox-api/pkg/logger"
	"github.com/1Storm3/flibox-api/pkg/sys"
)

type UserFilmController struct {
	userFilmService  controller.UserFilmService
	filmService      controller.FilmService
	recommendService controller.RecommendService
	kafkaProducer    *kafka.Producer
}

func NewUserFilmController(
	userFilmService controller.UserFilmService,
	filmService controller.FilmService,
	recommendService controller.RecommendService,
	kafkaProducer *kafka.Producer,
) *UserFilmController {
	return &UserFilmController{
		userFilmService:  userFilmService,
		filmService:      filmService,
		recommendService: recommendService,
		kafkaProducer:    kafkaProducer,
	}
}

func (g *UserFilmController) GetAll(c *fiber.Ctx) error {
	userID := c.Locals("userClaims").(*dto.Claims).UserID
	typeUserFilm := c.Query("type")

	ctx := c.Context()

	films, err := g.userFilmService.GetAll(ctx, userID, dto.TypeUserFilm(typeUserFilm), 20)
	if err != nil {
		return sys.HandleError(c, err)
	}

	return c.JSON(films)
}

func (g *UserFilmController) Add(c *fiber.Ctx) error {
	userID, filmID, typeUserFilm, err := helper.ExtractUserFilmParams(c)
	if err != nil {
		return err
	}

	ctx := c.Context()
	if err := g.checkFilmExistence(ctx, filmID); err != nil {
		return err
	}

	err = g.userFilmService.Add(ctx, dto.Params{
		UserID: userID,
		FilmID: filmID,
		Type:   typeUserFilm,
	})
	if err != nil {
		return sys.HandleError(c, err)
	}

	if typeUserFilm == dto.TypeUserFavourite {
		go func() {
			err := g.recommendService.CreateRecommendations(dto.RecommendationsParams{
				UserID: userID,
			})
			if err != nil {
				logger.Info("Произошла ошибка при создании рекомендаций")
				logger.Error(err.Error())
			}
		}()

		//go func() {
		//	err := g.kafkaProducer.Produce(map[string]interface{}{
		//		"user_id": userID,
		//		"film_id": filmID,
		//	})
		//	if err != nil {
		//		logger.Info("Произошла ошибка при отправке сообщения в Kafka")
		//		logger.Error(err.Error())
		//	}
		//}()
	}

	return c.JSON(fiber.Map{
		"data": "Фильм добавлен в избранное",
	})
}

func (g *UserFilmController) Delete(c *fiber.Ctx) error {
	userID, filmID, typeUserFilm, err := helper.ExtractUserFilmParams(c)
	if err != nil {
		return err
	}

	ctx := c.Context()
	err = g.userFilmService.Delete(ctx, dto.Params{
		UserID: userID,
		FilmID: filmID,
		Type:   typeUserFilm,
	})
	if err != nil {
		return sys.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"data": "Фильм удален из избранного",
	})
}

func (g *UserFilmController) checkFilmExistence(ctx context.Context, filmID string) error {
	isExist, err := g.filmService.GetOne(ctx, filmID)
	if err != nil {
		return err
	}
	if isExist.Description == nil {
		return sys.NewError(sys.ErrFilmNotFound, "")
	}
	return nil
}
