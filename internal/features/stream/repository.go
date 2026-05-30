package stream

import (
	"context"
	"cv_api/internal/models"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) VideoRepository {
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

func (r *PostgresRepository) GetAllData(ctx context.Context) ([]models.Video, error) {
	var videos []models.Video

	err := r.db.WithContext(ctx).Model(&models.Video{}).Order("created_at DESC").Find(&videos).Error
	if err != nil {
		return nil, err
	}

	return videos, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context, offset int, limit int, all bool) ([]models.Video, int64, error) {
	var videos []models.Video
	var total int64

	err := r.db.WithContext(ctx).Model(&models.Video{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []models.Video{}, 0, nil
	}

	if all {
		err = r.db.WithContext(ctx).Model(&models.Video{}).Order("created_at DESC").Find(&videos).Error
	} else {
		err = r.db.WithContext(ctx).Model(&models.Video{}).Offset(offset).Limit(limit).Order("created_at DESC").Find(&videos).Error
	}

	if err != nil {
		return nil, 0, err
	}

	return videos, total, nil
}

func (r *PostgresRepository) GetAllRecents(ctx context.Context, max int) ([]models.Video, error) {
	var videos []models.Video

	err := r.db.WithContext(ctx).Model(&models.Video{}).Limit(max).Order("created_at DESC").Find(&videos).Error
	if err != nil {
		return nil, err
	}

	return videos, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uint64) (*models.Video, error) {
	var video models.Video

	if err := r.db.WithContext(ctx).First(&video, id).Error; err != nil {
		return nil, err
	}

	return &video, nil
}

func (r *PostgresRepository) ExistsById(ctx context.Context, id uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Video{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *PostgresRepository) Create(ctx context.Context, data *models.Video) error {
	if err := r.db.WithContext(ctx).Create(data).Error; err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) Update(ctx context.Context, id uint64, data *models.Video) error {
	if err := r.db.WithContext(ctx).Model(&models.Video{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}
