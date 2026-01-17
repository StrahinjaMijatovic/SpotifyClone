package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"example.com/users-service/config"
	"example.com/users-service/models"
	"example.com/users-service/utils"
)

var (
	usersDB     *mongo.Database
	redisClient *redis.Client
)

func InitHandlers(db *mongo.Database, redis *redis.Client) {
	usersDB = db
	redisClient = redis
}

// Register handles user registration with validation
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogSecurityEvent("validation_failed", "register", c.ClientIP(), "Invalid request data")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Input validation
	req.Username = utils.SanitizeString(req.Username)
	req.Email = utils.SanitizeString(req.Email)
	req.FirstName = utils.SanitizeString(req.FirstName)
	req.LastName = utils.SanitizeString(req.LastName)

	// Validate username format
	if !utils.ValidateUsername(req.Username) {
		utils.LogSecurityEvent("validation_failed", "register", c.ClientIP(), "Invalid username format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be 3-50 alphanumeric characters or underscore"})
		return
	}

	// Validate email
	if !utils.ValidateEmail(req.Email) {
		utils.LogSecurityEvent("validation_failed", "register", c.ClientIP(), "Invalid email format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Validate names
	if !utils.ValidateName(req.FirstName) || !utils.ValidateName(req.LastName) {
		utils.LogSecurityEvent("validation_failed", "register", c.ClientIP(), "Invalid name format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Names must contain only letters and be 2-50 characters"})
		return
	}

	// Check for dangerous characters
	if utils.ContainsSpecialChars(req.Username) || utils.ContainsSpecialChars(req.FirstName) || utils.ContainsSpecialChars(req.LastName) {
		utils.LogSecurityEvent("validation_failed", "register", c.ClientIP(), "Special characters detected")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input contains invalid characters"})
		return
	}

	// Password strength validation
	if !utils.ValidatePasswordStrength(req.Password) {
		utils.LogSecurityEvent("validation_failed", "register", c.ClientIP(), "Weak password")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain uppercase, lowercase, number and special character"})
		return
	}

	// Check password blacklist
	if utils.IsPasswordBlacklisted(req.Password) {
		utils.LogSecurityEvent("validation_failed", "register", c.ClientIP(), "Blacklisted password")
		c.JSON(http.StatusBadRequest, gin.H{"error": "This password is too common. Please choose a more unique password"})
		return
	}

	ctx := c.Request.Context()

	// Check if username exists
	var existingUser models.User
	err := usersDB.Collection("users").FindOne(ctx, bson.M{"username": req.Username}).Decode(&existingUser)
	if err == nil {
		utils.LogSecurityEvent("validation_failed", "register", c.ClientIP(), "Username already exists")
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Check if email exists
	err = usersDB.Collection("users").FindOne(ctx, bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		utils.LogSecurityEvent("validation_failed", "register", c.ClientIP(), "Email already exists")
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Generate verification token
	verificationToken, err := utils.GenerateVerificationToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verification token"})
		return
	}

	role := models.RoleRegular
	if strings.HasPrefix(req.Email, "admin@") {
		role = models.RoleAdmin
	}

	user := models.User{
		ID:                        primitive.NewObjectID(),
		Username:                  req.Username,
		Email:                     req.Email,
		PasswordHash:              string(hashedPassword),
		FirstName:                 req.FirstName,
		LastName:                  req.LastName,
		Role:                      role,
		EmailVerified:             false,
		EmailVerificationToken:    verificationToken,
		EmailVerificationTokenExp: time.Now().Add(24 * time.Hour),
		PasswordChangedAt:         time.Now(),
		CreatedAt:                 time.Now(),
		UpdatedAt:                 time.Now(),
		FailedLoginAttempts:       0,
	}

	_, err = usersDB.Collection("users").InsertOne(ctx, user)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Send verification email
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:4200"
	}
	verificationLink := fmt.Sprintf("%s/verify-email?token=%s", frontendURL, verificationToken)
	emailBody := fmt.Sprintf("Hi %s,\n\nPlease verify your email by clicking this link:\n%s\n\nThis link expires in 24 hours.", req.FirstName, verificationLink)

	go utils.SendEmail(req.Email, "Verify your email", emailBody)

	utils.LogSecurityEvent("success", "register", c.ClientIP(), fmt.Sprintf("User %s registered", req.Username))

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful. Please check your email to verify your account.",
		"user_id": user.ID.Hex(),
	})
}

