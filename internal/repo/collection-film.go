package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/1Storm3/flibox-api/database/postgres"
	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type CollectionFilmRepo struct {
	storage *postgres.Storage
}

func NewCollectionFilmRepo(storage *postgres.Storage) *CollectionFilmRepo {
	return &CollectionFilmRepo{
		storage: storage,
	}
}

func (c *CollectionFilmRepo) GetFilmsByCollectionId(
	ctx context.Context,
	collectionID string,
	page int, pageSize int,
) ([]model.Film, int64, error) {
	var films []model.Film
	var totalRecords int64

	offset := (page - 1) * pageSize

	err := c.storage.DB().WithContext(ctx).
		Model(&model.CollectionFilm{}).
		Where("collection_id = ?", collectionID).
		Count(&totalRecords).Error
	if err != nil {
		return nil, 0, err
	}

	err = c.storage.DB().WithContext(ctx).
		Model(&model.CollectionFilm{}).
		Select("films.id, films.name_ru, films.name_original, films.year, films.poster_url, films.rating_kinopoisk, films.type").
		Joins("JOIN films ON films.id = collection_films.film_id").
		Where("collection_films.collection_id = ?", collectionID).
		Offset(offset).
		Limit(pageSize).
		Scan(&films).Error
	if err != nil {
		return nil, 0, err
	}
	return films, totalRecords, nil
}

func (c *CollectionFilmRepo) Add(
	ctx context.Context,
	collectionId string,
	filmId int,
) error {
	var collection model.Collection
	err := c.storage.DB().WithContext(ctx).
		Table("collections").
		Where("id = ?", collectionId).
		First(&collection).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httperror.New(
				http.StatusNotFound,
				"Коллекция не найдена",
			)
		}
		return httperror.New(
			http.StatusInternalServerError,
			fmt.Sprintf("Ошибка при получении коллекции: %v", err),
		)
	}

	newCollectionFilm := &model.CollectionFilm{
		CollectionID: collectionId,
		FilmID:       filmId,
	}
	err = c.storage.DB().WithContext(ctx).Create(newCollectionFilm).Error

	if err != nil {
		if strings.Contains(err.Error(), "violates unique constraint") {
			return httperror.New(
				http.StatusConflict,
				"Фильм уже добавлен в коллекцию",
			)
		}
		if strings.Contains(err.Error(), "collection_films_film_id_fkey") {
			return httperror.New(
				http.StatusNotFound,
				"Фильм не найден",
			)
		}
		return httperror.New(
			http.StatusInternalServerError,
			fmt.Sprintf("Ошибка при добавлении фильма в коллекцию: %v", err),
		)
	}

	return nil
}

func (c *CollectionFilmRepo) Delete(ctx context.Context, collectionId string, filmId int) error {
	err := c.storage.DB().WithContext(ctx).
		Where("collection_id = ? AND film_id = ?", collectionId, filmId).
		Delete(&model.CollectionFilm{}).Error

	if err != nil {
		return httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	return nil
}
