package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"example.com/users-service/config"
	"example.com/users-service/handlers"
	"example.com/users-service/middleware"
	"example.com/users-service/tracing"
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

	// Inicijalizuj distributed tracing
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "users-service"
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

	// Initialize password expiry config
	config.InitPasswordConfig()
	log.Printf("Password expiry configured: max age = %s", config.GetPasswordMaxAgeString())

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

	// Dodaj tracing middleware
	router.Use(tracing.TracingMiddleware(serviceName))

	// Initialize handlers
	handlers.InitHandlers(usersDB, redisClient)
	handlers.EnsureUserIndexes(usersDB)

	// Initialize middleware
	middleware.InitAuthMiddleware(redisClient)

	// Routes
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

	// Start server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			},
		},
	}

	go func() {
		if tlsEnabled == "true" {
			log.Printf("Users service starting on HTTPS port %s", port)
			if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Failed to start HTTPS server: %v", err)
			}
		} else {
			log.Printf("Users service starting on HTTP port %s (TLS disabled)", port)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Failed to start server: %v", err)
			}
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
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:4200"
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
