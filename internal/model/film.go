package model

type Film struct {
	ID              *int     `json:"kinopoiskId"`
	NameRU          *string  `json:"nameRu"`
	NameOriginal    *string  `json:"nameOriginal"`
	Year            *int     `json:"year"`
	PosterURL       *string  `json:"posterUrl"`
	RatingKinopoisk *float64 `json:"ratingKinopoisk"`
	Description     *string  `json:"description"`
	LogoURL         *string  `json:"logoUrl"`
	Type            *string  `json:"type"`
	CoverURL        *string  `json:"coverUrl"`
	TrailerURL      *string  `json:"trailerUrl"`
	Sequels         []*Film
	Similars        []*Film
	Genres          []string `json:"genres"`
}
