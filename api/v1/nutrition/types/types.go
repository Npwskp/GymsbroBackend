package types

type Nutrient struct {
	Name   string  `json:"name" validate:"required"`
	Amount float64 `json:"amount" validate:"required"`
	Unit   string  `json:"unit" validate:"required"`
}

type Ingredient struct {
	IngredientId  string  `json:"ingredientid" validate:"required" bson:"ingredientid"`
	Name          string  `json:"name" validate:"required" bson:"name"`
	Value         float64 `json:"value" validate:"required"`
	Uint          string  `json:"unit" validate:"required"`
	NumOfServings float64 `json:"numofservings" validate:"required"`
}
