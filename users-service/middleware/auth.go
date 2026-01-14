package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"example.com/users-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitAuthMiddleware(redis *redis.Client) {
	redisClient = redis
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.LogSecurityEvent("failed", "auth", c.ClientIP(), "Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.LogSecurityEvent("failed", "auth", c.ClientIP(), "Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			utils.LogSecurityEvent("failed", "auth", c.ClientIP(), "Invalid/expired JWT")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Check blacklist (logout)
		if redisClient != nil && claims.ID != "" && claims.ExpiresAt != nil {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
			defer cancel()

			key := "bl:" + claims.ID
			exists, err := redisClient.Exists(ctx, key).Result()
			if err == nil && exists > 0 {
				utils.LogSecurityEvent("failed", "auth", c.ClientIP(), "Blacklisted token used")
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token revoked"})
				c.Abort()
				return
			}
		}

		// Store claims for handlers
		c.Set("claims", claims)
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			utils.LogSecurityEvent("failed", "authz", c.ClientIP(), "Admin access denied")
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
