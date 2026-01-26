package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"example.com/subsriptions-service/handlers"
	"example.com/subsriptions-service/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8004"
	}

	redisURI := os.Getenv("REDIS_URI")
	if redisURI == "" {
		redisURI = "redis://localhost:6380"
	}

	opt, err := redis.ParseURL(redisURI)
	if err != nil {
		log.Fatal("Failed to parse Redis URI:", err)
	}

	redisClient = redis.NewClient(opt)

	ctx := context.Background()
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	log.Println("Connected to Redis")

	router := gin.Default()
	// CORS is handled by the API Gateway
	// router.Use(corsMiddleware())

	handlers.InitHandlers(redisClient)
	setupRoutes(router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("Subscriptions service starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	if err := redisClient.Close(); err != nil {
		log.Fatal("Failed to close Redis connection:", err)
	}

	log.Println("Server exited")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func setupRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		api.POST("/subscriptions", middleware.AuthMiddleware(), handlers.CreateSubscription)
		api.GET("/subscriptions", middleware.AuthMiddleware(), handlers.GetSubscriptions)
		api.GET("/subscriptions/:target_id", middleware.AuthMiddleware(), handlers.CheckSubscription)
		api.DELETE("/subscriptions/:id", middleware.AuthMiddleware(), handlers.DeleteSubscription)
		api.GET("/subscriptions/followers/:artist_id", handlers.GetFollowersByArtist) // Called by content-service
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "subscriptions-service"})
	})
}
