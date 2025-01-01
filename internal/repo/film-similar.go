package repo

import (
	"context"

	"gorm.io/gorm"

	"github.com/1Storm3/flibox-api/database/postgres"
	"github.com/1Storm3/flibox-api/internal/dto"
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

	err := s.storage.DB().
		WithContext(ctx).
		Where("film_id = ?", filmID).
		Preload("Film", func(db *gorm.DB) *gorm.DB {
			return db.Table("films")
		}).
		Table("film_similars").
		Find(&filmSimilars).Error

	if err != nil {
		return nil, err
	}

	return filmSimilars, nil
}

func (s *FilmSimilarRepo) Save(ctx context.Context, filmID int, similarID int) error {
	var existingSimilar dto.FilmSimilarRepoDTO

	err := s.storage.DB().WithContext(ctx).Where("film_id = ? AND similar_id = ?", filmID, similarID).
		Table("film_similars").First(&existingSimilar).Error

	if err != nil {
		return err
	}

	err = s.storage.DB().WithContext(ctx).
		Table("film_similars").Create(&dto.FilmSimilarRepoDTO{
		FilmID:    filmID,
		SimilarID: similarID,
	}).Error

	if err != nil {
		return err
	}

	return nil
}
