package dto

import "time"

type HistoryFilmsRepoDTO struct {
	ID        string      `json:"id" gorm:"primaryKey;default:uuid_generate_v4()"`
	UserID    string      `json:"userId" gorm:"column:user_id"`
	FilmID    int         `json:"filmId" gorm:"column:film_id"`
	CreatedAt time.Time   `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time   `json:"updatedAt" gorm:"column:updated_at"`
	Film      FilmRepoDTO `gorm:"foreignKey:FilmID"`
	User      UserRepoDTO `gorm:"foreignKey:UserID"`
}
