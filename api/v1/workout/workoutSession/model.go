package workoutSession

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkoutSession struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID      string             `json:"userid" bson:"userid" validate:"required"`
	WorkoutID   string             `json:"workoutid" bson:"workoutid"`
	Type        SessionType        `json:"type" bson:"type" validate:"required"`
	StartTime   time.Time          `json:"start_time" bson:"start_time"`
	EndTime     time.Time          `json:"end_time" bson:"end_time"`
	Status      SessionStatus      `json:"status" bson:"status"`
	TotalVolume float64            `json:"total_volume" bson:"total_volume"`
	Duration    int                `json:"duration" bson:"duration"`
	Exercises   []SessionExercise  `json:"exercises" validate:"dive"`
	Notes       string             `json:"notes" bson:"notes"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type SessionType string

const (
	PlannedSession SessionType = "planned"
	CustomSession  SessionType = "custom"
	LoggedSession  SessionType = "logged"
)

type SessionExercise struct {
	ExerciseID    string `json:"exerciseid" bson:"exerciseid" validate:"required"`
	ExerciseLogID string `json:"exerciselogid" bson:"exerciselogid"`
	Order         int    `json:"order" bson:"order" validate:"required,min=0"`
}
type SessionStatus string

const (
	StatusInProgress SessionStatus = "in_progress"
	StatusCompleted  SessionStatus = "completed"
)
