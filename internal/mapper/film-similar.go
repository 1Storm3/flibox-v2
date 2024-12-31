package mapper

import "github.com/1Storm3/flibox-api/internal/model"

func MapModelFilmSimilarToModelFilm(filmSimilar model.FilmSimilar) model.Film {
	return model.Film{
		ID:              filmSimilar.Film.ID,
		NameRU:          filmSimilar.Film.NameRU,
		Description:     filmSimilar.Film.Description,
		Year:            filmSimilar.Film.Year,
		NameOriginal:    filmSimilar.Film.NameOriginal,
		RatingKinopoisk: filmSimilar.Film.RatingKinopoisk,
		LogoURL:         filmSimilar.Film.LogoURL,
		PosterURL:       filmSimilar.Film.PosterURL,
		Type:            filmSimilar.Film.Type,
		Genres:          filmSimilar.Film.Genres,
	}
}
