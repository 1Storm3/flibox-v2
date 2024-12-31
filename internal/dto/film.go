package dto

import "github.com/lib/pq"

type FilmSearchResponseDTO struct {
	ID              *int     `json:"kinopoiskId"`
	NameRU          *string  `json:"nameRu"`
	Type            *string  `json:"type"`
	NameOriginal    *string  `json:"nameOriginal"`
	Year            *int     `json:"year"`
	RatingKinopoisk *float64 `json:"ratingKinopoisk" gorm:"column:rating_kinopoisk"`
	PosterURL       *string  `json:"posterUrl"`
}

type FilmRepoDTO struct {
	ID              *int           `json:"kinopoiskId" gorm:"column:id"`
	NameRU          *string        `json:"nameRu" gorm:"column:name_ru"`
	NameOriginal    *string        `json:"nameOriginal" gorm:"column:name_original"`
	Year            *int           `json:"year" gorm:"column:year"`
	PosterURL       *string        `json:"posterUrl" gorm:"column:poster_url"`
	RatingKinopoisk *float64       `json:"ratingKinopoisk" gorm:"column:rating_kinopoisk"`
	Description     *string        `json:"description" gorm:"column:description"`
	LogoURL         *string        `json:"logoUrl" gorm:"column:logo_url"`
	Type            *string        `json:"type" gorm:"column:type"`
	CoverURL        *string        `json:"coverUrl" gorm:"column:cover_url"`
	TrailerURL      *string        `json:"trailerUrl" gorm:"column:trailer_url"`
	Sequels         []*FilmRepoDTO `gorm:"many2many:film_sequels;joinForeignKey:film_id;JoinReferences:sequel_id"`
	Similars        []*FilmRepoDTO `gorm:"many2many:film_similars;joinForeignKey:film_id;JoinReferences:similar_id"`
	Genres          pq.StringArray `json:"genres" gorm:"type:text[];column:genres"`
}

type ResponseFilmDTO struct {
	ID              *int           `json:"kinopoiskId"`
	NameRU          *string        `json:"nameRu"`
	NameOriginal    *string        `json:"nameOriginal"`
	Year            *int           `json:"year"`
	RatingKinopoisk *float64       `json:"ratingKinopoisk"`
	PosterURL       *string        `json:"posterUrl"`
	Description     *string        `json:"description"`
	LogoURL         *string        `json:"logoUrl"`
	Type            *string        `json:"type"`
	CoverURL        *string        `json:"coverUrl"`
	TrailerURL      *string        `json:"trailerUrl"`
	Genres          pq.StringArray `json:"genres"`
}
