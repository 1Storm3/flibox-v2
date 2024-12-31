package model

import "github.com/1Storm3/flibox-api/internal/dto"

type UserFilm struct {
	UserID string           `json:"userId"`
	FilmID int              `json:"filmId"`
	Type   dto.TypeUserFilm `json:"type"`
	Film   Film
}
