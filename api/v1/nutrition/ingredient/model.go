package ingredient

import (
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/nutrition/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ingredient struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID      string             `json:"userid" validate:"required" bson:"userid"`
	Name        string             `json:"name" validate:"required" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Category    string             `json:"category" bson:"category"`
	Image       string             `json:"image" default:"null"`
	Calories    float64            `json:"calories" default:"0"`
	Nutrients   []types.Nutrient   `json:"nutrients,omitempty"`
	CreatedAt   time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdateAt    time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func CreateIngredientModel(dto *CreateIngredientDto) *Ingredient {
	return &Ingredient{
		UserID:      "",
		Name:        dto.Name,
		Description: dto.Category,
		Category:    dto.Category,
		Image:       dto.Image,
		Calories:    dto.Calories,
		Nutrients:   dto.Nutrients,
		CreatedAt:   time.Now(),
		UpdateAt:    time.Now(),
	}
}
