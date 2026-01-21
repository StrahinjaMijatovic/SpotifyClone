package proxy

import (
	"crypto/tls"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	usersServiceURL          = getEnv("USERS_SERVICE_URL", "https://localhost:8001")
	contentServiceURL        = getEnv("CONTENT_SERVICE_URL", "https://localhost:8002")
	ratingsServiceURL        = getEnv("RATINGS_SERVICE_URL", "https://localhost:8003")
	subscriptionsServiceURL  = getEnv("SUBSCRIPTIONS_SERVICE_URL", "https://localhost:8004")
	notificationsServiceURL  = getEnv("NOTIFICATIONS_SERVICE_URL", "https://localhost:8005")
	recommendationServiceURL = getEnv("RECOMMENDATION_SERVICE_URL", "https://localhost:8006")
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func proxyRequest(c *gin.Context, baseURL string) {
	// Skip TLS verification for self-signed certificates in development
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: 15 * time.Second, Transport: tr}

	// ✅ Ne diramo path — gateway je već /api/v1, i servisi su /api/v1
	url := baseURL + c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		url += "?" + c.Request.URL.RawQuery
	}

	req, err := http.NewRequestWithContext(c.Request.Context(), c.Request.Method, url, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Copy headers
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to connect to service"})
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}

	c.Writer.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(c.Writer, resp.Body)
}

func ProxyToUsersService(c *gin.Context)          { proxyRequest(c, usersServiceURL) }
func ProxyToContentService(c *gin.Context)        { proxyRequest(c, contentServiceURL) }
func ProxyToRatingsService(c *gin.Context)        { proxyRequest(c, ratingsServiceURL) }
func ProxyToSubscriptionsService(c *gin.Context)  { proxyRequest(c, subscriptionsServiceURL) }
func ProxyToNotificationsService(c *gin.Context)  { proxyRequest(c, notificationsServiceURL) }
func ProxyToRecommendationService(c *gin.Context) { proxyRequest(c, recommendationServiceURL) }

// DeleteSongCascade handles cascade deletion of a song across all services
// 1. Deletes the song from content-service
// 2. Deletes all ratings for the song from ratings-service
// 3. Deletes the song from recommendation graph
func DeleteSongCascade(c *gin.Context) {
	songID := c.Param("id")
	if songID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Song ID required"})
		return
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: 15 * time.Second, Transport: tr}

	// Copy authorization header for service calls
	authHeader := c.GetHeader("Authorization")

	var errors []string

	// 1. Delete song from content-service
	contentURL := contentServiceURL + "/api/v1/songs/" + songID
	contentReq, _ := http.NewRequestWithContext(c.Request.Context(), "DELETE", contentURL, nil)
	contentReq.Header.Set("Authorization", authHeader)

	contentResp, err := client.Do(contentReq)
	if err != nil {
		errors = append(errors, "Failed to connect to content service")
	} else {
		defer contentResp.Body.Close()
		if contentResp.StatusCode == http.StatusNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
			return
		}
		if contentResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(contentResp.Body)
			errors = append(errors, "Content service error: "+string(body))
		}
	}

	// 2. Delete all ratings for the song from ratings-service
	ratingsURL := ratingsServiceURL + "/api/v1/ratings/" + songID + "/all"
	ratingsReq, _ := http.NewRequestWithContext(c.Request.Context(), "DELETE", ratingsURL, nil)
	ratingsReq.Header.Set("Authorization", authHeader)

	ratingsResp, err := client.Do(ratingsReq)
	if err != nil {
		errors = append(errors, "Failed to connect to ratings service")
	} else {
		defer ratingsResp.Body.Close()
		if ratingsResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(ratingsResp.Body)
			errors = append(errors, "Ratings service error: "+string(body))
		}
	}

	// 3. Delete song from recommendation graph
	recsURL := recommendationServiceURL + "/api/v1/recommendations/songs/" + songID
	recsReq, _ := http.NewRequestWithContext(c.Request.Context(), "DELETE", recsURL, nil)
	recsReq.Header.Set("Authorization", authHeader)

	recsResp, err := client.Do(recsReq)
	if err != nil {
		errors = append(errors, "Failed to connect to recommendation service")
	} else {
		defer recsResp.Body.Close()
		if recsResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(recsResp.Body)
			errors = append(errors, "Recommendation service error: "+string(body))
		}
	}

	if len(errors) > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Song deleted with some errors",
			"errors":  errors,
			"song_id": songID,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Song and all related data deleted successfully",
		"song_id": songID,
	})
}
