package config

import (
	"os"
	"strconv"
	"time"
)

var (
	// PasswordMaxAgeDays je maksimalni broj dana koliko lozinka može biti aktivna
	// Default: 60 dana
	// Za simulaciju na odbrani: postaviti PASSWORD_MAX_AGE_MINUTES env varijablu
	PasswordMaxAgeDays int = 60

	// PasswordMaxAgeDuration je izračunata duration za proveru isteka
	PasswordMaxAgeDuration time.Duration
)

// InitPasswordConfig učitava konfiguraciju za istek lozinke
// Ako je postavljena PASSWORD_MAX_AGE_MINUTES, koristi se ta vrednost u minutama (za demo)
// Inače se koristi PASSWORD_MAX_AGE_DAYS (default 60)
func InitPasswordConfig() {
	// Prvo proveri da li je postavljena demo varijabla u minutama
	if minutesStr := os.Getenv("PASSWORD_MAX_AGE_MINUTES"); minutesStr != "" {
		if minutes, err := strconv.Atoi(minutesStr); err == nil && minutes > 0 {
			PasswordMaxAgeDuration = time.Duration(minutes) * time.Minute
			PasswordMaxAgeDays = 0 // Oznaka da se koristi minute mode
			return
		}
	}

	// Inače koristi dane
	if daysStr := os.Getenv("PASSWORD_MAX_AGE_DAYS"); daysStr != "" {
		if days, err := strconv.Atoi(daysStr); err == nil && days > 0 {
			PasswordMaxAgeDays = days
		}
	}

	PasswordMaxAgeDuration = time.Duration(PasswordMaxAgeDays) * 24 * time.Hour
}

// GetPasswordMaxAgeString vraća čitljivu reprezentaciju perioda isteka
func GetPasswordMaxAgeString() string {
	if PasswordMaxAgeDays == 0 {
		// Koristi se minute mode (demo)
		minutes := int(PasswordMaxAgeDuration.Minutes())
		return strconv.Itoa(minutes) + " minuta"
	}
	return strconv.Itoa(PasswordMaxAgeDays) + " dana"
}

// IsPasswordExpired proverava da li je lozinka istekla
func IsPasswordExpired(passwordChangedAt time.Time) bool {
	return time.Since(passwordChangedAt) > PasswordMaxAgeDuration
}

// GetPasswordAge vraća starost lozinke
func GetPasswordAge(passwordChangedAt time.Time) time.Duration {
	return time.Since(passwordChangedAt)
}

// GetDaysUntilExpiry vraća broj dana do isteka (ili minuta ako je demo mode)
func GetDaysUntilExpiry(passwordChangedAt time.Time) int {
	age := time.Since(passwordChangedAt)
	remaining := PasswordMaxAgeDuration - age

	if remaining <= 0 {
		return 0
	}

	if PasswordMaxAgeDays == 0 {
		// Minute mode - vraća minute
		return int(remaining.Minutes())
	}

	// Vraća dane
	return int(remaining.Hours() / 24)
}
