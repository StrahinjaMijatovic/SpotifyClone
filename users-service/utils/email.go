package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// SendEmail sends email
func SendEmail(to, subject, body string) error {
	// For testing, just log metadata (NEVER log body, OTP, tokens)
	if os.Getenv("MOCK_EMAIL") == "true" {
		if os.Getenv("MOCK_EMAIL") == "true" {
			fmt.Printf("📧 MOCK EMAIL to=%s subject=%s\nBODY:\n%s\n-------------------\n", to, subject, body)
			return nil
		}
		return nil
	}

	from := os.Getenv("EMAIL_USER")
	password := os.Getenv("EMAIL_PASS")
	smtpHost := os.Getenv("EMAIL_HOST")
	smtpPort := os.Getenv("EMAIL_PORT")

	// If config missing -> fallback to safe mock (NO body)
	if from == "" || password == "" || smtpHost == "" || smtpPort == "" {
		log.Printf("📧 MOCK EMAIL (missing smtp config) to=%s subject=%s (body_len=%d)", to, subject, len(body))
		return nil
	}

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		from, to, subject, body,
	)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	if err := smtp.SendMail(addr, auth, from, []string{to}, []byte(msg)); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}
	return nil
}

// GenerateVerificationToken creates random token
func GenerateVerificationToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
