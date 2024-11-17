package meal

import "github.com/Npwskp/GymsbroBackend/api/v1/nutrition/types"

type CreateMealDto struct {
	Name        string             `json:"name" validate:"required"`
	Description string             `json:"description"`
	Category    string             `json:"category"`
	UserID      string             `json:"userid"`
	Image       string             `json:"image"`
	Calories    float64            `json:"calories"`
	Nutrients   []types.Nutrient   `json:"nutrients"`
	Ingredients []types.Ingredient `json:"ingredients"`
}

type UpdateMealDto struct {
	Description string  `json:"description"`
	Carb        float64 `json:"carb"`
	Protein     float64 `json:"protein"`
	Fat         float64 `json:"fat"`
	Calories    float64 `json:"calories"`
}

type SearchFilters struct {
	Query       string `json:"q"` // Search query
	Category    string `json:"category"`
	MinCalories int    `json:"minCalories"`
	MaxCalories int    `json:"maxCalories"`
	Nutrients   string `json:"nutrients"`
	UserID      string `json:"userid"`
}
