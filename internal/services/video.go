package services

import (
	"context"
	"cv_api/internal/models"
	"cv_api/internal/repository"
)

type VideoService struct {
	repo *repository.VideoRepository
}

func NewVideoService(repo *repository.VideoRepository) *VideoService {
	return &VideoService{repo: repo}
}

func (s *VideoService) GetVideos(ctx context.Context) ([]models.Video, error) {
	return s.repo.FindAll(ctx)
}
