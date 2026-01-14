package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var logFile *os.File

// InitLogger initializes the log file
func InitLogger() error {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	logPath := filepath.Join(logDir, "security.log")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	logFile = file
	log.SetOutput(file)
	return nil
}

// CloseLogger closes the log file
func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}

// LogSecurityEvent logs security-related events
func LogSecurityEvent(eventType, action, ip, details string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] TYPE=%s ACTION=%s IP=%s DETAILS=%s\n",
		timestamp, eventType, action, ip, details)

	// Log to file
	if logFile != nil {
		logFile.WriteString(logMessage)
	}

	// Also log to stdout for Docker logs
	fmt.Print(logMessage)
}

// RotateLogs rotates log files if they exceed size limit (10MB)
func RotateLogs() {
	logPath := "logs/security.log"
	info, err := os.Stat(logPath)
	if err != nil {
		return
	}

	// If file size > 10MB, rotate
	if info.Size() > 10*1024*1024 {
		timestamp := time.Now().Format("20060102-150405")
		newPath := fmt.Sprintf("logs/security-%s.log", timestamp)

		// Close current file
		if logFile != nil {
			logFile.Close()
		}

		// Rename old file
		os.Rename(logPath, newPath)

		// Create new file
		InitLogger()
	}
}
