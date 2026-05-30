package post

import (
	"context"
	"cv_api/internal/models"
	"cv_api/internal/shared/service"
	"cv_api/internal/shared/utils"
	"fmt"
	"strings"
	"time"
)

type PostBasic struct {
	ID        uint64 `json:"id" gorm:"column:id"`
	Slug      string `json:"slug" gorm:"column:slug"`
	Title     string `json:"title" gorm:"column:title"`
	Banner    string `json:"banner" gorm:"column:banner"`
	StateID   uint64 `json:"state_id" gorm:"column:state_id"`
	Category  string `json:"category" gorm:"column:category"`
	CreatedAt string `json:"created_at" gorm:"column:created_at"`
}

type PostInfoBasic struct {
	ID    uint64 `json:"id" gorm:"column:id"`
	Slug  string `json:"slug" gorm:"column:slug"`
	Title string `json:"title" gorm:"column:title"`
}

type PotstRepository interface {
	GetAllStates(ctx context.Context) ([]models.StatePost, error)
	GetAll(ctx context.Context, offset int, limit int) ([]PostBasic, int64, error)
	GetAllPublic(ctx context.Context, offset int, limit int) ([]models.Post, int64, error)
	GetAllRecents(ctx context.Context, max int) ([]models.Post, error)
	GetBySlugPublic(ctx context.Context, slug string) (*models.Post, error)
	GetBySlugPrivate(ctx context.Context, slug string, userID uint64) (*models.Post, error)
	GetBasicInfoBySlug(ctx context.Context, slug string) (*PostInfoBasic, error)
	GetBySlug(ctx context.Context, slug string) (*models.Post, error)
	GetByID(ctx context.Context, id uint64) (*models.Post, error)
	ExistsBySlug(ctx context.Context, slug string) (bool, error)
	ExistsById(ctx context.Context, id uint64) (bool, error)
	Create(ctx context.Context, data *models.Post) error
	Update(ctx context.Context, slug string, data map[string]any) error
}

type PostService interface {
	GetAllStates(ctx context.Context) ([]models.StatePost, error)
	GetAll(ctx context.Context, page int, pageSize int) (*service.PaginatedResult[PostBasic], error)
	GetAllPublic(ctx context.Context, page int, pageSize int) (*service.PaginatedResult[models.Post], error)
	GetAllRecients(ctx context.Context, max int) (*service.PaginatedResult[models.Post], error)
	GetBySlugPublic(ctx context.Context, slug string) (*models.Post, error)
	GetBySlugPrivate(ctx context.Context, slug string, userID uint64) (*models.Post, error)
	GetBySlug(ctx context.Context, slug string) (*models.Post, error)
	GetByID(ctx context.Context, id uint64) (*models.Post, error)
	ExistsBySlug(ctx context.Context, slug string) (bool, error)
	ExistsById(ctx context.Context, id uint64) (bool, error)
	Create(ctx context.Context, authID uint64, data *models.Post) error
	Update(ctx context.Context, authID uint64, slug string, post *models.Post) error
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

func NewPost(data *models.Post) (*models.Post, error) {
	if data.AuthorID == 0 {
		return nil, fmt.Errorf("author_id is required")
	}

	if data.StateID == 0 {
		return nil, fmt.Errorf("state_id is required")
	}

	if strings.TrimSpace(data.Title) == "" {
		return nil, fmt.Errorf("title is required")
	}

	if strings.TrimSpace(data.Content) == "" {
		return nil, fmt.Errorf("content is required")
	}

	if strings.TrimSpace(data.Banner) == "" {
		return nil, fmt.Errorf("banner is required")
	}

	if len(data.Tags) == 0 {
		return nil, fmt.Errorf("tags is required")
	}

	if strings.TrimSpace(data.Category) == "" {
		return nil, fmt.Errorf("category is required")
	}

	data.Slug = generateSlug(data.Title)

	post := &models.Post{
		Slug:     data.Slug,
		Title:    utils.NormalizeText(data.Title),
		Content:  data.Content,
		Banner:   data.Banner,
		AuthorID: data.AuthorID,
		Tags:     data.Tags,
		Category: data.Category,
		StateID:  data.StateID,
	}

	return post, nil
}

func BuildPostUpdateData(data *models.Post, changeTitle bool) map[string]any {
	updateData := make(map[string]any)

	if strings.TrimSpace(data.Title) != "" && !changeTitle {
		updateData["title"] = utils.NormalizeText(data.Title)
		updateData["slug"] = generateSlug(data.Title)
	}

	if strings.TrimSpace(data.Content) != "" {
		updateData["content"] = data.Content
		updateData["content_clean"] = utils.CleanMarkdownForSearch(data.Content, false)
	}

	if strings.TrimSpace(data.Banner) != "" {
		updateData["banner"] = data.Banner
	}

	if len(data.Tags) > 0 {
		updateData["tags"] = data.Tags
	}

	if strings.TrimSpace(data.Category) != "" {
		updateData["category"] = data.Category
	}

	return updateData
}
