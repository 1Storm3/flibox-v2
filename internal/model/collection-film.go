package model

type CollectionFilm struct {
	CollectionID string `json:"collectionId"`
	FilmID       int    `json:"filmId"`
	Film         Film
	Collection   Collection
}
