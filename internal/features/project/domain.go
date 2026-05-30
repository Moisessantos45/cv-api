package project

import (
	"context"
	"cv_api/internal/models"
	"cv_api/internal/shared/service"
	"cv_api/internal/shared/utils"
	"errors"
	"fmt"
	"strings"
	"time"
)

type ProjectBasic struct {
	ID          uint64 `json:"id" gorm:"column:id"`
	Slug        string `json:"slug" gorm:"column:slug"`
	Title       string `json:"title" gorm:"column:title"`
	TypeProject string `json:"type_project" gorm:"column:type_project"`
	Banner      string `json:"banner" gorm:"column:banner"`
	StateID     uint64 `json:"state_id" gorm:"column:state_id"`
	CreatedAt   string `json:"created_at" gorm:"column:created_at"`
}

type ProjectInfoBasic struct {
	ID    uint64 `json:"id" gorm:"column:id"`
	Slug  string `json:"slug" gorm:"column:slug"`
	Title string `json:"title" gorm:"column:title"`
}

type ProjectRepository interface {
	WithTransaction(fn func(repo *PostgresRepository) error) error
	GetAllData(ctx context.Context) ([]models.Project, error)
	GetAll(ctx context.Context, offset int, limit int) ([]ProjectBasic, int64, error)
	GetAllPublic(ctx context.Context, offset int, limit int) ([]models.Project, int64, error)
	GetAllRecents(ctx context.Context, max int) ([]models.Project, error)
	GetByID(ctx context.Context, id uint64) (*models.Project, error)
	GetBySlug(ctx context.Context, slug string) (*models.Project, error)
	GetBySlugPublic(ctx context.Context, slug string) (*models.Project, error)
	GetBasicInfoBySlug(ctx context.Context, slug string) (*ProjectInfoBasic, error)
	GetCurrentCount(ctx context.Context, id uint64) (int64, error)
	ExistsById(ctx context.Context, id uint64) (bool, error)
	Create(ctx context.Context, data *models.Project) error
	Update(ctx context.Context, id uint64, data map[string]any) error
	UpdateState(ctx context.Context, id uint64, stateID uint64) error
	UpdateCounter(ctx context.Context, id uint64, counter int64) error
}

type ProjectService interface {
	CreateAll(ctx context.Context, data []models.Project) error
	GetAllData(ctx context.Context) ([]models.Project, error)
	GetAll(ctx context.Context, page int, pageSize int) (*service.PaginatedResult[ProjectBasic], error)
	GetAllPublic(ctx context.Context, page int, pageSize int) (*service.PaginatedResult[models.Project], error)
	GetAllRecents(ctx context.Context, max int) (*service.PaginatedResult[models.Project], error)
	GetBySlug(ctx context.Context, slug string) (*models.Project, error)
	GetBySlugPublic(ctx context.Context, slug string) (*models.Project, error)
	GetByID(ctx context.Context, id uint64) (*models.Project, error)
	Create(ctx context.Context, data *models.Project) error
	Update(ctx context.Context, slug string, data *models.Project) error
	UpdateState(ctx context.Context, id uint64, stateID uint64) error
	UpdateCounter(ctx context.Context, id uint64, counter int64) error
}

type PaginatedProjects struct {
	Data     []models.Project  `json:"data"`
	Paginate models.Pagination `json:"paginate"`
}

