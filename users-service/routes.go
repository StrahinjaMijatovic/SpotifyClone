package main

import (
	"example.com/users-service/handlers"
	"example.com/users-service/middleware"
	"github.com/gin-gonic/gin"
)

func setupRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		// Public routes
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)
		api.POST("/reset-password", handlers.ResetPassword)
		api.POST("/reset-password/confirm", handlers.ResetPasswordConfirm)

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/change-password", handlers.ChangePassword)
			protected.GET("/profile", handlers.GetProfile)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "users-service"})
	})
}
