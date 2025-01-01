package repo

import (
	"context"

	"github.com/1Storm3/flibox-api/database/postgres"
	"github.com/1Storm3/flibox-api/internal/dto"
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
	err := u.storage.DB().WithContext(ctx).Where("id = ?", userDTO.ID).Table("users").First(&user).Error
	if err != nil {
		return dto.UserRepoDTO{}, err
	}

	user.IsVerified = userDTO.IsVerified
	user.VerifiedToken = userDTO.VerifiedToken

	err = u.storage.DB().WithContext(ctx).Table("users").Save(&user).Error
	if err != nil {
		return dto.UserRepoDTO{}, err
	}
	return user, nil
}

func (u *UserRepo) GetOneByNickName(ctx context.Context, nickName string) (dto.UserRepoDTO, error) {
	var user dto.UserRepoDTO
	err := u.storage.DB().WithContext(ctx).
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
		First(&user).Error
	if err != nil {
		return dto.UserRepoDTO{}, err
	}
	return user, nil
}

func (u *UserRepo) GetOneById(_ context.Context, id string) (dto.UserRepoDTO, error) {
	var user dto.UserRepoDTO

	err := u.storage.DB().
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
		First(&user).Error
	if err != nil {
		return dto.UserRepoDTO{}, err
	}

	return user, nil
}

func (u *UserRepo) GetOneByEmail(ctx context.Context, email string) (dto.UserRepoDTO, error) {
	var user dto.UserRepoDTO

	err := u.storage.DB().WithContext(ctx).Where("email = ?", email).Table("users").First(&user).Error
	if err != nil {
		return dto.UserRepoDTO{}, err
	}

	return user, nil
}

func (u *UserRepo) Create(ctx context.Context, user dto.UserRepoDTO) (dto.UserRepoDTO, error) {
	err := u.storage.DB().WithContext(ctx).Table("users").Create(&user).Error
	if err != nil {
		return dto.UserRepoDTO{}, err
	}
	return user, nil
}

func (u *UserRepo) Update(ctx context.Context, userDTO dto.UserRepoDTO) (dto.UserRepoDTO, error) {
	tx := u.storage.DB().WithContext(ctx).Begin().Table("users")
	var user dto.UserRepoDTO
	if err := tx.Where("id = ?", userDTO.ID).First(&user).Error; err != nil {
		tx.Rollback()
		return dto.UserRepoDTO{}, err
	}

	if err := tx.Model(&user).Updates(userDTO).Error; err != nil {
		tx.Rollback()
		return dto.UserRepoDTO{}, err
	}

	tx.Commit()
	return user, nil
}
