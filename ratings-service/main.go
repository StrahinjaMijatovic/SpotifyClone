package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"example.com/ratings-service/handlers"
	"example.com/ratings-service/middleware"
	"example.com/ratings-service/tracing"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8003"
	}

	// Inicijalizuj distributed tracing
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "ratings-service"
	}

	tp, err := tracing.InitTracer(serviceName)
	if err != nil {
		log.Printf("Warning: Failed to initialize tracing: %v", err)
	} else {
		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				log.Printf("Error shutting down tracer: %v", err)
			}
		}()
		log.Println("Distributed tracing initialized")
	}

	redisURI := os.Getenv("REDIS_URI")
	if redisURI == "" {
		redisURI = "redis://localhost:6379"
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
	router.Use(corsMiddleware())
	// Dodaj tracing middleware
	router.Use(tracing.TracingMiddleware(serviceName))

	handlers.InitHandlers(redisClient)
	setupRoutes(router)

	// TLS Configuration
	tlsEnabled := os.Getenv("TLS_ENABLED")
	certFile := os.Getenv("TLS_CERT_FILE")
	keyFile := os.Getenv("TLS_KEY_FILE")

	if certFile == "" {
		certFile = "certs/cert.pem"
	}
	if keyFile == "" {
		keyFile = "certs/key.pem"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	go func() {
		if tlsEnabled == "true" {
			log.Printf("Ratings service starting on HTTPS port %s", port)
			if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Failed to start HTTPS server: %v", err)
			}
		} else {
			log.Printf("Ratings service starting on port %s", port)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Failed to start server: %v", err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctxShutdown)
	redisClient.Close()
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")

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
		api.POST("/ratings", middleware.AuthMiddleware(), handlers.CreateRating)
		api.GET("/ratings", middleware.AuthMiddleware(), handlers.GetRatings)
		api.GET("/ratings/:songId", handlers.GetSongRatings)
		api.DELETE("/ratings/:songId", middleware.AuthMiddleware(), handlers.DeleteRating)

		// Admin route - delete all ratings for a song (used when song is deleted)
		api.DELETE("/ratings/:songId/all", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.DeleteAllSongRatings)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
