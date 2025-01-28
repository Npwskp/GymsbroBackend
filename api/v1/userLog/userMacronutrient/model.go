package macronutrientLog

import (
	"time"

	userFitnessPreferenceEnums "github.com/Npwskp/GymsbroBackend/api/v1/user/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserMacronutrientLog struct {
	ID             primitive.ObjectID                        `json:"id,omitempty" bson:"_id,omitempty"`
	UserID         string                                    `json:"userid" bson:"userid"`
	Macronutrients userFitnessPreferenceEnums.Macronutrients `json:"macronutrients" bson:"macronutrients"`
	CreatedAt      time.Time                                 `json:"created_at,omitempty" bson:"created_at,omitempty" default:"null"`
	UpdatedAt      time.Time                                 `json:"updated_at,omitempty" bson:"updated_at,omitempty" default:"null"`
}
