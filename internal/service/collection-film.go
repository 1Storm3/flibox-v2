package service

import (
	"context"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/mapper"
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
	filmDto dto.CreateCollectionFilmDTO,
) error {
	return c.collectionFilmRepo.Add(ctx, collectionId, filmDto.FilmID)
}

func (c *CollectionFilmService) Delete(
	ctx context.Context,
	collectionId string,
	filmDto dto.DeleteCollectionFilmDTO,
) error {
	return c.collectionFilmRepo.Delete(ctx, collectionId, filmDto.FilmID)
}

func (c *CollectionFilmService) GetFilmsByCollectionId(
	ctx context.Context,
	collectionID string,
	page int,
	pageSize int,
) (films dto.FilmsByCollectionIdDTO, totalRecords int64, err error) {
	result, totalRecords, err := c.collectionFilmRepo.GetFilmsByCollectionId(ctx, collectionID, page, pageSize)

	return dto.FilmsByCollectionIdDTO{
		CollectionID: collectionID,
		Films:        mapper.MapModelFilmsToDTOs(result),
	}, totalRecords, err
}
