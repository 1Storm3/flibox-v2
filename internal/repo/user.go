package repo

import (
	"context"
	"errors"
	"github.com/1Storm3/flibox-api/database/postgres"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type UserRepo struct {
	storage *postgres.Storage
}

func NewUserRepo(storage *postgres.Storage) *UserRepo {
	return &UserRepo{
		storage: storage,
	}
}

func (u *UserRepo) UpdateForVerify(ctx context.Context, userDTO dto.UpdateForVerifyDTO) (model.User, error) {
	var user model.User
	result := u.storage.DB().WithContext(ctx).Where("id = ?", userDTO.ID).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return model.User{},
			httperror.New(
				http.StatusNotFound,
				"Пользователь не найден",
			)
	}

	user.IsVerified = userDTO.IsVerified
	user.VerifiedToken = userDTO.VerifiedToken

	result = u.storage.DB().WithContext(ctx).Save(&user)
	if result.Error != nil {
		return model.User{},
			httperror.New(
				http.StatusInternalServerError,
				result.Error.Error(),
			)
	}
	return user, nil
}

func (u *UserRepo) GetOneByNickName(ctx context.Context, nickName string) (model.User, error) {
	var user model.User
	result := u.storage.DB().WithContext(ctx).
		Select("id",
			"nick_name",
			"name",
			"email",
			"photo",
			"role",
			"is_verified",
			"updated_at",
			"created_at",
		).
		Where("nick_name = ?", nickName).
		First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return model.User{},
			httperror.New(
				http.StatusNotFound,
				"Пользователь не найден",
			)
	} else if result.Error != nil {
		return model.User{},
			httperror.New(
				http.StatusInternalServerError,
				result.Error.Error(),
			)
	}
	return user, nil
}

func (u *UserRepo) GetOneById(ctx context.Context, id string) (model.User, error) {
	var user model.User

	result := u.storage.DB().
		Select("id",
			"nick_name",
			"name",
			"email",
			"photo",
			"role",
			"is_verified",
			"is_blocked",
			"updated_at",
			"created_at",
		).
		Where("id = ?", id).
		First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return model.User{},
			httperror.New(
				http.StatusNotFound,
				"Пользователь не найден",
			)
	} else if result.Error != nil {
		return model.User{},
			httperror.New(
				http.StatusInternalServerError,
				result.Error.Error(),
			)
	}

	return user, nil
}

func (u *UserRepo) GetOneByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User

	result := u.storage.DB().WithContext(ctx).Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return model.User{}, nil
	} else if result.Error != nil {
		return model.User{},
			httperror.New(
				http.StatusInternalServerError,
				result.Error.Error(),
			)
	}

	return user, nil
}

func (u *UserRepo) Create(ctx context.Context, user model.User) (model.User, error) {
	result := u.storage.DB().WithContext(ctx).Create(&user)
	if result.Error != nil {
		return model.User{},
			httperror.New(
				http.StatusInternalServerError,
				result.Error.Error(),
			)
	}
	return user, nil
}

func (u *UserRepo) Update(ctx context.Context, userDTO dto.UpdateUserDTO) (model.User, error) {
	tx := u.storage.DB().WithContext(ctx).Begin()
	var user model.User
	if err := tx.Where("id = ?", userDTO.ID).First(&user).Error; err != nil {
		tx.Rollback()
		return model.User{}, httperror.New(http.StatusNotFound, "Пользователь не найден")
	}

	if err := tx.Model(&user).Updates(userDTO).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "duplicate key value") {
			return model.User{}, httperror.New(
				http.StatusConflict,
				"Пользователь с таким никнеймом или почтой уже существует",
			)
		}
		return model.User{}, httperror.New(http.StatusInternalServerError, err.Error())
	}

	tx.Commit()
	return user, nil
}
