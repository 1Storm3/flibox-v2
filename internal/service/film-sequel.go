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

type FilmSequelService struct {
	filmSequelRepo FilmSequelRepo
	filmService    controller.FilmService
	cfg            *config.Config
}

const baseUrlForAllSequels = "https://kinopoiskapiunofficial.tech/api/v2.1/films/%s/sequels_and_prequels"

func NewFilmSequelService(filmSequelRepo FilmSequelRepo, filmService controller.FilmService, cfg *config.Config) *FilmSequelService {
	return &FilmSequelService{
		filmSequelRepo: filmSequelRepo,
		filmService:    filmService,
		cfg:            cfg,
	}
}

func (s *FilmSequelService) GetAll(ctx context.Context, filmId string) ([]model.FilmSequel, error) {
	result, err := s.filmSequelRepo.GetAll(ctx, filmId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []model.FilmSequel{}, nil
		}
		return []model.FilmSequel{}, sys.NewError(sys.ErrDatabaseFailure, err.Error())
	}
	if len(result) > 0 {
		var sequels []model.FilmSequel
		for _, sequel := range result {
			res, err := s.filmService.GetOne(ctx, strconv.Itoa(sequel.SequelID))

			if err != nil {
				return []model.FilmSequel{}, err
			}
			sequels = append(sequels, model.FilmSequel{
				Film: res,
			})
		}
		return sequels, nil
	}

	sequels, err := s.FetchSequels(ctx, filmId)

	if err != nil {
		return []model.FilmSequel{}, err
	}

	return sequels, nil
}

func (s *FilmSequelService) FetchSequels(ctx context.Context, filmId string) ([]model.FilmSequel, error) {
	apiKey := s.cfg.DB.ApiKey
	baseUrlForAllSequels := fmt.Sprintf(baseUrlForAllSequels, filmId)
	req, err := http.NewRequest("GET", baseUrlForAllSequels, nil)

	if err != nil {
		return []model.FilmSequel{}, sys.NewError(sys.ErrUnknown, err.Error())
	}

	req.Header.Add("X-API-KEY", apiKey)

	client := &http.Client{}
	resAllSequels, err := client.Do(req)
	if err != nil {
		return []model.FilmSequel{},
			sys.NewError(sys.ErrUnknown, err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resAllSequels.Body)

	if resAllSequels.StatusCode != http.StatusOK {
		return []model.FilmSequel{}, sys.NewError(sys.ErrUnknown, "Сиквелы не найдены")
	}

	bodyAllSequels, err := io.ReadAll(resAllSequels.Body)
	if err != nil {
		return []model.FilmSequel{},
			sys.NewError(sys.ErrUnknown, err.Error())
	}

	var externalSequels []dto.GetExternalSequelResponseDTO

	err = json.Unmarshal(bodyAllSequels, &externalSequels)

	if err != nil {
		return []model.FilmSequel{},
			sys.NewError(sys.ErrUnknown, err.Error())
	}

	var sequels []model.FilmSequel
	for _, sequel := range externalSequels {
		filmExist, err := s.filmService.GetOne(ctx, strconv.Itoa(sequel.FilmId))

		if err != nil {
			return []model.FilmSequel{}, err
		}

		filmIdInt, err := strconv.Atoi(filmId)

		if err != nil {
			return []model.FilmSequel{},
				sys.NewError(sys.ErrUnknown, err.Error())
		}

		err = s.filmSequelRepo.Save(ctx, filmIdInt, sequel.FilmId)
		if err != nil {
			return nil, sys.NewError(sys.ErrDatabaseFailure, err.Error())
		}

		sequels = append(sequels, model.FilmSequel{
			Film: filmExist,
		})
	}

	return sequels, nil
}
