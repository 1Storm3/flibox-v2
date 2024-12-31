package service

import (
	"context"
	"net/http"

	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/mapper"
	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
)

type UserFilmService struct {
	userFilmRepo UserFilmRepo
}

func NewUserFilmService(userFilmRepo UserFilmRepo) *UserFilmService {
	return &UserFilmService{
		userFilmRepo: userFilmRepo,
	}
}

func (s *UserFilmService) AddMany(ctx context.Context, params []dto.Params) error {
	return s.userFilmRepo.AddMany(ctx, params)
}

func (s *UserFilmService) DeleteMany(ctx context.Context, userID string) error {
	return s.userFilmRepo.DeleteMany(ctx, userID)
}

func (s *UserFilmService) GetAll(ctx context.Context, userId string, typeUserFilm dto.TypeUserFilm, limit int) ([]model.UserFilm, error) {
	result, err := s.userFilmRepo.GetAllForRecommend(ctx, userId, typeUserFilm, limit)

	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		if typeUserFilm == dto.TypeUserFavourite {
			return []model.UserFilm{}, httperror.New(
				http.StatusNotFound,
				"Избранные фильмы не найдены у этого пользователя",
			)
		} else {
			return []model.UserFilm{},
				httperror.New(
					http.StatusNotFound,
					"Рекомендации не найдены у этого пользователя",
				)
		}
	}

	var res []model.UserFilm
	for _, film := range result {
		res = append(res, mapper.MapUserFilmRepoDTOToUserFilmModel(film))
	}

	return res, nil
}

func (s *UserFilmService) Add(ctx context.Context, params dto.Params) error {
	return s.userFilmRepo.Add(ctx, params)
}

func (s *UserFilmService) Delete(ctx context.Context, params dto.Params) error {
	return s.userFilmRepo.Delete(ctx, params)
}
