package middleware

import (
	"context"
	"net/http"
	"time"

	"example.com/users-service/utils"
	"github.com/gin-gonic/gin"
)

// RateLimitMiddleware limits requests based on IP address
// limit: number of requests allowed
// window: time window for the limit
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if redisClient == nil {
			// Fail open if redis is not connected, or log error
			// For security, maybe better to log and allow, or panic during init
			c.Next()
			return
		}

		ip := c.ClientIP()
		key := "rate_limit:" + ip

		ctx := context.Background()

		// Increment request count
		count, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			// Redis error, fail open but log
			utils.LogSecurityEvent("error", "rate_limit", ip, "Redis error")
			c.Next()
			return
		}

		// Set expiration on first request
		if count == 1 {
			redisClient.Expire(ctx, key, window)
		}

		if count > int64(limit) {
			utils.LogSecurityEvent("blocked", "rate_limit_exceeded", ip, "Too many requests")
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
