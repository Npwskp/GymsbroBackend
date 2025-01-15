package workoutPlan

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkoutPlan struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    string             `json:"userid" bson:"userid" validate:"required"`
	WorkoutID string             `json:"workoutid" bson:"workoutid" validate:"required"`
	Dates     []time.Time        `json:"dates" bson:"dates" validate:"required,min=1"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
