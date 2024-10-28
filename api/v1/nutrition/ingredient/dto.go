package ingredient

import (
	"github.com/Npwskp/GymsbroBackend/api/v1/nutrition/types"
)

type CreateIngredientDto struct {
	UserID    string            `json:"userid" validate:"required"`
	Name      string            `json:"name" validate:"required"`
	Image     string            `json:"image" default:"null"`
	Calories  float64           `json:"calories" default:"0"`
	Nutrients *[]types.Nutrient `json:"nutrients,omitempty"`
}

type UpdateIngredientDto struct {
	UserID    string            `json:"userid" validate:"required"`
	Name      string            `json:"name"`
	Image     string            `json:"image"`
	Calories  float64           `json:"calories"`
	Nutrients *[]types.Nutrient `json:"nutrients,omitempty"`
}
