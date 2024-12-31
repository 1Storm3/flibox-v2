package model

type Comment struct {
	ID        string  `json:"id"`
	FilmID    int     `json:"filmId"`
	UserID    string  `json:"userId"`
	ParentID  *string `json:"parentId"`
	Content   *string `json:"content"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
	User      User
	Parent    *Comment
}
