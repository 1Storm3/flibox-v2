package dto

import "time"

type CreateCommentDTO struct {
	Content  *string `json:"content" validate:"required,min=1"`
	UserID   string  `json:"userId"`
	FilmID   int     `json:"filmId" validate:"required,min=36"`
	ParentID *string `json:"parentId" validate:"omitempty"`
}

type CommentRepoDTO struct {
	ID        string          `json:"id" gorm:"primaryKey;default:uuid_generate_v4()"`
	FilmID    int             `json:"filmId" gorm:"column:film_id"`
	UserID    string          `json:"userId" gorm:"column:user_id"`
	ParentID  *string         `json:"parentId" gorm:"column:parent_id"`
	Content   *string         `json:"content" gorm:"column:content"`
	CreatedAt time.Time       `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time       `json:"updatedAt" gorm:"column:updated_at"`
	User      UserRepoDTO     `gorm:"foreignKey:UserID;references:ID"`
	Parent    *CommentRepoDTO `gorm:"foreignKey:ParentID;references:ID;onDelete:CASCADE"`
}

type CommentResponseDTO struct {
	ID        string  `json:"id"`
	Content   *string `json:"content"`
	FilmID    int     `json:"filmId"`
	User      User    `json:"author"`
	ParentID  *string `json:"parentId"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

type UpdateCommentDTO struct {
	ID      string  `json:"id"`
	Content *string `json:"content" validate:"required,min=1"`
}

type User struct {
	ID       string  `json:"id"`
	NickName string  `json:"nickName"`
	Photo    *string `json:"photo"`
}
