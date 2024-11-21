package foodlog

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FoodLog struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    string             `json:"userid" validate:"required" bson:"userid"`
	DateTime  time.Time          `json:"datetime" validate:"required" bson:"datetime"`
	Meals     []string           `json:"meals"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdateAt  time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func CreateFoodLogModel(dto *CreateFoodLogDto) *FoodLog {
	return &FoodLog{
		UserID:    "",
		DateTime:  dto.DateTime,
		Meals:     dto.Meals,
		CreatedAt: time.Now(),
		UpdateAt:  time.Now(),
	}
}
