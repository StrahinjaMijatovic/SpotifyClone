package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"example.com/content-service/models"
)

var contentDB *mongo.Database

func InitHandlers(db *mongo.Database) {
	contentDB = db
}

// Genre handlers
func GetGenres(c *gin.Context) {
	ctx := c.Request.Context()

	cursor, err := contentDB.Collection("genres").Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch genres"})
		return
	}
	defer cursor.Close(ctx)

	var genres []models.Genre
	if err := cursor.All(ctx, &genres); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode genres"})
		return
	}

	c.JSON(http.StatusOK, genres)
}

// Artist handlers
func CreateArtist(c *gin.Context) {
	var req models.CreateArtistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	ctx := c.Request.Context()

	// Convert genre IDs
	genreIDs := make([]primitive.ObjectID, len(req.Genres))
	for i, genreID := range req.Genres {
		objID, err := primitive.ObjectIDFromHex(genreID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
			return
		}
		genreIDs[i] = objID
	}

	artist := models.Artist{
		ID:        primitive.NewObjectID(),
		Name:      req.Name,
		Biography: req.Biography,
		Genres:    genreIDs,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := contentDB.Collection("artists").InsertOne(ctx, artist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create artist"})
		return
	}

	c.JSON(http.StatusCreated, artist)
}

func GetArtists(c *gin.Context) {
	ctx := c.Request.Context()

	filter := bson.M{}

	// Filter by genre_id if provided
	if genreID := c.Query("genre_id"); genreID != "" {
		objID, err := primitive.ObjectIDFromHex(genreID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
			return
		}
		filter["genres"] = objID
	}

	cursor, err := contentDB.Collection("artists").Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch artists"})
		return
	}
	defer cursor.Close(ctx)

	var artists []models.Artist
	if err := cursor.All(ctx, &artists); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode artists"})
		return
	}

	c.JSON(http.StatusOK, artists)
}

func GetArtist(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artist ID"})
		return
	}

	var artist models.Artist
	err = contentDB.Collection("artists").FindOne(ctx, bson.M{"_id": objID}).Decode(&artist)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Artist not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, artist)
}

func UpdateArtist(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artist ID"})
		return
	}

	var req models.UpdateArtistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	update := bson.M{"$set": bson.M{"updated_at": time.Now()}}

	if req.Name != "" {
		update["$set"].(bson.M)["name"] = req.Name
	}
	if req.Biography != "" {
		update["$set"].(bson.M)["biography"] = req.Biography
	}
	if len(req.Genres) > 0 {
		genreIDs := make([]primitive.ObjectID, len(req.Genres))
		for i, genreID := range req.Genres {
			gID, err := primitive.ObjectIDFromHex(genreID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
				return
			}
			genreIDs[i] = gID
		}
		update["$set"].(bson.M)["genres"] = genreIDs
	}

	result, err := contentDB.Collection("artists").UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update artist"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artist not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Artist updated successfully"})
}

// Album handlers
func CreateAlbum(c *gin.Context) {
	var req models.CreateAlbumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	ctx := c.Request.Context()

	genreID, err := primitive.ObjectIDFromHex(req.Genre)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
		return
	}

	artistIDs := make([]primitive.ObjectID, len(req.Artists))
	for i, artistID := range req.Artists {
		objID, err := primitive.ObjectIDFromHex(artistID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artist ID"})
			return
		}
		artistIDs[i] = objID
	}

	album := models.Album{
		ID:        primitive.NewObjectID(),
		Name:      req.Name,
		Date:      req.Date,
		Genre:     genreID,
		Artists:   artistIDs,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = contentDB.Collection("albums").InsertOne(ctx, album)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create album"})
		return
	}

	c.JSON(http.StatusCreated, album)
}

func GetAlbums(c *gin.Context) {
	ctx := c.Request.Context()

	filter := bson.M{}

	// Filter by artist_id if provided
	if artistID := c.Query("artist_id"); artistID != "" {
		objID, err := primitive.ObjectIDFromHex(artistID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artist ID"})
			return
		}
		filter["artists"] = objID
	}

	cursor, err := contentDB.Collection("albums").Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch albums"})
		return
	}
	defer cursor.Close(ctx)

	var albums []models.Album
	if err := cursor.All(ctx, &albums); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode albums"})
		return
	}

	c.JSON(http.StatusOK, albums)
}

func GetAlbum(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	var album models.Album
	err = contentDB.Collection("albums").FindOne(ctx, bson.M{"_id": objID}).Decode(&album)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, album)
}

// Song handlers
func CreateSong(c *gin.Context) {
	var req models.CreateSongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	ctx := c.Request.Context()

	albumID, err := primitive.ObjectIDFromHex(req.Album)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	// Check if album exists
	var album models.Album
	err = contentDB.Collection("albums").FindOne(ctx, bson.M{"_id": albumID}).Decode(&album)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Album does not exist. Create album first."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	genreID, err := primitive.ObjectIDFromHex(req.Genre)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
		return
	}

	artistIDs := make([]primitive.ObjectID, len(req.Artists))
	for i, artistID := range req.Artists {
		objID, err := primitive.ObjectIDFromHex(artistID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artist ID"})
			return
		}
		artistIDs[i] = objID
	}

	song := models.Song{
		ID:        primitive.NewObjectID(),
		Name:      req.Name,
		Duration:  req.Duration,
		Genre:     genreID,
		Album:     albumID,
		Artists:   artistIDs,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = contentDB.Collection("songs").InsertOne(ctx, song)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create song"})
		return
	}

	c.JSON(http.StatusCreated, song)
}

func GetSongs(c *gin.Context) {
	ctx := c.Request.Context()

	filter := bson.M{}

	// Filter by album_id if provided
	if albumID := c.Query("album_id"); albumID != "" {
		objID, err := primitive.ObjectIDFromHex(albumID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
			return
		}
		filter["album"] = objID
	}

	cursor, err := contentDB.Collection("songs").Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch songs"})
		return
	}
	defer cursor.Close(ctx)

	var songs []models.Song
	if err := cursor.All(ctx, &songs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode songs"})
		return
	}

	c.JSON(http.StatusOK, songs)
}

func DeleteSong(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	// Check if song exists first
	var song models.Song
	err = contentDB.Collection("songs").FindOne(ctx, bson.M{"_id": objID}).Decode(&song)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Delete the song
	result, err := contentDB.Collection("songs").DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully", "song_id": id})
}

func SearchContent(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query required"})
		return
	}

	ctx := c.Request.Context()

	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	// Search artists
	var artists []models.Artist
	artistCursor, err := contentDB.Collection("artists").Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search artists"})
		return
	}
	defer artistCursor.Close(ctx)

	if err := artistCursor.All(ctx, &artists); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode artists"})
		return
	}

	// Search albums
	var albums []models.Album
	albumCursor, err := contentDB.Collection("albums").Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search albums"})
		return
	}
	defer albumCursor.Close(ctx)

	if err := albumCursor.All(ctx, &albums); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode albums"})
		return
	}

	// Search songs
	var songs []models.Song
	songCursor, err := contentDB.Collection("songs").Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search songs"})
		return
	}
	defer songCursor.Close(ctx)

	if err := songCursor.All(ctx, &songs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode songs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"artists": artists,
		"albums":  albums,
		"songs":   songs,
	})
}
