package routes

import (
	"cv_api/config"
	"cv_api/config/db"
	"cv_api/internal/features/post"
	"cv_api/internal/shared/middleware"
	"cv_api/internal/shared/service"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func PostRoutes(rg *gin.RouterGroup) {
	rd := config.Rdb
	cache := service.NewCacheService(rd)

	postRepo := post.NewPostgresRepository(db.DB)
	postSvc := post.NewPostService(postRepo, cache, profileUC)
	ph := post.NewPostHandlers(postSvc, profileUC)

	public := rg.Group("")
	public.Use(middleware.RateLimiterMiddleware(rate.Every(time.Minute/60), 60))
	{
		public.GET("/post/public", ph.GetAllPublic)
		public.GET("/post/recent", ph.GetAllRecents)
		public.GET("/post/:slug", ph.GetBySlugPublic)
	}

	protected := rg.Group("/post")
	protected.Use(authMiddleware())
	protected.Use(idempotencyMiddleware())
	{
		protected.GET("/states", ph.GetAllStates)
		protected.GET("/user", ph.GetAll)
		protected.GET("/slug-private/:slug", ph.GetBySlugPrivate)
		protected.POST("", ph.Create)
		protected.PUT("/:slug", ph.Update)
	}
}
