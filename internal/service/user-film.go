package service

import (
	"context"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/mapper"
	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"net/http"
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

func (s *UserFilmService) GetAll(ctx context.Context, userId string, typeUserFilm model.TypeUserFilm, limit int) ([]dto.GetUserFilmResponseDTO, error) {
	result, err := s.userFilmRepo.GetAllForRecommend(ctx, userId, typeUserFilm, limit)

	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		if typeUserFilm == model.TypeUserFavourite {
			return []dto.GetUserFilmResponseDTO{}, httperror.New(
				http.StatusNotFound,
				"Избранные фильмы не найдены у этого пользователя",
			)
		} else {
			return []dto.GetUserFilmResponseDTO{},
				httperror.New(
					http.StatusNotFound,
					"Рекомендации не найдены у этого пользователя",
				)
		}
	}

	var res []dto.GetUserFilmResponseDTO
	for _, film := range result {
		res = append(res, mapper.MapDomainUserFilmToResponseDTO(film))
	}

	return res, nil
}

func (s *UserFilmService) Add(ctx context.Context, params dto.Params) error {
	return s.userFilmRepo.Add(ctx, params)
}

func (s *UserFilmService) Delete(ctx context.Context, params dto.Params) error {
	return s.userFilmRepo.Delete(ctx, params)
}
