package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitHandlers(client *redis.Client) {
	redisClient = client
}

type Subscription struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"` // "artist" or "genre"
	TargetID  string    `json:"target_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateSubscriptionRequest struct {
	Type     string `json:"type" binding:"required,oneof=artist genre"`
	TargetID string `json:"target_id" binding:"required"`
}

func scanKeys(ctx context.Context, pattern string) ([]string, error) {
	var cursor uint64
	var keys []string

	for {
		batch, next, err := redisClient.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, err
		}

		keys = append(keys, batch...)
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return keys, nil
}

func CreateSubscription(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	subscription := Subscription{
		ID:        userID + ":" + req.Type + ":" + req.TargetID,
		UserID:    userID,
		Type:      req.Type,
		TargetID:  req.TargetID,
		CreatedAt: time.Now(),
	}

	key := "subscription:" + subscription.ID

	data, err := json.Marshal(subscription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize subscription"})
		return
	}

	ctx := c.Request.Context()
	if err := redisClient.Set(ctx, key, data, 0).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

func GetSubscriptions(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := c.Request.Context()
	keys, err := scanKeys(ctx, "subscription:"+userID+":*")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscriptions"})
		return
	}

	var subscriptions []Subscription
	for _, key := range keys {
		data, err := redisClient.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var subscription Subscription
		if err := json.Unmarshal([]byte(data), &subscription); err == nil {
			subscriptions = append(subscriptions, subscription)
		}
	}

	c.JSON(http.StatusOK, subscriptions)
}

func DeleteSubscription(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing subscription id"})
		return
	}

	// 🔒 Ne dozvoli da briše tuđe pretplate
	if !strings.HasPrefix(id, userID+":") {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own subscriptions"})
		return
	}

	key := "subscription:" + id

	ctx := c.Request.Context()
	if err := redisClient.Del(ctx, key).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription deleted successfully"})
}
