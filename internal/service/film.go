package service

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/mapper"
	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/1Storm3/flibox-api/pkg/sys"
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

func (f *FilmService) GetOne(ctx context.Context, filmId string) (model.Film, error) {
	result, err := f.filmRepo.GetOne(ctx, filmId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Film{}, nil
		}
		return model.Film{}, sys.NewError(sys.ErrDatabaseFailure, err.Error())
	}

	if result.ID == nil {
		externalFilm, err := f.externalService.SetExternalRequest(filmId)
		if err != nil {
			return model.Film{}, err
		}
		var genres []string
		for _, genre := range externalFilm.Genres {
			genres = append(genres, genre.Genre)
		}

		film := mapper.MapExternalFilmDTOToFilmRepoDTO(externalFilm)

		if err := f.filmRepo.Save(ctx, film); err != nil {
			return model.Film{}, sys.NewError(sys.ErrDatabaseFailure, err.Error())
		}

		filmDTO := mapper.MapFilmRepoDTOToFilmModel(film)

		return filmDTO, nil
	}

	filmDTO := mapper.MapFilmRepoDTOToFilmModel(result)

	return filmDTO, nil
}

func (f *FilmService) Search(ctx context.Context, match string, genres []string, page int, pageSize int) ([]model.Film, int64, error) {
	films, totalRecords, err := f.filmRepo.Search(ctx, match, genres, page, pageSize)

	if err != nil {
		return []model.Film{}, 0, sys.NewError(sys.ErrDatabaseFailure, err.Error())
	}

	var filmsDTO []model.Film
	for _, film := range films {
		filmsDTO = append(filmsDTO, mapper.MapFilmRepoDTOToFilmModel(film))
	}

	return filmsDTO, totalRecords, nil
}

func (f *FilmService) GetOneByNameRu(ctx context.Context, nameRu string) (model.Film, error) {
	result, err := f.filmRepo.GetOneByNameRu(ctx, nameRu)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Film{}, nil
		}
		return model.Film{}, sys.NewError(sys.ErrDatabaseFailure, err.Error())
	}
	return mapper.MapFilmRepoDTOToFilmModel(result), nil
}
