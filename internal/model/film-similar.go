package model

type FilmSimilar struct {
	SimilarID int `json:"similarId"`
	FilmID    int `json:"filmId"`
	Film      Film
}
