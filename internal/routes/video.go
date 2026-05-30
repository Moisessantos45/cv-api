package routes

import (
	"cv_api/config"
	"cv_api/config/db"
	"cv_api/internal/features/stream"
	"cv_api/internal/shared/middleware"
	"cv_api/internal/shared/service"
	"cv_api/internal/shared/utils"

	"github.com/gin-gonic/gin"
)

func VideoRoutes(rg *gin.RouterGroup) {
	rd := config.Rdb
	cache := service.NewCacheService(rd)
	maker := utils.NewPasetoMaker()

	videoRepo := stream.NewPostgresRepository(db.DB)
	videoSvc := stream.NewVideoUseCase(videoRepo, cache)
	vh := stream.NewVideoHandlers(videoSvc)

	//rg.GET("/stream/data", vh.GetAllData)
	rg.GET("/stream", vh.GetAll)
	//rg.GET("/stream/all", vh.GetAll)
	rg.GET("/stream/recent", vh.GetAllRecents)
	rg.GET("/stream/:id", vh.GetByID)
	protected := rg.Group("/stream")

	protected.Use(middleware.AuthMiddleware(maker, rd))
	{
		protected.POST("", vh.Create)
		protected.POST("/all", vh.CreateAll)
		protected.PUT("", vh.Update)
	}
}
