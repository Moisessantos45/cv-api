package repository

import (
	"context"
	"cv_api/internal/models"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/supabase-community/supabase-go"
)

type VideoRepository struct {
	sp   *supabase.Client
	rd   *redis.Client
	cTTL time.Duration
}

const videosCacheKey = "videos:all"

func NewVideoRepository(sp *supabase.Client, rd *redis.Client, cTTL time.Duration) *VideoRepository {
	return &VideoRepository{
		sp:   sp,
		rd:   rd,
		cTTL: cTTL,
	}
}

func (r *VideoRepository) FindAll(ctx context.Context) ([]models.Video, error) {
	cachedData, err := r.rd.Get(ctx, videosCacheKey).Bytes()
	if err == nil {
		var videos []models.Video
		if json.Unmarshal(cachedData, &videos) == nil {
			return videos, nil
		}
	} else if err != redis.Nil {
		log.Printf("Error al obtener videos de Redis: %v", err)
	}

	data, _, err := r.sp.From("Link_video").Select("*", "exact", false).Execute()
	if err != nil {
		return nil, err
	}

	var videos []models.Video
	if err := json.Unmarshal(data, &videos); err != nil {
		return nil, err
	}

	go func() {
		videosBytes, _ := json.Marshal(videos)
		_ = r.rd.Set(ctx, videosCacheKey, videosBytes, r.cTTL).Err()
	}()

	return videos, nil

}
