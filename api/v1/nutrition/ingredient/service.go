package ingredient

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/function"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	SearchFilteredIngredients(filters SearchFilters) ([]*Ingredient, error)
}

func (is *IngredientService) CreateIngredient(ingredient *CreateIngredientDto, userId string) (*Ingredient, error) {
	ingredientModel := CreateIngredientModel(ingredient)
	ingredientModel.UserID = userId

	result, err := is.DB.Collection("ingredient").InsertOne(context.Background(), ingredientModel)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: result.InsertedID}}
	createdIngredient := &Ingredient{}
	if err := is.DB.Collection("ingredient").FindOne(context.Background(), filter).Decode(createdIngredient); err != nil {
		return nil, err
	}

	return createdIngredient, nil
}

func (is *IngredientService) GetAllIngredients(userId string) ([]*Ingredient, error) {
	filter := bson.D{
		{Key: "$or", Value: []bson.D{
			{{Key: "userid", Value: ""}},
			{{Key: "userid", Value: userId}},
			{{Key: "userid", Value: primitive.Null{}}},
		}},
	}
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
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: doc.Name},
			{Key: "category", Value: doc.Category},
			{Key: "calories", Value: doc.Calories},
			{Key: "nutrients", Value: doc.Nutrients},
			{Key: "updated_at", Value: time.Now()},
		}},
	}
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

func (is *IngredientService) SearchFilteredIngredients(filters SearchFilters) ([]*Ingredient, error) {
	filterQuery := bson.D{}
	andConditions := []bson.D{}

	// Add name search if query is provided
	if filters.Query != "" {
		andConditions = append(andConditions, bson.D{
			{Key: "name", Value: bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: filters.Query, Options: "i"}}}},
		})
	}

	// Add user filter
	andConditions = append(andConditions, bson.D{
		{Key: "$or", Value: []bson.D{
			{{Key: "userid", Value: ""}},
			{{Key: "userid", Value: filters.UserID}},
			{{Key: "userid", Value: primitive.Null{}}},
		}},
	})

	// Add category filter
	if filters.Category != "" {
		andConditions = append(andConditions, bson.D{
			{Key: "category", Value: filters.Category},
		})
	}

	// Add calories range filter
	if filters.MinCalories > 0 || filters.MaxCalories > 0 {
		caloriesFilter := bson.D{}
		if filters.MinCalories > 0 {
			caloriesFilter = append(caloriesFilter, bson.E{Key: "$gte", Value: filters.MinCalories})
		}
		if filters.MaxCalories > 0 {
			caloriesFilter = append(caloriesFilter, bson.E{Key: "$lte", Value: filters.MaxCalories})
		}
		andConditions = append(andConditions, bson.D{
			{Key: "calories", Value: caloriesFilter},
		})
	}

	// Add nutrients filter
	if filters.Nutrients != "" {
		nutrients := strings.Split(filters.Nutrients, ",")
		nutrientFilters := bson.A{}
		for _, nutrient := range nutrients {
			nutrient = strings.TrimSpace(nutrient)
			if nutrient != "" {
				nutrientFilters = append(nutrientFilters, bson.M{
					"nutrients.name": bson.M{
						"$regex":   nutrient,
						"$options": "i",
					},
				})
			}
		}
		if len(nutrientFilters) > 0 {
			andConditions = append(andConditions, bson.D{
				{Key: "$or", Value: nutrientFilters},
			})
		}
	}

	// Combine all conditions
	if len(andConditions) > 0 {
		filterQuery = bson.D{{Key: "$and", Value: andConditions}}
	}

	// Execute query with limit
	opts := options.Find().SetLimit(20)
	cursor, err := is.DB.Collection("ingredient").Find(context.Background(), filterQuery, opts)
	if err != nil {
		return nil, fmt.Errorf("error searching ingredients: %w", err)
	}
	defer cursor.Close(context.Background())

	var ingredients []*Ingredient
	if err := cursor.All(context.Background(), &ingredients); err != nil {
		return nil, fmt.Errorf("error decoding ingredients: %w", err)
	}

	return ingredients, nil
}
