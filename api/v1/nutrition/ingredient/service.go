package ingredient

import (
	"context"

	"github.com/Npwskp/GymsbroBackend/api/v1/function"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IngredientService struct {
	DB *mongo.Database
}

type IIngredientService interface {
	CreateIngredient(ingredient *CreateIngredientDto, userId string) (*Ingredient, error)
	GetAllIngredients(userId string) ([]*Ingredient, error)
	GetIngredient(id string, userId string) (*Ingredient, error)
	GetIngredientByUser(userId string) ([]*Ingredient, error)
	DeleteIngredient(id string, userId string) error
	UpdateIngredient(doc *UpdateIngredientDto, id string, userId string) (*Ingredient, error)
}

func (is *IngredientService) CreateIngredient(ingredient *CreateIngredientDto, userId string) (*Ingredient, error) {
	ingredientModel := CreateIngredientModel(ingredient)
	result, err := is.DB.Collection("ingredient").InsertOne(context.Background(), ingredientModel)
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

func (is *IngredientService) GetAllIngredients(userId string) ([]*Ingredient, error) {
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

func (is *IngredientService) GetIngredient(id string, userId string) (*Ingredient, error) {
	err := function.CheckOwnership(is.DB, id, userId, "ingredient", &Ingredient{})
	if err != nil {
		return nil, err
	}

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

func (is *IngredientService) GetIngredientByUser(userId string) ([]*Ingredient, error) {
	filter := bson.D{{Key: "userid", Value: userId}}
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

func (is *IngredientService) DeleteIngredient(id string, userId string) error {
	err := function.CheckOwnership(is.DB, id, userId, "ingredient", &Ingredient{})
	if err != nil {
		return err
	}

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

func (is *IngredientService) UpdateIngredient(doc *UpdateIngredientDto, id string, userId string) (*Ingredient, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objectID}}
	update := bson.D{{Key: "$set", Value: doc}}
	_, err = is.DB.Collection("ingredient").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}
	ingredient, err := is.GetIngredient(id, userId)
	if err != nil {
		return nil, err
	}
	return ingredient, nil
}
