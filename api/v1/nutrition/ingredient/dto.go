package ingredient

import (
	"github.com/Npwskp/GymsbroBackend/api/v1/nutrition/types"
)

type CreateIngredientDto struct {
	Name        string            `json:"name" validate:"required"`
	Description string            `json:"description"`
	Category    string            `json:"category"`
	Image       string            `json:"image" default:"null"`
	Calories    float64           `json:"calories" default:"0"`
	Nutrients   *[]types.Nutrient `json:"nutrients,omitempty"`
}

type UpdateIngredientDto struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Category    string            `json:"category"`
	Image       string            `json:"image"`
	Calories    float64           `json:"calories"`
	Nutrients   *[]types.Nutrient `json:"nutrients,omitempty"`
}
