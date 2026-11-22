package experience

import (
	"context"
	"cv_api/internal/models"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) ExperienceRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) GetAll(ctx context.Context, offset int, limit int, all bool) ([]models.Experience, int64, error) {
	var experiences []models.Experience
	var total int64

	err := r.db.WithContext(ctx).Model(&models.Experience{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []models.Experience{}, 0, nil
	}

	if all {
		err = r.db.WithContext(ctx).Model(&models.Experience{}).Order("created_at DESC").Find(&experiences).Error
	} else {
		err = r.db.WithContext(ctx).Model(&models.Experience{}).Offset(offset).Limit(limit).Order("created_at DESC").Find(&experiences).Error
	}

	if err != nil {
		return nil, 0, err
	}

	return experiences, total, nil
}

func (r *PostgresRepository) GetAllRecents(ctx context.Context, max int) ([]models.Experience, error) {
	var experiences []models.Experience

	err := r.db.WithContext(ctx).Model(&models.Experience{}).Limit(max).Order("created_at DESC").Find(&experiences).Error
	if err != nil {
		return nil, err
	}

	return experiences, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uint64) (*models.Experience, error) {
	var experience models.Experience

	if err := r.db.WithContext(ctx).First(&experience, id).Error; err != nil {
		return nil, err
	}

	return &experience, nil
}

func (r *PostgresRepository) ExistsById(ctx context.Context, id uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Experience{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *PostgresRepository) Create(ctx context.Context, data *models.Experience) error {
	if err := r.db.WithContext(ctx).Create(data).Error; err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) Update(ctx context.Context, id uint64, data *models.Experience) error {
	if err := r.db.WithContext(ctx).Model(&models.Experience{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}