// VerifyEmail verifies user email with token
func VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token required"})
		return
	}

	ctx := c.Request.Context()

	var user models.User
	err := usersDB.Collection("users").FindOne(ctx, bson.M{
		"email_verification_token":     token,
		"email_verification_token_exp": bson.M{"$gt": time.Now()},
	}).Decode(&user)

	if err != nil {
		utils.LogSecurityEvent("failed", "email_verification", c.ClientIP(), "Invalid or expired token")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired verification token"})
		return
	}

	// Update user
	update := bson.M{
		"$set": bson.M{
			"email_verified":               true,
			"email_verification_token":     "",
			"email_verification_token_exp": time.Time{},
			"updated_at":                   time.Now(),
		},
	}

	_, err = usersDB.Collection("users").UpdateOne(ctx, bson.M{"_id": user.ID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email"})
		return
	}

	utils.LogSecurityEvent("success", "email_verification", c.ClientIP(), fmt.Sprintf("User %s verified email", user.Username))

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully. You can now login."})
}

// Login handles user login with OTP
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogSecurityEvent("validation_failed", "login", c.ClientIP(), "Invalid request data")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Sanitize input
	req.Username = utils.SanitizeString(req.Username)

	ctx := c.Request.Context()

	var user models.User
	err := usersDB.Collection("users").FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
	if err != nil {
		utils.LogSecurityEvent("failed", "login", c.ClientIP(), fmt.Sprintf("User %s not found", req.Username))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if account is locked
	if time.Now().Before(user.LockedUntil) {
		utils.LogSecurityEvent("failed", "login", c.ClientIP(), fmt.Sprintf("User %s account locked", req.Username))
		c.JSON(http.StatusForbidden, gin.H{"error": "Account is temporarily locked due to multiple failed login attempts"})
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		// Increment failed attempts
		failedAttempts := user.FailedLoginAttempts + 1
		update := bson.M{
			"$set": bson.M{
				"failed_login_attempts": failedAttempts,
				"last_failed_login":     time.Now(),
			},
		}

		// Lock account after 5 failed attempts
		if failedAttempts >= 5 {
			update["$set"].(bson.M)["locked_until"] = time.Now().Add(15 * time.Minute)
			utils.LogSecurityEvent("failed", "login", c.ClientIP(), fmt.Sprintf("User %s account locked after 5 failed attempts", req.Username))
		}

		usersDB.Collection("users").UpdateOne(ctx, bson.M{"_id": user.ID}, update)

		utils.LogSecurityEvent("failed", "login", c.ClientIP(), fmt.Sprintf("User %s invalid password", req.Username))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if email is verified
	if !user.EmailVerified {
		utils.LogSecurityEvent("failed", "login", c.ClientIP(), fmt.Sprintf("User %s email not verified", req.Username))
		c.JSON(http.StatusForbidden, gin.H{"error": "Please verify your email before logging in"})
		return
	}

	// Check password age (configurable via PASSWORD_MAX_AGE_DAYS or PASSWORD_MAX_AGE_MINUTES env vars)
	if config.IsPasswordExpired(user.PasswordChangedAt) {
		passwordAge := config.GetPasswordAge(user.PasswordChangedAt)
		maxAge := config.GetPasswordMaxAgeString()
		utils.LogSecurityEvent("failed", "login", c.ClientIP(),
			fmt.Sprintf("User %s password expired (age: %v, max: %s)", req.Username, passwordAge.Round(time.Minute), maxAge))
		c.JSON(http.StatusForbidden, gin.H{
			"error":       "Lozinka je istekla. Molimo resetujte vašu lozinku.",
			"error_code":  "PASSWORD_EXPIRED",
			"password_age": passwordAge.String(),
			"max_age":     maxAge,
			"message":     fmt.Sprintf("Vaša lozinka je starija od dozvoljenog perioda od %s. Morate resetovati lozinku da biste nastavili.", maxAge),
		})
		return
	}

	// Send password expiry warning (at 80% of max age)
	warningThreshold := config.PasswordMaxAgeDuration * 80 / 100
	passwordAge := config.GetPasswordAge(user.PasswordChangedAt)
	if passwordAge > warningThreshold {
		daysLeft := config.GetDaysUntilExpiry(user.PasswordChangedAt)
		var timeUnit string
		if config.PasswordMaxAgeDays == 0 {
			timeUnit = "minuta"
		} else {
			timeUnit = "dana"
		}
		emailBody := fmt.Sprintf("Zdravo %s,\n\nVaša lozinka će isteći za %d %s. Molimo promenite je što pre.", user.FirstName, daysLeft, timeUnit)
		go utils.SendEmail(user.Email, "Upozorenje o isteku lozinke", emailBody)
	}

	// Generate OTP
	otp, err := utils.GenerateOTP()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
		return
	}

	// Save OTP to Redis
	err = utils.SaveOTPToRedis(ctx, redisClient, user.ID.Hex(), otp, 5*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save OTP"})
		return
	}

	// Generate temp token
	tempToken, err := utils.GenerateTempToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate temp token"})
		return
	}

	// Save temp token
	err = utils.SaveTempToken(ctx, redisClient, tempToken, user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save temp token"})
		return
	}

	// Send OTP email
	emailBody := fmt.Sprintf("Hi %s,\n\nYour OTP code is: %s\n\nThis code expires in 5 minutes.", user.FirstName, otp)
	go utils.SendEmail(user.Email, "Your OTP Code", emailBody)

	// Reset failed login attempts
	usersDB.Collection("users").UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{
		"$set": bson.M{
			"failed_login_attempts": 0,
		},
	})

	utils.LogSecurityEvent("success", "login_otp_sent", c.ClientIP(), fmt.Sprintf("User %s OTP sent", req.Username))

	c.JSON(http.StatusOK, gin.H{
		"message":    "OTP sent to your email",
		"temp_token": tempToken,
	})
}

