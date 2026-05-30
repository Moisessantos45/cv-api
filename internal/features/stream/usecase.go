package stream

import (
	"context"
	"cv_api/internal/models"
	"cv_api/internal/shared/service"
	"fmt"
)

type VideoUseCase struct {
	repo  VideoRepository
	cache *service.CacheService
}

const videosCacheKey = "videos:all"

func NewVideoUseCase(repo VideoRepository, cache *service.CacheService) VideoService {
	return &VideoUseCase{
		repo:  repo,
		cache: cache,
	}
}

func NewPaginationVideo(data []models.Video, total int64, totalPages int, page int, pageSize int) *service.PaginatedResult[models.Video] {
	return &service.PaginatedResult[models.Video]{
		Data: data,
		Paginate: models.Pagination{
			Total: total, TotalPages: totalPages, Page: page, PageSize: pageSize,
		},
	}
}

func (s *VideoUseCase) GetAllData(ctx context.Context) ([]models.Video, error) {
	streams, err := s.repo.GetAllData(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al obtener los videos: %w", err)
	}

	return streams, nil
}

func (s *VideoUseCase) GetAll(ctx context.Context, page int, pageSize int, all bool) (*service.PaginatedResult[models.Video], error) {
	cacheKey := videosCacheKey
	if all {
		cacheKey += ":all"
	}
	return service.Paginate[models.Video](
		ctx,
		s.cache,
		cacheKey,
		page,
		pageSize,
		func(ctx context.Context, offset, limit int) ([]models.Video, int64, error) {
			return s.repo.GetAll(ctx, offset, limit, all)
		},
		NewPaginationVideo,
	)
}

func (s *VideoUseCase) GetAllRecients(ctx context.Context, max int) (*service.PaginatedResult[models.Video], error) {
	return service.GetRecents[models.Video](
		ctx,
		s.cache,
		videosCacheKey,
		max,
		s.repo.GetAllRecents,
		NewPaginationVideo,
	)
}

func (s *VideoUseCase) GetByID(ctx context.Context, id uint64) (*models.Video, error) {
	if id == 0 {
		return nil, fmt.Errorf("Id invalido")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *VideoUseCase) CreateAll(ctx context.Context, data []models.Video) error {
	return s.repo.WithTransaction(func(repo *PostgresRepository) error {
		for i := range data {
			if err := s.Create(ctx, &data[i]); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *VideoUseCase) Create(ctx context.Context, data *models.Video) error {

	if data == nil {
		return fmt.Errorf("los datos de la video son requeridos")
	}

	validatedVideo, err := NewVideo(*data)
	if err != nil {
		return err
	}

	err = s.repo.Create(ctx, validatedVideo)
	if err != nil {
		return fmt.Errorf("error al crear el video: %w", err)
	}

	*data = *validatedVideo
	return nil
}

func (s *VideoUseCase) Update(ctx context.Context, id uint64, data *models.Video) error {
	if id == 0 {
		return fmt.Errorf("Id invalido")
	}

	if data == nil {
		return fmt.Errorf("los datos del video son requeridos")
	}

	exists, err := s.repo.ExistsById(ctx, id)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("El video no existe")
	}

	validatedVideo, err := NewVideo(*data)
	if err != nil {
		return err
	}

	err = s.repo.Update(ctx, id, validatedVideo)
	if err != nil {
		return fmt.Errorf("error al crear el proyecto: %w", err)
	}

	return nil
}
