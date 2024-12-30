package repo

import (
	"context"
	"errors"
	"github.com/1Storm3/flibox-api/database/postgres"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type UserFilmRepo struct {
	storage *postgres.Storage
}

func NewUserFilmRepo(storage *postgres.Storage) *UserFilmRepo {
	return &UserFilmRepo{
		storage: storage,
	}
}

func (u *UserFilmRepo) AddMany(ctx context.Context, params []dto.Params) error {
	var userFilms []model.UserFilm
	for _, param := range params {
		filmIdInt, err := strconv.Atoi(param.FilmID)
		if err != nil {
			return httperror.New(
				http.StatusBadRequest,
				err.Error(),
			)
		}
		userFilms = append(userFilms, model.UserFilm{
			UserID: param.UserID,
			FilmID: filmIdInt,
			Type:   param.Type,
		})
	}

	result := u.storage.DB().WithContext(ctx).Create(&userFilms)
	if result.Error != nil {
		return httperror.New(
			http.StatusInternalServerError,
			result.Error.Error(),
		)
	}
	return nil
}

func (u *UserFilmRepo) DeleteMany(ctx context.Context, userID string) error {
	result := u.storage.DB().
		WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, model.TypeUserRecommend).
		Delete(&model.UserFilm{})
	if result.Error != nil {
		return httperror.New(
			http.StatusInternalServerError,
			result.Error.Error(),
		)
	}
	return nil
}

func (u *UserFilmRepo) GetAllForRecommend(ctx context.Context, userId string, typeUserFilm model.TypeUserFilm, limit int) ([]model.UserFilm, error) {
	var userFilms []model.UserFilm
	result := u.storage.DB().
		WithContext(ctx).
		Where("user_id = ?", userId).
		Where("type = ?", typeUserFilm).
		Preload("Film").
		Limit(limit).
		Find(&userFilms)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []model.UserFilm{}, nil

	} else if result.Error != nil {
		return []model.UserFilm{},
			httperror.New(
				http.StatusInternalServerError,
				result.Error.Error(),
			)
	}
	return userFilms, nil
}

func (u *UserFilmRepo) Add(ctx context.Context, params dto.Params) error {
	filmIDInt, err := strconv.Atoi(params.FilmID)
	if err != nil {
		return httperror.New(
			http.StatusBadRequest,
			"Неверный формат ID фильма",
		)
	}

	isFavourite := u.storage.DB().WithContext(ctx).
		Where("user_id = ? AND film_id = ? AND type = ?", params.UserID, filmIDInt, params.Type).
		Find(&model.UserFilm{})
	if isFavourite.RowsAffected > 0 {
		return httperror.New(
			http.StatusConflict,
			"Фильм уже добавлен в избранное",
		)
	}

	result := u.storage.DB().WithContext(ctx).Create(&model.UserFilm{
		UserID: params.UserID,
		FilmID: filmIDInt,
		Type:   params.Type,
	})
	if result.Error != nil {
		return httperror.New(
			http.StatusInternalServerError,
			result.Error.Error(),
		)
	}
	return nil
}

func (u *UserFilmRepo) Delete(ctx context.Context, params dto.Params) error {
	isFavourite := u.storage.DB().WithContext(ctx).
		Where("user_id = ? AND film_id = ? AND type = ?", params.UserID, params.FilmID, params.Type).
		Find(&model.UserFilm{})
	if isFavourite.RowsAffected == 0 {
		return httperror.New(
			http.StatusNotFound,
			"Фильм не найден в избранном",
		)
	}
	result := u.storage.DB().WithContext(ctx).
		Where("user_id = ? AND film_id = ? AND type = ?", params.UserID, params.FilmID, params.Type).
		Delete(&model.UserFilm{})
	if result.Error != nil {
		return httperror.New(
			http.StatusInternalServerError,
			result.Error.Error(),
		)
	}
	return nil
}
