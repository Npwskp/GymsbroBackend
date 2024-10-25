package meal

import "github.com/Npwskp/GymsbroBackend/api/v1/nutrition/types"

type CreateMealDto struct {
	Name        string             `json:"name" validate:"required"`
	UserID      string             `json:"userid"`
	Image       string             `json:"image"`
	Calories    float64            `json:"calories"`
	Nutrients   []types.Nutrient   `json:"nutrients"`
	Ingredients []types.Ingredient `json:"ingredients"`
}

type UpdateMealDto struct {
	Carb     float64 `json:"carb"`
	Protein  float64 `json:"protein"`
	Fat      float64 `json:"fat"`
	Calories float64 `json:"calories"`
}
