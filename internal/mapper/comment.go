package mapper

import (
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/model"
)

func MapCreateCommentDTOToCommentModel(dto dto.CreateCommentDTO) model.Comment {
	return model.Comment{
		Content:  dto.Content,
		FilmID:   dto.FilmID,
		UserID:   dto.UserID,
		ParentID: dto.ParentID,
	}
}

func MapCommentModelToCommentRepoDTO(comment model.Comment) dto.CommentRepoDTO {
	var parentComment *dto.CommentRepoDTO
	if comment.Parent != nil {
		mappedParent := MapCommentModelToCommentRepoDTO(*comment.Parent)
		parentComment = &mappedParent
	}

	return dto.CommentRepoDTO{
		ID:        comment.ID,
		FilmID:    comment.FilmID,
		UserID:    comment.UserID,
		ParentID:  comment.ParentID,
		Content:   comment.Content,
		CreatedAt: parseTimeStringToTime(comment.CreatedAt),
		UpdatedAt: parseTimeStringToTime(comment.UpdatedAt),
		User:      MapUserModelToUserRepoDTO(comment.User),
		Parent:    parentComment,
	}
}

func MapCommentRepoDTOToCommentModel(comment dto.CommentRepoDTO) model.Comment {
	var parentComment *model.Comment
	if comment.Parent != nil {
		mappedParent := MapCommentRepoDTOToCommentModel(*comment.Parent)
		parentComment = &mappedParent
	}

	return model.Comment{
		ID:        comment.ID,
		FilmID:    comment.FilmID,
		UserID:    comment.UserID,
		ParentID:  comment.ParentID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt.String(),
		UpdatedAt: comment.UpdatedAt.String(),
		User:      MapUserRepoDTOToUserModel(comment.User),
		Parent:    parentComment,
	}
}

func MapUpdateCommentDTOToCommentModel(dto dto.UpdateCommentDTO) model.Comment {
	return model.Comment{
		ID:      dto.ID,
		Content: dto.Content,
	}
}

func MapCommentModelToCommentResponseDTO(comment model.Comment) dto.CommentResponseDTO {
	return dto.CommentResponseDTO{
		ID:      comment.ID,
		Content: comment.Content,
		User: dto.User{
			ID:       comment.User.ID,
			NickName: comment.User.NickName,
			Photo:    comment.User.Photo,
		},
		FilmID:    comment.FilmID,
		ParentID:  comment.ParentID,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}
}
