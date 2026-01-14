package main

import (
	"time"

	"example.com/users-service/handlers"
	"example.com/users-service/middleware"
	"github.com/gin-gonic/gin"
)

func setupRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")

	// Apply global rate limiting to all API routes
	api.Use(middleware.RateLimitMiddleware(60, time.Minute))

	{
		// Public routes
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)

		api.POST("/verify-otp", handlers.VerifyOTP)
		api.GET("/verify-email", handlers.VerifyEmail)
		api.POST("/magic-link", handlers.RequestMagicLink)
		api.GET("/magic-login", handlers.MagicLogin)
		api.POST("/reset-password", handlers.ResetPassword)
		api.POST("/reset-password/confirm", handlers.ResetPasswordConfirm)

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/change-password", handlers.ChangePassword)
			protected.GET("/profile", handlers.GetProfile)
			protected.PUT("/profile", handlers.UpdateProfile)
			protected.DELETE("/profile", handlers.DeleteAccount)
			protected.POST("/logout", handlers.Logout)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "users-service"})
	})
}
