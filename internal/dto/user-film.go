package dto

type GetUserFilmResponseDTO struct {
	UserID string       `json:"userId"`
	FilmID int          `json:"filmId"`
	Type   TypeUserFilm `json:"type"`
	Film   ResponseFilmDTO
}

type Params struct {
	UserID string
	FilmID string
	Type   TypeUserFilm
}

type UserFilmRepoDTO struct {
	UserID string       `json:"userId" gorm:"column:user_id"`
	FilmID int          `json:"filmId" gorm:"column:film_id"`
	Film   FilmRepoDTO  `gorm:"foreignKey:FilmID;references:ID"`
	Type   TypeUserFilm `json:"type" gorm:"column:type"`
}

type TypeUserFilm string

const (
	TypeUserFavourite TypeUserFilm = "favourite"
	TypeUserRecommend TypeUserFilm = "recommend"
)
