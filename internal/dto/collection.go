package dto

import (
	"github.com/lib/pq"
	"time"
)

type CreateCollectionDTO struct {
	Name        string   `json:"name" validate:"required,min=4"`
	Description *string  `json:"description" validate:"omitempty,min=4"`
	CoverUrl    *string  `json:"coverUrl" validate:"omitempty,url"`
	Tags        []string `json:"tags" validate:"omitempty,min=1"`
}

type ResponseCollectionDTO struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	CoverUrl    *string  `json:"coverUrl"`
	User        User     `json:"author"`
	Tags        []string `json:"tags"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}

type CollectionRepoDTO struct {
	ID          string          `json:"id" gorm:"primaryKey;default:uuid_generate_v4()"`
	Name        string          `json:"name" gorm:"column:name"`
	Description *string         `json:"description" gorm:"column:description"`
	CoverUrl    *string         `json:"coverUrl" gorm:"column:cover_url"`
	Tags        *pq.StringArray `json:"tags" gorm:"type:text[];column:tags"`
	CreatedAt   time.Time       `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt   time.Time       `json:"updatedAt" gorm:"column:updated_at"`
	UserId      *string         `json:"userId" gorm:"column:user_id"`
	User        UserRepoDTO     `gorm:"foreignKey:UserId;references:ID"`
	Films       []FilmRepoDTO   `gorm:"many2many:collection_films;joinForeignKey:collection_id;JoinReferences:film_id"`
}

type UpdateCollectionDTO struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	CoverUrl    *string  `json:"coverUrl"`
	Tags        []string `json:"tags"`
}

type Film struct {
	ID              *int     `json:"kinopoiskId"`
	NameRU          *string  `json:"nameRu"`
	NameOriginal    *string  `json:"nameOriginal"`
	Type            *string  `json:"type"`
	Year            *int     `json:"year"`
	PosterURL       *string  `json:"posterUrl"`
	RatingKinopoisk *float64 `json:"ratingKinopoisk"`
}
