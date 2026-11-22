package profile

import (
	"context"
	"cv_api/internal/models"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) ProfileRepository {
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

func (r *PostgresRepository) GetAll(ctx context.Context, offset int, limit int, all bool) ([]models.Profile, int64, error) {
	var profiles []models.Profile
	var total int64

	err := r.db.WithContext(ctx).Model(&models.Profile{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []models.Profile{}, 0, nil
	}

	if all {
		err = r.db.WithContext(ctx).Model(&models.Profile{}).Order("created_at DESC").Find(&profiles).Error
	} else {
		err = r.db.WithContext(ctx).Model(&models.Profile{}).Offset(offset).Limit(limit).Order("created_at DESC").Find(&profiles).Error
	}

	if err != nil {
		return nil, 0, err
	}

	return profiles, total, nil
}

func (r *PostgresRepository) GetAllRecents(ctx context.Context, max int) ([]models.Profile, error) {
	var profiles []models.Profile

	err := r.db.WithContext(ctx).Model(&models.Profile{}).Limit(max).Order("created_at DESC").Find(&profiles).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []models.Profile{}, nil
		}

		return nil, err
	}

	return profiles, nil
}

func (r *PostgresRepository) GetBasicByAuthID(ctx context.Context, authID uint64) (*ProfileBasic, error) {
	var profile ProfileBasic

	if err := r.db.WithContext(ctx).Model(&models.Profile{}).Select("id, auth_id").Where("auth_id = ?", authID).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &profile, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uint64) (*models.Profile, error) {
	var profile models.Profile

	if err := r.db.WithContext(ctx).Model(&models.Profile{}).Where("auth_id = ?", id).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &profile, nil
}

func (r *PostgresRepository) GetByAuthID(ctx context.Context, authID uint64) (*models.Profile, error) {
	var profile models.Profile

	if err := r.db.WithContext(ctx).Model(&models.Profile{}).Where("auth_id = ?", authID).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &profile, nil
}

func (r *PostgresRepository) ExistsById(ctx context.Context, id uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Profile{}).Where("auth_id = ?", id).Count(&count).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, err
	}

	return count > 0, nil
}

func (r *PostgresRepository) Create(ctx context.Context, data *models.Profile) error {
	if err := r.db.WithContext(ctx).Create(data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}

		return err
	}
	return nil
}

func (r *PostgresRepository) Update(ctx context.Context, id uint64, data *models.Profile) error {
	if err := r.db.WithContext(ctx).Model(&models.Profile{}).Where("auth_id = ?", id).Updates(data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}

		return err
	}

	return nil
}
