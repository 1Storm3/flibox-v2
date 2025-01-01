package service

import (
	"context"
	"strings"

	"github.com/1Storm3/flibox-api/internal/mapper"
	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/1Storm3/flibox-api/pkg/sys"
)

type HistoryFilmsService struct {
	historyFilmsRepo HistoryFilmsRepo
}

func NewHistoryFilmsService(historyFilmsRepo HistoryFilmsRepo) *HistoryFilmsService {
	return &HistoryFilmsService{
		historyFilmsRepo: historyFilmsRepo,
	}
}

func (s *HistoryFilmsService) Add(ctx context.Context, filmId, userId string) error {
	err := s.historyFilmsRepo.Add(ctx, filmId, userId)
	if err != nil {
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			return sys.NewError(
				sys.ErrFilmNotFound,
				"Фильм не существует с таким ID",
			)
		}
		return sys.NewError(sys.ErrDatabaseFailure, err.Error())
	}
	return nil
}

func (s *HistoryFilmsService) GetAll(ctx context.Context, userId string) ([]model.HistoryFilms, error) {
	result, err := s.historyFilmsRepo.GetAll(ctx, userId)
	if err != nil {
		return nil, sys.NewError(sys.ErrDatabaseFailure, err.Error())
	}

	var historyFilms []model.HistoryFilms
	for _, historyFilm := range result {
		historyFilms = append(historyFilms, model.HistoryFilms{
			Film: mapper.MapFilmRepoDTOToFilmModel(historyFilm.Film),
		})
	}
	return historyFilms, nil
}
