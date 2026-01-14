package utils

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/redis/go-redis/v9"
)

// GenerateOTP creates a 6-digit OTP code
func GenerateOTP() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// SaveOTPToRedis stores OTP with expiration
func SaveOTPToRedis(ctx context.Context, client *redis.Client, userID, otp string, expiration time.Duration) error {
	key := fmt.Sprintf("otp:%s", userID)
	return client.Set(ctx, key, otp, expiration).Err()
}

// VerifyOTP checks if OTP is valid
func VerifyOTP(ctx context.Context, client *redis.Client, userID, otp string) (bool, error) {
	key := fmt.Sprintf("otp:%s", userID)
	storedOTP, err := client.Get(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if storedOTP == otp {
		// Delete OTP after successful verification
		client.Del(ctx, key)
		return true, nil
	}

	return false, nil
}

// GenerateTempToken creates temporary token for OTP verification
func GenerateTempToken(userID string) (string, error) {
	token := fmt.Sprintf("temp_%s_%d", userID, time.Now().Unix())
	return token, nil
}

// SaveTempToken stores temp token in Redis
func SaveTempToken(ctx context.Context, client *redis.Client, token, userID string) error {
	key := fmt.Sprintf("temp_token:%s", token)
	return client.Set(ctx, key, userID, 10*time.Minute).Err()
}

// VerifyTempToken checks if temp token is valid and returns userID
func VerifyTempToken(ctx context.Context, client *redis.Client, token string) (string, error) {
	key := fmt.Sprintf("temp_token:%s", token)
	userID, err := client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return userID, nil
}
