package mapper

import (
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/model"
)

func MapUpdateUserDTOToUserModel(dto dto.UpdateUserDTO) model.User {
	return model.User{
		ID:           dto.ID,
		NickName:     *dto.NickName,
		Name:         *dto.Name,
		Email:        *dto.Email,
		LastActivity: *dto.LastActivity,
		Photo:        dto.Photo,
	}
}

func MapUserModelToUserResponseDto(user model.User) dto.ResponseUserDTO {
	return dto.ResponseUserDTO{
		ID:         user.ID,
		Name:       user.Name,
		NickName:   user.NickName,
		Email:      user.Email,
		Photo:      user.Photo,
		Role:       user.Role,
		CreatedAt:  parseTimeStringToTime(user.CreatedAt),
		UpdatedAt:  parseTimeStringToTime(user.UpdatedAt),
		IsVerified: user.IsVerified,
	}

}

func MapUserModelToUserRepoDTO(user model.User) dto.UserRepoDTO {
	return dto.UserRepoDTO{
		ID:            user.ID,
		NickName:      user.NickName,
		Name:          user.Name,
		Email:         user.Email,
		Password:      user.Password,
		IsVerified:    user.IsVerified,
		IsBlocked:     user.IsBlocked,
		VerifiedToken: user.VerifiedToken,
		LastActivity:  parseTimeStringToTime(user.LastActivity),
		Photo:         user.Photo,
		Role:          user.Role,
		CreatedAt:     parseTimeStringToTime(user.CreatedAt),
		UpdatedAt:     parseTimeStringToTime(user.UpdatedAt),
	}
}

func MapModelUserToResponseDTO(user model.User) dto.MeResponseDTO {
	return dto.MeResponseDTO{
		Id:         user.ID,
		Name:       user.Name,
		NickName:   user.NickName,
		Email:      user.Email,
		Photo:      user.Photo,
		Role:       user.Role,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
		IsVerified: user.IsVerified,
	}
}

func MapUserRepoDTOToUserModel(user dto.UserRepoDTO) model.User {
	return model.User{
		ID:            user.ID,
		NickName:      user.NickName,
		Name:          user.Name,
		Email:         user.Email,
		Password:      user.Password,
		IsVerified:    user.IsVerified,
		IsBlocked:     user.IsBlocked,
		VerifiedToken: user.VerifiedToken,
		LastActivity:  user.LastActivity.String(),
		Photo:         user.Photo,
		Role:          user.Role,
		CreatedAt:     user.CreatedAt.String(),
		UpdatedAt:     user.UpdatedAt.String(),
	}
}
