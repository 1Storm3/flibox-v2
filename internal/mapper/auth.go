package mapper

import (
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/model"
)

func MapLoginDTOToUserModel(dto dto.LoginDTO) model.User {
	return model.User{
		Email:    dto.Email,
		Password: dto.Password,
	}
}

func MapRegisterDTOToUserModel(dto dto.RegisterDTO) model.User {
	return model.User{
		NickName: dto.NickName,
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
		Photo:    dto.Photo,
	}
}
