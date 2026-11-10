package config

import (
	"os"
	"time"
)

// GetJWTSecret mengembalikan JWT secret key
func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "galuraswara-secret-key-ganti-di-production"
	}
	return secret
}

// GetJWTExpiry mengembalikan durasi expiry JWT token
func GetJWTExpiry() time.Duration {
	return 24 * time.Hour
}
