package models

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	GinMode         string
	MongoURI        string
	MongoDatabase   string
	MongoCollection string
}

func LoadConfig() (Config, error) {
	_ = godotenv.Load()

	config := Config{
		Port:            getEnv("PORT", "8080"),
		GinMode:         getEnv("GIN_MODE", "release"),
		MongoURI:        os.Getenv("MONGO_URI"),
		MongoDatabase:   getEnv("MONGO_DATABASE", "movie_booking"),
		MongoCollection: getEnv("MONGO_COLLECTION", "movies"),
	}

	if config.MongoURI == "" {
		return Config{}, fmt.Errorf("MONGO_URI is required")
	}

	return config, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
