package foodlog

import "github.com/Npwskp/GymsbroBackend/api/v1/nutrition/types"

type AddMealToFoodLogDto struct {
	Date  string   `json:"date" validate:"required"`
	Meals []string `json:"meals"`
}

type UpdateFoodLogDto struct {
	Date  string   `json:"date" validate:"required"`
	Meals []string `json:"meals"`
}

type DailyNutrientResponse struct {
	Date      string           `json:"date"`
	Calories  float64          `json:"calories"`
	Nutrients []types.Nutrient `json:"nutrients"`
}
