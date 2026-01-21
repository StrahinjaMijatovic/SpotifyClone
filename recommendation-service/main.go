package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"example.com/recommendation-service/handlers"
	"example.com/recommendation-service/middleware"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var neo4jDriver neo4j.DriverWithContext

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8006"
	}

	neo4jURI := os.Getenv("NEO4J_URI")
	if neo4jURI == "" {
		neo4jURI = "bolt://localhost:7687"
	}

	neo4jUser := os.Getenv("NEO4J_USER")
	if neo4jUser == "" {
		neo4jUser = "neo4j"
	}

	neo4jPassword := os.Getenv("NEO4J_PASSWORD")
	if neo4jPassword == "" {
		neo4jPassword = "password"
	}

	driver, err := neo4j.NewDriverWithContext(
		neo4jURI,
		neo4j.BasicAuth(neo4jUser, neo4jPassword, ""),
	)
	if err != nil {
		log.Fatal("Failed to create Neo4j driver:", err)
	}
	neo4jDriver = driver

	// ✅ Proveri konekciju odmah (auth/URI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := neo4jDriver.VerifyConnectivity(ctx); err != nil {
		_ = neo4jDriver.Close(ctx)
		log.Fatal("Failed to connect to Neo4j:", err)
	}
	log.Println("Connected to Neo4j")

	router := gin.Default()
	router.Use(corsMiddleware())

	handlers.InitHandlers(neo4jDriver)
	setupRoutes(router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
	srv.ListenAndServeTLS("cert.pem", "key.pem")

	go func() {
		log.Printf("Recommendation service starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	if err := neo4jDriver.Close(shutdownCtx); err != nil {
		log.Fatal("Failed to close Neo4j driver:", err)
	}

	log.Println("Server exited")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, DELETE, OPTIONS")

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
		api.GET("/recommendations", middleware.AuthMiddleware(), handlers.GetRecommendations)

		// Admin route - delete song from recommendation graph (used when song is deleted)
		api.DELETE("/recommendations/songs/:songId", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.DeleteSong)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "recommendation-service"})
	})
}
