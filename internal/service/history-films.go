package service

import (
	"context"

	"github.com/1Storm3/flibox-api/internal/mapper"
	"github.com/1Storm3/flibox-api/internal/model"
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
	return s.historyFilmsRepo.Add(ctx, filmId, userId)
}

func (s *HistoryFilmsService) GetAll(ctx context.Context, userId string) ([]model.HistoryFilms, error) {
	result, err := s.historyFilmsRepo.GetAll(ctx, userId)
	if err != nil {
		return nil, err
	}

	var historyFilms []model.HistoryFilms
	for _, historyFilm := range result {
		historyFilms = append(historyFilms, model.HistoryFilms{
			Film: mapper.MapFilmRepoDTOToFilmModel(historyFilm.Film),
		})
	}
	return historyFilms, nil
}
