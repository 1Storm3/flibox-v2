package model

type FilmSequel struct {
	SequelID int `json:"sequelId"`
	FilmID   int `json:"filmId"`
	Film     Film
}
