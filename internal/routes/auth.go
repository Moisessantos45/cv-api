package routes

import (
	"cv_api/internal/features/auth"
	"cv_api/internal/shared/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func AuthRoutes(rg *gin.RouterGroup) {
	s := authUc
	h := auth.NewAuthHandler(s)

	public := rg.Group("")
	public.Use(middleware.RateLimiterMiddleware(rate.Every(time.Minute/20), 20))
	{
		public.POST("/login", h.Login)
		public.POST("/forward-email-verification", h.ForwardEmailVerification)
		public.POST("/forgot-password", h.SendPasswordReset)
		public.POST("/register", h.Register)
		public.POST("/logout", h.Logout)
		public.POST("/refresh-token", h.RefreshToken)
	}

	protected := rg.Group("/")
	protected.Use(authMiddleware())
	protected.Use(middleware.RateLimiterMiddleware(rate.Every(time.Minute/300), 300))
	{
		protected.GET("/confirm-account", h.ConfirmAccount)
		protected.GET("/session", h.GetSession)
		protected.POST("/verify-email", h.VerifyEmail)
		protected.PATCH("/reset-password", h.ResetPassword)
		protected.PATCH("/change-password", h.UpdatePassword)

		protected.GET("/generate-two-factor", h.TestingCrateTWOFA)
	}
}
