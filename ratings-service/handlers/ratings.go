package handlers

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitHandlers(client *redis.Client) {
	redisClient = client
}

type Rating struct {
	UserID  string    `json:"user_id"`
	SongID  string    `json:"song_id"`
	Rating  int       `json:"rating"`
	Created time.Time `json:"created"`
}

type CreateRatingRequest struct {
	SongID string `json:"song_id" binding:"required"`
	Rating int    `json:"rating" binding:"required,min=1,max=5"`
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

func CreateRating(c *gin.Context) {
	userID := c.GetString("user_id")
	var req CreateRatingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid data"})
		return
	}

	rating := Rating{
		UserID:  userID,
		SongID:  req.SongID,
		Rating:  req.Rating,
		Created: time.Now(),
	}

	key := "rating:" + req.SongID + ":" + userID

	data, err := json.Marshal(rating)
	if err != nil {
		c.JSON(500, gin.H{"error": "JSON error"})
		return
	}

	ctx := c.Request.Context()
	if err := redisClient.Set(ctx, key, data, 0).Err(); err != nil {
		c.JSON(500, gin.H{"error": "Redis error"})
		return
	}

	c.JSON(201, rating)
}

func GetRatings(c *gin.Context) {
	userID := c.GetString("user_id")
	ctx := c.Request.Context()

	keys, err := scanKeys(ctx, "rating:*:"+userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Redis error"})
		return
	}

	var ratings []Rating
	for _, key := range keys {
		val, err := redisClient.Get(ctx, key).Result()
		if err != nil {
			continue
		}
		var r Rating
		if err := json.Unmarshal([]byte(val), &r); err == nil {
			ratings = append(ratings, r)
		}
	}

	c.JSON(200, ratings)
}

func GetSongRatings(c *gin.Context) {
	songID := strings.TrimSpace(c.Param("songId"))
	if songID == "" {
		c.JSON(400, gin.H{"error": "songId required"})
		return
	}

	ctx := c.Request.Context()
	keys, err := scanKeys(ctx, "rating:"+songID+":*")
	if err != nil {
		c.JSON(500, gin.H{"error": "Redis error"})
		return
	}

	var ratings []Rating
	total := 0

	for _, key := range keys {
		val, err := redisClient.Get(ctx, key).Result()
		if err != nil {
			continue
		}
		var r Rating
		if err := json.Unmarshal([]byte(val), &r); err == nil {
			ratings = append(ratings, r)
			total += r.Rating
		}
	}

	avg := 0.0
	if len(ratings) > 0 {
		avg = float64(total) / float64(len(ratings))
	}

	c.JSON(200, gin.H{
		"ratings": ratings,
		"average": avg,
		"count":   len(ratings),
	})
}

func DeleteRating(c *gin.Context) {
	userID := c.GetString("user_id")
	songID := strings.TrimSpace(c.Param("songId"))

	if songID == "" {
		c.JSON(400, gin.H{"error": "songId required"})
		return
	}

	key := "rating:" + songID + ":" + userID
	ctx := c.Request.Context()

	result, err := redisClient.Del(ctx, key).Result()
	if err != nil {
		c.JSON(500, gin.H{"error": "Redis error"})
		return
	}

	if result == 0 {
		c.JSON(404, gin.H{"error": "Rating not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Rating deleted"})
}
