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

type FilmSequelRepo struct {
	storage *postgres.Storage
}

func NewFilmSequelRepo(storage *postgres.Storage) *FilmSequelRepo {
	return &FilmSequelRepo{
		storage: storage,
	}
}

func (s *FilmSequelRepo) GetAll(ctx context.Context, filmID string) ([]dto.FilmSequelRepoDTO, error) {
	var filmSequels []dto.FilmSequelRepoDTO
	result := s.storage.DB().
		WithContext(ctx).
		Where("film_id = ?", filmID).
		Preload("Film", func(db *gorm.DB) *gorm.DB {
			return db.Table("films")
		}).
		Table("film_sequels").
		Find(&filmSequels)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []dto.FilmSequelRepoDTO{}, nil
	} else if result.Error != nil {
		return []dto.FilmSequelRepoDTO{}, httperror.New(
			http.StatusInternalServerError,
			result.Error.Error(),
		)
	}
	return filmSequels, nil
}

func (s *FilmSequelRepo) Save(ctx context.Context, filmID int, sequelID int) error {
	var existingSequel dto.FilmSequelRepoDTO

	result := s.storage.DB().WithContext(ctx).Where("film_id = ? AND sequel_id = ?", filmID, sequelID).
		Table("film_sequels").First(&existingSequel)

	if result.Error == nil {
		return nil
	} else if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {

		return httperror.New(
			http.StatusInternalServerError,
			result.Error.Error(),
		)
	}

	createdResult := s.storage.DB().
		Table("film_sequels").Create(&dto.FilmSequelRepoDTO{
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
