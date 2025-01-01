package repo

import (
	"context"
	"errors"
	"net/http"

	"gorm.io/gorm"

	"github.com/1Storm3/flibox-api/database/postgres"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
)

type FilmSimilarRepo struct {
	storage *postgres.Storage
}

func NewFilmSimilarRepo(storage *postgres.Storage) *FilmSimilarRepo {
	return &FilmSimilarRepo{
		storage: storage,
	}
}

func (s *FilmSimilarRepo) GetAll(ctx context.Context, filmID string) ([]dto.FilmSimilarRepoDTO, error) {
	var filmSimilars []dto.FilmSimilarRepoDTO

	result := s.storage.DB().
		WithContext(ctx).
		Where("film_id = ?", filmID).
		Preload("Film", func(db *gorm.DB) *gorm.DB {
			return db.Table("films")
		}).
		Table("film_similars").
		Find(&filmSimilars)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []dto.FilmSimilarRepoDTO{}, nil
	} else if result.Error != nil {
		return []dto.FilmSimilarRepoDTO{},
			httperror.New(
				http.StatusInternalServerError,
				result.Error.Error())
	}

	return filmSimilars, nil
}

func (s *FilmSimilarRepo) Save(ctx context.Context, filmID int, similarID int) error {
	var existingSimilar dto.FilmSimilarRepoDTO

	result := s.storage.DB().WithContext(ctx).Where("film_id = ? AND similar_id = ?", filmID, similarID).
		Table("film_similars").First(&existingSimilar)

	if result.Error == nil {
		return nil
	} else if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {

		return httperror.New(
			http.StatusInternalServerError,
			result.Error.Error(),
		)
	}

	createdResult := s.storage.DB().WithContext(ctx).
		Table("film_similars").Create(&dto.FilmSimilarRepoDTO{
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
