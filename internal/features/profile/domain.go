package profile

import (
	"context"
	"cv_api/internal/models"
	"cv_api/internal/shared/service"
	"cv_api/internal/shared/utils"
	"errors"
	"strings"
)

type ProfileBasic struct {
	ID     uint64 `json:"id" gorm:"column:id"`
	AuthID uint64 `json:"auth_id" gorm:"column:auth_id"`
}

type ProfileRepository interface {
	WithTransaction(fn func(repo *PostgresRepository) error) error
	GetAll(ctx context.Context, offset int, limit int, all bool) ([]models.Profile, int64, error)
	GetAllRecents(ctx context.Context, max int) ([]models.Profile, error)
	GetByID(ctx context.Context, id uint64) (*models.Profile, error)
	GetBasicByAuthID(ctx context.Context, authID uint64) (*ProfileBasic, error)
	GetByAuthID(ctx context.Context, authID uint64) (*models.Profile, error)
	ExistsById(ctx context.Context, id uint64) (bool, error)
	Create(ctx context.Context, data *models.Profile) error
	Update(ctx context.Context, id uint64, data *models.Profile) error
}

type ProfileService interface {
	GetAll(ctx context.Context, offset int, limit int, all bool) (*service.PaginatedResult[models.Profile], error)
	GetAllRecents(ctx context.Context, max int) (*service.PaginatedResult[models.Profile], error)
	GetBasicByAuthID(ctx context.Context, authID uint64) (*ProfileBasic, error)
	GetByID(ctx context.Context, id uint64) (*models.Profile, error)
	GetByAuthID(ctx context.Context, authID uint64) (*models.Profile, error)
	Create(ctx context.Context, authID uint64, data *models.Profile) error
	Update(ctx context.Context, authID uint64, data *models.Profile) error
}

func NewProfile(data models.Profile) (*models.Profile, error) {
	if data.AuthID == 0 {
		return nil, errors.New("el ID de autenticación es obligatorio")
	}

	if strings.TrimSpace(data.Description) == "" {
		return nil, errors.New("la descripción del perfil es obligatoria")
	}

	if strings.TrimSpace(data.CVLink) != "" {
		if err := utils.ValidateURL(data.CVLink, "cv_link"); err != nil {
			return nil, err
		}
	}

	if strings.TrimSpace(data.LinkedInLink) != "" {
		if err := utils.ValidateURL(data.LinkedInLink, "linkedin_link"); err != nil {
			return nil, err
		}
	}

	if strings.TrimSpace(data.GitHubLink) != "" {
		if err := utils.ValidateURL(data.GitHubLink, "github_link"); err != nil {
			return nil, err
		}
	}

	if len(data.SocialLinks) == 0 {
		return nil, errors.New("al menos un enlace social es obligatorio")
	}

	return &models.Profile{
		Name:         data.Name,
		LastName:     data.LastName,
		CVLink:       data.CVLink,
		LinkedInLink: data.LinkedInLink,
		GitHubLink:   data.GitHubLink,
		SocialLinks:  data.SocialLinks,
		Description:  data.Description,
		AuthID:       data.AuthID,
	}, nil
}
