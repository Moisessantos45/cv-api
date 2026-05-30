package experience

import (
	"context"
	"cv_api/internal/models"
	"cv_api/internal/shared/service"
	"fmt"
)

type ExperienceUseCase struct {
	repo  ExperienceRepository
	cache *service.CacheService
}

const experiencesCacheKey = "experiences:all"

func NewExperienceService(repo ExperienceRepository, cache *service.CacheService) ExperienceService {
	return &ExperienceUseCase{
		repo:  repo,
		cache: cache,
	}
}

func newPaginationExperience(data []models.Experience, total int64, totalPages int, page int, pageSize int) *service.PaginatedResult[models.Experience] {
	return &service.PaginatedResult[models.Experience]{
		Data: data,
		Paginate: models.Pagination{
			Total: total, TotalPages: totalPages, Page: page, PageSize: pageSize,
		},
	}
}

func (s *ExperienceUseCase) GetAllExperiences(ctx context.Context, page int, pageSize int, all bool) (*service.PaginatedResult[models.Experience], error) {
	cacheKey := experiencesCacheKey
	if all {
		cacheKey += ":all"
	}
	return service.Paginate[models.Experience](
		ctx,
		s.cache,
		cacheKey,
		page,
		pageSize,
		func(ctx context.Context, offset, limit int) ([]models.Experience, int64, error) {
			return s.repo.GetAll(ctx, offset, limit, all)
		},
		newPaginationExperience,
	)
}

func (s *ExperienceUseCase) GetAllExperiencesRecients(ctx context.Context, max int) (*service.PaginatedResult[models.Experience], error) {
	return service.GetRecents[models.Experience](
		ctx,
		s.cache,
		experiencesCacheKey,
		max,
		s.repo.GetAllRecents,
		newPaginationExperience,
	)
}

func (s *ExperienceUseCase) GetExperienceById(ctx context.Context, id uint64) (*models.Experience, error) {
	if id == 0 {
		return nil, fmt.Errorf("Id invalido")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *ExperienceUseCase) AddExperience(ctx context.Context, data *models.Experience) error {
	if data == nil {
		return fmt.Errorf("los datos de la experiencia son requeridos")
	}

	validatedExperience, err := NewExperience(*data)
	if err != nil {
		return err
	}

	if err = s.repo.Create(ctx, validatedExperience); err != nil {
		return fmt.Errorf("error al crear la experiencia: %w", err)
	}

	*data = *validatedExperience
	return nil
}

func (s *ExperienceUseCase) Update(ctx context.Context, id uint64, data *models.Experience) error {
	if id == 0 {
		return fmt.Errorf("Id invalido")
	}

	if data == nil {
		return fmt.Errorf("los datos de la experiencia son requeridos")
	}

	exists, err := s.repo.ExistsById(ctx, id)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("la experiencia no existe")
	}

	validatedExperience, err := NewExperience(*data)
	if err != nil {
		return err
	}

	if err = s.repo.Update(ctx, id, validatedExperience); err != nil {
		return fmt.Errorf("error al actualizar la experiencia: %w", err)
	}

	return nil
}
