package bodyCompositionLog

import (
	"time"

	userFitnessPreferenceEnums "github.com/Npwskp/GymsbroBackend/api/v1/user/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserBodyCompositionLog struct {
	ID              primitive.ObjectID                             `json:"id,omitempty" bson:"_id,omitempty"`
	UserID          string                                         `json:"userid" bson:"userid"`
	Weight          float64                                        `json:"weight" default:"0"`
	BodyComposition userFitnessPreferenceEnums.BodyCompositionInfo `json:"body_composition" bson:"body_composition"`
	CreatedAt       time.Time                                      `json:"created_at,omitempty" bson:"created_at,omitempty" default:"null"`
	UpdatedAt       time.Time                                      `json:"updated_at,omitempty" bson:"updated_at,omitempty" default:"null"`
}
