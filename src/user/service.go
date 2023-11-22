package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID            int     `json:"id,omitempty" bson:"_id,omitempty"`
	Username      string  `json:"username" validate:"required,min=3,max=20"`
	Password      string  `json:"password" validate:"required"`
	Weight        float64 `json:"weight" default:"0"` // default:"0" is not working
	Height        float64 `json:"height" default:"0"` // default:"0" is not working
	Age           int     `json:"age" validate:"required,min=1,max=120"`
	Gender        string  `json:"gender" validate:"required"`
	Neck          float64 `json:"neck" default:"0"`          // default:"0" is not working
	Waist         float64 `json:"waist" default:"0"`         // default:"0" is not working
	Hip           float64 `json:"hip" default:"0"`           // default:"0" is not working
	ActivityLevel int     `json:"activityLevel" default:"0"` // default:"0" is not working
	CreatedAt     string  `json:"createdAt" default:"CURRENT_TIMESTAMP"`
}

type UserService struct {
	DB *mongo.Database
}

type IUserService interface {
	CreateUser(user *CreateUserDto) (*User, error)
}

func (us *UserService) CreateUser(user *CreateUserDto) (*User, error) {
	res := new(User)
	result, err := us.DB.Collection("users").InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: result.InsertedID}}
	createdRecord := us.DB.Collection("users").FindOne(context.Background(), filter)
	createdRecord.Decode(res)
	return res, nil
}
