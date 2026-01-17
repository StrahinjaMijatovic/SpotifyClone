package utils

import (
	"bufio"
	"html"
	"os"
	"regexp"
	"strings"
	"sync"
	"unicode"
)

var (
	passwordBlacklist     map[string]bool
	passwordBlacklistOnce sync.Once
)

// loadPasswordBlacklist loads the password blacklist from file
func loadPasswordBlacklist() {
	passwordBlacklist = make(map[string]bool)

	file, err := os.Open("data/password_blacklist.txt")
	if err != nil {
		// Try alternative path (when running from different directory)
		file, err = os.Open("users-service/data/password_blacklist.txt")
		if err != nil {
			return
		}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Store lowercase for case-insensitive matching
		passwordBlacklist[strings.ToLower(line)] = true
	}
}

// IsPasswordBlacklisted checks if password is in the blacklist
func IsPasswordBlacklisted(password string) bool {
	passwordBlacklistOnce.Do(loadPasswordBlacklist)
	return passwordBlacklist[strings.ToLower(password)]
}

// ValidatePasswordStrength checks if password meets requirements
func ValidatePasswordStrength(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// SanitizeString removes dangerous characters and limits length
func SanitizeString(s string) string {
	// HTML escape
	s = html.EscapeString(s)
	// Remove null bytes
	s = strings.ReplaceAll(s, "\x00", "")
	// Remove control characters
	s = strings.Map(func(r rune) rune {
		if r < 32 && r != '\n' && r != '\r' && r != '\t' {
			return -1
		}
		return r
	}, s)
	// Trim whitespace
	s = strings.TrimSpace(s)
	return s
}

// ValidateUsername checks username format
func ValidateUsername(username string) bool {
	// Only alphanumeric and underscore, 3-50 chars
	if len(username) < 3 || len(username) > 50 {
		return false
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return usernameRegex.MatchString(username)
}

// ValidateName checks if name contains only letters and spaces
func ValidateName(name string) bool {
	if len(name) < 2 || len(name) > 50 {
		return false
	}
	nameRegex := regexp.MustCompile(`^[a-zA-ZÀ-ÿ\s'-]+$`)
	return nameRegex.MatchString(name)
}

// BoundaryCheck ensures string is within limits
func BoundaryCheck(s string, min, max int) bool {
	length := len(s)
	return length >= min && length <= max
}

// ContainsSpecialChars checks for SQL injection characters
func ContainsSpecialChars(s string) bool {
	dangerousChars := []string{"'", "\"", ";", "--", "/*", "*/", "xp_", "sp_", "<script", "</script>"}
	lower := strings.ToLower(s)
	for _, char := range dangerousChars {
		if strings.Contains(lower, char) {
			return true
		}
	}
	return false
}

// ValidateNumeric checks if string contains only numbers
func ValidateNumeric(s string) bool {
	numericRegex := regexp.MustCompile(`^[0-9]+$`)
	return numericRegex.MatchString(s)
}

// Whitelist checks if string contains only allowed characters
func Whitelist(s string, allowedChars string) bool {
	pattern := "^[" + regexp.QuoteMeta(allowedChars) + "]+$"
	whitelistRegex := regexp.MustCompile(pattern)
	return whitelistRegex.MatchString(s)
}
