package workout

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Workout struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID      string             `json:"userid" bson:"userid"`
	Name        string             `json:"name" bson:"name" validate:"required"`
	Description string             `json:"description" bson:"description"`
	Exercises   []WorkoutExercise  `json:"exercises" bson:"exercises" validate:"required"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type WorkoutExercise struct {
	ExerciseID string `json:"exerciseid" bson:"exerciseid" validate:"required"`
	Order      int    `json:"order" bson:"order" validate:"min=0"`
}
