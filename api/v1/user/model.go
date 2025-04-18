package user

import (
	"time"

	authEnums "github.com/Npwskp/GymsbroBackend/api/v1/auth/enums"
	userFitnessPreferenceEnums "github.com/Npwskp/GymsbroBackend/api/v1/user/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID                             `json:"id,omitempty" bson:"_id,omitempty"`
	Username        string                                         `json:"username" validate:"required,min=3,max=20"`
	Email           string                                         `json:"email" validate:"required,email"`
	Password        string                                         `json:"password" validate:"required"`
	Weight          float64                                        `json:"weight" default:"0"`
	Height          float64                                        `json:"height" default:"0"`
	Age             int                                            `json:"age" validate:"required,min=1,max=120"`
	Gender          authEnums.GenderType                           `json:"gender" validate:"required"`
	NutritionInfo   userFitnessPreferenceEnums.NutritionInfo       `json:"nutrition_info" bson:"nutrition_info"`
	BodyComposition userFitnessPreferenceEnums.BodyCompositionInfo `json:"body_composition" bson:"body_composition"`
	Macronutrients  userFitnessPreferenceEnums.Macronutrients      `json:"macronutrients" bson:"macronutrients"`
	OAuthProvider   string                                         `json:"oauth_provider,omitempty" bson:"oauth_provider,omitempty"`
	OAuthID         string                                         `json:"oauth_id,omitempty" bson:"oauth_id,omitempty"`
	Picture         string                                         `json:"picture,omitempty" bson:"picture,omitempty"`
	CreatedAt       time.Time                                      `json:"created_at,omitempty" bson:"created_at,omitempty" default:"null"`
	UpdatedAt       time.Time                                      `json:"updated_at,omitempty" bson:"updated_at,omitempty" default:"null"`
	IsFirstLogin    bool                                           `json:"is_first_login" bson:"is_first_login" default:"true"`
}

func CreateUserModel(user *CreateUserDto) *User {
	model := &User{
		Username:     user.Username,
		Email:        user.Email,
		Password:     user.Password,
		Age:          user.Age,
		Gender:       user.Gender,
		IsFirstLogin: true,
		CreatedAt:    time.Now(),
	}

	model.NutritionInfo.ActivityLevel = user.ActivityLevel
	model.NutritionInfo.Goal = user.Goal

	// Calculate initial BMR if possible
	if model.Weight > 0 && model.Height > 0 && model.Age > 0 &&
		(model.Gender == "male" || model.Gender == "female") {
		model.NutritionInfo.BMR = userFitnessPreferenceEnums.CalculateBMR(
			model.Weight,
			model.Height,
			model.Age,
			model.Gender,
		)
	}

	// Calculate BMI if possible
	if model.Weight > 0 && model.Height > 0 {
		model.BodyComposition.BMI = userFitnessPreferenceEnums.CalculateBMI(
			model.Weight,
			model.Height,
		)
	}

	return model
}
