package routes

import (
	"cv_api/internal/features/profile"
	"cv_api/internal/shared/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func ProfileRoutes(rg *gin.RouterGroup) {

	ph := profile.NewProfileHandler(profileUC)

	protected := rg.Group("/profile")
	protected.Use(authMiddleware())
	protected.Use(middleware.RateLimiterMiddleware(rate.Every(time.Minute/300), 300))
	{
		// protected.GET("/dashboard-metrics", h.GetDashboardMetrics)
		protected.POST("", ph.Create)
		protected.GET("", ph.GetById)
		protected.PUT("", ph.Update)
	}
}
