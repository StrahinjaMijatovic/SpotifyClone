package handlers

import (
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

var session *gocql.Session

func InitHandlers(s *gocql.Session) {
	session = s
}

type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
}

func GetNotifications(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := c.Request.Context()

	var notifications []Notification
	iter := session.Query(`
		SELECT id, user_id, message, type, read, created_at
		FROM notifications
		WHERE user_id = ?
	`, userID).WithContext(ctx).Iter()

	var (
		id        gocql.UUID
		uid       string
		message   string
		ntype     string
		read      bool
		createdAt time.Time
	)

	for iter.Scan(&id, &uid, &message, &ntype, &read, &createdAt) {
		notifications = append(notifications, Notification{
			ID:        id.String(),
			UserID:    uid,
			Message:   message,
			Type:      ntype,
			Read:      read,
			CreatedAt: createdAt,
		})
	}

	if err := iter.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	// (Opcionalno) sort po created_at desc u aplikaciji
	sort.Slice(notifications, func(i, j int) bool {
		return notifications[i].CreatedAt.After(notifications[j].CreatedAt)
	})

	c.JSON(http.StatusOK, notifications)
}

func MarkAsRead(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")
	notificationID, err := gocql.ParseUUID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	ctx := c.Request.Context()

	// ✅ Pošto je PRIMARY KEY (user_id, id), WHERE mora imati user_id + id
	if err := session.Query(`
		UPDATE notifications
		SET read = true
		WHERE user_id = ? AND id = ?
	`, userID, notificationID).WithContext(ctx).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}
