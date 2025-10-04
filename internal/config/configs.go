package config

import (
	"log/slog"
	"os"
)

type Config struct {
	Port           string
	DSN            string
	SpoonacularKey string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		slog.Info("DB_DSN not set, using default database connection string.")
		dsn = "recipe_user:password@tcp(127.0.0.1:3306)/recipe_db?charset=utf8mb4&parseTime=True&loc=Local"
	}

	apiKey := os.Getenv("SPOONACULAR_API_KEY")

	return &Config{
		Port:           port,
		DSN:            dsn,
		SpoonacularKey: apiKey,
	}
}
