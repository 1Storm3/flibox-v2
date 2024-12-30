package dto

type CreateCollectionFilmDTO struct {
	FilmID int `json:"filmId" validate:"required"`
}

type DeleteCollectionFilmDTO struct {
	FilmID int `json:"filmId" validate:"required"`
}

type FilmsByCollectionIdDTO struct {
	CollectionID string `json:"collectionId" validate:"required"`
	Films        []Film `json:"films" validate:"required"`
}
