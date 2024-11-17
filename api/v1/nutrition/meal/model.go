package meal

import (
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/nutrition/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Meal struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name            string             `json:"name" validate:"required" bson:"name"`
	Description     string             `json:"description" bson:"description"`
	Category        string             `json:"category" bson:"category"`
	UserID          string             `json:"userid" bson:"userid"`
	Image           string             `json:"image" default:"null"`
	Calories        float64            `json:"calories" default:"0"`
	Nutrients       []types.Nutrient   `json:"nutrients,omitempty" bson:"nutrients"`
	Ingredients     []types.Ingredient `json:"ingredients,omitempty" bson:"ingredients"`
	BrandOwner      string             `json:"brandOwner" bson:"brand_owner"`
	BrandName       string             `json:"brandName" bson:"brand_name"`
	ServingSize     float64            `json:"servingSize" bson:"serving_size"`
	ServingSizeUnit string             `json:"servingSizeUnit" bson:"serving_size_unit"`
	CreatedAt       time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty" default:"null"`
	UpdatedAt       time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty" default:"null"`
}

func CreateMealModel(dto *CreateMealDto) *Meal {
	return &Meal{
		Name:        dto.Name,
		Description: dto.Description,
		Category:    dto.Category,
		UserID:      dto.UserID,
		Image:       dto.Image,
		Calories:    dto.Calories,
		Nutrients:   dto.Nutrients,
		Ingredients: dto.Ingredients,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
