package repository

import (
	"context"
	"cv_api/config"
	"cv_api/internal/models"
	"cv_api/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/supabase-community/supabase-go"
)

type ProjectRepository struct {
	sp   *supabase.Client
	rd   *redis.Client
	cTTL time.Duration
}

const projectsCacheKey = "projects:all"

func NewProjectRepository(sp *supabase.Client, rd *redis.Client, cTTL time.Duration) *ProjectRepository {
	return &ProjectRepository{
		sp:   config.Client,
		rd:   rd,
		cTTL: cTTL,
	}
}

func (r *ProjectRepository) FindAll(ctx context.Context) ([]models.Project, error) {
	cachedData, err := r.rd.Get(ctx, projectsCacheKey).Bytes()
	if err == nil {
		var projects []models.Project
		if json.Unmarshal(cachedData, &projects) == nil {
			return projects, nil
		}
	} else if err != redis.Nil {
		log.Printf("Error al obtener proyectos de Redis: %v", err)
	}

	data, _, err := r.sp.From("Proyectos").Select("*", "exact", false).Execute()
	if err != nil {
		return nil, err
	}

	projects, err := utils.ParseJSONToStringArray(data)
	if err != nil {
		return nil, err
	}

	go func() {
		projectBytes, _ := json.Marshal(projects)
		_ = r.rd.Set(ctx, projectsCacheKey, projectBytes, r.cTTL).Err()
	}()

	return projects, nil
}

func (r *ProjectRepository) FindByID(ctx context.Context, id string) (*models.Project, error) {
	cacheKey := fmt.Sprintf("project:%s", id)

	cachedData, err := r.rd.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var project models.Project
		if json.Unmarshal(cachedData, &project) == nil {
			return &project, nil
		}
	} else if err != redis.Nil {
		log.Printf("Error al obtener proyecto %s de Redis: %v", id, err)
	}

	data, _, err := r.sp.From("Proyectos").Select("*", "exact", false).Eq("id", id).Execute()
	if err != nil {
		return nil, err
	}

	projects, err := utils.ParseJSONToStringArray(data)
	if err != nil {
		return nil, err
	}

	if len(projects) == 0 {
		return nil, nil
	}

	project := &projects[0]

	go func() {
		projectBytes, _ := json.Marshal(project)
		_ = r.rd.Set(ctx, cacheKey, projectBytes, r.cTTL).Err()
	}()

	return project, nil
}

func (r *ProjectRepository) UpdateCounter(id string, counter int) (models.Project, error) {
	res, count, err := r.sp.
		From("Proyectos").
		Update(map[string]any{
			"counter_likes": counter,
		}, "", "exact").
		Eq("id", id).
		Execute()

	if err != nil {
		return models.Project{}, err
	}

	if count == 0 {
		return models.Project{}, nil
	}

	projects, err := utils.ParseJSONToStringArray(res)
	if err != nil {
		return models.Project{}, err
	}

	if len(projects) == 0 {
		return models.Project{}, nil
	}

	return projects[0], nil
}
