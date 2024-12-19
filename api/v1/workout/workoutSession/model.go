package workoutSession

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkoutSession struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID      string             `json:"userId" bson:"userId" validate:"required"`
	WorkoutID   string             `json:"workoutId" bson:"workoutId" validate:"required"`
	StartTime   time.Time          `json:"startTime" bson:"startTime"`
	EndTime     *time.Time         `json:"endTime" bson:"endTime"`
	Status      SessionStatus      `json:"status" bson:"status"`
	TotalVolume float64            `json:"totalVolume" bson:"totalVolume"`
	Duration    int                `json:"duration" bson:"duration"` // in seconds
	Exercises   []ExerciseEntry    `json:"exercises" validate:"required,dive"`
	Notes       string             `json:"notes"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type SessionStatus string

const (
	StatusInProgress SessionStatus = "in_progress"
	StatusCompleted  SessionStatus = "completed"
	StatusCancelled  SessionStatus = "cancelled"
)

type ExerciseEntry struct {
	ExerciseID    string    `json:"exerciseId" validate:"required"`
	ExerciseLogID string    `json:"exerciseLogId" bson:"exerciseLogId"`
	PlannedSets   int       `json:"plannedSets"`
	CompletedSets int       `json:"completedSets"`
	TotalVolume   float64   `json:"totalVolume"`
	StartTime     time.Time `json:"startTime"`
	EndTime       time.Time `json:"endTime"`
}
