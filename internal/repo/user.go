package repo

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"gorm.io/gorm"

	"github.com/1Storm3/flibox-api/database/postgres"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
)

type UserRepo struct {
	storage *postgres.Storage
}

func NewUserRepo(storage *postgres.Storage) *UserRepo {
	return &UserRepo{
		storage: storage,
	}
}

func (u *UserRepo) UpdateForVerify(ctx context.Context, userDTO dto.UserRepoDTO) (dto.UserRepoDTO, error) {
	var user dto.UserRepoDTO
	result := u.storage.DB().WithContext(ctx).Where("id = ?", userDTO.ID).Table("users").First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return dto.UserRepoDTO{},
			httperror.New(
				http.StatusNotFound,
				"Пользователь не найден",
			)
	}

	user.IsVerified = userDTO.IsVerified
	user.VerifiedToken = userDTO.VerifiedToken

	result = u.storage.DB().WithContext(ctx).Table("users").Save(&user)
	if result.Error != nil {
		return dto.UserRepoDTO{},
			httperror.New(
				http.StatusInternalServerError,
				result.Error.Error(),
			)
	}
	return user, nil
}

func (u *UserRepo) GetOneByNickName(ctx context.Context, nickName string) (dto.UserRepoDTO, error) {
	var user dto.UserRepoDTO
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
		Table("users").
		First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return dto.UserRepoDTO{},
			httperror.New(
				http.StatusNotFound,
				"Пользователь не найден",
			)
	} else if result.Error != nil {
		return dto.UserRepoDTO{},
			httperror.New(
				http.StatusInternalServerError,
				result.Error.Error(),
			)
	}
	return user, nil
}

func (u *UserRepo) GetOneById(_ context.Context, id string) (dto.UserRepoDTO, error) {
	var user dto.UserRepoDTO

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
		Table("users").
		First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return dto.UserRepoDTO{},
			httperror.New(
				http.StatusNotFound,
				"Пользователь не найден",
			)
	} else if result.Error != nil {
		return dto.UserRepoDTO{},
			httperror.New(
				http.StatusInternalServerError,
				result.Error.Error(),
			)
	}

	return user, nil
}

func (u *UserRepo) GetOneByEmail(ctx context.Context, email string) (dto.UserRepoDTO, error) {
	var user dto.UserRepoDTO

	result := u.storage.DB().WithContext(ctx).Where("email = ?", email).Table("users").First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return dto.UserRepoDTO{}, nil
	} else if result.Error != nil {
		return dto.UserRepoDTO{},
			httperror.New(
				http.StatusInternalServerError,
				result.Error.Error(),
			)
	}

	return user, nil
}

func (u *UserRepo) Create(ctx context.Context, user dto.UserRepoDTO) (dto.UserRepoDTO, error) {
	result := u.storage.DB().WithContext(ctx).Table("users").Create(&user)
	if result.Error != nil {
		return dto.UserRepoDTO{},
			httperror.New(
				http.StatusInternalServerError,
				result.Error.Error(),
			)
	}
	return user, nil
}

func (u *UserRepo) Update(ctx context.Context, userDTO dto.UserRepoDTO) (dto.UserRepoDTO, error) {
	tx := u.storage.DB().WithContext(ctx).Begin().Table("users")
	var user dto.UserRepoDTO
	if err := tx.Where("id = ?", userDTO.ID).First(&user).Error; err != nil {
		tx.Rollback()
		return dto.UserRepoDTO{}, httperror.New(http.StatusNotFound, "Пользователь не найден")
	}

	if err := tx.Model(&user).Updates(userDTO).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "duplicate key value") {
			return dto.UserRepoDTO{}, httperror.New(
				http.StatusConflict,
				"Пользователь с таким никнеймом или почтой уже существует",
			)
		}
		return dto.UserRepoDTO{}, httperror.New(http.StatusInternalServerError, err.Error())
	}

	tx.Commit()
	return user, nil
}
