package profile

import (
	"context"
	"cv_api/internal/features/auth"
	"cv_api/internal/models"
	"cv_api/internal/shared/service"
	"fmt"
)

type ProfileUseCase struct {
	repo   ProfileRepository
	ucAuth auth.AuthService
	cache  *service.CacheService
}

const profilesCacheKey = "profiles:all"

func NewProfileUseCase(repo ProfileRepository, ucAuth auth.AuthService, cache *service.CacheService) ProfileService {
	return &ProfileUseCase{
		repo:   repo,
		ucAuth: ucAuth,
		cache:  cache,
	}
}

func NewPaginationProfile(data []models.Profile, total int64, totalPages int, page int, pageSize int) *service.PaginatedResult[models.Profile] {
	return &service.PaginatedResult[models.Profile]{
		Data: data,
		Paginate: models.Pagination{
			Total: total, TotalPages: totalPages, Page: page, PageSize: pageSize,
		},
	}
}

func (s *ProfileUseCase) GetAll(ctx context.Context, page int, pageSize int, all bool) (*service.PaginatedResult[models.Profile], error) {
	cacheKey := profilesCacheKey
	if all {
		cacheKey += ":all"
	}
	return service.Paginate[models.Profile](
		ctx,
		s.cache,
		cacheKey,
		page,
		pageSize,
		func(ctx context.Context, offset, limit int) ([]models.Profile, int64, error) {
			return s.repo.GetAll(ctx, offset, limit, all)
		},
		NewPaginationProfile,
	)
}

func (s *ProfileUseCase) GetAllRecents(ctx context.Context, max int) (*service.PaginatedResult[models.Profile], error) {
	return service.GetRecents[models.Profile](
		ctx,
		s.cache,
		profilesCacheKey,
		max,
		s.repo.GetAllRecents,
		NewPaginationProfile,
	)
}

func (s *ProfileUseCase) GetBasicByAuthID(ctx context.Context, authID uint64) (*ProfileBasic, error) {
	if authID == 0 {
		return nil, fmt.Errorf("authID invalido")
	}

	return s.repo.GetBasicByAuthID(ctx, authID)
}

func (s *ProfileUseCase) GetByID(ctx context.Context, authID uint64) (*models.Profile, error) {
	if authID == 0 {
		return nil, fmt.Errorf("Id invalido")
	}

	return s.repo.GetByID(ctx, authID)
}

func (s *ProfileUseCase) GetByAuthID(ctx context.Context, authID uint64) (*models.Profile, error) {
	if authID == 0 {
		return nil, fmt.Errorf("authID invalido")
	}

	return s.repo.GetByAuthID(ctx, authID)
}

func (s *ProfileUseCase) Create(ctx context.Context, authID uint64, data *models.Profile) error {

	if data == nil {
		return fmt.Errorf("los datos del perfil son requeridos")
	}

	data.AuthID = authID

	validatedProfile, err := NewProfile(*data)
	if err != nil {
		return err
	}

	err = s.repo.WithTransaction(func(repo *PostgresRepository) error {
		if err := s.repo.Create(ctx, validatedProfile); err != nil {
			return err
		}

		if err := s.ucAuth.ChangeCompletProfile(authID, true); err != nil {
			return err
		}

		return nil

	})

	if err != nil {
		return err
	}

	*data = *validatedProfile
	return nil
}

func (s *ProfileUseCase) Update(ctx context.Context, authID uint64, data *models.Profile) error {
	if authID == 0 {
		return fmt.Errorf("Id de autorization invalido")
	}

	if data == nil {
		return fmt.Errorf("los datos del perfil son requeridos")
	}

	exists, err := s.repo.ExistsById(ctx, authID)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("El perfil no existe")
	}

	validatedProfile, err := NewProfile(*data)
	if err != nil {
		return err
	}

	err = s.repo.Update(ctx, authID, validatedProfile)
	if err != nil {
		return fmt.Errorf("error al actualizar el perfil: %w", err)
	}

	return nil
}
