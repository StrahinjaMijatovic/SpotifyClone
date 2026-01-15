package handlers

import (
	"example.com/api-gateway/proxy"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		// Users service routes
		api.POST("/register", proxy.ProxyToUsersService)
		api.POST("/login", proxy.ProxyToUsersService)
		api.POST("/verify-otp", proxy.ProxyToUsersService)
		api.GET("/verify-email", proxy.ProxyToUsersService)
		api.POST("/magic-link", proxy.ProxyToUsersService)
		api.GET("/magic-login", proxy.ProxyToUsersService)
		api.POST("/reset-password", proxy.ProxyToUsersService)
		api.POST("/reset-password/confirm", proxy.ProxyToUsersService)
		api.POST("/change-password", proxy.ProxyToUsersService)
		api.GET("/profile", proxy.ProxyToUsersService)
		api.PUT("/profile", proxy.ProxyToUsersService)
		api.DELETE("/profile", proxy.ProxyToUsersService)
		api.POST("/logout", proxy.ProxyToUsersService)

		// Legacy auth routes (keeping for compatibility)
		api.POST("/auth/register", proxy.ProxyToUsersService)
		api.POST("/auth/login", proxy.ProxyToUsersService)
		api.POST("/auth/reset-password", proxy.ProxyToUsersService)
		api.POST("/auth/reset-password/confirm", proxy.ProxyToUsersService)
		api.POST("/auth/change-password", proxy.ProxyToUsersService)
		api.GET("/auth/profile", proxy.ProxyToUsersService)

		// Content service routes
		api.GET("/genres", proxy.ProxyToContentService)
		api.GET("/artists", proxy.ProxyToContentService)
		api.GET("/artists/:id", proxy.ProxyToContentService)
		api.GET("/albums", proxy.ProxyToContentService)
		api.GET("/albums/:id", proxy.ProxyToContentService)
		api.GET("/songs", proxy.ProxyToContentService)
		api.GET("/search", proxy.ProxyToContentService)

		// Admin content routes
		api.POST("/artists", proxy.ProxyToContentService)
		api.PUT("/artists/:id", proxy.ProxyToContentService)
		api.POST("/albums", proxy.ProxyToContentService)
		api.POST("/songs", proxy.ProxyToContentService)

		// Ratings service routes
		api.POST("/ratings", proxy.ProxyToRatingsService)
		api.GET("/ratings", proxy.ProxyToRatingsService)
		api.GET("/ratings/:songId", proxy.ProxyToRatingsService)
		api.DELETE("/ratings/:songId", proxy.ProxyToRatingsService)

		// Subscriptions service routes
		api.POST("/subscriptions", proxy.ProxyToSubscriptionsService)
		api.GET("/subscriptions", proxy.ProxyToSubscriptionsService)
		api.DELETE("/subscriptions/:id", proxy.ProxyToSubscriptionsService)

		// Notifications service routes
		api.GET("/notifications", proxy.ProxyToNotificationsService)
		api.PUT("/notifications/:id/read", proxy.ProxyToNotificationsService)

		// Recommendations service routes
		api.GET("/recommendations", proxy.ProxyToRecommendationService)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "api-gateway"})
	})
}
