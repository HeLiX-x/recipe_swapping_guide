package api

import (
	"encoding/json"
	"net/http"
	"recipe/internal/models"
	"recipe/internal/services"

	"gorm.io/gorm"
)

type Handler struct {
	db           *gorm.DB
	spoonService *services.SpoonacularService
}

func NewHandler(db *gorm.DB, spoonService *services.SpoonacularService) *Handler {
	return &Handler{
		db:           db,
		spoonService: spoonService,
	}
}

func (h *Handler) UploadRecipeHandler(w http.ResponseWriter, r *http.Request) {
	var recipe models.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	comparison, err := h.spoonService.CreateRecipeAndSuggestSwaps(h.db, &recipe)
	if err != nil {
		http.Error(w, "Failed to process recipe", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comparison)
}

func (h *Handler) GetAllRecipesHandler(w http.ResponseWriter, r *http.Request) {
	var recipes []models.Recipe
	if err := h.db.Preload("Ingredients").Find(&recipes).Error; err != nil {
		http.Error(w, "Could not retrieve recipes", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}
