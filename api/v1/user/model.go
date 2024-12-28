package user

import (
	"time"

	userFitnessPreferenceEnums "github.com/Npwskp/GymsbroBackend/api/v1/user/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID                           `json:"id,omitempty" bson:"_id,omitempty"`
	Username       string                                       `json:"username" validate:"required,min=3,max=20"`
	Email          string                                       `json:"email" validate:"required,email"`
	Password       string                                       `json:"password" validate:"required"`
	Weight         float64                                      `json:"weight" default:"0"`
	Height         float64                                      `json:"height" default:"0"`
	Age            int                                          `json:"age" validate:"required,min=1,max=120"`
	Gender         string                                       `json:"gender" validate:"required"`
	Neck           float64                                      `json:"neck" default:"0"`
	Waist          float64                                      `json:"waist" default:"0"`
	Hip            float64                                      `json:"hip" default:"0"`
	ActivityLevel  userFitnessPreferenceEnums.ActivityLevelType `json:"activitylevel" default:"0"`
	Goal           userFitnessPreferenceEnums.GoalType          `json:"goal" default:"maintain"`
	BMR            float64                                      `json:"bmr" default:"0"`
	Macronutrients userFitnessPreferenceEnums.Macronutrients    `json:"macronutrients" default:"0"`
	OAuthProvider  string                                       `json:"oauth_provider,omitempty" bson:"oauth_provider,omitempty"`
	OAuthID        string                                       `json:"oauth_id,omitempty" bson:"oauth_id,omitempty"`
	Picture        string                                       `json:"picture,omitempty" bson:"picture,omitempty"`
	CreatedAt      time.Time                                    `json:"created_at,omitempty" bson:"created_at,omitempty" default:"null"`
	UpdatedAt      time.Time                                    `json:"updated_at,omitempty" bson:"updated_at,omitempty" default:"null"`
	IsFirstLogin   bool                                         `json:"is_first_login" bson:"is_first_login" default:"true"`
}

func CreateUserModel(user *CreateUserDto) *User {
	return &User{
		Username:      user.Username,
		Email:         user.Email,
		Password:      user.Password,
		Weight:        user.Weight,
		Height:        user.Height,
		Age:           user.Age,
		Gender:        user.Gender,
		Neck:          user.Neck,
		Waist:         user.Waist,
		Hip:           user.Hip,
		ActivityLevel: user.ActivityLevel,
		Goal:          user.Goal,
		IsFirstLogin:  true,
		CreatedAt:     time.Now(),
	}
}
