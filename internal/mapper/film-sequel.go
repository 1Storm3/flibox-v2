package mapper

import "github.com/1Storm3/flibox-api/internal/model"

func MapModelFilmSequelToModelFilm(filmSequel model.FilmSequel) model.Film {
	return model.Film{
		ID:              filmSequel.Film.ID,
		NameRU:          filmSequel.Film.NameRU,
		Description:     filmSequel.Film.Description,
		Year:            filmSequel.Film.Year,
		NameOriginal:    filmSequel.Film.NameOriginal,
		RatingKinopoisk: filmSequel.Film.RatingKinopoisk,
		LogoURL:         filmSequel.Film.LogoURL,
		PosterURL:       filmSequel.Film.PosterURL,
		Type:            filmSequel.Film.Type,
		Genres:          filmSequel.Film.Genres,
	}
}
