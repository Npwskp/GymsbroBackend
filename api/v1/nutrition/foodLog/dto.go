package foodlog

import "time"

type CreateFoodLogDto struct {
	DateTime time.Time `json:"date" validate:"required"`
	Meals    []string  `json:"meals"`
}

type UpdateFoodLogDto struct {
	DateTime time.Time `json:"date" validate:"required"`
	Meals    []string  `json:"meals"`
}
