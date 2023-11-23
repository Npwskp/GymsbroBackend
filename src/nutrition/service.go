package nutrition

import "go.mongodb.org/mongo-driver/mongo"

type Nutrition struct {
	ID        string `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    string `json:"userid" validate:"required" bson:"userid"`
	Carb      string `json:"carb" default:"0"` // default:"0" is not working
	Protein   string `json:"protein" default:"0"`
	Fat       string `json:"fat" default:"0"`
	Calories  string `json:"calories" default:"0"`
	CreatedAt string `json:"created_at" default:"CURRENT" bson:"created_at"`
}

type NutritionService struct {
	DB *mongo.Database
}

type INutritionService interface {
	CreateNutrition(nutrition *CreateNutritionDto) (*Nutrition, error)
	GetAllNutritions() ([]*Nutrition, error)
	GetNutrition(id string) (*Nutrition, error)
	DeleteNutrition(id string) error
	UpdateNutrition(doc *UpdateNutritionDto, id string) (*Nutrition, error)
}
