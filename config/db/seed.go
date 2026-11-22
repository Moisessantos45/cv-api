package db

import "cv_api/internal/models"

func SeedDatabase() error {
	var count int64
	if err := DB.Model(&models.StatePost{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		if err := DB.Create(&models.StatePosts).Error; err != nil {
			return err
		}
	}

	return nil
}
