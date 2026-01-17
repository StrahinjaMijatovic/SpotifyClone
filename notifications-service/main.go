package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"example.com/notifications-service/handlers"
	"example.com/notifications-service/middleware"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

const (
	maxRetries     = 30              // Maksimalan broj pokušaja
	initialBackoff = 2 * time.Second // Početno čekanje
	maxBackoff     = 30 * time.Second // Maksimalno čekanje između pokušaja
)

var cassandraSession *gocql.Session

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8005"
	}

	hostsEnv := os.Getenv("CASSANDRA_HOSTS")
	if hostsEnv == "" {
		hostsEnv = "localhost"
	}
	hosts := parseHosts(hostsEnv)

	keyspace := os.Getenv("CASSANDRA_KEYSPACE")
	if keyspace == "" {
		keyspace = "notifications"
	}

	// 1) Prvo napravi keyspace/tabelu (pre konekcije na keyspace)
	if err := createKeyspaceAndTable(hosts, keyspace); err != nil {
		log.Fatal("Failed to init Cassandra schema:", err)
	}

	// 2) Onda konekcija na keyspace (sa retry logikom)
	session, err := connectToKeyspace(hosts, keyspace)
	if err != nil {
		log.Fatal("Failed to connect to Cassandra:", err)
	}
	cassandraSession = session
	defer cassandraSession.Close()

	log.Println("Connected to Cassandra keyspace:", keyspace)

	router := gin.Default()
	router.Use(corsMiddleware())

	handlers.InitHandlers(cassandraSession)
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
		log.Printf("Notifications service starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func parseHosts(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return []string{"localhost"}
	}
	return out
}

// waitForCassandra čeka da Cassandra postane dostupna sa retry logikom
func waitForCassandra(hosts []string) (*gocql.Session, error) {
	var session *gocql.Session
	var err error
	backoff := initialBackoff

	for i := 0; i < maxRetries; i++ {
		cluster := gocql.NewCluster(hosts...)
		cluster.Consistency = gocql.Quorum
		cluster.Timeout = 10 * time.Second
		cluster.ConnectTimeout = 10 * time.Second

		session, err = cluster.CreateSession()
		if err == nil {
			log.Println("Successfully connected to Cassandra")
			return session, nil
		}

		log.Printf("Waiting for Cassandra (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(backoff)

		// Eksponencijalni backoff sa maksimumom
		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
		}
	}

	return nil, fmt.Errorf("failed to connect to Cassandra after %d attempts: %v", maxRetries, err)
}

// connectToKeyspace kreira sesiju na specifičan keyspace sa retry logikom
func connectToKeyspace(hosts []string, keyspace string) (*gocql.Session, error) {
	var session *gocql.Session
	var err error
	backoff := initialBackoff

	for i := 0; i < maxRetries; i++ {
		cluster := gocql.NewCluster(hosts...)
		cluster.Keyspace = keyspace
		cluster.Consistency = gocql.Quorum
		cluster.Timeout = 10 * time.Second
		cluster.ConnectTimeout = 10 * time.Second

		session, err = cluster.CreateSession()
		if err == nil {
			return session, nil
		}

		log.Printf("Waiting for Cassandra keyspace %s (attempt %d/%d): %v", keyspace, i+1, maxRetries, err)
		time.Sleep(backoff)

		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
		}
	}

	return nil, fmt.Errorf("failed to connect to keyspace %s after %d attempts: %v", keyspace, maxRetries, err)
}

func createKeyspaceAndTable(hosts []string, keyspace string) error {
	// Čekaj da Cassandra postane dostupna
	s, err := waitForCassandra(hosts)
	if err != nil {
		return err
	}
	defer s.Close()

	// Keyspace
	if err := s.Query(
		"CREATE KEYSPACE IF NOT EXISTS " + keyspace +
			" WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}",
	).Exec(); err != nil {
		return err
	}

	// Table u keyspace-u (koristi retry logiku)
	s2, err := connectToKeyspace(hosts, keyspace)
	if err != nil {
		return err
	}
	defer s2.Close()

	// ✅ user_id je partition key, id clustering key
	return s2.Query(`
		CREATE TABLE IF NOT EXISTS notifications (
			user_id TEXT,
			id UUID,
			message TEXT,
			type TEXT,
			read BOOLEAN,
			created_at TIMESTAMP,
			PRIMARY KEY (user_id, id)
		)
	`).Exec()
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, PUT, OPTIONS")

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
		api.GET("/notifications", middleware.AuthMiddleware(), handlers.GetNotifications)
		api.PUT("/notifications/:id/read", middleware.AuthMiddleware(), handlers.MarkAsRead)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "notifications-service"})
	})
}
