package controller

import (
	"context"
	"time"

	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/1Storm3/flibox-api/pkg/proto/gengrpc"
)

type AuthService interface {
	Login(ctx context.Context, req model.User) (string, error)
	Register(ctx context.Context, user model.User) (bool, error)
	Me(ctx context.Context, userId string) (model.User, error)
	Verify(ctx context.Context, token string) error
}

type CollectionService interface {
	Create(ctx context.Context, collection model.Collection, userID string) (model.Collection, error)
	GetAll(ctx context.Context, page, pageSize int) ([]model.Collection, int64, error)
	GetOne(ctx context.Context, collectionId string) (model.Collection, error)
	Update(ctx context.Context, collection model.Collection) (model.Collection, error)
	Delete(ctx context.Context, collectionId string) error
	GetAllMy(ctx context.Context, page, pageSize int, userID string) ([]model.Collection, int64, error)
}

type CollectionFilmService interface {
	Add(ctx context.Context, collectionId, filmId string) error
	Delete(ctx context.Context, collectionId, filmId string) error
	GetFilmsByCollectionId(ctx context.Context, collectionID string, page int, pageSize int) (films dto.FilmsByCollectionIdDTO, totalRecords int64, err error)
}

type CommentService interface {
	Create(ctx context.Context, comment model.Comment) (model.Comment, error)
	Update(ctx context.Context, comment model.Comment, commentID string) (model.Comment, error)
	Delete(ctx context.Context, commentID string) error
	GetAllByFilmId(ctx context.Context, filmID string, page int, pageSize int) ([]model.Comment, int64, error)
	GetOne(ctx context.Context, commentID string) (model.Comment, error)
}

type HistoryFilmsService interface {
	Add(ctx context.Context, filmId, userId string) error
	GetAll(ctx context.Context, userId string) ([]model.HistoryFilms, error)
}

type RecommendService interface {
	CreateRecommendations(params dto.RecommendationsParams) error
	GetFilmNamesForRecommendations(ctx context.Context, userID string) ([]*gengrpc.Film, error)
	GetFilmName(film *model.Film) *gengrpc.Film
	GetUniqueFilmIDsForRecommendations(ctx context.Context, recommendations []string) ([]*int, error)
	AddFilmRecommendations(ctx context.Context, userID string, filmIds []*int) error
}

type UserService interface {
	GetOneByNickName(ctx context.Context, nickName string) (model.User, error)
	GetOneByEmail(ctx context.Context, email string) (model.User, error)
	CheckPassword(ctx context.Context, user *model.User, password string) bool
	HashPassword(ctx context.Context, password string) (string, error)
	Create(ctx context.Context, user model.User) (model.User, error)
	GetOneById(ctx context.Context, id string) (model.User, error)
	Update(ctx context.Context, userDTO model.User) (model.User, error)
	UpdateForVerify(ctx context.Context, userDTO model.User) (model.User, error)
}

type UserFilmService interface {
	GetAll(ctx context.Context, userId string, typeUserFilm dto.TypeUserFilm, limit int) ([]model.UserFilm, error)
	Add(ctx context.Context, params dto.Params) error
	Delete(ctx context.Context, params dto.Params) error
	AddMany(ctx context.Context, params []dto.Params) error
	DeleteMany(ctx context.Context, userID string) error
}

type FilmSequelService interface {
	GetAll(ctx context.Context, filmId string) ([]model.FilmSequel, error)
	FetchSequels(ctx context.Context, filmId string) ([]model.FilmSequel, error)
}

type FilmSimilarService interface {
	GetAll(ctx context.Context, filmId string) ([]model.FilmSimilar, error)
	FetchSimilar(ctx context.Context, filmId string) ([]model.FilmSimilar, error)
}
type EmailService interface {
	SendEmail(email, body, subject string) error
}

type FilmService interface {
	GetOne(ctx context.Context, filmId string) (model.Film, error)
	Search(ctx context.Context, match string, genres []string, page, pageSize int) ([]model.Film, int64, error)
	GetOneByNameRu(ctx context.Context, nameRu string) (model.Film, error)
}

type ExternalService interface {
	SetExternalRequest(filmId string) (dto.GetExternalFilmDTO, error)
}

type TokenService interface {
	GenerateToken(jwtKey []byte, userID, role string, duration time.Duration) (string, error)
	GenerateEmailToken(email string, jwtKey []byte, duration time.Duration) (*string, error)
	ValidateEmailToken(tokenString string, jwtKey []byte) (string, error)
	ParseToken(tokenString string, jwtKey []byte) (*dto.Claims, error)
}

type S3Service interface {
	UploadFile(ctx context.Context, key string, file []byte) (string, error)
	DeleteFile(ctx context.Context, key string) error
}
