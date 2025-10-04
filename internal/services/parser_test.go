package services_test

import (
	"recipe/internal/models"
	"recipe/internal/services"
	"testing"
)

func TestParseIngredients(t *testing.T) {
	testCases := []struct {
		name           string
		rawIngredients []models.Ingredient
		expected       []models.Ingredient
	}{
		{
			name:           "Simple case",
			rawIngredients: []models.Ingredient{{Name: "1 cup whole milk"}},
			expected: []models.Ingredient{
				{Quantity: 1, Unit: "cup", Name: "Whole Milk"},
			},
		},
		{
			name:           "No unit",
			rawIngredients: []models.Ingredient{{Name: "2 large eggs"}},
			expected: []models.Ingredient{
				{Quantity: 2, Unit: "", Name: "Large Eggs"},
			},
		},
		{
			name:           "No quantity or unit",
			rawIngredients: []models.Ingredient{{Name: "salt to taste"}},
			expected: []models.Ingredient{
				{Quantity: 1, Unit: "", Name: "Salt To Taste"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := services.ParseIngredients(tc.rawIngredients)
			if len(result) != len(tc.expected) {
				t.Fatalf("Expected %d ingredients, but got %d", len(tc.expected), len(result))
			}
			for i := range result {
				if result[i].Name != tc.expected[i].Name || result[i].Quantity != tc.expected[i].Quantity || result[i].Unit != tc.expected[i].Unit {
					t.Errorf("Expected %+v, but got %+v", tc.expected[i], result[i])
				}
			}
		})
	}
}
