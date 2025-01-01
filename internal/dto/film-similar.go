package dto

type GetExternalSimilarDTO struct {
	FilmId       int     `json:"filmId"`
	NameRu       *string `json:"nameRu"`
	NameOriginal *string `json:"nameOriginal"`
	PosterUrl    *string `json:"posterUrl"`
}

type FilmSimilarRepoDTO struct {
	FilmID    int         `json:"filmId" gorm:"column:film_id"`
	SimilarID int         `json:"similarId" gorm:"column:similar_id"`
	Film      FilmRepoDTO `gorm:"foreignKey:FilmID;references:ID"`
}
