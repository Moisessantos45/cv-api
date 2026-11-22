package experience

import (
	"context"
	"cv_api/internal/models"
	"cv_api/internal/shared/service"
	"errors"
	"strings"
)

type ExperienceRepository interface {
	GetAll(ctx context.Context, offset int, limit int, all bool) ([]models.Experience, int64, error)
	GetAllRecents(ctx context.Context, max int) ([]models.Experience, error)
	GetByID(ctx context.Context, id uint64) (*models.Experience, error)
	ExistsById(ctx context.Context, id uint64) (bool, error)
	Create(ctx context.Context, data *models.Experience) error
	Update(ctx context.Context, id uint64, data *models.Experience) error
}

type ExperienceService interface {
	GetAllExperiences(ctx context.Context, page int, pageSize int, all bool) (*service.PaginatedResult[models.Experience], error)
	GetAllExperiencesRecients(ctx context.Context, max int) (*service.PaginatedResult[models.Experience], error)
	GetExperienceById(ctx context.Context, id uint64) (*models.Experience, error)
	AddExperience(ctx context.Context, data *models.Experience) error
	Update(ctx context.Context, id uint64, data *models.Experience) error
}

func NewExperience(data models.Experience) (*models.Experience, error) {
	if strings.TrimSpace(data.Title) == "" {
		return nil, errors.New("el título del puesto es obligatorio")
	}

	if strings.TrimSpace(data.Company) == "" {
		return nil, errors.New("el nombre de la empresa es obligatorio")
	}

	if strings.TrimSpace(data.StartDate) == "" {
		return nil, errors.New("la fecha de inicio es obligatoria")
	}

	if strings.TrimSpace(data.Description) == "" {
		return nil, errors.New("la descripción de la experiencia es obligatoria")
	}

	if len(data.SkillsLearned) != 0 {
		return nil, errors.New("la lista de habilidades aprendidas no puede estar vacía")
	}

	return &models.Experience{
		Title:         data.Title,
		Company:       data.Company,
		Location:      data.Location,
		StartDate:     data.StartDate,
		EndDate:       data.EndDate,
		Description:   data.Description,
		SkillsLearned: data.SkillsLearned,
	}, nil
}
