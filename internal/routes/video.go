package routes

import (
	"cv_api/config"
	"cv_api/internal/handlers"
	"cv_api/internal/repository"
	"cv_api/internal/services"
	"time"

	"github.com/gin-gonic/gin"
)

func VideoRoutes(api *gin.RouterGroup) {
	client := config.Client
	clientRd := config.Rdb
	videoRepo := repository.NewVideoRepository(client, clientRd, time.Minute*15)
	videoSvc := services.NewVideoService(videoRepo)
	videoHandler := handlers.NewVideoHandlers(videoSvc)

	video := api.Group("/video")
	{
		video.GET("", videoHandler.GetVideos)
		// auth.POST("/register", registerHandler)
	}
}
