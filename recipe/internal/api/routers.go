package api

import (
	"log/slog"
	"net/http"
	"recipe/internal/services"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterRoutes(r *mux.Router, db *gorm.DB, apiKey string) {
	spoonService := services.NewSpoonacularService(apiKey)

	h := NewHandler(db, spoonService)

	r.HandleFunc("/upload", h.UploadRecipeHandler).Methods("POST")
	r.HandleFunc("/recipes", h.GetAllRecipesHandler).Methods("GET")

	r.Use(LoggingMiddleware)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("request handled",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
		)
	})
}
