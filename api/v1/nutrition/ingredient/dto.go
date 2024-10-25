package ingredient

import "time"

type CreateIngredientDto struct {
	UserID    string    `json:"userid" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Image     string    `json:"image" default:"null"`
	Calories  float64   `json:"calories" default:"0"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty" default:"null"`
}

type UpdateIngredientDto struct {
	UserID   string  `json:"userid" validate:"required"`
	Name     string  `json:"name"`
	Image    string  `json:"image"`
	Calories float64 `json:"calories"`
}
