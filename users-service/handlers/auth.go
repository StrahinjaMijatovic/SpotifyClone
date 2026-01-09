package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"example.com/users-service/models"
	"example.com/users-service/utils"
)

var usersDB *mongo.Database

func InitHandlers(db *mongo.Database) {
	usersDB = db
}

func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if !utils.ValidatePasswordStrength(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password does not meet strength requirements"})
		return
	}

	ctx := c.Request.Context()

	var existingUser models.User

	// Check username
	err := usersDB.Collection("users").
		FindOne(ctx, bson.M{"username": req.Username}).
		Decode(&existingUser)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	} else if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check email
	err = usersDB.Collection("users").
		FindOne(ctx, bson.M{"email": req.Email}).
		Decode(&existingUser)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	} else if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		ID:                primitive.NewObjectID(),
		Username:          req.Username,
		Email:             req.Email,
		PasswordHash:      string(hashedPassword),
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		Role:              models.RoleRegular,
		EmailVerified:     false,
		PasswordChangedAt: time.Now(),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	_, err = usersDB.Collection("users").
		InsertOne(ctx, user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully. Please verify your email.",
		"user_id": user.ID.Hex(),
	})
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	ctx := c.Request.Context()

	var user models.User
	err := usersDB.Collection("users").
		FindOne(ctx, bson.M{"username": req.Username}).
		Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if time.Since(user.PasswordChangedAt) > 60*24*time.Hour {
		c.JSON(http.StatusForbidden, gin.H{"error": "Password expired. Please reset your password."})
		return
	}

	token, err := utils.GenerateJWT(user.ID.Hex(), user.Username, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":         user.ID.Hex(),
			"username":   user.Username,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"role":       user.Role,
		},
	})
}

func ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	ctx := c.Request.Context()

	var user models.User
	objID, _ := primitive.ObjectIDFromHex(userID.(string))

	err := usersDB.Collection("users").
		FindOne(ctx, bson.M{"_id": objID}).
		Decode(&user)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid current password"})
		return
	}

	if time.Since(user.PasswordChangedAt) < 24*time.Hour {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password can only be changed after 1 day"})
		return
	}

	if !utils.ValidatePasswordStrength(req.NewPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New password does not meet strength requirements"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"password_hash":       string(hashedPassword),
			"password_changed_at": time.Now(),
			"updated_at":          time.Now(),
		},
	}

	_, err = usersDB.Collection("users").
		UpdateOne(ctx, bson.M{"_id": objID}, update)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	ctx := c.Request.Context()

	var user models.User
	_ = usersDB.Collection("users").
		FindOne(ctx, bson.M{"email": req.Email}).
		Decode(&user)

	c.JSON(http.StatusOK, gin.H{"message": "If email exists, reset link has been sent"})
}

func ResetPasswordConfirm(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// Čak i ako se dve registracije pošalju u istoj milisekundi, Mongo će dozvoliti samo jednu, a druga će dobiti grešku.
func EnsureUserIndexes(db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	users := db.Collection("users")

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "username", Value: 1}},
			Options: options.Index().
				SetUnique(true).
				SetName("unique_username"),
		},
		{
			Keys: bson.D{{Key: "email", Value: 1}},
			Options: options.Index().
				SetUnique(true).
				SetName("unique_email"),
		},
	}

	_, err := users.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Fatalf("Failed to create MongoDB indexes: %v", err)
	}

	log.Println("MongoDB indexes for users ensured")
}
