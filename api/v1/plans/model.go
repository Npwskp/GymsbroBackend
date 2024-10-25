package plans

import "go.mongodb.org/mongo-driver/bson/primitive"

type Plan struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID     string             `json:"userid" validate:"required" bson:"userid"`
	TypeOfPlan string             `json:"typeofplan" validate:"required"`
	DayOfWeek  string             `json:"dayofweek" validate:"required"`
	Exercise   []string           `json:"exercise"`
}
