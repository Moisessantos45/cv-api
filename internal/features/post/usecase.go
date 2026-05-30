package post

import (
	"context"
	"cv_api/internal/features/profile"
	"cv_api/internal/models"
	"cv_api/internal/shared/service"
	"fmt"
	"strings"
)

type PostUseCase struct {
	repo  PotstRepository
	cache *service.CacheService
	pS    profile.ProfileService
}

const postsCacheKey = "posts:all"

func NewPostService(repo PotstRepository, cache *service.CacheService, profileService profile.ProfileService) PostService {
	return &PostUseCase{
		repo:  repo,
		cache: cache,
		pS:    profileService,
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

func (uc *PostUseCase) GetAllStates(ctx context.Context) ([]models.StatePost, error) {
	return uc.repo.GetAllStates(ctx)
}

func (s *PostUseCase) GetAll(ctx context.Context, page int, pageSize int) (*service.PaginatedResult[PostBasic], error) {
	return service.Paginate[PostBasic](
		ctx,
		s.cache,
		postsCacheKey,
		page,
		pageSize,
		s.repo.GetAll,
		NewPaginationResult[PostBasic],
	)
}

func (s *PostUseCase) GetAllPublic(ctx context.Context, page int, pageSize int) (*service.PaginatedResult[models.Post], error) {
	return service.Paginate[models.Post](
		ctx,
		s.cache,
		postsCacheKey+":public",
		page,
		pageSize,
		s.repo.GetAllPublic,
		NewPaginationResult[models.Post],
	)
}

func (s *PostUseCase) GetAllRecients(ctx context.Context, max int) (*service.PaginatedResult[models.Post], error) {
	return service.GetRecents[models.Post](
		ctx,
		s.cache,
		postsCacheKey,
		max,
		s.repo.GetAllRecents,
		NewPaginationResult[models.Post],
	)
}

func (s *PostUseCase) GetBySlugPublic(ctx context.Context, slug string) (*models.Post, error) {
	if strings.TrimSpace(slug) == "" {
		return nil, fmt.Errorf("slug invalido")
	}

	return s.repo.GetBySlugPublic(ctx, slug)
}

func (s *PostUseCase) GetBySlugPrivate(ctx context.Context, slug string, userID uint64) (*models.Post, error) {
	if strings.TrimSpace(slug) == "" {
		return nil, fmt.Errorf("slug invalido")
	}

	if userID == 0 {
		return nil, fmt.Errorf("Id de usuario invalido")
	}

	return s.repo.GetBySlugPrivate(ctx, slug, userID)
}

func (s *PostUseCase) GetBySlug(ctx context.Context, slug string) (*models.Post, error) {
	if strings.TrimSpace(slug) == "" {
		return nil, fmt.Errorf("slug invalido")
	}

	return s.repo.GetBySlug(ctx, slug)
}

func (s *PostUseCase) GetByID(ctx context.Context, id uint64) (*models.Post, error) {
	if id == 0 {
		return nil, fmt.Errorf("Id invalido")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *PostUseCase) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	return s.repo.ExistsBySlug(ctx, slug)
}

func (s *PostUseCase) ExistsById(ctx context.Context, id uint64) (bool, error) {
	return s.repo.ExistsById(ctx, id)
}

func (s *PostUseCase) Create(ctx context.Context, authID uint64, data *models.Post) error {
	if data == nil {
		return fmt.Errorf("los datos del post son requeridos")
	}

	if authID == 0 {
		return fmt.Errorf("Id de autenticacion invalido")
	}

	profile, err := s.pS.GetBasicByAuthID(ctx, authID)
	if err != nil {
		return fmt.Errorf("error al obtener el perfil del usuario: %w", err)
	}

	if profile == nil {
		return fmt.Errorf("perfil del usuario no encontrado")
	}

	data.AuthorID = profile.ID

	validatedPost, err := NewPost(data)
	if err != nil {
		return err
	}

	err = s.repo.Create(ctx, validatedPost)
	if err != nil {
		return fmt.Errorf("error al crear el post: %w", err)
	}

	*data = *validatedPost
	return nil
}

func (s *PostUseCase) Update(ctx context.Context, authID uint64, slug string, post *models.Post) error {
	if strings.TrimSpace(slug) == "" {
		return fmt.Errorf("slug invalido")
	}

	findPost, err := s.repo.GetBasicInfoBySlug(ctx, slug)
	if err != nil {
		return fmt.Errorf("error fetching post: %v", err)
	}

	isEqualTitle := compareStrings(findPost.Title, post.Title)

	updateData := BuildPostUpdateData(post, isEqualTitle)

	err = s.repo.Update(ctx, slug, updateData)
	if err != nil {
		return fmt.Errorf("error al actualizar el post: %w", err)
	}

	return nil
}
