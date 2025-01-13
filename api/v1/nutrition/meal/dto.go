package meal

import "github.com/Npwskp/GymsbroBackend/api/v1/nutrition/types"

type CreateMealDto struct {
	Name        string             `json:"name" validate:"required"`
	Description string             `json:"description"`
	Category    string             `json:"category"`
	Image       string             `json:"image"`
	Calories    float64            `json:"calories"`
	Nutrients   []types.Nutrient   `json:"nutrients"`
	Ingredients []types.Ingredient `json:"ingredients"`
	IsQuickAdd  bool               `json:"isQuickAdd" bson:"is_quick_add" default:"false"`
}

type UpdateMealDto struct {
	Description string             `json:"description"`
	Category    string             `json:"category"`
	Image       string             `json:"image"`
	Calories    float64            `json:"calories"`
	Nutrients   []types.Nutrient   `json:"nutrients"`
	Ingredients []types.Ingredient `json:"ingredients"`
}

type SearchFilters struct {
	Query       string `json:"q"` // Search query
	Category    string `json:"category"`
	MinCalories int    `json:"minCalories"`
	MaxCalories int    `json:"maxCalories"`
	Nutrients   string `json:"nutrients"`
	UserID      string `json:"userid"`
}

type CalculateNutrientBody struct {
	Ingredients []types.Ingredient `json:"ingredients"`
}

type CalculateNutrientResponse struct {
	Nutrients []types.Nutrient `json:"nutrients"`
	Calories  float64          `json:"calories"`
}
