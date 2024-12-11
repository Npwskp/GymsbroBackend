package foodlog

import (
	"context"
	"time"

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
		UpdateAt:  time.Now(),
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
	var foodlogs []*FoodLog
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
			{Key: "update_at", Value: time.Now()},
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
