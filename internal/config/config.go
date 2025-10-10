package config

import (
	"log"
	"os"
)

type Config struct {
	ServerPort    string
	SessionKey    []byte
	AdminUsername string
	AdminPassword string
	DatabasePath  string
	UploadDir     string
	PublicDir     string
	TemplatesDir  string
}

func Load() *Config {
	// Load from environment variables or use defaults
	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		log.Println("WARNING: Using default session key. Set SESSION_KEY environment variable in production!")
		sessionKey = "super-secret-key-that-is-32-bytes-long-so-it-is-secure"
	}

	return &Config{
		ServerPort:    getEnv("PORT", "8080"),
		SessionKey:    []byte(sessionKey),
		AdminUsername: getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "password123"),
		DatabasePath:  getEnv("DB_PATH", "school.db"),
		UploadDir:     "public/uploads",
		PublicDir:     "public",
		TemplatesDir:  "templates",
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
