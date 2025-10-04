package models

import "gorm.io/gorm"

type Ingredient struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	RecipeID uint    `json:"-"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
	Name     string  `json:"name" gorm:"index"`
}

type Recipe struct {
	gorm.Model
	ID           uint         `gorm:"primaryKey" json:"id"`
	Title        string       `json:"title" gorm:"index"`
	Ingredients  []Ingredient `json:"ingredients"`
	Instructions string       `json:"instructions"`
}

type IngredientSwap struct {
	OriginalName      string  `json:"original_name"`
	SuggestedName     string  `json:"suggested_name"`
	OriginalCalories  float64 `json:"original_calories"`
	SuggestedCalories float64 `json:"suggested_calories"`
	Quantity          float64 `json:"quantity"`
	Unit              string  `json:"unit"`
	CaloriesSaved     float64 `json:"calories_saved"`
}

type RecipeComparison struct {
	OriginalRecipe     Recipe           `json:"original_recipe"`
	HealthierRecipe    []IngredientSwap `json:"healthier_recipe"`
	TotalCaloriesSaved float64          `json:"total_calories_saved"`
}

type SpoonIngredientSearch struct {
	Results []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"results"`
}

type SpoonIngredientNutrition struct {
	Nutrition struct {
		Nutrients []struct {
			Name   string  `json:"name"`
			Amount float64 `json:"amount"`
			Unit   string  `json:"unit"`
		} `json:"nutrients"`
	} `json:"nutrition"`
}
