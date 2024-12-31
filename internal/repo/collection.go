package repo

import (
	"context"
	"errors"
	"net/http"

	"gorm.io/gorm"

	"github.com/1Storm3/flibox-api/database/postgres"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
)

type CollectionRepo struct {
	storage *postgres.Storage
}

func NewCollectionRepo(storage *postgres.Storage) *CollectionRepo {
	return &CollectionRepo{
		storage: storage,
	}
}

func (c *CollectionRepo) GetAllMy(ctx context.Context, page, pageSize int, userID string) ([]dto.CollectionRepoDTO, int64, error) {
	var collections []dto.CollectionRepoDTO
	var totalRecords int64
	err := c.storage.DB().WithContext(ctx).Model(&dto.CollectionRepoDTO{}).Where("user_id = ?", userID).Table("collections").Count(&totalRecords).Error
	if err != nil {
		return []dto.CollectionRepoDTO{}, 0, err
	}
	err = c.storage.DB().WithContext(ctx).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Table("users")
		}).
		Where("user_id = ?", userID).Table("collections").Order("created_at DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&collections).Error
	if err != nil {
		return []dto.CollectionRepoDTO{}, 0, err
	}
	return collections, totalRecords, nil
}

func (c *CollectionRepo) GetAll(ctx context.Context, page, pageSize int) ([]dto.CollectionRepoDTO, int64, error) {
	var collections []dto.CollectionRepoDTO
	var totalRecords int64
	err := c.storage.DB().WithContext(ctx).Table("collections").Model(&dto.CollectionRepoDTO{}).Count(&totalRecords).Error
	if err != nil {
		return []dto.CollectionRepoDTO{}, 0, err
	}
	err = c.storage.DB().WithContext(ctx).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Table("users")
		}).
		Table("collections").Order("created_at DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&collections).Error
	if err != nil {
		return []dto.CollectionRepoDTO{}, 0, err
	}
	return collections, totalRecords, nil
}

func (c *CollectionRepo) Delete(ctx context.Context, collectionId string) error {
	err := c.storage.DB().WithContext(ctx).
		Where("id = ?", collectionId).
		Table("collections").
		Delete(&dto.CollectionRepoDTO{}).Error

	if err != nil {
		return httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	return nil
}

func (c *CollectionRepo) Update(ctx context.Context, collection dto.CollectionRepoDTO) (dto.CollectionRepoDTO, error) {
	err := c.storage.DB().WithContext(ctx).Table("collections").Model(&collection).Where("id = ?", collection.ID).Updates(collection).Error
	if err != nil {
		return dto.CollectionRepoDTO{}, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	err = c.storage.DB().WithContext(ctx).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Table("users")
		}).
		Table("collections").First(&collection, "id = ?", collection.ID).Error
	if err != nil {
		return dto.CollectionRepoDTO{}, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	return collection, nil
}

func (c *CollectionRepo) GetOne(ctx context.Context, collectionId string) (dto.CollectionRepoDTO, error) {
	var collection dto.CollectionRepoDTO
	err := c.storage.DB().WithContext(ctx).
		Where("id = ?", collectionId).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Table("users")
		}).
		Table("collections").
		First(&collection).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.CollectionRepoDTO{}, httperror.New(
				http.StatusNotFound,
				"Коллекция не найдена",
			)
		}
		return dto.CollectionRepoDTO{}, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	return collection, nil
}

func (c *CollectionRepo) Create(ctx context.Context, collection dto.CollectionRepoDTO) (dto.CollectionRepoDTO, error) {
	result := c.storage.DB().WithContext(ctx).Table("collections").Create(&collection)

	if result.Error != nil {
		return dto.CollectionRepoDTO{}, result.Error
	}
	return collection, nil
}
