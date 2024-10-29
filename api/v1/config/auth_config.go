package config

import (
	"os"
	"time"
)

const (
	JWTExpirationTime = 72 * time.Hour
	BcryptCost        = 12 // Higher than default (10)
	CookieSecure      = true
	CookieHTTPOnly    = true
	CookieSameSite    = "Strict"
)

// GetJWTSecret retrieves JWT secret from environment variable
func GetJWTSecret() string {
	return os.Getenv("JWT_SECRET")
}
