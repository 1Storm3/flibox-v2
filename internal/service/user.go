package service

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/mapper"
	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/1Storm3/flibox-api/internal/shared/helper"
	"github.com/1Storm3/flibox-api/pkg/sys"
)

type UserService struct {
	userRepo  UserRepo
	s3Service controller.S3Service
}

func NewUserService(userRepo UserRepo, s3Service controller.S3Service) *UserService {
	return &UserService{
		userRepo:  userRepo,
		s3Service: s3Service,
	}
}

func (s *UserService) UpdateForVerify(ctx context.Context, user model.User) (model.User, error) {
	userRepoDTO := mapper.MapUserModelToUserRepoDTO(user)
	result, err := s.userRepo.UpdateForVerify(ctx, userRepoDTO)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, sys.NewError(sys.ErrUserNotFound, err.Error())
		}
	}

	return mapper.MapUserRepoDTOToUserModel(result), nil
}

func (s *UserService) GetOneByNickName(ctx context.Context, nickName string) (model.User, error) {
	result, err := s.userRepo.GetOneByNickName(ctx, nickName)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, sys.NewError(sys.ErrUserNotFound, err.Error())
		}
	}
	return mapper.MapUserRepoDTOToUserModel(result), nil
}

func (s *UserService) GetOneById(ctx context.Context, id string) (model.User, error) {

	user, err := s.userRepo.GetOneById(ctx, id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, sys.NewError(sys.ErrUserNotFound, err.Error())
		}
	}
	return mapper.MapUserRepoDTOToUserModel(user), nil
}

func (s *UserService) GetOneByEmail(ctx context.Context, email string) (model.User, error) {
	result, err := s.userRepo.GetOneByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, sys.NewError(sys.ErrUserNotFound, err.Error())
		}
	}
	return mapper.MapUserRepoDTOToUserModel(result), nil
}

func (s *UserService) CheckPassword(_ context.Context, user *model.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (s *UserService) HashPassword(_ context.Context, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", sys.NewError(sys.ErrPasswordHashGeneration, err.Error())
	}
	return string(hashedPassword), nil
}

func (s *UserService) Create(ctx context.Context, user model.User) (model.User, error) {
	userRepo := mapper.MapUserModelToUserRepoDTO(user)
	result, err := s.userRepo.Create(ctx, userRepo)

	if err != nil {
		return model.User{}, sys.NewError(sys.ErrCreateUser, err.Error())
	}

	return mapper.MapUserRepoDTOToUserModel(result), nil
}

func (s *UserService) Update(ctx context.Context, userDTO model.User) (model.User, error) {
	if userDTO.Photo != nil {
		user, err := s.GetOneById(ctx, userDTO.ID)
		if err != nil {
			return model.User{}, err
		}

		if user.Photo != nil && *user.Photo != "" {
			key, err := helper.ExtractS3Key(*user.Photo)
			if err != nil {
				return model.User{}, err
			}

			err = s.s3Service.DeleteFile(ctx, key)
			if err != nil {
				return model.User{}, err
			}
		}
	}
	userRepo := mapper.MapUserModelToUserRepoDTO(userDTO)
	result, err := s.userRepo.Update(ctx, userRepo)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return model.User{}, sys.NewError(sys.ErrUserAlreadyExists, err.Error())
		}
		return model.User{}, sys.NewError(sys.ErrUpdateUser, err.Error())
	}

	return mapper.MapUserRepoDTOToUserModel(result), nil

}
