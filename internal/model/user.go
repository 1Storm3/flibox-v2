package model

type User struct {
	ID            string  `json:"id"`
	NickName      string  `json:"nickName"`
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	Password      string  `json:"password"`
	Photo         *string `json:"photo"`
	Role          string  `json:"role"`
	VerifiedToken *string `json:"verifiedToken"`
	IsVerified    bool    `json:"isVerified"`
	IsBlocked     bool    `json:"isBlocked"`
	LastActivity  string  `json:"lastActivity"`
	UpdatedAt     string  `json:"updatedAt"`
	CreatedAt     string  `json:"createdAt"`
}
