package main

import (
	"example.com/content-service/handlers"
	"example.com/content-service/middleware"
	"github.com/gin-gonic/gin"
)

func setupRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		// Public routes
		api.GET("/genres", handlers.GetGenres)
		api.GET("/artists", handlers.GetArtists)
		api.GET("/artists/:id", handlers.GetArtist)
		api.GET("/albums", handlers.GetAlbums)
		api.GET("/albums/:id", handlers.GetAlbum)
		api.GET("/songs", handlers.GetSongs)
		api.GET("/search", handlers.SearchContent)

		// Admin routes
		admin := api.Group("/")
		admin.Use(middleware.AuthMiddleware())
		admin.Use(middleware.AdminMiddleware())
		{
			admin.POST("/artists", handlers.CreateArtist)
			admin.PUT("/artists/:id", handlers.UpdateArtist)
			admin.POST("/albums", handlers.CreateAlbum)
			admin.POST("/songs", handlers.CreateSong)
			admin.DELETE("/songs/:id", handlers.DeleteSong)
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "content-service"})
	})
}
