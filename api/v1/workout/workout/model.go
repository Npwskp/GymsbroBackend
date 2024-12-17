package workout

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Workout struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID      string             `json:"userId" bson:"userId"`
	Name        string             `json:"name" validate:"required"`
	Description string             `json:"description"`
	Exercises   []WorkoutExercise  `json:"exercises" validate:"required"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type WorkoutExercise struct {
	ExerciseID string `json:"exerciseId" validate:"required"`
	Sets       int    `json:"sets" validate:"required,min=1"`
	Reps       int    `json:"reps" validate:"required,min=1"`
	Order      int    `json:"order" validate:"required,min=1"`
}
