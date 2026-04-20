package config

import (
	"log"
	"os"
	"time"
)

type AuthConfig struct {
	JWTSecret            []byte
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	DatabaseUrl          string
}

var Env *AuthConfig

func InitEnv() {
	Env = &AuthConfig{
		JWTSecret:            []byte(os.Getenv("JWT_SECRET")),
		AccessTokenDuration:  getDuration("JWT_ACCESS_DURATION", 15*time.Minute),
		RefreshTokenDuration: getDuration("JWT_REFRESH_DURATION", 30*24*time.Hour),
		DatabaseUrl:          os.Getenv("DATABASE_URL"),
	}
}

func getDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if dur, err := time.ParseDuration(value); err == nil {
			return dur
		}
		log.Printf("Invalid duration for %s, using fallback", key)
	}
	return fallback
}
