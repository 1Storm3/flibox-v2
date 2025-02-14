package service

import (
	"context"
	"errors"

	"github.com/1Storm3/flibox-api/internal/mapper"
	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/1Storm3/flibox-api/pkg/sys"
	"gorm.io/gorm"
)

type CollectionService struct {
	collectionRepo CollectionRepo
}

func NewCollectionService(collectionRepo CollectionRepo) *CollectionService {
	return &CollectionService{
		collectionRepo: collectionRepo,
	}
}

func (c *CollectionService) Update(ctx context.Context, collection model.Collection) (model.Collection, error) {
	collectionRepo := mapper.MapCollectionModelToCollectionRepoDTO(collection)

	isExist, err := c.GetOne(ctx, collection.ID)

	if err != nil {
		return model.Collection{}, err
	}

	if isExist.ID == "" {
		return model.Collection{}, sys.NewError(sys.ErrCollectionNotFound, "collection not found")
	}

	result, err := c.collectionRepo.Update(ctx, collectionRepo)

	if err != nil {
		return model.Collection{}, sys.NewError(sys.ErrDatabaseFailure, err.Error())
	}

	return mapper.MapCollectionRepoDTOToCollectionModel(result), nil
}

func (c *CollectionService) Delete(ctx context.Context, collectionId string) error {
	err := c.collectionRepo.Delete(ctx, collectionId)

	if err != nil {
		return sys.NewError(sys.ErrDatabaseFailure, err.Error())
	}
	return nil
}

func (c *CollectionService) Create(ctx context.Context, collection model.Collection, userID string) (model.Collection, error) {
	collectionRepo := mapper.MapCollectionModelToCollectionRepoDTO(collection)

	collectionRepo.UserId = &userID

	result, err := c.collectionRepo.Create(ctx, collectionRepo)
	if err != nil {
		return model.Collection{}, err
	}
	return mapper.MapCollectionRepoDTOToCollectionModel(result), nil
}

func (c *CollectionService) GetAll(ctx context.Context, page, pageSize int) ([]model.Collection, int64, error) {
	result, totalRecords, err := c.collectionRepo.GetAll(ctx, page, pageSize)

	if err != nil {
		return []model.Collection{}, 0, sys.NewError(sys.ErrDatabaseFailure, err.Error())
	}
	var resultDTO []model.Collection

	for _, collection := range result {
		resultDTO = append(resultDTO, mapper.MapCollectionRepoDTOToCollectionModel(collection))
	}

	return resultDTO, totalRecords, nil
}

func (c *CollectionService) GetOne(ctx context.Context, collectionId string) (model.Collection, error) {
	result, err := c.collectionRepo.GetOne(ctx, collectionId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Collection{}, sys.NewError(sys.ErrCollectionNotFound, err.Error())
		}
		return model.Collection{}, sys.NewError(sys.ErrDatabaseFailure, err.Error())
	}

	return mapper.MapCollectionRepoDTOToCollectionModel(result), nil
}

func (c *CollectionService) GetAllMy(ctx context.Context, page, pageSize int, userID string) ([]model.Collection, int64, error) {
	result, totalRecords, err := c.collectionRepo.GetAllMy(ctx, page, pageSize, userID)
	if err != nil {
		return []model.Collection{}, 0, sys.NewError(sys.ErrDatabaseFailure, err.Error())
	}
	var resultDTO []model.Collection
	for _, collection := range result {
		resultDTO = append(resultDTO, mapper.MapCollectionRepoDTOToCollectionModel(collection))
	}

	return resultDTO, totalRecords, nil
}
