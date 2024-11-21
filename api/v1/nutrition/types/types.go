package types

type Nutrient struct {
	Name   string  `json:"name" validate:"required" bson:"name"`
	Amount float64 `json:"amount" validate:"required" bson:"amount"`
	Unit   string  `json:"unit" validate:"required" bson:"unit"`
}

type Ingredient struct {
	IngredientId  string  `json:"ingredientId" validate:"required" bson:"ingredientid"`
	Name          string  `json:"name" validate:"required" bson:"name"`
	Value         float64 `json:"value" validate:"required" bson:"value"`
	Uint          string  `json:"unit" validate:"required" bson:"unit"`
	NumOfServings float64 `json:"numOfServings" validate:"required" bson:"num_of_servings"`
}
