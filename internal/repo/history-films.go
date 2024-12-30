package repo

import (
	"context"
	"github.com/1Storm3/flibox-api/database/postgres"
	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"net/http"
	"strconv"
	"strings"
)

type HistoryFilmsRepo struct {
	storage *postgres.Storage
}

func NewHistoryFilmsRepo(storage *postgres.Storage) *HistoryFilmsRepo {
	return &HistoryFilmsRepo{
		storage: storage,
	}
}

func (r *HistoryFilmsRepo) GetAll(ctx context.Context, userId string) ([]model.HistoryFilms, error) {
	var historyFilms []model.HistoryFilms
	res := r.storage.DB().WithContext(ctx).
		Where("user_id = ?", userId).
		Order("created_at DESC").
		Preload("Film").Limit(5).
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
	isExist := r.storage.DB().WithContext(ctx).Where("user_id = ? AND film_id = ?", userId, filmId).Find(&model.HistoryFilms{})
	if isExist.RowsAffected > 0 {
		return httperror.New(
			http.StatusConflict,
			"Фильм уже добавлен в историю просмотров",
		)
	}
	filmIdInt, _ := strconv.Atoi(filmId)
	res := r.storage.DB().WithContext(ctx).Create(&model.HistoryFilms{
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
