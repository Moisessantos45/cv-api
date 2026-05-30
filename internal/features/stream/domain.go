package stream

import (
	"context"
	"cv_api/internal/models"
	"cv_api/internal/shared/service"
	"errors"
	"strings"
)

type VideoRepository interface {
	WithTransaction(fn func(repo *PostgresRepository) error) error
	GetAllData(ctx context.Context) ([]models.Video, error)
	GetAll(ctx context.Context, offset int, limit int, all bool) ([]models.Video, int64, error)
	GetAllRecents(ctx context.Context, max int) ([]models.Video, error)
	GetByID(ctx context.Context, id uint64) (*models.Video, error)
	ExistsById(ctx context.Context, id uint64) (bool, error)
	Create(ctx context.Context, data *models.Video) error
	Update(ctx context.Context, id uint64, data *models.Video) error
}

type VideoService interface {
	CreateAll(ctx context.Context, data []models.Video) error
	GetAllData(ctx context.Context) ([]models.Video, error)
	GetAll(ctx context.Context, page int, pageSize int, all bool) (*service.PaginatedResult[models.Video], error)
	GetAllRecients(ctx context.Context, max int) (*service.PaginatedResult[models.Video], error)
	GetByID(ctx context.Context, id uint64) (*models.Video, error)
	Create(ctx context.Context, data *models.Video) error
	Update(ctx context.Context, id uint64, data *models.Video) error
}

func NewVideo(data models.Video) (*models.Video, error) {
	if strings.TrimSpace(data.Title) == "" {
		return nil, errors.New("el título del puesto es obligatorio")
	}
	if strings.TrimSpace(data.URL) == "" {
		return nil, errors.New("el tipo de contratación es obligatorio")
	}

	if strings.TrimSpace(data.Description) == "" {
		return nil, errors.New("la descripción es obligatoria")
	}
	if strings.TrimSpace(data.CreatedAt) == "" {
		return nil, errors.New("las responsabilidades son obligatorias")
	}

	return &models.Video{
		Title:       data.Title,
		URL:         strings.TrimSpace(data.URL),
		Description: strings.TrimSpace(data.Description),
		CreatedAt:   strings.TrimSpace(data.CreatedAt),
	}, nil
}
