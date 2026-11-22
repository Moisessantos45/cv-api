package routes

import (
	"cv_api/config"
	"cv_api/internal/features/auth"
	"cv_api/internal/shared/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func AuthRoutes(rg *gin.RouterGroup) {
	s := authUc
	h := auth.NewAuthHandler(s)
	strict := middleware.RateLimiterRedis(config.Rdb, "auth", rate.Every(time.Minute/10), 5)

	public := rg.Group("")
	public.Use(idempotencyMiddleware())
	{
		public.POST("/login", strict, h.Login)
		public.POST("/forward-email-verification", strict, h.ForwardEmailVerification)
		public.POST("/forgot-password", strict, h.SendPasswordReset)
		public.POST("/register", strict, h.Register)
		public.POST("/logout", h.Logout)
		public.POST("/refresh-token", h.RefreshToken)
	}

	protected := rg.Group("/")
	protected.Use(authMiddleware())
	protected.Use(idempotencyMiddleware())
	{
		protected.GET("/confirm-account", h.ConfirmAccount)
		protected.GET("/session", h.GetSession)
		protected.POST("/verify-email", h.VerifyEmail)
		protected.PATCH("/reset-password", h.ResetPassword)
		protected.PATCH("/change-password", h.UpdatePassword)

		protected.GET("/generate-two-factor", h.TestingCrateTWOFA)
	}
}
