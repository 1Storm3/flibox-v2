package mapper

import "github.com/1Storm3/flibox-api/internal/dto"

func MapExternalFilmDTOToFilmRepoDTO(external dto.GetExternalFilmDTO) dto.FilmRepoDTO {
	return dto.FilmRepoDTO{
		ID:              external.ID,
		NameRU:          external.NameRU,
		Description:     external.Description,
		Year:            external.Year,
		NameOriginal:    external.NameOriginal,
		RatingKinopoisk: external.RatingKinopoisk,
		LogoURL:         external.LogoURL,
		PosterURL:       external.PosterURL,
		Type:            external.Type,
		Genres:          MapExternalGenreToString(external.Genres),
	}
}

func MapExternalGenreToString(genre []dto.Genre) []string {
	var genres []string
	for _, g := range genre {
		genres = append(genres, g.Genre)
	}
	return genres

}
