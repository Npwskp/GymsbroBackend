package meal

import (
	"time"

	"github.com/Npwskp/GymsbroBackend/src/nutrition/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Meal struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" validate:"required" bson:"name"`
	UserID      string             `json:"userid" bson:"userid"`
	Image       string             `json:"image" default:"null"`
	Calories    float64            `json:"calories" default:"0"`
	Nutrients   []types.Nutrient   `json:"nutrients,omitempty"`
	Ingredients []types.Ingredient `json:"ingredients,omitempty"`
	CreatedAt   time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty" default:"null"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty" default:"null"`
}

func CreateMealModel(dto *CreateMealDto) *Meal {
	return &Meal{
		Name:        dto.Name,
		UserID:      dto.UserID,
		Image:       dto.Image,
		Calories:    dto.Calories,
		Nutrients:   dto.Nutrients,
		Ingredients: dto.Ingredients,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
