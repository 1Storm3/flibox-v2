package dto

import "time"

type CreateUserDTO struct {
	Name     string  `json:"name" validate:"required,min=4"`
	NickName string  `json:"nickName" validate:"required,min=4"`
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=6"`
	Photo    *string `json:"photo" validate:"omitempty,url"`
}

type UserRepoDTO struct {
	ID            string    `json:"id" gorm:"column:id;primaryKey;default:uuid_generate_v4()"`
	NickName      string    `json:"nickName" gorm:"column:nick_name"`
	Name          string    `json:"name" gorm:"column:name"`
	Email         string    `json:"email" gorm:"column:email"`
	Password      string    `json:"password" gorm:"column:password"`
	Photo         *string   `json:"photo" gorm:"column:photo"`
	Role          string    `json:"role" gorm:"column:role"`
	VerifiedToken *string   `json:"verifiedToken" gorm:"column:verified_token"`
	IsVerified    bool      `json:"isVerified" gorm:"column:is_verified"`
	IsBlocked     bool      `json:"isBlocked" gorm:"column:is_blocked"`
	LastActivity  time.Time `json:"lastActivity" gorm:"column:last_activity"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"column:updated_at"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:created_at"`
}

type ResponseUserDTO struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	NickName   string    `json:"nickName"`
	Email      string    `json:"email"`
	Photo      *string   `json:"photo"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	IsVerified bool      `json:"isVerified"`
}

type UpdateUserDTO struct {
	ID           string  `json:"id" validate:"required"`
	NickName     *string `json:"nickName" validate:"omitempty,min=4"`
	Name         *string `json:"name" validate:"omitempty,min=4"`
	Email        *string `json:"email" validate:"omitempty,email"`
	LastActivity *string `json:"lastActivity" validate:"omitempty"`
	Photo        *string `json:"photo" validate:"omitempty,url"`
}

type UpdateForVerifyDTO struct {
	ID            string  `json:"id" validate:"required"`
	IsVerified    bool    `json:"isVerified" validate:"required"`
	VerifiedToken *string `json:"verifiedToken" validate:"required"`
}
