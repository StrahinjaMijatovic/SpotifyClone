package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var driver neo4j.DriverWithContext

func InitHandlers(d neo4j.DriverWithContext) {
	driver = d
}

func GetRecommendations(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := c.Request.Context()

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	// ✅ Ne preporučuj pesme koje je user već ocenio
	// ✅ Koristi ctx iz request-a
	result, err := session.Run(ctx,
		`MATCH (u:User {id: $userID})-[:RATED]->(:Song)-[:HAS_GENRE]->(g:Genre)
		 MATCH (similar:User)-[:RATED]->(similarSong:Song)-[:HAS_GENRE]->(g)
		 WHERE similar.id <> $userID
		   AND NOT (u)-[:RATED]->(similarSong)
		 RETURN similarSong.id as songId, COUNT(*) as score
		 ORDER BY score DESC
		 LIMIT 10`,
		map[string]any{"userID": userID},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recommendations"})
		return
	}

	recommendations := make([]map[string]any, 0, 10)
	for result.Next(ctx) {
		record := result.Record()
		songID, _ := record.Get("songId")
		score, _ := record.Get("score")

		recommendations = append(recommendations, map[string]any{
			"song_id": songID,
			"score":   score,
		})
	}

	if err := result.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process recommendations"})
		return
	}

	c.JSON(http.StatusOK, recommendations)
}
