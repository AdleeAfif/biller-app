package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	MongoDBURI   string
	DatabaseName string
	JWTSecret    string
	Port         string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Try to load .env file (ignore error if not found)
	_ = godotenv.Load()

	config := &Config{
		MongoDBURI:   getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		DatabaseName: getEnv("DATABASE_NAME", "biller_app"),
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		Port:         getEnv("PORT", "8080"),
	}

	log.Printf("Config loaded: Port=%s, Database=%s", config.Port, config.DatabaseName)
	return config
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
