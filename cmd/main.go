package main

import (
	"log/slog"
	"net/http"
	"os"
	"recipe/internal/api"
	"recipe/internal/config"
	"recipe/internal/database"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists.
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found, using environment variables.")
	}

	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load configuration
	cfg := config.Load()
	if cfg.SpoonacularKey == "" {
		slog.Error("SPOONACULAR_API_KEY environment variable not set. Aborting.")
		os.Exit(1)
	}

	// Connect to the database
	db, err := database.ConnectDB(cfg.DSN)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	slog.Info("Successfully connected to the database.")

	// Create and register routes
	router := mux.NewRouter()
	api.RegisterRoutes(router, db, cfg.SpoonacularKey)

	// Start the server
	addr := ":" + cfg.Port
	slog.Info("Server starting", "address", addr)
	err = http.ListenAndServe(addr, router)
	if err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
