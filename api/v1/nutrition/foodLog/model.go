package foodlog

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FoodLog struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    string             `json:"userid" validate:"required" bson:"userid"`
	Date      string             `json:"date" validate:"required" bson:"date"`
	Meals     []string           `json:"meals"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdateAt  time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
