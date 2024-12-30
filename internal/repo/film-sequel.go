package repo

import (
	"context"
	"errors"
	"github.com/1Storm3/flibox-api/database/postgres"
	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"gorm.io/gorm"
	"net/http"
)

type FilmSequelRepo struct {
	storage *postgres.Storage
}

func NewFilmSequelRepo(storage *postgres.Storage) *FilmSequelRepo {
	return &FilmSequelRepo{
		storage: storage,
	}
}

func (s *FilmSequelRepo) GetAll(ctx context.Context, filmID string) ([]model.FilmSequel, error) {
	var filmSequels []model.FilmSequel
	result := s.storage.DB().
		WithContext(ctx).
		Where("film_id = ?", filmID).
		Preload("Film").
		Find(&filmSequels)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []model.FilmSequel{}, nil
	} else if result.Error != nil {
		return []model.FilmSequel{}, httperror.New(
			http.StatusInternalServerError,
			result.Error.Error(),
		)
	}
	return filmSequels, nil
}

func (s *FilmSequelRepo) Save(ctx context.Context, filmID int, sequelID int) error {
	var existingSequel model.FilmSequel

	result := s.storage.DB().WithContext(ctx).Where("film_id = ? AND sequel_id = ?", filmID, sequelID).First(&existingSequel)

	if result.Error == nil {
		return nil
	} else if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {

		return httperror.New(
			http.StatusInternalServerError,
			result.Error.Error(),
		)
	}

	createdResult := s.storage.DB().Create(&model.FilmSequel{
		FilmID:   filmID,
		SequelID: sequelID,
	})

	if createdResult.Error != nil {
		return httperror.New(
			http.StatusInternalServerError,
			createdResult.Error.Error(),
		)
	}

	return nil
}