// VerifyOTP verifies OTP and returns JWT token
func VerifyOTP(c *gin.Context) {
	var req models.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogSecurityEvent("validation_failed", "verify_otp", c.ClientIP(), "Invalid request data")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	ctx := c.Request.Context()

	// Verify temp token
	userID, err := utils.VerifyTempToken(ctx, redisClient, req.TempToken)
	if err != nil {
		utils.LogSecurityEvent("failed", "verify_otp", c.ClientIP(), "Invalid temp token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired temp token"})
		return
	}

	// Verify OTP
	valid, err := utils.VerifyOTP(ctx, redisClient, userID, req.OTPCode)
	if err != nil || !valid {
		utils.LogSecurityEvent("failed", "verify_otp", c.ClientIP(), "Invalid OTP")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP code"})
		return
	}

	// Get user
	objID, _ := primitive.ObjectIDFromHex(userID)
	var user models.User
	err = usersDB.Collection("users").FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID.Hex(), user.Username, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	utils.LogSecurityEvent("success", "login", c.ClientIP(), fmt.Sprintf("User %s logged in", user.Username))

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

// RequestMagicLink sends magic link for account recovery
func RequestMagicLink(c *gin.Context) {
	var req models.MagicLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	req.Email = utils.SanitizeString(req.Email)

	ctx := c.Request.Context()

	var user models.User
	err := usersDB.Collection("users").FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		// Don't reveal if email exists
		c.JSON(http.StatusOK, gin.H{"message": "If email exists, magic link has been sent"})
		return
	}

	// Generate magic link token
	token, err := utils.GenerateVerificationToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Save token to user
	update := bson.M{
		"$set": bson.M{
			"verification_token":     token,
			"verification_token_exp": time.Now().Add(15 * time.Minute),
			"updated_at":             time.Now(),
		},
	}
	usersDB.Collection("users").UpdateOne(ctx, bson.M{"_id": user.ID}, update)

	// Send magic link email
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:4200"
	}
	magicLink := fmt.Sprintf("%s/magic-login?token=%s", frontendURL, token)
	emailBody := fmt.Sprintf("Hi %s,\n\nClick this link to login:\n%s\n\nThis link expires in 15 minutes.", user.FirstName, magicLink)

	go utils.SendEmail(user.Email, "Magic Login Link", emailBody)

	utils.LogSecurityEvent("success", "magic_link_sent", c.ClientIP(), fmt.Sprintf("Magic link sent to %s", req.Email))

	c.JSON(http.StatusOK, gin.H{"message": "If email exists, magic link has been sent"})
}

func ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogSecurityEvent("validation_failed", "reset_password", c.ClientIP(), "Invalid request data")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	req.Email = utils.SanitizeString(req.Email)

	if !utils.ValidateEmail(req.Email) {
		utils.LogSecurityEvent("validation_failed", "reset_password", c.ClientIP(), "Invalid email format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	ctx := c.Request.Context()

	var user models.User
	err := usersDB.Collection("users").FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		// Don't reveal if email exists or not
		c.JSON(http.StatusOK, gin.H{"message": "If email exists, reset link has been sent"})
		return
	}

	// Generate reset token
	resetToken, err := utils.GenerateVerificationToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset token"})
		return
	}

	// Save token to user
	update := bson.M{
		"$set": bson.M{
			"verification_token":     resetToken,
			"verification_token_exp": time.Now().Add(1 * time.Hour),
			"updated_at":             time.Now(),
		},
	}

	_, err = usersDB.Collection("users").UpdateOne(ctx, bson.M{"_id": user.ID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save reset token"})
		return
	}

	// Send reset email
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:4200"
	}
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, resetToken)
	emailBody := fmt.Sprintf("Hi %s,\n\nClick this link to reset your password:\n%s\n\nThis link expires in 1 hour.\n\nIf you didn't request this, please ignore this email.",
		user.FirstName, resetLink)

	go utils.SendEmail(user.Email, "Password Reset Request", emailBody)

	utils.LogSecurityEvent("success", "reset_password_sent", c.ClientIP(), fmt.Sprintf("Reset link sent to %s", req.Email))

	c.JSON(http.StatusOK, gin.H{"message": "If email exists, reset link has been sent"})
}

// ResetPasswordConfirm resets password with token
func ResetPasswordConfirm(c *gin.Context) {
	var req models.ResetPasswordConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogSecurityEvent("validation_failed", "reset_password_confirm", c.ClientIP(), "Invalid request data")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Validate new password
	if !utils.ValidatePasswordStrength(req.NewPassword) {
		utils.LogSecurityEvent("validation_failed", "reset_password_confirm", c.ClientIP(), "Weak password")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain uppercase, lowercase, number and special character"})
		return
	}

	// Check password blacklist
	if utils.IsPasswordBlacklisted(req.NewPassword) {
		utils.LogSecurityEvent("validation_failed", "reset_password_confirm", c.ClientIP(), "Blacklisted password")
		c.JSON(http.StatusBadRequest, gin.H{"error": "This password is too common. Please choose a more unique password"})
		return
	}

	ctx := c.Request.Context()

	var user models.User
	err := usersDB.Collection("users").FindOne(ctx, bson.M{
		"verification_token":     req.Token,
		"verification_token_exp": bson.M{"$gt": time.Now()},
	}).Decode(&user)

	if err != nil {
		utils.LogSecurityEvent("failed", "reset_password_confirm", c.ClientIP(), "Invalid or expired token")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired reset token"})
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Update password and clear token
	update := bson.M{
		"$set": bson.M{
			"password_hash":          string(hashedPassword),
			"password_changed_at":    time.Now(),
			"verification_token":     "",
			"verification_token_exp": time.Time{},
			"updated_at":             time.Now(),
			"failed_login_attempts":  0,
			"locked_until":           time.Time{},
		},
	}

	_, err = usersDB.Collection("users").UpdateOne(ctx, bson.M{"_id": user.ID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}

	// Send confirmation email
	emailBody := fmt.Sprintf("Hi %s,\n\nYour password has been successfully reset.\n\nIf you didn't make this change, please contact support immediately.",
		user.FirstName)
	go utils.SendEmail(user.Email, "Password Reset Successful", emailBody)

	utils.LogSecurityEvent("success", "reset_password_confirm", c.ClientIP(), fmt.Sprintf("User %s reset password", user.Username))

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully. You can now login with your new password."})
}

// EnsureUserIndexes creates unique indexes for username and email
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
		{
			Keys: bson.D{{Key: "verification_token", Value: 1}},
			Options: options.Index().
				SetName("verification_token_idx").
				SetSparse(true),
		},
	}

	_, err := users.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Fatalf("Failed to create MongoDB indexes: %v", err)
	}

	log.Println("MongoDB indexes for users ensured")
}

