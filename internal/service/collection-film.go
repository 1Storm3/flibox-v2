package service

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/mapper"
	"github.com/1Storm3/flibox-api/pkg/sys"
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
		return sys.NewError(sys.ErrInvalidRequestData, err.Error())
	}
	err = c.collectionFilmRepo.Add(ctx, collectionId, filmIdInt)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sys.NewError(sys.ErrCollectionNotFound, err.Error())
		}
		if strings.Contains(err.Error(), "violates unique constraint") {
			return sys.NewError(sys.ErrFilmAlreadyAdded, err.Error())
		}
		if strings.Contains(err.Error(), "collection_films_film_id_fkey") {
			return sys.NewError(sys.ErrFilmNotFound, err.Error())
		}
		return sys.NewError(sys.ErrDatabaseFailure, err.Error())
	}
	return nil
}

func (c *CollectionFilmService) Delete(
	ctx context.Context,
	collectionId string,
	filmDto string,
) error {
	filmIdInt, err := strconv.Atoi(filmDto)
	if err != nil {
		return sys.NewError(sys.ErrInvalidRequestData, err.Error())
	}
	return c.collectionFilmRepo.Delete(ctx, collectionId, filmIdInt)
}

func (c *CollectionFilmService) GetFilmsByCollectionId(
	ctx context.Context,
	collectionID string,
	page int,
	pageSize int,
) (films dto.FilmsByCollectionIdDTO, totalRecords int64, err error) {
	result, totalRecords, err := c.collectionFilmRepo.GetFilmsByCollectionID(ctx, collectionID, page, pageSize)

	if err != nil {
		return dto.FilmsByCollectionIdDTO{}, 0, sys.NewError(sys.ErrDatabaseFailure, err.Error())
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
