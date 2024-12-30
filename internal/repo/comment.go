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

type CommentRepo struct {
	storage *postgres.Storage
}

func NewCommentRepo(storage *postgres.Storage) *CommentRepo {
	return &CommentRepo{
		storage: storage,
	}
}

func (c *CommentRepo) Create(ctx context.Context, comment model.Comment) (model.Comment, error) {
	result := c.storage.DB().WithContext(ctx).Create(&comment)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "violates foreign key") {
			return model.Comment{}, httperror.New(
				http.StatusConflict,
				"Родительского комментария не существует с таким ID",
			)
		}
		return model.Comment{}, httperror.New(
			http.StatusInternalServerError,
			result.Error.Error(),
		)
	}

	err := c.storage.DB().WithContext(ctx).Preload("User").First(&comment, "id = ?", comment.ID).Error
	if err != nil {
		return model.Comment{}, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	return comment, nil
}

func (c *CommentRepo) GetCountByParentId(ctx context.Context, parentId string) (int64, error) {
	var count int64
	err := c.storage.DB().WithContext(ctx).
		Model(&model.Comment{}).
		Where("parent_id = ?", parentId).
		Count(&count).Error

	if err != nil {
		return 0, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	return count, nil
}

func (c *CommentRepo) GetOne(ctx context.Context, commentID string) (model.Comment, error) {
	var comment model.Comment

	err := c.storage.DB().WithContext(ctx).Where("id = ?", commentID).Preload("User").First(&comment).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Comment{}, httperror.New(
			http.StatusNotFound,
			"Комментарий не найден",
		)
	} else if err != nil {
		return model.Comment{}, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	return comment, nil
}

func (c *CommentRepo) Delete(ctx context.Context, commentID string) error {
	err := c.storage.DB().WithContext(ctx).
		Where("id = ?", commentID).
		Delete(&model.Comment{}).
		Error

	if err != nil {
		return httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	return nil
}

func (c *CommentRepo) GetAllByFilmId(ctx context.Context, filmID string, page, pageSize int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var totalRecords int64

	offset := (page - 1) * pageSize

	err := c.storage.DB().WithContext(ctx).Model(&model.Comment{}).Where("film_id = ?", filmID).Count(&totalRecords).Error
	if err != nil {
		return []model.Comment{}, 0, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	err = c.storage.DB().WithContext(ctx).
		Where("film_id = ?", filmID).
		Preload("User").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&comments).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []model.Comment{}, 0, nil
	} else if err != nil {
		return []model.Comment{}, 0, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	return comments, totalRecords, nil
}

func (c *CommentRepo) Update(ctx context.Context, commentDTO dto.UpdateCommentDTO, commentID string) (model.Comment, error) {
	var comment model.Comment
	if commentDTO.Content == nil {
		err := c.storage.DB().WithContext(ctx).Model(&comment).Where("id = ?", commentID).Update("content", nil).Error
		if err != nil {
			return model.Comment{}, httperror.New(http.StatusInternalServerError, err.Error())
		}
	} else {
		err := c.storage.DB().WithContext(ctx).Model(&comment).Where("id = ?", commentID).Updates(commentDTO).Error
		if err != nil {
			return model.Comment{}, httperror.New(http.StatusInternalServerError, err.Error())
		}
	}
	err := c.storage.DB().WithContext(ctx).Preload("User").First(&comment, "id = ?", commentID).Error
	if err != nil {
		return model.Comment{}, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	return comment, nil
}
