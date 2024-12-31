package service

import (
	"context"
	"net/http"
	"strconv"

	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/mapper"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
)

type CollectionFilmService struct {
	collectionFilmRepo CollectionFilmRepo
}

func NewCollectionFilmService(collectionFilmRepo CollectionFilmRepo) *CollectionFilmService {
	return &CollectionFilmService{
		collectionFilmRepo: collectionFilmRepo,
	}
}

func (c *CollectionFilmService) Add(
	ctx context.Context,
	collectionId string,
	filmDto string,
) error {
	filmIdInt, err := strconv.Atoi(filmDto)
	if err != nil {
		return httperror.New(
			http.StatusBadRequest,
			err.Error(),
		)
	}
	return c.collectionFilmRepo.Add(ctx, collectionId, filmIdInt)
}

func (c *CollectionFilmService) Delete(
	ctx context.Context,
	collectionId string,
	filmDto string,
) error {
	filmIdInt, err := strconv.Atoi(filmDto)
	if err != nil {
		return httperror.New(
			http.StatusBadRequest,
			err.Error(),
		)
	}
	return c.collectionFilmRepo.Delete(ctx, collectionId, filmIdInt)
}

func (c *CollectionFilmService) GetFilmsByCollectionId(
	ctx context.Context,
	collectionID string,
	page int,
	pageSize int,
) (films dto.FilmsByCollectionIdDTO, totalRecords int64, err error) {
	result, totalRecords, err := c.collectionFilmRepo.GetFilmsByCollectionId(ctx, collectionID, page, pageSize)

	if err != nil {
		return dto.FilmsByCollectionIdDTO{}, 0, err
	}

	var filmsDTO []dto.ResponseFilmDTO
	for _, film := range result {
		filmsDTO = append(filmsDTO,
			mapper.MapModelFilmToResponseDTO(
				mapper.MapFilmRepoDTOToFilmModel(film),
			),
		)
	}

	return dto.FilmsByCollectionIdDTO{
		CollectionID: collectionID,
		Films:        filmsDTO,
	}, totalRecords, err
}