func NewProject(data *models.Project) (*models.Project, error) {
	if strings.TrimSpace(data.Title) == "" {
		return nil, errors.New("el título del puesto es obligatorio")
	}

	if strings.TrimSpace(data.TypeProject) == "" {
		return nil, errors.New("el tipo de contratación es obligatorio")
	}

	if len(data.Technologies) == 0 {
		return nil, errors.New("las tecnologías utilizadas son obligatorias")
	}

	if len(data.Characteristics) == 0 {
		return nil, errors.New("las características del proyecto son obligatorias")
	}

	if len(data.Learning) == 0 {
		return nil, errors.New("lo aprendido en el proyecto es obligatorio")
	}

	if len(data.Images) == 0 {
		return nil, errors.New("las imágenes del proyecto son obligatorias")
	}

	if strings.TrimSpace(data.Banner) == "" {
		return nil, errors.New("el banner del proyecto es obligatorio")
	}

	if strings.TrimSpace(data.Description) == "" {
		return nil, errors.New("la descripción del proyecto es obligatoria")
	}

	if strings.TrimSpace(data.CreatedAt) == "" {
		return nil, errors.New("la fecha de creación del proyecto es obligatoria")
	}

	if data.StateID == 0 {
		return nil, errors.New("el estado del proyecto es obligatorio")
	}

	data.Slug = generateSlug(data.Title)

	return &models.Project{
		Slug:            data.Slug,
		Title:           data.Title,
		TypeProject:     strings.TrimSpace(data.TypeProject),
		Technologies:    data.Technologies,
		Characteristics: data.Characteristics,
		Banner:          strings.TrimSpace(data.Banner),
		Description:     strings.TrimSpace(data.Description),
		CreatedAt:       strings.TrimSpace(data.CreatedAt),
		StateID:         data.StateID,
		Learning:        data.Learning,
		Images:          data.Images,
		Link:            strings.TrimSpace(data.Link),
		LinkFrontend:    data.LinkFrontend,
		LinkBackend:     data.LinkBackend,
	}, nil
}

func BuildProjectUpdateData(data *models.Project, changeTitle bool) map[string]any {
	updateData := make(map[string]any)

	if strings.TrimSpace(data.Title) != "" && !changeTitle {
		updateData["title"] = utils.NormalizeText(data.Title)
		updateData["slug"] = generateSlug(data.Title)
	}

	if strings.TrimSpace(data.TypeProject) != "" {
		updateData["type_project"] = strings.TrimSpace(data.TypeProject)
	}

	if len(data.Technologies) > 0 {
		updateData["technologies"] = data.Technologies
	}

	if len(data.Characteristics) > 0 {
		updateData["characteristics"] = data.Characteristics
	}

	if strings.TrimSpace(data.Banner) != "" {
		updateData["banner"] = strings.TrimSpace(data.Banner)
	}

	if strings.TrimSpace(data.Description) != "" {
		updateData["description"] = strings.TrimSpace(data.Description)
	}

	if strings.TrimSpace(data.CreatedAt) != "" {
		updateData["created_at"] = strings.TrimSpace(data.CreatedAt)
	}

	if len(data.Learning) > 0 {
		updateData["learning"] = data.Learning
	}

	if len(data.Images) > 0 {
		updateData["images"] = data.Images
	}

	if strings.TrimSpace(data.Link) != "" {
		updateData["link"] = strings.TrimSpace(data.Link)
	}

	if len(data.LinkFrontend) > 0 {
		updateData["link_frontend"] = data.LinkFrontend
	}

	if len(data.LinkBackend) > 0 {
		updateData["link_backend"] = data.LinkBackend
	}

	return updateData
}

func generateSlug(title string) string {
	slugNormalized := utils.NormalizeText(title)
	slugBase := strings.ToLower(strings.TrimSpace(slugNormalized))
	slugBase = strings.ReplaceAll(slugBase, " ", "-")

	suffix := fmt.Sprintf("-%d", time.Now().UnixNano())

	const maxSlugLen = 55
	maxBaseLen := max(1, maxSlugLen-len(suffix))

	if len(slugBase) > maxBaseLen {
		slugBase = slugBase[:maxBaseLen]
		slugBase = strings.TrimRight(slugBase, "-")
	}

	return slugBase + suffix
}

func compareStrings(a string, b string) bool {
	if strings.TrimSpace(a) == "" && strings.TrimSpace(b) == "" {
		return true
	}

	return utils.NormalizeText(a) == utils.NormalizeText(b)
}
