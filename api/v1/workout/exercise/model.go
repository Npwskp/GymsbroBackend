package exercise

import "time"

type Exercise struct {
	ID          string    `json:"id,omitempty" bson:"_id,omitempty"`
	UserID      string    `json:"userid" bson:"userid" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Type        []string  `json:"type" validate:"required"`
	Muscle      []string  `json:"muscle" validate:"required"`
	Image       string    `json:"image" validate:"required"`
	CreatedAt   time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
