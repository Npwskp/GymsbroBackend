package types

type Nutrient struct {
	Name     string  `json:"name" validate:"required"`
	Type     string  `json:"type" validate:"required"`
	Category string  `json:"category"`
	Value    float64 `json:"value" validate:"required"`
	Unit     string  `json:"unit" validate:"required"`
}

type Ingredient struct {
	IngredientId  string  `json:"ingredientid" validate:"required" bson:"ingredientid"`
	Value         float64 `json:"value" validate:"required"`
	Uint          string  `json:"unit" validate:"required"`
	NumOfServings float64 `json:"numofservings" validate:"required"`
}
