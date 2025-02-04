package foodlog

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/nutrition/meal"
	"github.com/Npwskp/GymsbroBackend/api/v1/nutrition/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FoodLogService struct {
	DB *mongo.Database
}

type IFoodLogService interface {
	AddMealToFoodLog(foodlog *AddMealToFoodLogDto, userid string) (*FoodLog, error)
	GetFoodLog(id string, userid string) (*FoodLog, error)
	GetFoodLogByUser(userid string) ([]*FoodLog, error)
	GetFoodLogByUserDate(userid string, date string) (*FoodLog, error)
	DeleteFoodLog(id string, userid string) error
	UpdateFoodLog(doc *UpdateFoodLogDto, id string, userid string) (*FoodLog, error)
	CalculateDailyNutrients(date string, userid string) (*DailyNutrientResponse, error)
}

func (fs *FoodLogService) AddMealToFoodLog(foodlog *AddMealToFoodLogDto, userid string) (*FoodLog, error) {
	// Check if a food log already exists for this user and date
	existingLog, err := fs.GetFoodLogByUserDate(userid, foodlog.Date)
	if err == nil && existingLog != nil {
		// Update existing log by adding the new meal
		existingLog.Meals = append(existingLog.Meals, foodlog.Meals...)

		// Create update document
		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "meals", Value: existingLog.Meals},
				{Key: "update_at", Value: time.Now()},
			}},
		}

		// Update the existing document
		filter := bson.D{{Key: "_id", Value: existingLog.ID}}
		_, err = fs.DB.Collection("foodlog").UpdateOne(context.Background(), filter, update)
		if err != nil {
			return nil, err
		}

		return existingLog, nil
	}

	// Create new food log if none exists
	newFoodLog := &FoodLog{
		UserID:    userid,
		Date:      foodlog.Date,
		Meals:     foodlog.Meals,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert the new food log
	result, err := fs.DB.Collection("foodlog").InsertOne(context.Background(), newFoodLog)
	if err != nil {
		return nil, err
	}

	// Fetch and return the created document
	filter := bson.D{{Key: "_id", Value: result.InsertedID}}
	createdFoodLog := &FoodLog{}
	if err := fs.DB.Collection("foodlog").FindOne(context.Background(), filter).Decode(createdFoodLog); err != nil {
		return nil, err
	}

	return createdFoodLog, nil
}

func (fs *FoodLogService) GetFoodLog(id string, userid string) (*FoodLog, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: objectID}}
	foodlog := &FoodLog{}
	err = fs.DB.Collection("foodlog").FindOne(context.Background(), filter).Decode(foodlog)
	if err != nil {
		return nil, err
	}
	return foodlog, nil
}

func (fs *FoodLogService) GetFoodLogByUser(userid string) ([]*FoodLog, error) {
	filter := bson.D{{Key: "userid", Value: userid}}
	cursor, err := fs.DB.Collection("foodlog").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	foodlogs := make([]*FoodLog, 0)
	if err := cursor.All(context.Background(), &foodlogs); err != nil {
		return nil, err
	}
	return foodlogs, nil
}

func (fs *FoodLogService) GetFoodLogByUserDate(userid string, date string) (*FoodLog, error) {
	// Create a filter for user and date
	filter := bson.D{
		{Key: "userid", Value: userid},
		{Key: "date", Value: date},
	}

	// Find one document
	foodlog := &FoodLog{}
	err := fs.DB.Collection("foodlog").FindOne(context.Background(), filter).Decode(foodlog)
	if err != nil {
		return nil, err
	}

	return foodlog, nil
}

func (fs *FoodLogService) DeleteFoodLog(id string, userid string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: objectID}}
	_, err = fs.DB.Collection("foodlog").DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (fs *FoodLogService) UpdateFoodLog(doc *UpdateFoodLogDto, id string, userid string) (*FoodLog, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: objectID}}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "date", Value: doc.Date},
			{Key: "meals", Value: doc.Meals},
			{Key: "updated_at", Value: time.Now()},
		}},
	}

	_, err = fs.DB.Collection("foodlog").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	// Fetch and return updated document
	foodlog := &FoodLog{}
	err = fs.DB.Collection("foodlog").FindOne(context.Background(), filter).Decode(foodlog)
	if err != nil {
		return nil, err
	}
	return foodlog, nil
}

func (fs *FoodLogService) CalculateDailyNutrients(date string, userid string) (*DailyNutrientResponse, error) {
	// Get food log for the specified date
	foodLog, err := fs.GetFoodLogByUserDate(userid, date)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Return zero values if no food log exists
			return &DailyNutrientResponse{
				Date:      date,
				Calories:  0,
				Nutrients: []types.Nutrient{},
			}, nil
		}
		return nil, fmt.Errorf("error fetching food log: %w", err)
	}

	// Initialize meal service
	mealService := &meal.MealService{DB: fs.DB}

	totalCalories := 0.0
	totalNutrients := make(map[string]types.Nutrient)

	// Calculate nutrients for each meal in the food log
	for _, mealID := range foodLog.Meals {
		meal, err := mealService.GetMeal(mealID, userid)
		if err != nil {
			return nil, fmt.Errorf("error fetching meal %s: %w", mealID, err)
		}

		// Add calories
		totalCalories += meal.Calories

		// Combine nutrients
		for _, nutrient := range meal.Nutrients {
			if existing, ok := totalNutrients[nutrient.Name]; ok {
				// Update existing nutrient
				roundedAmount, _ := strconv.ParseFloat(fmt.Sprintf("%.5f", existing.Amount+nutrient.Amount), 64)
				existing.Amount = roundedAmount
				totalNutrients[nutrient.Name] = existing
			} else {
				// Add new nutrient
				totalNutrients[nutrient.Name] = types.Nutrient{
					Name:   nutrient.Name,
					Amount: nutrient.Amount,
					Unit:   nutrient.Unit,
				}
			}
		}
	}

	// Convert nutrients map to slice
	nutrientValues := make([]types.Nutrient, 0, len(totalNutrients))
	for _, nutrient := range totalNutrients {
		nutrientValues = append(nutrientValues, nutrient)
	}

	return &DailyNutrientResponse{
		Date:      date,
		Calories:  totalCalories,
		Nutrients: nutrientValues,
	}, nil
}
