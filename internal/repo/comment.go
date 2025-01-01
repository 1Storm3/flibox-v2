package repo

import (
	"context"

	"gorm.io/gorm"

	"github.com/1Storm3/flibox-api/database/postgres"
	"github.com/1Storm3/flibox-api/internal/dto"
)

type CommentRepo struct {
	storage *postgres.Storage
}

func NewCommentRepo(storage *postgres.Storage) *CommentRepo {
	return &CommentRepo{
		storage: storage,
	}
}

func (c *CommentRepo) Create(ctx context.Context, comment dto.CommentRepoDTO) (dto.CommentRepoDTO, error) {
	result := c.storage.DB().WithContext(ctx).Table("comments").Create(&comment)

	if result.Error != nil {
		return dto.CommentRepoDTO{}, result.Error
	}

	err := c.storage.DB().WithContext(ctx).Table("comments").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Table("users")
		}).
		First(&comment, "id = ?", comment.ID).Error
	if err != nil {
		return dto.CommentRepoDTO{}, err
	}

	return comment, nil
}

func (c *CommentRepo) GetCountByParentId(ctx context.Context, parentId string) (int64, error) {
	var count int64
	err := c.storage.DB().WithContext(ctx).
		Table("comments").
		Model(&dto.CommentRepoDTO{}).
		Where("parent_id = ?", parentId).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *CommentRepo) GetOne(ctx context.Context, commentID string) (dto.CommentRepoDTO, error) {
	var comment dto.CommentRepoDTO

	err := c.storage.DB().WithContext(ctx).Where("id = ?", commentID).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Table("users")
		}).
		Table("comments").First(&comment).Error
	if err != nil {
		return dto.CommentRepoDTO{}, err
	}

	return comment, nil
}

func (c *CommentRepo) Delete(ctx context.Context, commentID string) error {
	err := c.storage.DB().WithContext(ctx).
		Table("comments").
		Where("id = ?", commentID).
		Delete(&dto.CommentRepoDTO{}).
		Error

	if err != nil {
		return err
	}
	return nil
}

func (c *CommentRepo) GetAllByFilmId(ctx context.Context, filmID string, page, pageSize int) ([]dto.CommentRepoDTO, int64, error) {
	var comments []dto.CommentRepoDTO
	var totalRecords int64

	offset := (page - 1) * pageSize

	err := c.storage.DB().WithContext(ctx).Table("comments").Model(&dto.CommentRepoDTO{}).Where("film_id = ?", filmID).Count(&totalRecords).Error
	if err != nil {
		return []dto.CommentRepoDTO{}, 0, err
	}

	err = c.storage.DB().WithContext(ctx).
		Where("film_id = ?", filmID).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Table("users")
		}).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Table("comments").
		Find(&comments).Error

	if err != nil {
		return []dto.CommentRepoDTO{}, 0, err
	}

	return comments, totalRecords, nil
}

func (c *CommentRepo) Update(ctx context.Context, commentDTO dto.CommentRepoDTO, commentID string) (dto.CommentRepoDTO, error) {
	var comment dto.CommentRepoDTO
	if commentDTO.Content == nil {
		err := c.storage.DB().WithContext(ctx).Table("comments").Model(&comment).Where("id = ?", commentID).Update("content", nil).Error
		if err != nil {
			return dto.CommentRepoDTO{}, err
		}
	} else {
		err := c.storage.DB().WithContext(ctx).Table("comments").Model(&comment).Where("id = ?", commentID).Updates(commentDTO).Error
		if err != nil {
			return dto.CommentRepoDTO{}, err
		}
	}
	err := c.storage.DB().WithContext(ctx).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Table("users")
		}).
		Table("comments").First(&comment, "id = ?", commentID).Error
	if err != nil {
		return dto.CommentRepoDTO{}, err
	}

	return comment, nil
}
