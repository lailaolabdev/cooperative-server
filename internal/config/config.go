package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port, MongoURI, Database, JWTSecret string
	AllowedOrigins                      []string
}

func Load() (Config, error) {
	// Load local development values when .env exists. Existing process
	// environment variables keep priority, which is required in production.
	_ = godotenv.Load()

	c := Config{
		Port: env("PORT", "8080"), MongoURI: env("MONGODB_URI", "mongodb://localhost:27017"),
		Database: env("MONGODB_DATABASE", "cooperative_db"), JWTSecret: os.Getenv("JWT_SECRET"),
		AllowedOrigins: strings.Split(env("ALLOWED_ORIGINS", "http://localhost:5173"), ","),
	}
	if len(c.JWTSecret) < 32 {
		return c, fmt.Errorf("JWT_SECRET must contain at least 32 characters")
	}
	return c, nil
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
