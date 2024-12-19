package exerciseLog

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExerciseLog struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID     string             `json:"userId" bson:"userId" validate:"required"`
	ExerciseID string             `json:"exerciseId" bson:"exerciseId" validate:"required"`
	Date       time.Time          `json:"date" bson:"date"`
	Sets       []SetLog           `json:"sets" validate:"required,dive"`
	Notes      string             `json:"notes"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
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
