package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"github.com/1Storm3/flibox-api/pkg/logger"
)

type UserFilmController struct {
	userFilmService  controller.UserFilmService
	filmService      controller.FilmService
	recommendService controller.RecommendService
}

func NewUserFilmController(
	userFilmService controller.UserFilmService,
	filmService controller.FilmService,
	recommendService controller.RecommendService,
) *UserFilmController {
	return &UserFilmController{
		userFilmService:  userFilmService,
		filmService:      filmService,
		recommendService: recommendService,
	}
}

func (g *UserFilmController) GetAll(c *fiber.Ctx) error {
	userID := c.Locals("userClaims").(*dto.Claims).UserID
	typeUserFilm := c.Query("type")

	ctx := c.Context()

	films, err := g.userFilmService.GetAll(ctx, userID, dto.TypeUserFilm(typeUserFilm), 20)
	if err != nil {
		return httperror.HandleError(c, err)
	}

	return c.JSON(films)
}

func (g *UserFilmController) Add(c *fiber.Ctx) error {
	userID, filmID, typeUserFilm, err := extractUserFilmParams(c)
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
		return httperror.HandleError(c, err)
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
	}

	return c.JSON(fiber.Map{
		"data": "Фильм добавлен в избранное",
	})
}

func (g *UserFilmController) Delete(c *fiber.Ctx) error {
	userID, filmID, typeUserFilm, err := extractUserFilmParams(c)
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
		return httperror.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"data": "Фильм удален из избранного",
	})
}

func extractUserFilmParams(c *fiber.Ctx) (userID, filmID string, typeUserFilm dto.TypeUserFilm, err error) {
	userID = c.Locals("userClaims").(*dto.Claims).UserID
	filmID = c.Params("filmId")

	typeUserFilmReq := c.Query("type")
	if err := ParseTypeUserFilm(typeUserFilmReq, &typeUserFilm); err != nil {
		return "", "", "", httperror.New(
			http.StatusBadRequest,
			"Недопустимый тип фильма: "+typeUserFilmReq,
		)
	}
	return userID, filmID, typeUserFilm, nil
}

func (g *UserFilmController) checkFilmExistence(ctx context.Context, filmID string) error {
	isExist, err := g.filmService.GetOne(ctx, filmID)
	if err != nil {
		return err
	}
	if isExist.Description == nil {
		return httperror.New(
			http.StatusNotFound,
			"Фильм не найден",
		)
	}
	return nil
}

func ParseTypeUserFilm(s string, t *dto.TypeUserFilm) error {
	switch s {
	case string(dto.TypeUserFavourite):
		*t = dto.TypeUserFavourite
	case string(dto.TypeUserRecommend):
		*t = dto.TypeUserRecommend
	default:
		return errors.New("Неверный тип фильма: " + s)
	}
	return nil
}
