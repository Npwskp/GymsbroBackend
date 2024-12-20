package foodlog

type AddMealToFoodLogDto struct {
	Date  string   `json:"date" validate:"required"`
	Meals []string `json:"meals"`
}

type UpdateFoodLogDto struct {
	Date  string   `json:"date" validate:"required"`
	Meals []string `json:"meals"`
}
