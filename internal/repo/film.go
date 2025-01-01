package repo

import (
	"context"

	"github.com/lib/pq"

	"github.com/1Storm3/flibox-api/database/postgres"
	"github.com/1Storm3/flibox-api/internal/dto"
)

type FilmRepo struct {
	storage *postgres.Storage
}

func NewFilmRepo(storage *postgres.Storage) *FilmRepo {
	return &FilmRepo{
		storage: storage,
	}
}

func (f *FilmRepo) GetOneByNameRu(ctx context.Context, nameRu string) (dto.FilmRepoDTO, error) {
	var film dto.FilmRepoDTO

	err := f.storage.DB().
		WithContext(ctx).
		Table("films").
		Where("name_ru ILIKE ?", "%"+nameRu+"%").
		First(&film).Error

	if err != nil {
		return dto.FilmRepoDTO{}, err
	}

	return film, nil
}

func (f *FilmRepo) GetOne(ctx context.Context, filmID string) (dto.FilmRepoDTO, error) {
	var film dto.FilmRepoDTO

	err := f.storage.DB().WithContext(ctx).Table("films").Where("id = ?", filmID).First(&film).Error

	if err != nil {
		return dto.FilmRepoDTO{}, err
	}

	return film, nil
}

func (f *FilmRepo) Save(ctx context.Context, film dto.FilmRepoDTO) error {
	err := f.storage.DB().WithContext(ctx).Table("films").Create(&film).Error

	if err != nil {
		return err
	}

	return nil
}

func (f *FilmRepo) Search(
	ctx context.Context,
	match string,
	genres []string,
	limit, pageSize int,
) ([]dto.FilmRepoDTO, int64, error) {
	var films []dto.FilmRepoDTO
	var totalRecords int64

	offset := (limit - 1) * pageSize

	query := f.storage.DB().WithContext(ctx).Table("films")

	query = query.Where("name_ru ILIKE ? OR name_original ILIKE ?", "%"+match+"%", "%"+match+"%")

	if len(genres) > 0 {
		query = query.Where("genres && ?", pq.Array(genres))
	}

	err := query.Count(&totalRecords).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.
		Limit(pageSize).
		Offset(offset).
		Find(&films).Error

	if err != nil {
		return nil, 0, err
	}

	return films, totalRecords, nil
}
