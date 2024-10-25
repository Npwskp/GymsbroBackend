package ingredient

import (
	"context"
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/function"
	"github.com/Npwskp/GymsbroBackend/api/v1/nutrition/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Ingredient struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    string             `json:"userid" validate:"required" bson:"userid"`
	Name      string             `json:"name" validate:"required" bson:"name"`
	Image     string             `json:"image" default:"null"`
	Calories  float64            `json:"calories" default:"0"`
	Nutrients []types.Nutrient   `json:"nutrients" default:"null"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty" default:"null"`
}

type IngredientService struct {
	DB *mongo.Database
}

type IIngredientService interface {
	CreateIngredient(ingredient *CreateIngredientDto) (*Ingredient, error)
	GetAllIngredients() ([]*Ingredient, error)
	GetIngredient(id string) (*Ingredient, error)
	GetIngredientByUser(userid string) ([]*Ingredient, error)
	DeleteIngredient(id string) error
	UpdateIngredient(doc *UpdateIngredientDto, id string) (*Ingredient, error)
}

func (is *IngredientService) CreateIngredient(ingredient *CreateIngredientDto) (*Ingredient, error) {
	ingredient.CreatedAt = time.Now()
	result, err := is.DB.Collection("ingredient").InsertOne(context.Background(), ingredient)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: result.InsertedID}}
	createdRecord := is.DB.Collection("ingredient").FindOne(context.Background(), filter)
	createdIngredient := &Ingredient{}
	if err := createdRecord.Decode(createdIngredient); err != nil {
		return nil, err
	}
	return createdIngredient, nil
}

func (is *IngredientService) GetAllIngredients() ([]*Ingredient, error) {
	cursor, err := is.DB.Collection("ingredient").Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	var ingredients []*Ingredient
	if err := cursor.All(context.Background(), &ingredients); err != nil {
		return nil, err
	}
	return ingredients, nil
}

func (is *IngredientService) GetIngredient(id string) (*Ingredient, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: objectID}}
	record := is.DB.Collection("ingredient").FindOne(context.Background(), filter)
	ingredient := &Ingredient{}
	if err := record.Decode(ingredient); err != nil {
		return nil, err
	}
	return ingredient, nil
}

func (is *IngredientService) GetIngredientByUser(userid string) ([]*Ingredient, error) {
	filter := bson.D{{Key: "userid", Value: userid}}
	cursor, err := is.DB.Collection("ingredient").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	var ingredients []*Ingredient
	if err := cursor.All(context.Background(), &ingredients); err != nil {
		return nil, err
	}
	return ingredients, nil
}

func (is *IngredientService) DeleteIngredient(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: objectID}}
	_, err = is.DB.Collection("ingredient").DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (is *IngredientService) UpdateIngredient(doc *UpdateIngredientDto, id string) (*Ingredient, error) {
	objectID, err := function.CheckOwnership(is.DB, id, doc.UserID, "ingredient", Ingredient{}) // TODO: Use userId from token
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objectID}}
	update := bson.D{{Key: "$set", Value: doc}}
	_, err = is.DB.Collection("ingredient").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}
	ingredient, err := is.GetIngredient(id)
	if err != nil {
		return nil, err
	}
	return ingredient, nil
}
