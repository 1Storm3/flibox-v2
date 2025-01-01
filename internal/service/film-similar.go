package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/1Storm3/flibox-api/internal/config"
	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/1Storm3/flibox-api/pkg/sys"
)

type FilmsSimilarService struct {
	filmSimilarRepo FilmSimilarRepo
	filmService     controller.FilmService
	cfg             *config.Config
}

const baseUrlForAllSimilar = "https://kinopoiskapiunofficial.tech/api/v2.2/films/%s/similars"

func NewFilmSimilarService(filmSimilarRepo FilmSimilarRepo, filmService controller.FilmService, cfg *config.Config) *FilmsSimilarService {
	return &FilmsSimilarService{
		filmSimilarRepo: filmSimilarRepo,
		filmService:     filmService,
		cfg:             cfg,
	}
}

func (s *FilmsSimilarService) GetAll(ctx context.Context, filmId string) ([]model.FilmSimilar, error) {
	result, err := s.filmSimilarRepo.GetAll(ctx, filmId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []model.FilmSimilar{}, nil
		}
		return []model.FilmSimilar{}, sys.NewError(sys.ErrDatabaseFailure, err.Error())
	}

	if len(result) > 0 {
		var similars []model.FilmSimilar
		for _, similar := range result {
			res, err := s.filmService.GetOne(ctx, strconv.Itoa(similar.SimilarID))

			if err != nil {
				return []model.FilmSimilar{}, err
			}
			similars = append(similars, model.FilmSimilar{
				Film: res,
			})
		}
		return similars, nil
	}
	similars, err := s.FetchSimilar(ctx, filmId)
	if err != nil {
		return []model.FilmSimilar{}, err
	}
	return similars, nil
}

func (s *FilmsSimilarService) FetchSimilar(ctx context.Context, filmId string) ([]model.FilmSimilar, error) {
	apikey := s.cfg.DB.ApiKey
	baseUrlForAllSimilar := fmt.Sprintf(baseUrlForAllSimilar, filmId)
	req, err := http.NewRequest("GET", baseUrlForAllSimilar, nil)

	if err != nil {
		return []model.FilmSimilar{}, sys.NewError(sys.ErrUnknown, err.Error())
	}

	req.Header.Add("X-API-KEY", apikey)

	client := &http.Client{}
	resAllSimilars, err := client.Do(req)

	if err != nil {
		return []model.FilmSimilar{},
			sys.NewError(sys.ErrUnknown, err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resAllSimilars.Body)

	if resAllSimilars.StatusCode != http.StatusOK {
		return []model.FilmSimilar{}, sys.NewError(sys.ErrUnknown, "Похожие фильмы не найдены")
	}
	bodyAllSimilars, err := io.ReadAll(resAllSimilars.Body)

	if err != nil {
		return []model.FilmSimilar{},
			sys.NewError(sys.ErrUnknown, err.Error())
	}

	var apiResponse struct {
		Total int                         `json:"total"`
		Items []dto.GetExternalSimilarDTO `json:"items"`
	}

	if apiResponse.Total == 0 || len(apiResponse.Items) == 0 {
		return nil, sys.NewError(sys.ErrUnknown, "Похожие фильмы не найдены")
	}

	err = json.Unmarshal(bodyAllSimilars, &apiResponse)

	if err != nil {
		return []model.FilmSimilar{},
			sys.NewError(sys.ErrUnknown, err.Error())
	}
	var similars []model.FilmSimilar
	for _, similar := range apiResponse.Items {
		filmIsExist, err := s.filmService.GetOne(ctx, strconv.Itoa(similar.FilmId))

		if err != nil {
			return []model.FilmSimilar{}, err
		}

		filmIdInt, err := strconv.Atoi(filmId)

		if err != nil {
			return []model.FilmSimilar{},
				sys.NewError(sys.ErrInvalidRequestData, err.Error())
		}

		err = s.filmSimilarRepo.Save(ctx, filmIdInt, similar.FilmId)
		if err != nil {
			return nil, sys.NewError(sys.ErrDatabaseFailure, err.Error())
		}

		similars = append(similars, model.FilmSimilar{
			Film: filmIsExist,
		})
	}
	return similars, nil
}
