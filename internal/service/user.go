package service

import (
	"context"
	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/1Storm3/flibox-api/internal/shared/helper"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"golang.org/x/crypto/bcrypt"
	"net/http"
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

func (s *UserService) UpdateForVerify(ctx context.Context, userDTO dto.UpdateForVerifyDTO) (model.User, error) {
	return s.userRepo.UpdateForVerify(ctx, userDTO)
}

func (s *UserService) GetOneByNickName(ctx context.Context, nickName string) (model.User, error) {
	return s.userRepo.GetOneByNickName(ctx, nickName)
}

func (s *UserService) GetOneById(ctx context.Context, id string) (model.User, error) {

	user, err := s.userRepo.GetOneById(ctx, id)

	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s *UserService) GetOneByEmail(ctx context.Context, email string) (model.User, error) {
	return s.userRepo.GetOneByEmail(ctx, email)
}

func (s *UserService) CheckPassword(_ context.Context, user *model.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (s *UserService) HashPassword(_ context.Context, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	return string(hashedPassword), nil
}

func (s *UserService) Create(ctx context.Context, user model.User) (model.User, error) {
	return s.userRepo.Create(ctx, user)
}

func (s *UserService) Update(ctx context.Context, userDTO dto.UpdateUserDTO) (model.User, error) {
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

	return s.userRepo.Update(context.Background(), userDTO)
}
