package main

import (
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"recipe/internal/api"
	"recipe/internal/config"
	"recipe/internal/database"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found, using environment variables.")
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// --- Add this block to start the pprof server ---
	go func() {
		slog.Info("Starting pprof server on localhost:6060")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			slog.Error("pprof server failed to start", "error", err)
		}
	}()
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

	// Create and register routes for your main application
	router := mux.NewRouter()
	api.RegisterRoutes(router, db, cfg.SpoonacularKey)

	addr := ":" + cfg.Port
	slog.Info("Server starting", "address", addr)
	err = http.ListenAndServe(addr, router)
	if err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
