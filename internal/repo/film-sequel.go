package repo

import (
	"context"

	"gorm.io/gorm"

	"github.com/1Storm3/flibox-api/database/postgres"
	"github.com/1Storm3/flibox-api/internal/dto"
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
	err := s.storage.DB().
		WithContext(ctx).
		Where("film_id = ?", filmID).
		Preload("Film", func(db *gorm.DB) *gorm.DB {
			return db.Table("films")
		}).
		Table("film_sequels").
		Find(&filmSequels).Error

	if err != nil {
		return nil, err
	}
	return filmSequels, nil
}

func (s *FilmSequelRepo) Save(ctx context.Context, filmID int, sequelID int) error {
	var existingSequel dto.FilmSequelRepoDTO

	err := s.storage.DB().WithContext(ctx).Where("film_id = ? AND sequel_id = ?", filmID, sequelID).
		Table("film_sequels").First(&existingSequel).Error

	if err != nil {
		return err
	}

	err = s.storage.DB().
		Table("film_sequels").Create(&dto.FilmSequelRepoDTO{
		FilmID:   filmID,
		SequelID: sequelID,
	}).Error

	if err != nil {
		return err
	}

	return nil
}
