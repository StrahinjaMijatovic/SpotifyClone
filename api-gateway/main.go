package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"example.com/api-gateway/handlers"
	"example.com/api-gateway/middleware"
	"example.com/api-gateway/tracing"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Inicijalizuj distributed tracing
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "api-gateway"
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
		log.Println("Distributed tracing initialized with Jaeger")
	}

	router := gin.Default()

	router.Use(corsMiddleware())
	router.Use(middleware.RateLimitMiddleware())

	// Dodaj tracing middleware
	router.Use(tracing.TracingMiddleware(serviceName))

	handlers.SetupRoutes(router)

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

	if tlsEnabled == "true" {
		log.Printf("API Gateway starting on HTTPS port %s", port)
		if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil {
			log.Fatal("Failed to start HTTPS server:", err)
		}
	} else {
		log.Printf("API Gateway starting on HTTP port %s (TLS disabled)", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:4200"
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", frontendURL)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