// MagicLogin logs in user with magic link token
func MagicLogin(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token required"})
		return
	}

	ctx := c.Request.Context()

	var user models.User
	err := usersDB.Collection("users").FindOne(ctx, bson.M{
		"verification_token":     token,
		"verification_token_exp": bson.M{"$gt": time.Now()},
	}).Decode(&user)

	if err != nil {
		utils.LogSecurityEvent("failed", "magic_login", c.ClientIP(), "Invalid or expired token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired magic link"})
		return
	}

	// Check if email is verified
	if !user.EmailVerified {
		utils.LogSecurityEvent("failed", "magic_login", c.ClientIP(), fmt.Sprintf("User %s email not verified", user.Username))
		c.JSON(http.StatusForbidden, gin.H{"error": "Please verify your email first"})
		return
	}

	// Clear verification token
	update := bson.M{
		"$set": bson.M{
			"verification_token":     "",
			"verification_token_exp": time.Time{},
			"failed_login_attempts":  0,
			"locked_until":           time.Time{},
			"updated_at":             time.Now(),
		},
	}

	_, err = usersDB.Collection("users").UpdateOne(ctx, bson.M{"_id": user.ID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process login"})
		return
	}

	// Generate JWT token
	jwtToken, err := utils.GenerateJWT(user.ID.Hex(), user.Username, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	utils.LogSecurityEvent("success", "magic_login", c.ClientIP(), fmt.Sprintf("User %s logged in via magic link", user.Username))

	c.JSON(http.StatusOK, gin.H{
		"token": jwtToken,
		"user": gin.H{
			"id":         user.ID.Hex(),
			"username":   user.Username,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"role":       user.Role,
		},
		"message": "Login successful",
	})
}

// ResetPasswordConfirmValidate validates reset token from GET link
func ResetPasswordConfirmValidate(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token required"})
		return
	}

	ctx := c.Request.Context()

	count, err := usersDB.Collection("users").CountDocuments(ctx, bson.M{
		"password_reset_token":     token,
		"password_reset_token_exp": bson.M{"$gt": time.Now()},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if count == 0 {
		utils.LogSecurityEvent("failed", "reset_password_validate", c.ClientIP(), "Invalid/expired token")
		c.JSON(http.StatusBadRequest, gin.H{"valid": false, "error": "Invalid or expired token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true, "message": "Token is valid"})
}

// ChangePassword allows authenticated users to change password,
// requires last password change to be >= 24h ago + strength validation
func ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogSecurityEvent("validation_failed", "change_password", c.ClientIP(), "Invalid request data")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Strength validation
	if !utils.ValidatePasswordStrength(req.NewPassword) {
		utils.LogSecurityEvent("validation_failed", "change_password", c.ClientIP(), "Weak password")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain uppercase, lowercase, number and special character"})
		return
	}

	// Check password blacklist
	if utils.IsPasswordBlacklisted(req.NewPassword) {
		utils.LogSecurityEvent("validation_failed", "change_password", c.ClientIP(), "Blacklisted password")
		c.JSON(http.StatusBadRequest, gin.H{"error": "This password is too common. Please choose a more unique password"})
		return
	}

	userIDAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDHex := userIDAny.(string)
	objID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	ctx := c.Request.Context()

	var user models.User
	if err := usersDB.Collection("users").FindOne(ctx, bson.M{"_id": objID}).Decode(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// Must be at least 1 day since last change
	if time.Since(user.PasswordChangedAt) < 24*time.Hour {
		utils.LogSecurityEvent("failed", "change_password", c.ClientIP(), "Password change too soon (<24h)")
		c.JSON(http.StatusForbidden, gin.H{"error": "Password can be changed only if it is at least 1 day old"})
		return
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		utils.LogSecurityEvent("failed", "change_password", c.ClientIP(), "Invalid current password")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid current password"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"password_hash":       string(hashedPassword),
			"password_changed_at": time.Now(),
			"updated_at":          time.Now(),
		},
	}

	if _, err := usersDB.Collection("users").UpdateOne(ctx, bson.M{"_id": user.ID}, update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change password"})
		return
	}

	utils.LogSecurityEvent("success", "change_password", c.ClientIP(), "Password changed")

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// Logout revokes current JWT by blacklisting its JTI in Redis until expiration
func Logout(c *gin.Context) {
	claimsAny, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims := claimsAny.(*utils.Claims)
	if claims.ID == "" || claims.ExpiresAt == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token claims"})
		return
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Already expired"})
		return
	}

	ctx := c.Request.Context()
	key := "bl:" + claims.ID

	if err := redisClient.Set(ctx, key, "1", ttl).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	utils.LogSecurityEvent("success", "logout", c.ClientIP(), "Token revoked (blacklisted)")
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}
