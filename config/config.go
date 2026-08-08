// for all env vars and config in one place
package config

import (
		"github.com/joho/godotenv"
		"log"
		"os"
)

//adding all the config params used throughout the codebase here
type Config struct {
    GeminiAPIKey string
    Port         string
	VALID_API_KEYS string
}

var C Config

func Init(){
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	C = Config{
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
        Port:         getEnv("PORT", "8080"),
		VALID_API_KEYS: os.Getenv("VALID_API_KEYS"),
	}
}

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}