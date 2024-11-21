package foodlog

import "time"

type CreateFoodLogDto struct {
	UserID    string    `json:"userid" validate:"required"`
	Date      string    `json:"date" validate:"required"`
	Meals     []string  `json:"meals"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type UpdateFoodLogDto struct {
	UserID string   `json:"userid" validate:"required"`
	Date   string   `json:"date" validate:"required"`
	Meals  []string `json:"meals"`
}
