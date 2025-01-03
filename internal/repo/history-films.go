package repo

import (
	"context"
	"strconv"

	"gorm.io/gorm"

	"github.com/1Storm3/flibox-api/database/postgres"
	"github.com/1Storm3/flibox-api/internal/dto"
)

type HistoryFilmsRepo struct {
	storage *postgres.Storage
}

func NewHistoryFilmsRepo(storage *postgres.Storage) *HistoryFilmsRepo {
	return &HistoryFilmsRepo{
		storage: storage,
	}
}

func (r *HistoryFilmsRepo) GetAll(ctx context.Context, userId string) ([]dto.HistoryFilmsRepoDTO, error) {
	var historyFilms []dto.HistoryFilmsRepoDTO
	err := r.storage.DB().WithContext(ctx).
		Where("user_id = ?", userId).
		Order("created_at DESC").
		Preload("Film", func(db *gorm.DB) *gorm.DB {
			return db.Table("films")
		}).
		Limit(5).
		Table("history_films").
		Find(&historyFilms).Error
	if err != nil {
		return nil, err
	}
	return historyFilms, nil
}

func (r *HistoryFilmsRepo) Add(ctx context.Context, filmId, userId string) error {
	isExist := r.storage.DB().WithContext(ctx).Where("user_id = ? AND film_id = ?", userId, filmId).
		Table("history_films").Find(&dto.HistoryFilmsRepoDTO{})
	if isExist.RowsAffected > 0 {
		return nil
	}
	filmIdInt, _ := strconv.Atoi(filmId)

	err := r.storage.DB().WithContext(ctx).
		Table("history_films").Create(&dto.HistoryFilmsRepoDTO{
		UserID: userId,
		FilmID: filmIdInt,
	}).Error
	if err != nil {
		return err
	}
	return nil
}
