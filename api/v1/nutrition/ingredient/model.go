package ingredient

import (
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/nutrition/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ingredient struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    string             `json:"userid" validate:"required" bson:"userid"`
	Name      string             `json:"name" validate:"required" bson:"name"`
	Image     string             `json:"image" default:"null"`
	Calories  float64            `json:"calories" default:"0"`
	Nutrients *[]types.Nutrient  `json:"nutrients,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty" default:"null"`
	UpdateAt  time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty" default:"null"`
}

func CreateIngredientModel(dto *CreateIngredientDto) *Ingredient {
	return &Ingredient{
		UserID:    dto.UserID,
		Name:      dto.Name,
		Image:     dto.Image,
		Calories:  dto.Calories,
		Nutrients: dto.Nutrients,
		CreatedAt: time.Now(),
		UpdateAt:  time.Now(),
	}
}

func GetMockIngredient() *Ingredient {
	return &Ingredient{
		ID:        primitive.NewObjectID(),
		UserID:    "mockuserid",
		Name:      "mockname",
		Image:     "mockimage",
		Calories:  100,
		Nutrients: nil,
		CreatedAt: time.Now(),
		UpdateAt:  time.Now(),
	}
}
