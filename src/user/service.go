package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Npwskp/GymsbroBackend/src/function"
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
	ActivityLevel int                `json:"activitylevel" default:"0"` // default:"0" is not working
	CreatedAt     time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type UserService struct {
	DB *mongo.Database
}

type IUserService interface {
	CreateUser(user *CreateUserDto) (*User, error)
	GetAllUsers() ([]*User, error)
	GetUser(id string) (*User, error)
	DeleteUser(id string) error
	UpdateUsernamePassword(doc *UpadateUsernamePasswordDto, id string) (*User, error)
	UpdateBody(doc *UpdateBodyDto, id string) (*User, error)
}

func (us *UserService) CreateUser(user *CreateUserDto) (*User, error) {
	user.CreatedAt = time.Now()
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

func (us *UserService) DeleteUser(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: oid}}
	if _, err := us.DB.Collection("users").DeleteOne(context.Background(), filter); err != nil {
		return err
	}
	return nil
}

func (us *UserService) UpdateUsernamePassword(doc *UpadateUsernamePasswordDto, id string) (*User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: oid}}
	user, err := us.GetUser(id)
	if err != nil {
		return nil, err
	}
	if strings.Compare(user.Password, doc.Password) != 0 {
		return nil, errors.New("password is not correct")
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "username", Value: function.Coalesce(doc.Username, user.Username)},
			{Key: "password", Value: function.Coalesce(doc.NewPassword, user.Password)},
		}},
	}

	// Perform the update
	result, err := us.DB.Collection("users").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	// Check if any document was modified
	if result.ModifiedCount == 0 {
		return nil, errors.New("no user found for the given ID")
	}

	// Retrieve the updated document
	filter = bson.D{{Key: "_id", Value: oid}}
	UpdatedUser := &User{}
	updatedRecord := us.DB.Collection("users").FindOne(context.Background(), filter)
	if err := updatedRecord.Decode(&UpdatedUser); err != nil {
		return nil, err
	}

	return UpdatedUser, nil
}

func (us *UserService) UpdateBody(doc *UpdateBodyDto, id string) (*User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: oid}}
	user, err := us.GetUser(id)
	if err != nil {
		return nil, err
	}
	fmt.Println(user.Weight, "(", doc.Weight, ")")
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "weight", Value: function.Coalesce(doc.Weight, user.Weight)},
			{Key: "height", Value: function.Coalesce(doc.Height, user.Height)},
			{Key: "age", Value: function.Coalesce(doc.Age, user.Age)},
			{Key: "gender", Value: function.Coalesce(doc.Gender, user.Gender)},
			{Key: "neck", Value: function.Coalesce(doc.Neck, user.Neck)},
			{Key: "waist", Value: function.Coalesce(doc.Waist, user.Waist)},
			{Key: "hip", Value: function.Coalesce(doc.Hip, user.Hip)},
			{Key: "activitylevel", Value: function.Coalesce(doc.ActivityLevel, user.ActivityLevel)},
		}},
	}
	result, err := us.DB.Collection("users").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	// Check if any document was modified
	if result.ModifiedCount == 0 {
		return nil, errors.New("no user found for the given ID")
	}

	// Retrieve the updated document
	filter = bson.D{{Key: "_id", Value: oid}}
	UpdatedUser := &User{}
	updatedRecord := us.DB.Collection("users").FindOne(context.Background(), filter)
	if err := updatedRecord.Decode(&UpdatedUser); err != nil {
		return nil, err
	}

	return UpdatedUser, nil
}
