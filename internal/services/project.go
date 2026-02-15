package services

import (
	"context"
	"cv_api/internal/models"
	"cv_api/internal/repository"
)

type ProjectService struct {
	repo *repository.ProjectRepository
}

func NewProjectService(repo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) GetProjects(ctx context.Context) ([]models.Project, error) {
	return s.repo.FindAll(ctx)
}

func (s *ProjectService) GetProject(ctx context.Context, id string) (*models.Project, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ProjectService) UpdateCounter(ctx context.Context, id string, value int) (models.Project, error) {
	return s.repo.UpdateCounter(id, value)
}
