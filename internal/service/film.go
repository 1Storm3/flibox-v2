package service

import (
	"context"
	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/mapper"
	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/lib/pq"
)

type FilmService struct {
	filmRepo        FilmRepo
	externalService controller.ExternalService
}

func NewFilmService(filmRepo FilmRepo, externalService controller.ExternalService) *FilmService {
	return &FilmService{
		filmRepo:        filmRepo,
		externalService: externalService,
	}
}

func (f *FilmService) GetOne(ctx context.Context, filmId string) (dto.ResponseFilmDTO, error) {
	result, err := f.filmRepo.GetOne(ctx, filmId)
	if err != nil {
		return dto.ResponseFilmDTO{}, err
	}

	if result.ID == nil {
		externalFilm, err := f.externalService.SetExternalRequest(filmId)
		if err != nil {
			return dto.ResponseFilmDTO{}, err
		}
		var genres []string
		for _, genre := range externalFilm.Genres {
			genres = append(genres, genre.Genre)
		}

		film := model.Film{
			ID:              externalFilm.ID,
			NameRU:          externalFilm.NameRU,
			NameOriginal:    externalFilm.NameOriginal,
			Year:            externalFilm.Year,
			PosterURL:       externalFilm.PosterURL,
			RatingKinopoisk: externalFilm.RatingKinopoisk,
			Description:     externalFilm.Description,
			LogoURL:         externalFilm.LogoURL,
			Type:            externalFilm.Type,
			Genres:          pq.StringArray(genres),
		}

		if err := f.filmRepo.Save(ctx, film); err != nil {
			return dto.ResponseFilmDTO{}, err
		}

		filmDTO := mapper.MapModelFilmToResponseDTO(film)

		return filmDTO, nil
	}

	filmDTO := mapper.MapModelFilmToResponseDTO(result)
	return filmDTO, nil
}

func (f *FilmService) Search(ctx context.Context, match string, genres []string, page int, pageSize int) ([]dto.SearchResponseDTO, int64, error) {
	films, totalRecords, err := f.filmRepo.Search(ctx, match, genres, page, pageSize)

	if err != nil {
		return []dto.SearchResponseDTO{}, 0, err
	}

	var filmsDTO []dto.SearchResponseDTO
	for _, film := range films {
		filmsDTO = append(filmsDTO, mapper.MapModelFilmToResponseSearchDTO(film))
	}

	return filmsDTO, totalRecords, nil
}

func (f *FilmService) GetOneByNameRu(ctx context.Context, nameRu string) (dto.ResponseFilmDTO, error) {
	result, err := f.filmRepo.GetOneByNameRu(ctx, nameRu)
	if err != nil {
		return dto.ResponseFilmDTO{}, err
	}
	return mapper.MapModelFilmToResponseDTO(result), nil
}
