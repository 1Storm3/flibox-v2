package repo

import (
	"context"
	"strconv"

	"gorm.io/gorm"

	"github.com/1Storm3/flibox-api/database/postgres"
	"github.com/1Storm3/flibox-api/internal/dto"
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
	var userFilms []dto.UserFilmRepoDTO
	for _, param := range params {
		filmIdInt, _ := strconv.Atoi(param.FilmID)
		userFilms = append(userFilms, dto.UserFilmRepoDTO{
			UserID: param.UserID,
			FilmID: filmIdInt,
			Type:   param.Type,
		})
	}

	err := u.storage.DB().WithContext(ctx).Table("user_films").Create(&userFilms).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *UserFilmRepo) DeleteMany(ctx context.Context, userID string) error {
	err := u.storage.DB().
		WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, dto.TypeUserRecommend).
		Table("user_films").
		Delete(&dto.UserFilmRepoDTO{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *UserFilmRepo) GetAllForRecommend(ctx context.Context, userId string, typeUserFilm dto.TypeUserFilm, limit int) ([]dto.UserFilmRepoDTO, error) {
	var userFilms []dto.UserFilmRepoDTO
	err := u.storage.DB().
		WithContext(ctx).
		Where("user_id = ?", userId).
		Where("type = ?", typeUserFilm).
		Preload("Film", func(db *gorm.DB) *gorm.DB {
			return db.Table("films")
		}).
		Limit(limit).
		Table("user_films").
		Find(&userFilms).Error

	if err != nil {
		return nil, err
	}
	return userFilms, nil
}

func (u *UserFilmRepo) Add(ctx context.Context, params dto.Params) error {
	filmIDInt, _ := strconv.Atoi(params.FilmID)

	isFavourite := u.storage.DB().WithContext(ctx).
		Where("user_id = ? AND film_id = ? AND type = ?", params.UserID, filmIDInt, params.Type).
		Table("user_films").
		Find(&dto.UserFilmRepoDTO{})
	if isFavourite.RowsAffected > 0 {
		return nil
	}

	err := u.storage.DB().WithContext(ctx).Table("user_films").Create(&dto.UserFilmRepoDTO{
		UserID: params.UserID,
		FilmID: filmIDInt,
		Type:   params.Type,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *UserFilmRepo) Delete(ctx context.Context, params dto.Params) error {
	isFavourite := u.storage.DB().WithContext(ctx).
		Where("user_id = ? AND film_id = ? AND type = ?", params.UserID, params.FilmID, params.Type).
		Table("user_films").
		Find(&dto.UserFilmRepoDTO{})
	if isFavourite.RowsAffected == 0 {
		return nil
	}

	err := u.storage.DB().WithContext(ctx).
		Where("user_id = ? AND film_id = ? AND type = ?", params.UserID, params.FilmID, params.Type).
		Table("user_films").
		Delete(&dto.UserFilmRepoDTO{}).Error
	if err != nil {
		return err
	}
	return nil
}
