package service

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/mapper"
	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/1Storm3/flibox-api/pkg/sys"
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
	err := s.userFilmRepo.AddMany(ctx, params)
	if err != nil {
		return sys.NewError(sys.ErrDatabaseFailure, err.Error())
	}
	return nil
}

func (s *UserFilmService) DeleteMany(ctx context.Context, userID string) error {
	err := s.userFilmRepo.DeleteMany(ctx, userID)
	if err != nil {
		return sys.NewError(sys.ErrDatabaseFailure, err.Error())
	}
	return nil
}

func (s *UserFilmService) GetAll(ctx context.Context, userId string, typeUserFilm dto.TypeUserFilm, limit int) ([]model.UserFilm, error) {
	result, err := s.userFilmRepo.GetAllForRecommend(ctx, userId, typeUserFilm, limit)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []model.UserFilm{}, nil
		}
		return []model.UserFilm{}, sys.NewError(sys.ErrDatabaseFailure, err.Error())
	}
	if len(result) == 0 {
		if typeUserFilm == dto.TypeUserFavourite {
			return []model.UserFilm{}, sys.NewError(sys.ErrFavouriteNotFound, "Пользователь не добавил ни одного фильма в избранное")
		} else {
			return []model.UserFilm{},
				sys.NewError(sys.ErrRecommendationsNotFound, "Рекомендации не найдены")
		}
	}

	var res []model.UserFilm
	for _, film := range result {
		res = append(res, mapper.MapUserFilmRepoDTOToUserFilmModel(film))
	}

	return res, nil
}

func (s *UserFilmService) Add(ctx context.Context, params dto.Params) error {
	err := s.userFilmRepo.Add(ctx, params)
	if err != nil {
		return sys.NewError(sys.ErrDatabaseFailure, err.Error())
	}
	return nil
}

func (s *UserFilmService) Delete(ctx context.Context, params dto.Params) error {
	err := s.userFilmRepo.Delete(ctx, params)
	if err != nil {
		return sys.NewError(sys.ErrDatabaseFailure, err.Error())
	}
	return nil
}
