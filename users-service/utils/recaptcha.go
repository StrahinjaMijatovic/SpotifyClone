package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

var recaptchaSecretKey = "6Le-GU0sAAAAAH5edSdsYVW9QTHKBjy7Z4SEKfn8"

func init() {
	if secret := os.Getenv("RECAPTCHA_SECRET_KEY"); secret != "" {
		recaptchaSecretKey = secret
	}
}

// RecaptchaResponse represents Google's reCAPTCHA verification response
type RecaptchaResponse struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes,omitempty"`
	Score       float64   `json:"score,omitempty"`       // For reCAPTCHA v3
	Action      string    `json:"action,omitempty"`      // For reCAPTCHA v3
}

// VerifyRecaptcha verifies the reCAPTCHA token with Google's API
func VerifyRecaptcha(token string, remoteIP string) (bool, error) {
	if token == "" {
		return false, fmt.Errorf("reCAPTCHA token is empty")
	}

	// Check if reCAPTCHA is disabled (for testing)
	if os.Getenv("RECAPTCHA_DISABLED") == "true" {
		LogSecurityEvent("info", "recaptcha_disabled", remoteIP, "reCAPTCHA verification skipped (disabled)")
		return true, nil
	}

	// Prepare the verification request
	verifyURL := "https://www.google.com/recaptcha/api/siteverify"

	resp, err := http.PostForm(verifyURL, url.Values{
		"secret":   {recaptchaSecretKey},
		"response": {token},
		"remoteip": {remoteIP},
	})
	if err != nil {
		LogSecurityEvent("error", "recaptcha_verify", remoteIP, fmt.Sprintf("Failed to verify: %v", err))
		return false, fmt.Errorf("failed to verify reCAPTCHA: %v", err)
	}
	defer resp.Body.Close()

	// Parse the response
	var result RecaptchaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		LogSecurityEvent("error", "recaptcha_parse", remoteIP, fmt.Sprintf("Failed to parse response: %v", err))
		return false, fmt.Errorf("failed to parse reCAPTCHA response: %v", err)
	}

	// Log the result
	if result.Success {
		LogSecurityEvent("success", "recaptcha_verify", remoteIP, fmt.Sprintf("reCAPTCHA verified for hostname: %s", result.Hostname))
	} else {
		LogSecurityEvent("failed", "recaptcha_verify", remoteIP, fmt.Sprintf("reCAPTCHA failed: %v", result.ErrorCodes))
	}

	return result.Success, nil
}

// IsRecaptchaEnabled returns whether reCAPTCHA verification is enabled
func IsRecaptchaEnabled() bool {
	return os.Getenv("RECAPTCHA_DISABLED") != "true"
}
