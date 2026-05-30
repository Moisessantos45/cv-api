package routes

import (
	"cv_api/internal/features/profile"

	"github.com/gin-gonic/gin"
)

func ProfileRoutes(rg *gin.RouterGroup) {

	ph := profile.NewProfileHandler(profileUC)

	protected := rg.Group("/profile")
	protected.Use(authMiddleware())
	{
		// protected.GET("/dashboard-metrics", h.GetDashboardMetrics)
		protected.POST("", ph.Create)
		protected.GET("", ph.GetById)
		protected.PUT("", ph.Update)
	}
}
