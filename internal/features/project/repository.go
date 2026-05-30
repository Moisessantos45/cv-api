package project

import (
	"context"
	"cv_api/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) ProjectRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) WithTransaction(fn func(repo *PostgresRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &PostgresRepository{db: tx}
		return fn(txRepo)
	})
}

func (r *PostgresRepository) GetAllData(ctx context.Context) ([]models.Project, error) {
	var projects []models.Project

	err := r.db.WithContext(ctx).Find(&projects).Error
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context, offset int, limit int) ([]ProjectBasic, int64, error) {
	var projects []ProjectBasic
	var total int64

	db := r.db.WithContext(ctx).Select("id", "slug", "title", "type_project", "banner", "created_at", "state_id").Order("created_at DESC")

	err := db.Model(&models.Project{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []ProjectBasic{}, 0, nil
	}

	err = db.WithContext(ctx).Offset(offset).Limit(limit).Find(&projects).Error

	if err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (r *PostgresRepository) GetAllPublic(ctx context.Context, offset int, limit int) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	db := r.db.WithContext(ctx).Select("id", "slug", "title", "type_project", "description", "technologies", "banner", "created_at", "state_id").Order("created_at DESC")

	err := db.Model(&models.Project{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []models.Project{}, 0, nil
	}

	err = db.WithContext(ctx).Offset(offset).Limit(limit).Find(&projects).Error

	if err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (r *PostgresRepository) GetAllRecents(ctx context.Context, max int) ([]models.Project, error) {
	var projects []models.Project

	db := r.db.WithContext(ctx).Select("id", "slug", "title", "type_project", "description", "technologies", "banner", "created_at", "state_id").Order("created_at DESC")

	err := db.WithContext(ctx).Limit(max).Find(&projects).Error
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uint64) (*models.Project, error) {
	var project models.Project

	if err := r.db.WithContext(ctx).First(&project, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("proyecto no encontrado para el id: %d", id)
		}

		return nil, err
	}

	return &project, nil
}

func (r *PostgresRepository) GetBySlug(ctx context.Context, slug string) (*models.Project, error) {
	var project models.Project

	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&project).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("proyecto no encontrado para el slug: %s", slug)
		}

		return nil, err
	}

	return &project, nil
}

func (r *PostgresRepository) GetBySlugPublic(ctx context.Context, slug string) (*models.Project, error) {
	var project models.Project

	err := r.db.WithContext(ctx).
		Select(
			"id", "slug", "title", "type_project", "description",
			"technologies", "characteristics", "learning",
			"banner", "images", "link", "created_at",
			"link_frontend", "link_backend",
		).
		Where("slug = ? AND state_id = ?", slug, 2).
		First(&project).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("proyecto no encontrado para el slug: %s", slug)
		}
		return nil, err
	}

	return &project, nil
}

func (r *PostgresRepository) GetBasicInfoBySlug(ctx context.Context, slug string) (*ProjectInfoBasic, error) {
	var info ProjectInfoBasic

	if err := r.db.WithContext(ctx).Select("id", "slug", "title").Where("slug = ?", slug).First(&info).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("proyecto no encontrado para el slug: %s", slug)
		}

		return nil, err
	}

	return &info, nil
}

func (r *PostgresRepository) GetCurrentCount(ctx context.Context, id uint64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Project{}).Select("counter_likes").Where("id = ?", id).Scan(&count).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("proyecto no encontrado para el id: %d", id)
		}

		return 0, fmt.Errorf("error al obtener counter_likes: %w", err)
	}

	if id == 0 {
		return 0, fmt.Errorf("empresa no encontrada para el id: %d", id)
	}

	return count, nil
}

func (r *PostgresRepository) ExistsById(ctx context.Context, id uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Project{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, err
	}

	return count > 0, nil
}

func (r *PostgresRepository) Create(ctx context.Context, data *models.Project) error {
	if err := r.db.WithContext(ctx).Create(data).Error; err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) Update(ctx context.Context, id uint64, data map[string]any) error {
	if err := r.db.WithContext(ctx).Model(&models.Project{}).Where("id = ?", id).Updates(data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("proyecto no encontrado para el id: %d", id)
		}

		return err
	}

	return nil
}

func (r *PostgresRepository) UpdateState(ctx context.Context, id uint64, stateID uint64) error {
	if err := r.db.WithContext(ctx).Model(&models.Project{}).Where("id = ?", id).Update("state_id", stateID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("proyecto no encontrado para el id: %d", id)
		}

		return err
	}

	return nil
}

func (r *PostgresRepository) UpdateCounter(ctx context.Context, id uint64, counter int64) error {
	if err := r.db.WithContext(ctx).Model(&models.Project{}).Where("id = ?", id).Update("counter_likes", counter).Error; err != nil {
		return err
	}

	return nil
}
