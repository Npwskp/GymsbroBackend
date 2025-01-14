package exerciseLog

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExerciseLog struct {
	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID           string             `json:"userid" bson:"userid" validate:"required"`
	ExerciseID       string             `json:"exerciseid" validate:"required"`
	WorkoutSessionID string             `json:"workoutsessionid" bson:"workoutsessionid" validate:"required"`
	CompletedSets    int                `json:"completed_sets" bson:"completed_sets"`
	TotalVolume      float64            `json:"total_volume" bson:"total_volume"`
	Notes            string             `json:"notes" bson:"notes"`
	Duration         int                `json:"duration" bson:"duration"`
	DateTime         time.Time          `json:"datetime" bson:"datetime"`
	Sets             []SetLog           `json:"sets" validate:"dive"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
}

type SetType string

const (
	WarmUpSet  SetType = "warm_up"
	WorkingSet SetType = "working"
	DropSet    SetType = "drop"
	FailureSet SetType = "failure"
)

type SetLog struct {
	Weight    float64 `json:"weight" validate:"required,min=1"`
	Reps      int     `json:"reps" validate:"required,min=1"`
	SetNumber int     `json:"setNumber" validate:"required,min=1"`
	Type      SetType `json:"type" validate:"required,oneof=warm_up working drop failure"`
	RPE       *int    `json:"rpe" validate:"omitempty,min=1,max=10"`
}
