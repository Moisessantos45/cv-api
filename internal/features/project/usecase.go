package project

import (
	"context"
	"cv_api/internal/models"
	"cv_api/internal/shared/service"
	"fmt"
)

type ProjectUseCase struct {
	repo  ProjectRepository
	cache *service.CacheService
}

const projectsCacheKey = "projects:all"

func NewProjectService(repo ProjectRepository, cache *service.CacheService) ProjectService {
	return &ProjectUseCase{
		repo:  repo,
		cache: cache,
	}
}

func NewPaginationResult[T any](data []T, total int64, totalPages int, page int, pageSize int) *service.PaginatedResult[T] {
	return &service.PaginatedResult[T]{
		Data: data,
		Paginate: models.Pagination{
			Total:      total,
			TotalPages: totalPages,
			Page:       page,
			PageSize:   pageSize,
		},
	}
}

func (s *ProjectUseCase) GetAllData(ctx context.Context) ([]models.Project, error) {
	return s.repo.GetAllData(ctx)
}

func (s *ProjectUseCase) GetAll(ctx context.Context, page int, pageSize int) (*service.PaginatedResult[ProjectBasic], error) {
	return service.Paginate[ProjectBasic](
		ctx,
		s.cache,
		projectsCacheKey,
		page,
		pageSize,
		s.repo.GetAll,
		NewPaginationResult[ProjectBasic],
	)
}

func (s *ProjectUseCase) GetAllPublic(ctx context.Context, page int, pageSize int) (*service.PaginatedResult[models.Project], error) {
	return service.Paginate[models.Project](
		ctx,
		s.cache,
		projectsCacheKey+":public",
		page,
		pageSize,
		s.repo.GetAllPublic,
		NewPaginationResult[models.Project],
	)
}

func (s *ProjectUseCase) GetAllRecents(ctx context.Context, max int) (*service.PaginatedResult[models.Project], error) {
	return service.GetRecents[models.Project](
		ctx,
		s.cache,
		projectsCacheKey,
		max,
		s.repo.GetAllRecents,
		NewPaginationResult[models.Project],
	)
}

func (s *ProjectUseCase) GetByID(ctx context.Context, id uint64) (*models.Project, error) {
	if id == 0 {
		return nil, fmt.Errorf("Id invalido")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *ProjectUseCase) GetBySlug(ctx context.Context, slug string) (*models.Project, error) {
	if slug == "" {
		return nil, fmt.Errorf("el slug del proyecto es requerido")
	}

	return s.repo.GetBySlug(ctx, slug)
}

func (s *ProjectUseCase) GetBySlugPublic(ctx context.Context, slug string) (*models.Project, error) {
	if slug == "" {
		return nil, fmt.Errorf("el slug del proyecto es requerido")
	}

	return s.repo.GetBySlugPublic(ctx, slug)
}

func (s *ProjectUseCase) CreateAll(ctx context.Context, data []models.Project) error {
	return s.repo.WithTransaction(func(repo *PostgresRepository) error {
		for i := range data {
			// log.Printf("Creating project: %d", data[i].StateID)
			if err := s.Create(ctx, &data[i]); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *ProjectUseCase) Create(ctx context.Context, data *models.Project) error {

	if data == nil {
		return fmt.Errorf("los datos de la proyecto son requeridos")
	}

	validatedProject, err := NewProject(data)
	if err != nil {
		return err
	}

	err = s.repo.Create(ctx, validatedProject)
	if err != nil {
		return fmt.Errorf("error al crear el proyecto: %w", err)
	}

	*data = *validatedProject
	return nil
}

func (s *ProjectUseCase) Update(ctx context.Context, slug string, data *models.Project) error {
	if slug == "" {
		return fmt.Errorf("el slug del proyecto es requerido")
	}

	if data == nil {
		return fmt.Errorf("los datos del proyecto son requeridos")
	}

	findPost, err := s.repo.GetBasicInfoBySlug(ctx, slug)
	if err != nil {
		return fmt.Errorf("error fetching post: %v", err)
	}

	isEqualTitle := compareStrings(findPost.Title, data.Title)

	updateData := BuildProjectUpdateData(data, isEqualTitle)

	err = s.repo.Update(ctx, findPost.ID, updateData)
	if err != nil {
		return fmt.Errorf("error al crear el proyecto: %w", err)
	}

	return nil
}

func (s *ProjectUseCase) UpdateState(ctx context.Context, id uint64, stateID uint64) error {
	if id == 0 {
		return fmt.Errorf("el ID del proyecto es requerido")
	}

	exists, err := s.repo.ExistsById(ctx, id)
	if err != nil {
		return fmt.Errorf("error fetching post: %v", err)
	}

	if !exists {
		return fmt.Errorf("el proyecto no existe")
	}

	err = s.repo.UpdateState(ctx, id, stateID)
	if err != nil {
		return fmt.Errorf("error al actualizar el estado del proyecto: %w", err)
	}

	return nil
}

func (s *ProjectUseCase) UpdateCounter(ctx context.Context, id uint64, value int64) error {
	if id == 0 {
		return fmt.Errorf("Id invalido")
	}

	currentCount, err := s.repo.GetCurrentCount(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.UpdateCounter(ctx, id, currentCount+1)
}
