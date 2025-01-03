package service

import (
	"context"

	"github.com/1Storm3/flibox-api/internal/dto"
)

type UserRepo interface {
	GetOneById(ctx context.Context, id string) (dto.UserRepoDTO, error)
	GetOneByEmail(ctx context.Context, email string) (dto.UserRepoDTO, error)
	Create(ctx context.Context, user dto.UserRepoDTO) (dto.UserRepoDTO, error)
	GetOneByNickName(ctx context.Context, nickName string) (dto.UserRepoDTO, error)
	Update(ctx context.Context, userDTO dto.UserRepoDTO) (dto.UserRepoDTO, error)
	UpdateForVerify(ctx context.Context, user dto.UserRepoDTO) (dto.UserRepoDTO, error)
}

type FilmRepo interface {
	GetOne(ctx context.Context, filmId string) (dto.FilmRepoDTO, error)
	Save(ctx context.Context, film dto.FilmRepoDTO) error
	Search(ctx context.Context, match string, genres []string, page, pageSize int) ([]dto.FilmRepoDTO, int64, error)
	GetOneByNameRu(ctx context.Context, nameRu string) (dto.FilmRepoDTO, error)
}

type CollectionFilmRepo interface {
	Add(ctx context.Context, collectionId string, filmId int) error
	Delete(ctx context.Context, collectionId string, filmId int) error
	GetFilmsByCollectionID(ctx context.Context, collectionID string, page int, pageSize int) ([]dto.FilmRepoDTO, int64, error)
}

type CollectionRepo interface {
	GetOne(ctx context.Context, collectionId string) (dto.CollectionRepoDTO, error)
	Create(ctx context.Context, collection dto.CollectionRepoDTO) (dto.CollectionRepoDTO, error)
	GetAll(ctx context.Context, page, pageSize int) ([]dto.CollectionRepoDTO, int64, error)
	Delete(ctx context.Context, collectionId string) error
	Update(ctx context.Context, collection dto.CollectionRepoDTO) (dto.CollectionRepoDTO, error)
	GetAllMy(ctx context.Context, page, pageSize int, userID string) ([]dto.CollectionRepoDTO, int64, error)
}

type CommentRepo interface {
	Create(ctx context.Context, comment dto.CommentRepoDTO) (dto.CommentRepoDTO, error)
	Delete(ctx context.Context, commentID string) error
	GetAllByFilmId(ctx context.Context, filmID string, page, pageSize int) ([]dto.CommentRepoDTO, int64, error)
	Update(ctx context.Context, comment dto.CommentRepoDTO, commentID string) (dto.CommentRepoDTO, error)
	GetOne(ctx context.Context, commentID string) (dto.CommentRepoDTO, error)
	GetCountByParentId(ctx context.Context, parentId string) (int64, error)
}

type HistoryFilmsRepo interface {
	GetAll(ctx context.Context, userId string) ([]dto.HistoryFilmsRepoDTO, error)
	Add(ctx context.Context, filmId, userId string) error
}

type UserFilmRepo interface {
	GetAllForRecommend(ctx context.Context, userId string, typeUserFilm dto.TypeUserFilm, limit int) ([]dto.UserFilmRepoDTO, error)
	Add(ctx context.Context, params dto.Params) error
	Delete(ctx context.Context, params dto.Params) error
	AddMany(ctx context.Context, params []dto.Params) error
	DeleteMany(ctx context.Context, userID string) error
}

type FilmSequelRepo interface {
	GetAll(ctx context.Context, filmId string) ([]dto.FilmSequelRepoDTO, error)
	Save(ctx context.Context, filmId int, sequelId int) error
}

type FilmSimilarRepo interface {
	GetAll(ctx context.Context, filmId string) ([]dto.FilmSimilarRepoDTO, error)
	Save(ctx context.Context, filmId int, similarId int) error
}
