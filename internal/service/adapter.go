package service

import (
	"context"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/model"
)

type UserRepo interface {
	GetOneById(ctx context.Context, id string) (model.User, error)
	GetOneByEmail(ctx context.Context, email string) (model.User, error)
	Create(ctx context.Context, user model.User) (model.User, error)
	GetOneByNickName(ctx context.Context, nickName string) (model.User, error)
	Update(ctx context.Context, userDTO dto.UpdateUserDTO) (model.User, error)
	UpdateForVerify(ctx context.Context, userDTO dto.UpdateForVerifyDTO) (model.User, error)
}

type FilmRepo interface {
	GetOne(ctx context.Context, filmId string) (model.Film, error)
	Save(ctx context.Context, film model.Film) error
	Search(ctx context.Context, match string, genres []string, page, pageSize int) ([]model.Film, int64, error)
	GetOneByNameRu(ctx context.Context, nameRu string) (model.Film, error)
}

type CollectionFilmRepo interface {
	Add(ctx context.Context, collectionId string, filmId int) error
	Delete(ctx context.Context, collectionId string, filmId int) error
	GetFilmsByCollectionId(ctx context.Context, collectionID string, page int, pageSize int) ([]model.Film, int64, error)
}

type CollectionRepo interface {
	GetOne(ctx context.Context, collectionId string) (model.Collection, error)
	Create(ctx context.Context, collection model.Collection) (model.Collection, error)
	GetAll(ctx context.Context, page, pageSize int) ([]model.Collection, int64, error)
	Delete(ctx context.Context, collectionId string) error
	Update(ctx context.Context, collection model.Collection, collectionId string) (model.Collection, error)
	GetAllMy(ctx context.Context, page, pageSize int, userID string) ([]model.Collection, int64, error)
}

type CommentRepo interface {
	Create(ctx context.Context, comment model.Comment) (model.Comment, error)
	Delete(ctx context.Context, commentID string) error
	GetAllByFilmId(ctx context.Context, filmID string, page, pageSize int) ([]model.Comment, int64, error)
	Update(ctx context.Context, comment dto.UpdateCommentDTO, commentID string) (model.Comment, error)
	GetOne(ctx context.Context, commentID string) (model.Comment, error)
	GetCountByParentId(ctx context.Context, parentId string) (int64, error)
}

type HistoryFilmsRepo interface {
	GetAll(ctx context.Context, userId string) ([]model.HistoryFilms, error)
	Add(ctx context.Context, filmId, userId string) error
}

type UserFilmRepo interface {
	GetAllForRecommend(ctx context.Context, userId string, typeUserFilm model.TypeUserFilm, limit int) ([]model.UserFilm, error)
	Add(ctx context.Context, params dto.Params) error
	Delete(ctx context.Context, params dto.Params) error
	AddMany(ctx context.Context, params []dto.Params) error
	DeleteMany(ctx context.Context, userID string) error
}

type FilmSequelRepo interface {
	GetAll(ctx context.Context, filmId string) ([]model.FilmSequel, error)
	Save(ctx context.Context, filmId int, sequelId int) error
}

type FilmSimilarRepo interface {
	GetAll(ctx context.Context, filmId string) ([]model.FilmSimilar, error)
	Save(ctx context.Context, filmId int, similarId int) error
}
