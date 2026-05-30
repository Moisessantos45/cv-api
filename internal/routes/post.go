package routes

import (
	"cv_api/config"
	"cv_api/config/db"
	"cv_api/internal/features/post"
	"cv_api/internal/shared/middleware"
	"cv_api/internal/shared/service"
	"cv_api/internal/shared/utils"

	"github.com/gin-gonic/gin"
)

func PostRoutes(rg *gin.RouterGroup) {
	rd := config.Rdb
	cache := service.NewCacheService(rd)
	maker := utils.NewPasetoMaker()

	postRepo := post.NewPostgresRepository(db.DB)
	postSvc := post.NewPostService(postRepo, cache, profileUC)
	ph := post.NewPostHandlers(postSvc, profileUC)

	rg.GET("/post/public", ph.GetAllPublic)
	rg.GET("/post/recent", ph.GetAllRecents)
	rg.GET("/post/:slug", ph.GetBySlugPublic)
	protected := rg.Group("/post")
	protected.Use(middleware.AuthMiddleware(maker, rd))
	{
		protected.GET("/states", ph.GetAllStates)
		protected.GET("/user", ph.GetAll)
		protected.GET("/slug-private/:slug", ph.GetBySlugPrivate)
		protected.POST("", ph.Create)
		protected.PUT("/:slug", ph.Update)
	}
}
