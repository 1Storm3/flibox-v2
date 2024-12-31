package mapper

import (
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/model"
)

func MapUserFilmRepoDTOToUserFilmModel(dto dto.UserFilmRepoDTO) model.UserFilm {
	return model.UserFilm{
		UserID: dto.UserID,
		FilmID: dto.FilmID,
		Type:   dto.Type,
	}
}
