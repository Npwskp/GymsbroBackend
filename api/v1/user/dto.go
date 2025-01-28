package user

import (
	authEnums "github.com/Npwskp/GymsbroBackend/api/v1/auth/enums"
	userFitnessPreferenceEnums "github.com/Npwskp/GymsbroBackend/api/v1/user/enums"
)

type CreateUserDto struct {
	Username      string                                       `json:"username" validate:"required,min=3,max=20"`
	Email         string                                       `json:"email" validate:"required,email"`
	Password      string                                       `json:"password" validate:"required"`
	Weight        float64                                      `json:"weight" default:"0"`
	Height        float64                                      `json:"height" default:"0"`
	Age           int                                          `json:"age" validate:"required,min=1,max=120"`
	Gender        authEnums.GenderType                         `json:"gender" validate:"required"`
	ActivityLevel userFitnessPreferenceEnums.ActivityLevelType `json:"activityLevel" default:"sedentary"`
	Goal          userFitnessPreferenceEnums.GoalType          `json:"goal" default:"maintain"`
	OAuthProvider string                                       `json:"oauth_provider,omitempty" default:""`
	OAuthID       string                                       `json:"oauth_id,omitempty" default:""`
	Picture       string                                       `json:"picture,omitempty" default:""`
}

type UpdateUsernamePasswordDto struct {
	Username    string `json:"username"`
	Password    string `json:"password" validate:"required"`
	NewPassword string `json:"newPassword"`
}

type UpdateBodyDto struct {
	Weight          float64                                        `json:"weight"`
	Height          float64                                        `json:"height"`
	Age             int                                            `json:"age"`
	Gender          authEnums.GenderType                           `json:"gender"`
	NutritionInfo   userFitnessPreferenceEnums.NutritionInfo       `json:"nutrition_info"`
	BodyComposition userFitnessPreferenceEnums.BodyCompositionInfo `json:"body_composition"`
	Macronutrients  userFitnessPreferenceEnums.Macronutrients      `json:"macronutrients"`
}
