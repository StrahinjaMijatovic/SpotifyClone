package proxy

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	usersServiceURL          = getEnv("USERS_SERVICE_URL", "http://localhost:8001")
	contentServiceURL        = getEnv("CONTENT_SERVICE_URL", "http://localhost:8002")
	ratingsServiceURL        = getEnv("RATINGS_SERVICE_URL", "http://localhost:8003")
	subscriptionsServiceURL  = getEnv("SUBSCRIPTIONS_SERVICE_URL", "http://localhost:8004")
	notificationsServiceURL  = getEnv("NOTIFICATIONS_SERVICE_URL", "http://localhost:8005")
	recommendationServiceURL = getEnv("RECOMMENDATION_SERVICE_URL", "http://localhost:8006")
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func proxyRequest(c *gin.Context, baseURL string) {
	client := &http.Client{Timeout: 15 * time.Second}

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
