package meal

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MealService struct {
	DB *mongo.Database
}

type IMealService interface {
	CreateMeal(meal *CreateMealDto) (*Meal, error)
	GetAllMeals() ([]*Meal, error)
	GetMeal(id string) (*Meal, error)
	GetMealByUser(userid string) ([]*Meal, error)
	DeleteMeal(id string) error
	UpdateMeal(doc *UpdateMealDto, id string) (*Meal, error)
	SearchFilteredMeals(filters SearchFilters) ([]*Meal, error)
}

func (ns *MealService) CreateMeal(meal *CreateMealDto) (*Meal, error) {
	mealModel := CreateMealModel(meal)

	result, err := ns.DB.Collection("meal").InsertOne(context.Background(), mealModel)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: result.InsertedID}}
	createdRecord := ns.DB.Collection("meal").FindOne(context.Background(), filter)
	createdMeal := &Meal{}
	if err := createdRecord.Decode(createdMeal); err != nil {
		return nil, err
	}
	return createdMeal, nil
}

func (ns *MealService) GetAllMeals() ([]*Meal, error) {
	cursor, err := ns.DB.Collection("meal").Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	var Meals []*Meal
	if err := cursor.All(context.Background(), &Meals); err != nil {
		return nil, err
	}
	return Meals, nil
}

func (ns *MealService) GetMeal(id string) (*Meal, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: oid}}
	meal := &Meal{}
	if err := ns.DB.Collection("meal").FindOne(context.Background(), filter).Decode(meal); err != nil {
		return nil, err
	}
	return meal, nil
}

func (ns *MealService) GetMealByUser(userid string) ([]*Meal, error) {
	filter := bson.M{"userid": bson.M{"$in": []string{userid}}}
	cursor, err := ns.DB.Collection("meal").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	var Meals []*Meal
	if err := cursor.All(context.Background(), &Meals); err != nil {
		return nil, err
	}
	return Meals, nil
}

func (ns *MealService) DeleteMeal(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: oid}}
	if _, err := ns.DB.Collection("meal").DeleteOne(context.Background(), filter); err != nil {
		return err
	}
	return nil
}

func (ns *MealService) UpdateMeal(doc *UpdateMealDto, id string) (*Meal, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: oid}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "carb", Value: doc.Carb},
			{Key: "protein", Value: doc.Protein},
			{Key: "fat", Value: doc.Fat},
			{Key: "calories", Value: doc.Calories},
			{Key: "updated_at", Value: time.Now()},
		}},
	}
	if _, err := ns.DB.Collection("meal").UpdateOne(context.Background(), filter, update); err != nil {
		return nil, err
	}
	return ns.GetMeal(id)
}

func (ns *MealService) SearchFilteredMeals(filters SearchFilters) ([]*Meal, error) {
	filterQuery := bson.D{}
	andConditions := []bson.D{}

	// Add name search if query is provided
	if filters.Query != "" {
		andConditions = append(andConditions, bson.D{
			{Key: "name", Value: bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: filters.Query, Options: "i"}}}},
		})
	}

	// Add user filter
	if filters.UserID != "" {
		andConditions = append(andConditions, bson.D{
			{Key: "$or", Value: []bson.D{
				{{Key: "userid", Value: ""}},
				{{Key: "userid", Value: filters.UserID}},
				{{Key: "userid", Value: primitive.Null{}}},
			}},
		})
	}

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
	cursor, err := ns.DB.Collection("meal").Find(context.Background(), filterQuery, opts)
	if err != nil {
		return nil, fmt.Errorf("error searching meals: %w", err)
	}
	defer cursor.Close(context.Background())

	var meals []*Meal
	if err := cursor.All(context.Background(), &meals); err != nil {
		return nil, fmt.Errorf("error decoding meals: %w", err)
	}

	return meals, nil
}
