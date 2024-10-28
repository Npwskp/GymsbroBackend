package ingredient

import (
	"context"
	"errors"

	"github.com/Npwskp/GymsbroBackend/api/v1/function"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

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
	objectID, err := is.CheckOwnershipAndGetObjectId(id, "1") // TODO: Use userId from token
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
	objectID, err := is.CheckOwnershipAndGetObjectId(id, doc.UserID)
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

func (is *IngredientService) CheckOwnershipAndGetObjectId(id string, userid string) (primitive.ObjectID, error) {
	isEntityOwner, err := function.CheckOwnership(is.DB, id, userid, "ingredient", Ingredient{}) // TODO: Use userId from token
	var objectID primitive.ObjectID
	if err != nil {
		return objectID, err
	}

	if !isEntityOwner {
		return objectID, errors.New("update failed: user does not own this entity")
	} else {
		objectID, err = primitive.ObjectIDFromHex(id)
		if err != nil {
			return objectID, err
		}
	}
	return objectID, nil
}
