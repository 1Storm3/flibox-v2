package repo

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/1Storm3/flibox-api/database/postgres"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
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

	result := u.storage.DB().WithContext(ctx).Table("user_films").Create(&userFilms)
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
		Where("user_id = ? AND type = ?", userID, dto.TypeUserRecommend).
		Table("user_films").
		Delete(&dto.UserFilmRepoDTO{})
	if result.Error != nil {
		return httperror.New(
			http.StatusInternalServerError,
			result.Error.Error(),
		)
	}
	return nil
}

func (u *UserFilmRepo) GetAllForRecommend(ctx context.Context, userId string, typeUserFilm dto.TypeUserFilm, limit int) ([]dto.UserFilmRepoDTO, error) {
	var userFilms []dto.UserFilmRepoDTO
	result := u.storage.DB().
		WithContext(ctx).
		Where("user_id = ?", userId).
		Where("type = ?", typeUserFilm).
		Preload("Film", func(db *gorm.DB) *gorm.DB {
			return db.Table("films")
		}).
		Limit(limit).
		Table("user_films").
		Find(&userFilms)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []dto.UserFilmRepoDTO{}, nil

	} else if result.Error != nil {
		return []dto.UserFilmRepoDTO{},
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
		Table("user_films").
		Find(&dto.UserFilmRepoDTO{})
	if isFavourite.RowsAffected > 0 {
		return httperror.New(
			http.StatusConflict,
			"Фильм уже добавлен в избранное",
		)
	}

	result := u.storage.DB().WithContext(ctx).Table("user_films").Create(&model.UserFilm{
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
		Table("user_films").
		Find(&dto.UserFilmRepoDTO{})
	if isFavourite.RowsAffected == 0 {
		return httperror.New(
			http.StatusNotFound,
			"Фильм не найден в избранном",
		)
	}

	result := u.storage.DB().WithContext(ctx).
		Where("user_id = ? AND film_id = ? AND type = ?", params.UserID, params.FilmID, params.Type).
		Table("user_films").
		Delete(&dto.UserFilmRepoDTO{})
	if result.Error != nil {
		return httperror.New(
			http.StatusInternalServerError,
			result.Error.Error(),
		)
	}
	return nil
}
