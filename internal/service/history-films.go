package service

import (
	"context"
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
	return s.historyFilmsRepo.GetAll(ctx, userId)
}
