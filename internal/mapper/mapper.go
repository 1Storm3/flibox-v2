package mapper

import (
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/model"
)

func MapModelCollectionToResponseDTO(collection model.Collection) dto.ResponseCollectionDTO {
	return dto.ResponseCollectionDTO{
		ID:          collection.ID,
		Name:        collection.Name,
		Description: collection.Description,
		CoverUrl:    collection.CoverUrl,
		User: dto.User{
			ID:       collection.User.ID,
			NickName: collection.User.NickName,
			Photo:    collection.User.Photo,
		},
		Tags:      collection.Tags,
		CreatedAt: collection.CreatedAt,
		UpdatedAt: collection.UpdatedAt,
	}
}

func MapModelFilmToDTO(film model.Film) dto.Film {
	return dto.Film{
		ID:              film.ID,
		NameRU:          film.NameRU,
		NameOriginal:    film.NameOriginal,
		Type:            film.Type,
		Year:            film.Year,
		PosterURL:       film.PosterURL,
		RatingKinopoisk: film.RatingKinopoisk,
	}
}

func MapModelFilmsToDTOs(films []model.Film) []dto.Film {
	dtoFilms := make([]dto.Film, len(films))
	for i, film := range films {
		dtoFilms[i] = MapModelFilmToDTO(film)
	}
	return dtoFilms
}

func MapModelUserToResponseDTO(user model.User) dto.MeResponseDTO {
	return dto.MeResponseDTO{
		Id:         user.ID,
		Name:       user.Name,
		NickName:   user.NickName,
		Email:      user.Email,
		Photo:      user.Photo,
		Role:       user.Role,
		CreatedAt:  user.CreatedAt.String(),
		UpdatedAt:  user.UpdatedAt.String(),
		IsVerified: user.IsVerified,
	}
}

func MapModelCommentToResponseDTO(comment model.Comment) dto.ResponseDTO {
	return dto.ResponseDTO{
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

func MapModelFilmToResponseSearchDTO(film model.Film) dto.SearchResponseDTO {
	return dto.SearchResponseDTO{
		ID:              film.ID,
		NameRU:          film.NameRU,
		NameOriginal:    film.NameOriginal,
		Type:            film.Type,
		Year:            film.Year,
		RatingKinopoisk: film.RatingKinopoisk,
		PosterURL:       film.PosterURL,
	}
}

func MapModelFilmToResponseDTO(film model.Film) dto.ResponseFilmDTO {
	return dto.ResponseFilmDTO{
		ID:              film.ID,
		NameRU:          film.NameRU,
		NameOriginal:    film.NameOriginal,
		Year:            film.Year,
		RatingKinopoisk: film.RatingKinopoisk,
		PosterURL:       film.PosterURL,
		Description:     film.Description,
		LogoURL:         film.LogoURL,
		Type:            film.Type,
		Genres:          film.Genres,
		CoverURL:        film.CoverURL,
		TrailerURL:      film.TrailerURL,
	}
}

func MapDomainUserFilmToResponseDTO(userFilm model.UserFilm) dto.GetUserFilmResponseDTO {
	return dto.GetUserFilmResponseDTO{
		UserID: userFilm.UserID,
		FilmID: userFilm.FilmID,
		Film:   userFilm.Film,
		Type:   userFilm.Type,
	}
}
