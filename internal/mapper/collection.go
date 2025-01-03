package mapper

import (
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/lib/pq"
)

func MapCollectionRepoDTOToCollectionModel(dto dto.CollectionRepoDTO) model.Collection {
	return model.Collection{
		ID:          dto.ID,
		Name:        dto.Name,
		Description: dto.Description,
		CoverUrl:    dto.CoverUrl,
		Tags:        *dto.Tags,
		UserId:      dto.UserId,
		CreatedAt:   dto.CreatedAt.String(),
		UpdatedAt:   dto.UpdatedAt.String(),
		User:        MapUserRepoDTOToUserModel(dto.User),
		Films:       MapManyFilmsRepoDTOToManyFilmsModel(dto.Films),
	}
}

func MapCreateCollectionDTOToCollectionModel(dto dto.CreateCollectionDTO) model.Collection {
	return model.Collection{
		Name:        dto.Name,
		Description: dto.Description,
		CoverUrl:    dto.CoverUrl,
		Tags:        dto.Tags,
	}
}

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

func MapUpdateCollectionDTOToCollectionModel(dto dto.UpdateCollectionDTO) model.Collection {
	return model.Collection{
		ID:          dto.ID,
		Name:        dto.Name,
		Description: dto.Description,
		CoverUrl:    dto.CoverUrl,
		Tags:        dto.Tags,
	}
}

func MapCollectionModelToCollectionRepoDTO(collection model.Collection) dto.CollectionRepoDTO {
	tags := pq.StringArray(collection.Tags)
	return dto.CollectionRepoDTO{
		ID:          collection.ID,
		Name:        collection.Name,
		Description: collection.Description,
		CoverUrl:    collection.CoverUrl,
		Tags:        &tags,
		CreatedAt:   parseTimeStringToTime(collection.CreatedAt),
		UpdatedAt:   parseTimeStringToTime(collection.UpdatedAt),
		UserId:      collection.UserId,
	}
}
