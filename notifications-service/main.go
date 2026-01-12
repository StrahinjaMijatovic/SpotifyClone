package main

import (
	"context"
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

	// 2) Onda konekcija na keyspace
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 10 * time.Second
	cluster.ConnectTimeout = 10 * time.Second

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Failed to connect to Cassandra:", err)
	}
	cassandraSession = session
	defer cassandraSession.Close()

	log.Println("Connected to Cassandra")

	router := gin.Default()
	router.Use(corsMiddleware())

	handlers.InitHandlers(cassandraSession)
	setupRoutes(router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

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

func createKeyspaceAndTable(hosts []string, keyspace string) error {
	// Session bez keyspace-a (system)
	cluster := gocql.NewCluster(hosts...)
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 10 * time.Second
	cluster.ConnectTimeout = 10 * time.Second

	s, err := cluster.CreateSession()
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

	// Table u keyspace-u
	cluster2 := gocql.NewCluster(hosts...)
	cluster2.Keyspace = keyspace
	cluster2.Consistency = gocql.Quorum
	cluster2.Timeout = 10 * time.Second
	cluster2.ConnectTimeout = 10 * time.Second

	s2, err := cluster2.CreateSession()
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
