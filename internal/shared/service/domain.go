package service

import (
	"context"
	"cv_api/internal/models"

	"github.com/redis/go-redis/v9"
)

type CacheRepository[T any] interface {
	GetAll(ctx context.Context, offset int, limit int) ([]T, int64, error)
	GetAllPublic(ctx context.Context, offset int, limit int) ([]T, int64, error)
	GetAllRecents(ctx context.Context, max int) ([]T, error)
	GetByID(ctx context.Context, id uint64) (*T, error)
	ExistsById(ctx context.Context, id uint64) (bool, error)
	Create(ctx context.Context, data *T) error
	Update(ctx context.Context, id uint64, data *T) error
}

type PaginatedResult[T any] struct {
	Data     []T               `json:"data"`
	Paginate models.Pagination `json:"paginate"`
}

type CacheService struct {
	rd *redis.Client
}

func NewCacheService(rd *redis.Client) *CacheService {
	return &CacheService{rd: rd}
}
