package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/1Storm3/flibox-api/internal/config"
	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"io"
	"net/http"
	"strconv"
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

func (s *FilmSequelService) GetAll(ctx context.Context, filmId string) ([]dto.ResponseFilmDTO, error) {
	result, err := s.filmSequelRepo.GetAll(ctx, filmId)

	if err != nil {
		return []dto.ResponseFilmDTO{}, err
	}
	if len(result) > 0 {
		var sequels []dto.ResponseFilmDTO
		for _, sequel := range result {
			res, err := s.filmService.GetOne(ctx, strconv.Itoa(sequel.SequelID))

			if err != nil {
				return []dto.ResponseFilmDTO{}, err
			}
			sequels = append(sequels, res)
		}
		return sequels, nil
	}

	sequels, err := s.FetchSequels(ctx, filmId)

	if err != nil {
		return []dto.ResponseFilmDTO{}, err
	}

	return sequels, nil
}

func (s *FilmSequelService) FetchSequels(ctx context.Context, filmId string) ([]dto.ResponseFilmDTO, error) {
	apiKey := s.cfg.DB.ApiKey
	baseUrlForAllSequels := fmt.Sprintf(baseUrlForAllSequels, filmId)
	req, err := http.NewRequest("GET", baseUrlForAllSequels, nil)

	if err != nil {
		return []dto.ResponseFilmDTO{}, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	req.Header.Add("X-API-KEY", apiKey)

	client := &http.Client{}
	resAllSequels, err := client.Do(req)
	if err != nil {
		return []dto.ResponseFilmDTO{},
			httperror.New(
				http.StatusInternalServerError,
				err.Error(),
			)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resAllSequels.Body)

	if resAllSequels.StatusCode != http.StatusOK {
		return []dto.ResponseFilmDTO{}, httperror.New(
			http.StatusNotFound,
			"Сиквелы не найдены",
		)
	}

	bodyAllSequels, err := io.ReadAll(resAllSequels.Body)
	if err != nil {
		return []dto.ResponseFilmDTO{},
			httperror.New(
				http.StatusInternalServerError,
				err.Error(),
			)
	}

	var externalSequels []dto.GetExternalSequelResponseDTO

	err = json.Unmarshal(bodyAllSequels, &externalSequels)
	var sequels []dto.ResponseFilmDTO
	for _, sequel := range externalSequels {
		filmExist, err := s.filmService.GetOne(ctx, strconv.Itoa(sequel.FilmId))

		if err != nil {
			return []dto.ResponseFilmDTO{}, err
		}

		filmIdInt, err := strconv.Atoi(filmId)

		if err != nil {
			return []dto.ResponseFilmDTO{},
				httperror.New(
					http.StatusInternalServerError,
					err.Error(),
				)
		}

		err = s.filmSequelRepo.Save(ctx, filmIdInt, sequel.FilmId)
		if err != nil {
			return nil, err
		}

		sequels = append(sequels, filmExist)
	}

	if err != nil {
		return []dto.ResponseFilmDTO{},
			httperror.New(
				http.StatusInternalServerError,
				err.Error(),
			)
	}
	return sequels, nil
}
