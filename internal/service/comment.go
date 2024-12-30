package service

import (
	"context"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/mapper"
	"github.com/1Storm3/flibox-api/internal/model"
)

type CommentService struct {
	commentRepo CommentRepo
}

func NewCommentService(commentRepo CommentRepo) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
	}
}

func (c *CommentService) Create(ctx context.Context, comment dto.CreateCommentDTO, userID string) (dto.ResponseDTO, error) {
	result, err := c.commentRepo.Create(ctx, model.Comment{
		Content:  comment.Content,
		FilmID:   comment.FilmID,
		UserID:   userID,
		ParentID: comment.ParentID,
	})

	if err != nil {
		return dto.ResponseDTO{}, err
	}
	return mapper.MapModelCommentToResponseDTO(result), nil
}

func (c *CommentService) GetOne(ctx context.Context, commentID string) (dto.ResponseDTO, error) {
	result, err := c.commentRepo.GetOne(ctx, commentID)
	if err != nil {
		return dto.ResponseDTO{}, err
	}

	return mapper.MapModelCommentToResponseDTO(result), nil
}

func (c *CommentService) Update(ctx context.Context, comment dto.UpdateCommentDTO, commentID string) (dto.ResponseDTO, error) {
	result, err := c.commentRepo.Update(ctx, comment, commentID)
	if err != nil {
		return dto.ResponseDTO{}, err
	}
	return mapper.MapModelCommentToResponseDTO(result), nil
}

func (c *CommentService) Delete(ctx context.Context, commentID string) error {
	comment, err := c.commentRepo.GetOne(ctx, commentID)
	if err != nil {
		return err
	}
	if comment.ParentID == nil {
		countChildComments, err := c.commentRepo.GetCountByParentId(ctx, commentID)
		if err != nil {
			return err
		}
		if countChildComments != 0 {
			_, err := c.commentRepo.Update(ctx, dto.UpdateCommentDTO{Content: nil}, commentID)
			if err != nil {
				return err
			}
		} else {
			err := c.commentRepo.Delete(ctx, commentID)
			if err != nil {
				return err
			}
		}
	} else {
		countSiblingComment, err := c.commentRepo.GetCountByParentId(ctx, *comment.ParentID)
		if err != nil {
			return err
		}
		if countSiblingComment == 1 {
			parentComment, err := c.commentRepo.GetOne(ctx, *comment.ParentID)
			if err != nil {
				return err
			}
			if parentComment.Content == nil {
				err := c.commentRepo.Delete(ctx, *comment.ParentID)
				if err != nil {
					return err
				}
			}
		}

		err = c.commentRepo.Delete(ctx, commentID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CommentService) GetAllByFilmId(ctx context.Context, filmID string, page int, pageSize int) ([]dto.ResponseDTO, int64, error) {
	comments, totalRecords, err := c.commentRepo.GetAllByFilmId(ctx, filmID, page, pageSize)
	if err != nil {
		return []dto.ResponseDTO{}, 0, err
	}
	var commentsDTO []dto.ResponseDTO
	for _, comment := range comments {
		commentsDTO = append(commentsDTO, mapper.MapModelCommentToResponseDTO(comment))
	}
	if len(commentsDTO) == 0 {
		return []dto.ResponseDTO{}, totalRecords, nil
	}
	return commentsDTO, totalRecords, nil
}
