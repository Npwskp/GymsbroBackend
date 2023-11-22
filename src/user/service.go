package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username      string             `json:"username" validate:"required,min=3,max=20"`
	Password      string             `json:"password" validate:"required"`
	Weight        float64            `json:"weight" default:"0"` // default:"0" is not working
	Height        float64            `json:"height" default:"0"` // default:"0" is not working
	Age           int                `json:"age" validate:"required,min=1,max=120"`
	Gender        string             `json:"gender" validate:"required"`
	Neck          float64            `json:"neck" default:"0"`          // default:"0" is not working
	Waist         float64            `json:"waist" default:"0"`         // default:"0" is not working
	Hip           float64            `json:"hip" default:"0"`           // default:"0" is not working
	ActivityLevel int                `json:"activityLevel" default:"0"` // default:"0" is not working
	CreatedAt     string             `json:"createdAt" default:"CURRENT_TIMESTAMP"`
}

type UserService struct {
	DB *mongo.Database
}

type IUserService interface {
	CreateUser(user *CreateUserDto) (*User, error)
	GetAllUsers() ([]*User, error)
	GetUser(id string) (*User, error)
}

func (us *UserService) CreateUser(user *CreateUserDto) (*User, error) {
	result, err := us.DB.Collection("users").InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: result.InsertedID}}
	createdRecord := us.DB.Collection("users").FindOne(context.Background(), filter)
	createdUser := &User{}
	if err := createdRecord.Decode(createdUser); err != nil {
		return nil, err
	}
	return createdUser, nil
}

func (us *UserService) GetAllUsers() ([]*User, error) {
	cursor, err := us.DB.Collection("users").Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	var users []*User
	if err := cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (us *UserService) GetUser(id string) (*User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: oid}}
	user := &User{}
	if err := us.DB.Collection("users").FindOne(context.Background(), filter).Decode(user); err != nil {
		return nil, err
	}
	return user, nil
}
