package user

import "time"

type CreateUserDto struct {
	Username      string    `json:"username" validate:"required,min=3,max=20"`
	Email         string    `json:"email" validate:"required,email"`
	Password      string    `json:"password" validate:"required"`
	Weight        float64   `json:"weight" default:"0"` // default:"0" is not working
	Height        float64   `json:"height" default:"0"` // default:"0" is not working
	Age           int       `json:"age" validate:"required,min=1,max=120"`
	Gender        string    `json:"gender" validate:"required"`
	Neck          float64   `json:"neck" default:"0"`          // default:"0" is not working
	Waist         float64   `json:"waist" default:"0"`         // default:"0" is not working
	Hip           float64   `json:"hip" default:"0"`           // default:"0" is not working
	ActivityLevel int       `json:"activityLevel" default:"0"` // default:"0" is not working
	CreatedAt     time.Time `json:"created_at,omitempty" bson:"created_at,omitempty" set:"omitempty" default:"null"`
}

type UpadateUsernamePasswordDto struct {
	Username    string `json:"username"`
	Password    string `json:"password" validate:"required"`
	NewPassword string `json:"newPassword"`
}

type UpdateBodyDto struct {
	Weight        float64 `json:"weight"`
	Height        float64 `json:"height"`
	Age           int     `json:"age"`
	Gender        string  `json:"gender"`
	Neck          float64 `json:"neck"`
	Waist         float64 `json:"waist"`
	Hip           float64 `json:"hip"`
	ActivityLevel int     `json:"activityLevel"`
}
