package mapper

import (
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/model"
)

func MapManyFilmsRepoDTOToManyFilmsModel(dto []dto.FilmRepoDTO) []model.Film {
	films := make([]model.Film, len(dto))
	for i, film := range dto {
		films[i] = MapFilmRepoDTOToFilmModel(film)
	}
	return films
}

func MapModelFilmToResponseDTO(film model.Film) dto.ResponseFilmDTO {
	return dto.ResponseFilmDTO{
		ID:              film.ID,
		NameRU:          film.NameRU,
		NameOriginal:    film.NameOriginal,
		Year:            film.Year,
		RatingKinopoisk: film.RatingKinopoisk,
		PosterURL:       film.PosterURL,
		Description:     film.Description,
		LogoURL:         film.LogoURL,
		Type:            film.Type,
		Genres:          film.Genres,
		CoverURL:        film.CoverURL,
		TrailerURL:      film.TrailerURL,
	}
}

func MapModelFilmToResponseSearchDTO(film model.Film) dto.FilmSearchResponseDTO {
	return dto.FilmSearchResponseDTO{
		ID:              film.ID,
		NameRU:          film.NameRU,
		NameOriginal:    film.NameOriginal,
		Type:            film.Type,
		Year:            film.Year,
		RatingKinopoisk: film.RatingKinopoisk,
		PosterURL:       film.PosterURL,
	}
}

func MapFilmRepoDTOToFilmModel(film dto.FilmRepoDTO) model.Film {
	return model.Film{
		ID:              film.ID,
		NameRU:          film.NameRU,
		NameOriginal:    film.NameOriginal,
		Type:            film.Type,
		Year:            film.Year,
		PosterURL:       film.PosterURL,
		RatingKinopoisk: film.RatingKinopoisk,
		Description:     film.Description,
		LogoURL:         film.LogoURL,
		CoverURL:        film.CoverURL,
		TrailerURL:      film.TrailerURL,
		Genres:          film.Genres,
	}
}
