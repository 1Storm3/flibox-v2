package dto

type GetExternalSequelResponseDTO struct {
	FilmId       int     `json:"filmId"`
	NameRu       *string `json:"nameRu"`
	NameOriginal *string `json:"nameOriginal"`
	PosterUrl    *string `json:"posterUrl"`
}

type FilmSequelRepoDTO struct {
	SequelID int         `json:"sequelId" gorm:"column:sequel_id"`
	FilmID   int         `json:"filmId" gorm:"column:film_id"`
	Film     FilmRepoDTO `gorm:"foreignKey:FilmID;references:ID"`
}
