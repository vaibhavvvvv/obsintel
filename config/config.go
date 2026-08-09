// for all env vars and config in one place
package config

import (
		"github.com/joho/godotenv"
		"log"
		"os"
)

//adding all the config params used throughout the codebase here
type Config struct {
    GEMINI_API_KEY string
	AI_MODEL		string
    Port         string
	VALID_API_KEYS string
	DBUser        string
    DBPassword    string
    DBHost        string
    DBPort        string
    DBName        string
}

var C Config

func Init(){
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	C = Config{
		GEMINI_API_KEY: os.Getenv("GEMINI_API_KEY"),
		AI_MODEL		: os.Getenv("AI_MODEL"),
        Port:         getEnv("PORT", "8080"),
		VALID_API_KEYS: os.Getenv("VALID_API_KEYS"),
		DBUser:       os.Getenv("DB_USER"),
		DBPassword:   os.Getenv("DB_PASSWORD"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBName:       os.Getenv("DB_NAME"),
	}
}

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}