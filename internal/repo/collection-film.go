package repo

import (
	"context"

	"github.com/1Storm3/flibox-api/database/postgres"
	"github.com/1Storm3/flibox-api/internal/dto"
)

type CollectionFilmRepo struct {
	storage *postgres.Storage
}

func NewCollectionFilmRepo(storage *postgres.Storage) *CollectionFilmRepo {
	return &CollectionFilmRepo{
		storage: storage,
	}
}

func (c *CollectionFilmRepo) GetFilmsByCollectionID(
	ctx context.Context,
	collectionID string,
	page int, pageSize int,
) ([]dto.FilmRepoDTO, int64, error) {
	var films []dto.FilmRepoDTO
	var totalRecords int64

	offset := (page - 1) * pageSize

	err := c.storage.DB().WithContext(ctx).
		Table("collection_films").
		Model(&dto.FilmRepoDTO{}).
		Where("collection_id = ?", collectionID).
		Count(&totalRecords).Error
	if err != nil {
		return nil, 0, err
	}

	err = c.storage.DB().WithContext(ctx).
		Table("collection_films").
		Model(&dto.FilmRepoDTO{}).
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
	collectionID string,
	filmID int,
) error {
	var collection dto.CollectionRepoDTO
	err := c.storage.DB().WithContext(ctx).
		Table("collections").
		Where("id = ?", collectionID).
		First(&collection).Error

	if err != nil {
		return err
	}

	newCollectionFilm := &dto.CollectionFilmRepoDTO{
		CollectionID: collectionID,
		FilmID:       filmID,
	}
	err = c.storage.DB().WithContext(ctx).Table("collection_films").Create(newCollectionFilm).Error

	if err != nil {
		return err
	}

	return nil
}

func (c *CollectionFilmRepo) Delete(ctx context.Context, collectionId string, filmId int) error {
	err := c.storage.DB().WithContext(ctx).
		Table("collection_films").
		Where("collection_id = ? AND film_id = ?", collectionId, filmId).
		Delete(&dto.CollectionRepoDTO{}).Error

	if err != nil {
		return err
	}

	return nil
}
