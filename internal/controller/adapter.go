package controller

import (
	"context"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/1Storm3/flibox-api/pkg/proto/gengrpc"
	"time"
)

type AuthService interface {
	Login(ctx context.Context, dto dto.LoginDTO) (string, error)
	Register(ctx context.Context, user dto.RegisterDTO) (bool, error)
	Me(ctx context.Context, userId string) (dto.MeResponseDTO, error)
	Verify(ctx context.Context, tokenDto string) error
}

type CollectionService interface {
	Create(ctx context.Context, collection dto.CreateCollectionDTO, userID string) (dto.ResponseCollectionDTO, error)
	GetAll(ctx context.Context, page, pageSize int) ([]dto.ResponseCollectionDTO, int64, error)
	GetOne(ctx context.Context, collectionId string) (dto.ResponseCollectionDTO, error)
	Update(ctx context.Context, collection dto.UpdateCollectionDTO, collectionId string) (dto.ResponseCollectionDTO, error)
	Delete(ctx context.Context, collectionId string) error
	GetAllMy(ctx context.Context, page, pageSize int, userID string) ([]dto.ResponseCollectionDTO, int64, error)
}

type CollectionFilmService interface {
	Add(ctx context.Context, collectionId string, filmDto dto.CreateCollectionFilmDTO) error
	Delete(ctx context.Context, collectionId string, filmDto dto.DeleteCollectionFilmDTO) error
	GetFilmsByCollectionId(ctx context.Context, collectionID string, page int, pageSize int) (films dto.FilmsByCollectionIdDTO, totalRecords int64, err error)
}

type CommentService interface {
	Create(ctx context.Context, comment dto.CreateCommentDTO, userID string) (dto.ResponseDTO, error)
	Update(ctx context.Context, comment dto.UpdateCommentDTO, commentID string) (dto.ResponseDTO, error)
	Delete(ctx context.Context, commentID string) error
	GetAllByFilmId(ctx context.Context, filmID string, page int, pageSize int) ([]dto.ResponseDTO, int64, error)
	GetOne(ctx context.Context, commentID string) (dto.ResponseDTO, error)
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
	Update(ctx context.Context, userDTO dto.UpdateUserDTO) (model.User, error)
	UpdateForVerify(ctx context.Context, userDTO dto.UpdateForVerifyDTO) (model.User, error)
}

type UserFilmService interface {
	GetAll(ctx context.Context, userId string, typeUserFilm model.TypeUserFilm, limit int) ([]dto.GetUserFilmResponseDTO, error)
	Add(ctx context.Context, params dto.Params) error
	Delete(ctx context.Context, params dto.Params) error
	AddMany(ctx context.Context, params []dto.Params) error
	DeleteMany(ctx context.Context, userID string) error
}

type FilmSequelService interface {
	GetAll(ctx context.Context, filmId string) ([]dto.ResponseFilmDTO, error)
	FetchSequels(ctx context.Context, filmId string) ([]dto.ResponseFilmDTO, error)
}

type FilmSimilarService interface {
	GetAll(ctx context.Context, filmId string) ([]dto.ResponseFilmDTO, error)
	FetchSimilar(ctx context.Context, filmId string) ([]dto.ResponseFilmDTO, error)
}
type EmailService interface {
	SendEmail(email, body, subject string) error
}

type FilmService interface {
	GetOne(ctx context.Context, filmId string) (dto.ResponseFilmDTO, error)
	Search(ctx context.Context, match string, genres []string, page, pageSize int) ([]dto.SearchResponseDTO, int64, error)
	GetOneByNameRu(ctx context.Context, nameRu string) (dto.ResponseFilmDTO, error)
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
