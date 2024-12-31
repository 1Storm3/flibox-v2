package model

type Collection struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	CoverUrl    *string  `json:"coverUrl"`
	Tags        []string `json:"tags"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
	UserId      *string  `json:"userId"`
	User        User
	Films       []Film
}
