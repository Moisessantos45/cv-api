package routes

import (
	"cv_api/config"
	"cv_api/config/db"
	"cv_api/internal/features/stream"
	"cv_api/internal/shared/middleware"
	"cv_api/internal/shared/service"
	"cv_api/internal/shared/utils"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func VideoRoutes(rg *gin.RouterGroup) {
	rd := config.Rdb
	cache := service.NewCacheService(rd)
	maker := utils.NewPasetoMaker()

	videoRepo := stream.NewPostgresRepository(db.DB)
	videoSvc := stream.NewVideoUseCase(videoRepo, cache)
	vh := stream.NewVideoHandlers(videoSvc)

	public := rg.Group("")
	public.Use(middleware.RateLimiterMiddleware(rate.Every(time.Minute/60), 60))
	{
		public.GET("/stream", vh.GetAll)
		public.GET("/stream/recent", vh.GetAllRecents)
		public.GET("/stream/:id", vh.GetByID)
	}

	protected := rg.Group("/stream")
	protected.Use(middleware.AuthMiddleware(maker, rd))
	protected.Use(middleware.RateLimiterMiddleware(rate.Every(time.Minute/300), 300))
	{
		protected.POST("", vh.Create)
		protected.POST("/all", vh.CreateAll)
		protected.PUT("", vh.Update)
	}
}
