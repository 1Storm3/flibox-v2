package repo

import (
	"context"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"

	"github.com/1Storm3/flibox-api/database/postgres"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
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
	res := r.storage.DB().WithContext(ctx).
		Where("user_id = ?", userId).
		Order("created_at DESC").
		Preload("Film", func(db *gorm.DB) *gorm.DB {
			return db.Table("films")
		}).
		Limit(5).
		Table("history_films").
		Find(&historyFilms)
	if res.Error != nil {
		return nil, httperror.New(
			http.StatusInternalServerError,
			res.Error.Error(),
		)
	}
	return historyFilms, nil
}

func (r *HistoryFilmsRepo) Add(ctx context.Context, filmId, userId string) error {
	isExist := r.storage.DB().WithContext(ctx).Where("user_id = ? AND film_id = ?", userId, filmId).
		Table("history_films").Find(&dto.HistoryFilmsRepoDTO{})
	if isExist.RowsAffected > 0 {
		return httperror.New(
			http.StatusConflict,
			"Фильм уже добавлен в историю просмотров",
		)
	}
	filmIdInt, _ := strconv.Atoi(filmId)
	res := r.storage.DB().WithContext(ctx).
		Table("history_films").Create(&dto.HistoryFilmsRepoDTO{
		UserID: userId,
		FilmID: filmIdInt,
	})
	if res.Error != nil {
		if strings.Contains(res.Error.Error(), "violates foreign key constraint") {
			return httperror.New(
				http.StatusConflict,
				"Фильм не существует с таким ID",
			)
		}
		return httperror.New(
			http.StatusInternalServerError,
			res.Error.Error(),
		)
	}
	return nil
}
