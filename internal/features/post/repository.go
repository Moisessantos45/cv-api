package post

import (
	"context"
	"cv_api/internal/models"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) PotstRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) GetAllStates(ctx context.Context) ([]models.StatePost, error) {
	var states []models.StatePost
	err := r.db.WithContext(ctx).Find(&states).Error
	if err != nil {
		return nil, err
	}
	return states, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context, offset int, limit int) ([]PostBasic, int64, error) {
	var posts []PostBasic
	var total int64

	db := r.db.WithContext(ctx).Select("id", "slug", "title", "banner", "state_id", "category", "created_at").Order("created_at DESC")

	err := db.Model(&models.Post{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []PostBasic{}, 0, nil
	}

	err = db.Offset(offset).Limit(limit).Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *PostgresRepository) GetAllPublic(ctx context.Context, offset int, limit int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	db := r.db.WithContext(ctx).Select("id", "slug", "title", "banner", "content", "tags", "category", "created_at", "updated_at").Where("state_id = ?", 2).Order("created_at DESC")

	err := db.Model(&models.Post{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []models.Post{}, 0, nil
	}

	err = db.Offset(offset).Limit(limit).Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *PostgresRepository) GetAllRecents(ctx context.Context, max int) ([]models.Post, error) {
	var posts []models.Post

	db := r.db.WithContext(ctx).Select("id", "slug", "title", "content", "author_id", "tags", "category", "state_id", "created_at", "updated_at").Where("state_id = ?", 2).Order("created_at DESC")

	err := db.Limit(max).Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostgresRepository) GetBasicInfoBySlug(ctx context.Context, slug string) (*PostInfoBasic, error) {
	var post PostInfoBasic

	if err := r.db.WithContext(ctx).Select("id", "slug", "title").Where("slug = ?", slug).First(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostgresRepository) GetBySlugPublic(ctx context.Context, slug string) (*models.Post, error) {
	var post models.Post

	if err := r.db.WithContext(ctx).Where("slug = ? AND state_id = ?", slug, 2).First(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostgresRepository) GetBySlugPrivate(ctx context.Context, slug string, userID uint64) (*models.Post, error) {
	var post models.Post

	if err := r.db.WithContext(ctx).Where("slug = ? AND author_id = ?", slug, userID).First(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostgresRepository) GetBySlug(ctx context.Context, slug string) (*models.Post, error) {
	var post models.Post

	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uint64) (*models.Post, error) {
	var post models.Post

	if err := r.db.WithContext(ctx).First(&post, id).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostgresRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Post{}).Where("slug = ?", slug).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *PostgresRepository) ExistsById(ctx context.Context, id uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Post{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *PostgresRepository) Create(ctx context.Context, data *models.Post) error {
	if err := r.db.WithContext(ctx).Create(data).Error; err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) Update(ctx context.Context, slug string, data map[string]any) error {
	if err := r.db.WithContext(ctx).Model(&models.Post{}).Where("slug = ?", slug).Updates(data).Error; err != nil {
		return err
	}

	return nil
}
