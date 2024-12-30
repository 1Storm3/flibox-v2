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

type FilmSimilarRepo struct {
	storage *postgres.Storage
}

func NewFilmSimilarRepo(storage *postgres.Storage) *FilmSimilarRepo {
	return &FilmSimilarRepo{
		storage: storage,
	}
}

func (s *FilmSimilarRepo) GetAll(ctx context.Context, filmID string) ([]model.FilmSimilar, error) {
	var filmSimilars []model.FilmSimilar

	result := s.storage.DB().
		WithContext(ctx).
		Where("film_id = ?", filmID).
		Preload("Film").
		Find(&filmSimilars)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []model.FilmSimilar{}, nil
	} else if result.Error != nil {
		return []model.FilmSimilar{},
			httperror.New(
				http.StatusInternalServerError,
				result.Error.Error())
	}

	return filmSimilars, nil
}

func (s *FilmSimilarRepo) Save(ctx context.Context, filmID int, similarID int) error {
	var existingSimilar model.FilmSimilar

	result := s.storage.DB().WithContext(ctx).Where("film_id = ? AND similar_id = ?", filmID, similarID).First(&existingSimilar)

	if result.Error == nil {
		return nil
	} else if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {

		return httperror.New(
			http.StatusInternalServerError,
			result.Error.Error(),
		)
	}

	createdResult := s.storage.DB().WithContext(ctx).Create(&model.FilmSimilar{
		FilmID:    filmID,
		SimilarID: similarID,
	})

	if createdResult.Error != nil {
		return httperror.New(
			http.StatusInternalServerError,
			createdResult.Error.Error(),
		)
	}

	return nil
}
