package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"recipe/internal/models"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
)

type SpoonacularService struct {
	apiKey       string
	client       *http.Client
	cache        *sync.Map
	DefaultSwaps map[string]string
}

func NewSpoonacularService(apiKey string) *SpoonacularService {
	return &SpoonacularService{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		cache: new(sync.Map),
		DefaultSwaps: map[string]string{
			"whole milk":     "almond milk",
			"butter":         "olive oil spray",
			"cheddar cheese": "reduced-fat cheddar",
			"white bread":    "whole wheat bread",
			"sour cream":     "greek yogurt",
			"regular pasta":  "whole grain pasta",
		},
	}
}

var unitRegex = regexp.MustCompile(`(?i)^(teaspoons?|tbsps?|cups?|grams?|g|kg|ml|liters?|oz|lbs?|pounds?|tablespoons?|tsps?|fl oz)$`)

func ParseIngredients(rawIngredients []models.Ingredient) []models.Ingredient {
	var ingredients []models.Ingredient
	for _, item := range rawIngredients {
		line := strings.TrimSpace(item.Name)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		var quantity float64 = 1.0
		var unit string = ""
		nameParts := parts

		if q, err := strconv.ParseFloat(parts[0], 64); err == nil {
			quantity = q
			nameParts = parts[1:]
		}

		if len(nameParts) > 0 && unitRegex.MatchString(nameParts[0]) {
			unit = nameParts[0]
			nameParts = nameParts[1:]
		}

		name := strings.Join(nameParts, " ")
		finalName := strings.Title(strings.ToLower(strings.TrimSpace(name)))

		ingredients = append(ingredients, models.Ingredient{
			Quantity: quantity,
			Unit:     strings.ToLower(unit),
			Name:     finalName,
		})
	}
	return ingredients
}

func (s *SpoonacularService) CreateRecipeAndSuggestSwaps(db *gorm.DB, recipe *models.Recipe) (*models.RecipeComparison, error) {
	recipe.Ingredients = ParseIngredients(recipe.Ingredients)

	if err := db.Create(recipe).Error; err != nil {
		return nil, fmt.Errorf("failed to save recipe: %w", err)
	}

	swaps, total := s.substituteIngredients(recipe.Ingredients)

	return &models.RecipeComparison{
		OriginalRecipe:     *recipe,
		HealthierRecipe:    swaps,
		TotalCaloriesSaved: total,
	}, nil
}

func (s *SpoonacularService) substituteIngredients(ingredients []models.Ingredient) ([]models.IngredientSwap, float64) {
	var swaps []models.IngredientSwap
	var total float64 = 0
	var wg sync.WaitGroup
	resultChan := make(chan models.IngredientSwap, len(ingredients))
	for _, ing := range ingredients {
		wg.Add(1)
		go s.processIngredient(ing, &wg, resultChan)
	}
	wg.Wait()
	close(resultChan)
	for swap := range resultChan {
		total += swap.CaloriesSaved
		swaps = append(swaps, swap)
	}
	return swaps, total
}

func (s *SpoonacularService) processIngredient(ing models.Ingredient, wg *sync.WaitGroup, resultChan chan<- models.IngredientSwap) {
	defer wg.Done()
	handleErr := func(err error, stage string) {
		slog.Warn("Error processing ingredient", "ingredient", ing.Name, "stage", stage, "error", err)
		resultChan <- models.IngredientSwap{OriginalName: ing.Name, Quantity: ing.Quantity, Unit: ing.Unit}
	}
	searchResult, err := s.getIngredientInfo(ing.Name)
	if err != nil {
		handleErr(err, "search")
		return
	}
	var searchRes models.SpoonIngredientSearch
	if err := json.Unmarshal(searchResult, &searchRes); err != nil || len(searchRes.Results) == 0 {
		handleErr(fmt.Errorf("no results found or unmarshal error: %w", err), "unmarshal search")
		return
	}
	origID := searchRes.Results[0].ID
	nutriData, err := s.getIngredientNutrition(origID)
	if err != nil {
		handleErr(err, "get nutrition")
		return
	}
	var nutriRes models.SpoonIngredientNutrition
	if err := json.Unmarshal(nutriData, &nutriRes); err != nil {
		handleErr(err, "unmarshal nutrition")
		return
	}
	var origCal float64
	for _, n := range nutriRes.Nutrition.Nutrients {
		if n.Name == "Calories" {
			origCal = n.Amount
			break
		}
	}
	var suggCal float64
	suggestedName := s.DefaultSwaps[strings.ToLower(ing.Name)]
	if suggestedName != "" {
		suggInfo, err := s.getIngredientInfo(suggestedName)
		if err == nil {
			var suggSearchRes models.SpoonIngredientSearch
			if json.Unmarshal(suggInfo, &suggSearchRes) == nil && len(suggSearchRes.Results) > 0 {
				suggID := suggSearchRes.Results[0].ID
				suggNutriDetail, err := s.getIngredientNutrition(suggID)
				if err == nil {
					var suggNutriRes models.SpoonIngredientNutrition
					if json.Unmarshal(suggNutriDetail, &suggNutriRes) == nil {
						for _, n := range suggNutriRes.Nutrition.Nutrients {
							if n.Name == "Calories" {
								suggCal = n.Amount
								break
							}
						}
					}
				}
			}
		}
	}
	saved := (origCal - suggCal) * ing.Quantity
	if saved < 0 {
		saved = 0
	}
	resultChan <- models.IngredientSwap{
		OriginalName:      ing.Name,
		SuggestedName:     suggestedName,
		OriginalCalories:  origCal,
		SuggestedCalories: suggCal,
		Quantity:          ing.Quantity,
		Unit:              ing.Unit,
		CaloriesSaved:     saved,
	}
}

func (s *SpoonacularService) getFromAPI(url, cacheKey string) ([]byte, error) {
	if cached, ok := s.cache.Load(cacheKey); ok {
		return cached.([]byte), nil
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	s.cache.Store(cacheKey, body)
	return body, nil
}

func (s *SpoonacularService) getIngredientInfo(query string) ([]byte, error) {
	url := fmt.Sprintf("https://api.spoonacular.com/food/ingredients/search?query=%s&apiKey=%s", query, s.apiKey)
	return s.getFromAPI(url, "search_"+query)
}

func (s *SpoonacularService) getIngredientNutrition(id int) ([]byte, error) {
	url := fmt.Sprintf("https://api.spoonacular.com/food/ingredients/%d/information?amount=1&unit=cup&apiKey=%s", id, s.apiKey)
	return s.getFromAPI(url, fmt.Sprintf("nutrition_%d", id))
}
