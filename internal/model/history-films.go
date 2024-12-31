package model

type HistoryFilms struct {
	ID        string `json:"id"`
	UserID    string `json:"userId" `
	FilmID    int    `json:"filmId" `
	CreatedAt string `json:"createdAt" `
	UpdatedAt string `json:"updatedAt" `
	Film      Film
	User      User
}
