package service

import (
	"context"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/mapper"
	"github.com/1Storm3/flibox-api/internal/model"
)

type CollectionService struct {
	collectionRepo CollectionRepo
}

func NewCollectionService(collectionRepo CollectionRepo) *CollectionService {
	return &CollectionService{
		collectionRepo: collectionRepo,
	}
}

func (c *CollectionService) Update(ctx context.Context, collection dto.UpdateCollectionDTO, collectionId string) (dto.ResponseCollectionDTO, error) {
	result, err := c.collectionRepo.Update(ctx, model.Collection{
		Name:        collection.Name,
		Description: collection.Description,
		CoverUrl:    collection.CoverUrl,
		Tags:        collection.Tags,
	}, collectionId)
	if err != nil {
		return dto.ResponseCollectionDTO{}, err
	}
	return mapper.MapModelCollectionToResponseDTO(result), nil
}

func (c *CollectionService) Delete(ctx context.Context, collectionId string) error {
	return c.collectionRepo.Delete(ctx, collectionId)
}

func (c *CollectionService) Create(ctx context.Context, collection dto.CreateCollectionDTO, userID string) (dto.ResponseCollectionDTO, error) {
	result, err := c.collectionRepo.Create(ctx, model.Collection{
		Name:        collection.Name,
		Description: collection.Description,
		CoverUrl:    collection.CoverUrl,
		Tags:        collection.Tags,
		UserId:      &userID,
	})
	if err != nil {
		return dto.ResponseCollectionDTO{}, err
	}
	return mapper.MapModelCollectionToResponseDTO(result), nil
}

func (c *CollectionService) GetAll(ctx context.Context, page, pageSize int) ([]dto.ResponseCollectionDTO, int64, error) {
	result, totalRecords, err := c.collectionRepo.GetAll(ctx, page, pageSize)
	if err != nil {
		return []dto.ResponseCollectionDTO{}, 0, err
	}
	var resultDTO []dto.ResponseCollectionDTO
	for _, collection := range result {
		resultDTO = append(resultDTO, mapper.MapModelCollectionToResponseDTO(collection))
	}
	return resultDTO, totalRecords, nil
}

func (c *CollectionService) GetOne(ctx context.Context, collectionId string) (dto.ResponseCollectionDTO, error) {
	result, err := c.collectionRepo.GetOne(ctx, collectionId)
	if err != nil {
		return dto.ResponseCollectionDTO{}, err
	}
	return mapper.MapModelCollectionToResponseDTO(result), nil
}

func (c *CollectionService) GetAllMy(ctx context.Context, page, pageSize int, userID string) ([]dto.ResponseCollectionDTO, int64, error) {
	result, totalRecords, err := c.collectionRepo.GetAllMy(ctx, page, pageSize, userID)
	if err != nil {
		return []dto.ResponseCollectionDTO{}, 0, err
	}
	var resultDTO []dto.ResponseCollectionDTO
	for _, collection := range result {
		resultDTO = append(resultDTO, mapper.MapModelCollectionToResponseDTO(collection))
	}
	return resultDTO, totalRecords, nil
}
