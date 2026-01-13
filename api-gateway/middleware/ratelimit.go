package middleware

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func init() {
	redisURI := os.Getenv("REDIS_URI")
	if redisURI == "" {
		redisURI = "redis://localhost:6379"
	}

	opt, err := redis.ParseURL(redisURI)
	if err != nil {
		return
	}

	redisClient = redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		redisClient = nil // fallback
	}
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		rate := int64(100)
		window := time.Minute

		if redisClient != nil {
			ctx := c.Request.Context()
			key := "ratelimit:" + ip

			count, err := redisClient.Incr(ctx, key).Result()
			if err != nil {
				c.Next()
				return
			}

			if count == 1 {
				_ = redisClient.Expire(ctx, key, window).Err()
			}

			if count > rate {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error": "Rate limit exceeded. Please try again later.",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
