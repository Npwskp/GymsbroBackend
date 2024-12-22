package exerciseLog

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExerciseLog struct {
	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID           string             `json:"userId" bson:"userId" validate:"required"`
	ExerciseID       string             `json:"exerciseId" validate:"required"`
	WorkoutSessionID string             `json:"workoutSessionId" bson:"workoutSessionId" validate:"required"`
	CompletedSets    int                `json:"completedSets"`
	TotalVolume      float64            `json:"totalVolume"`
	Notes            string             `json:"notes"`
	TimeUsedInSec    int                `json:"timeUsedInSec"`
	Date             time.Time          `json:"date"`
	Sets             []SetLog           `json:"sets" validate:"dive"`
	CreatedAt        time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt        time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type SetType string

const (
	WarmUpSet  SetType = "warm_up"
	WorkingSet SetType = "working"
	DropSet    SetType = "drop"
	FailureSet SetType = "failure"
	BackOffSet SetType = "back_off"
)

type SetLog struct {
	Weight    float64 `json:"weight" validate:"required,min=0"`
	Reps      int     `json:"reps" validate:"required,min=0"`
	SetNumber int     `json:"setNumber" validate:"required,min=1"`
	Type      SetType `json:"type" validate:"required,oneof=warm_up working drop failure back_off"`
	RPE       *int    `json:"rpe" validate:"omitempty,min=1,max=10"`
}
