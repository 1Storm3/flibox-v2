package service

import (
	"context"
	"strconv"

	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/delivery/grpc"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/1Storm3/flibox-api/pkg/logger"
	"github.com/1Storm3/flibox-api/pkg/proto/gengrpc"
)

type RecommendService struct {
	historyFilmsService controller.HistoryFilmsService
	filmService         controller.FilmService
	userFilmService     controller.UserFilmService
	grpcClient          grpc.ClientConnInterface
}

func NewRecommendService(
	historyFilmsService controller.HistoryFilmsService,
	filmService controller.FilmService,
	userFilmService controller.UserFilmService,
	grpcClient grpc.ClientConnInterface,
) *RecommendService {
	return &RecommendService{
		grpcClient:          grpcClient,
		filmService:         filmService,
		userFilmService:     userFilmService,
		historyFilmsService: historyFilmsService,
	}
}

func (s *RecommendService) CreateRecommendations(params dto.RecommendationsParams) error {
	ctx := context.Background()

	filmNames, err := s.GetFilmNamesForRecommendations(ctx, params.UserID)
	if err != nil {
		return err
	}

	if len(filmNames) == 0 {
		logger.Info("Нет фильмов для рекомендаций")
		return nil
	}

	err = s.userFilmService.DeleteMany(ctx, params.UserID)
	if err != nil {
		return err
	}

	recommendations, err := s.grpcClient.GetRecommendations(ctx, filmNames)
	if err != nil {
		return err
	}

	filmIds, err := s.GetUniqueFilmIDsForRecommendations(ctx, recommendations)
	if err != nil {
		return err
	}

	if len(filmIds) == 0 {
		logger.Info("Нет рекомендаций")
		return nil
	}

	err = s.AddFilmRecommendations(ctx, params.UserID, filmIds)
	if err != nil {
		return err
	}

	logger.Info("Рекомендации созданы")
	return nil
}

func (s *RecommendService) GetFilmNamesForRecommendations(ctx context.Context, userID string) ([]*gengrpc.Film, error) {
	var filmNames []*gengrpc.Film

	favouriteFilms, err := s.userFilmService.GetAll(ctx, userID, dto.TypeUserFavourite, 5)
	if err != nil {
		return nil, err
	}

	historyFilms, err := s.historyFilmsService.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, historyFilm := range historyFilms {
		film := &historyFilm.Film
		filmNames = append(filmNames, s.GetFilmName(film))
	}

	for _, favouriteFilm := range favouriteFilms {
		film := &favouriteFilm.Film
		filmNames = append(filmNames, s.GetFilmName(film))
	}

	return filmNames, nil
}

func (s *RecommendService) GetFilmName(film *model.Film) *gengrpc.Film {
	var filmName string
	if film.NameOriginal != nil {
		filmName = *film.NameOriginal
	} else if film.NameRU != nil {
		filmName = *film.NameRU
	}
	return &gengrpc.Film{NameOriginal: filmName}
}

func (s *RecommendService) GetUniqueFilmIDsForRecommendations(ctx context.Context, recommendations []string) ([]*int, error) {
	var filmIds []*int
	seenFilmIDs := make(map[int]struct{})

	for _, film := range recommendations {
		filmExist, err := s.filmService.GetOneByNameRu(ctx, film)
		if err != nil {
			return nil, err
		}

		if filmExist.ID == nil {
			// Запрос во внешний апи
			continue
		}
		if _, exists := seenFilmIDs[*filmExist.ID]; exists {
			continue
		}
		seenFilmIDs[*filmExist.ID] = struct{}{}
		filmIds = append(filmIds, filmExist.ID)
	}

	return filmIds, nil
}

func (s *RecommendService) AddFilmRecommendations(ctx context.Context, userID string, filmIds []*int) error {
	var recommendFilms []dto.Params
	for _, id := range filmIds {
		recommendFilms = append(recommendFilms, dto.Params{
			UserID: userID,
			FilmID: strconv.Itoa(*id),
			Type:   dto.TypeUserRecommend,
		})
	}

	return s.userFilmService.AddMany(ctx, recommendFilms)
}
