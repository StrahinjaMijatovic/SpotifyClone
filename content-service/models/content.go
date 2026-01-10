package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Genre struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

type Artist struct {
	ID        primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name      string               `json:"name" bson:"name"`
	Biography string               `json:"biography" bson:"biography"`
	Genres    []primitive.ObjectID `json:"genres" bson:"genres"`
	CreatedAt time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time            `json:"updated_at" bson:"updated_at"`
}

type Album struct {
	ID        primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name      string               `json:"name" bson:"name"`
	Date      time.Time            `json:"date" bson:"date"`
	Genre     primitive.ObjectID   `json:"genre" bson:"genre"`
	Artists   []primitive.ObjectID `json:"artists" bson:"artists"`
	CreatedAt time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time            `json:"updated_at" bson:"updated_at"`
}

type Song struct {
	ID        primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name      string               `json:"name" bson:"name"`
	Duration  int                  `json:"duration" bson:"duration"` // in seconds
	Genre     primitive.ObjectID   `json:"genre" bson:"genre"`
	Album     primitive.ObjectID   `json:"album" bson:"album"`
	Artists   []primitive.ObjectID `json:"artists" bson:"artists"`
	CreatedAt time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time            `json:"updated_at" bson:"updated_at"`
}

type CreateArtistRequest struct {
	Name      string   `json:"name" binding:"required,min=1,max=100"`
	Biography string   `json:"biography" binding:"required,min=10"`
	Genres    []string `json:"genres" binding:"required,min=1"`
}

type UpdateArtistRequest struct {
	Name      string   `json:"name" binding:"omitempty,min=1,max=100"`
	Biography string   `json:"biography" binding:"omitempty,min=10"`
	Genres    []string `json:"genres" binding:"omitempty,min=1"`
}

type CreateAlbumRequest struct {
	Name    string    `json:"name" binding:"required,min=1,max=100"`
	Date    time.Time `json:"date" binding:"required"`
	Genre   string    `json:"genre" binding:"required"`
	Artists []string  `json:"artists" binding:"required,min=1"`
}

type CreateSongRequest struct {
	Name     string   `json:"name" binding:"required,min=1,max=100"`
	Duration int      `json:"duration" binding:"required,min=1"`
	Genre    string   `json:"genre" binding:"required"`
	Album    string   `json:"album" binding:"required"`
	Artists  []string `json:"artists" binding:"required,min=1"`
}
