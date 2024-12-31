package dto

type CreateCollectionFilmDTO struct {
	FilmID int `json:"filmId" validate:"required"`
}

type DeleteCollectionFilmDTO struct {
	FilmID int `json:"filmId" validate:"required"`
}

type FilmsByCollectionIdDTO struct {
	CollectionID string            `json:"collectionId" validate:"required"`
	Films        []ResponseFilmDTO `json:"films" validate:"required"`
}

type CollectionFilmRepoDTO struct {
	CollectionID string            `json:"collectionId" gorm:"column:collection_id"`
	FilmID       int               `json:"filmId" gorm:"column:film_id"`
	Film         FilmRepoDTO       `gorm:"foreignKey:FilmID;references:ID"`
	Collection   CollectionRepoDTO `gorm:"foreignKey:CollectionID;references:ID"`
}
