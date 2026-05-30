package routes

import (
	"cv_api/config"
	"cv_api/config/db"
	"cv_api/internal/features/project"
	"cv_api/internal/shared/middleware"
	"cv_api/internal/shared/service"
	"cv_api/internal/shared/utils"

	"github.com/gin-gonic/gin"
)

func ProjectRoutes(rg *gin.RouterGroup) {
	rd := config.Rdb
	cache := service.NewCacheService(rd)
	maker := utils.NewPasetoMaker()

	projectRepo := project.NewPostgresRepository(db.DB)
	projectSvc := project.NewProjectService(projectRepo, cache)
	ph := project.NewProjectHandlers(projectSvc)

	// rg.GET("/project/data", ph.GetAllData)
	rg.GET("/project/public", ph.GetAllPublic)
	rg.GET("/project/recent", ph.GetALLRecents)
	rg.GET("/project/:slug", ph.GetBySlugPublic)
	protected := rg.Group("/project")

	protected.Use(middleware.AuthMiddleware(maker, rd))
	{
		protected.GET("", ph.GetAll)
		protected.GET("/slug/:slug", ph.GetBySlug)
		protected.POST("", ph.Create)
		protected.POST("/all", ph.CreateAll)
		protected.PUT("/:id", ph.Update)
		protected.PATCH("/:id/state", ph.UpdateState)
		protected.PUT("/likes/:id", ph.UpdateCounter)
	}
}
