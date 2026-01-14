package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"example.com/users-service/handlers"
	"example.com/users-service/middleware"
	"example.com/users-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	usersDB     *mongo.Database
	redisClient *redis.Client
)

func main() {
	// Initialize logger
	if err := utils.InitLogger(); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer utils.CloseLogger()

	// Start log rotation goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			utils.RotateLogs()
		}
	}()

	// Load configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}

	mongodbURI := os.Getenv("MONGODB_URI")
	if mongodbURI == "" {
		mongodbURI = "mongodb://localhost:27017/users"
	}

	redisURI := os.Getenv("REDIS_URI")
	if redisURI == "" {
		redisURI = "redis://localhost:6379"
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongodbURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	mongoClient = client
	usersDB = client.Database("users")

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}
	log.Println("Connected to MongoDB")

	// Connect to Redis
	opt, err := redis.ParseURL(redisURI)
	if err != nil {
		log.Fatal("Failed to parse Redis URI:", err)
	}
	redisClient = redis.NewClient(opt)

	redisCtx := context.Background()
	if _, err := redisClient.Ping(redisCtx).Result(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	log.Println("Connected to Redis")

	// Setup router
	router := gin.Default()
	router.Use(corsMiddleware())

	// Initialize handlers
	handlers.InitHandlers(usersDB, redisClient)
	handlers.EnsureUserIndexes(usersDB)
	
	// Initialize middleware
	middleware.InitAuthMiddleware(redisClient)

	// Routes
	setupRoutes(router)

	// Start server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("Users service starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	if err := mongoClient.Disconnect(ctx); err != nil {
		log.Fatal("Failed to disconnect from MongoDB:", err)
	}

	if err := redisClient.Close(); err != nil {
		log.Fatal("Failed to close Redis connection:", err)
	}

	log.Println("Server exited")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "false")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
