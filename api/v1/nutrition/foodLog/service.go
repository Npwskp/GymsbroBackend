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
	CreateFoodLog(foodlog *CreateFoodLogDto) (*FoodLog, error)
	GetFoodLog(id string) (*FoodLog, error)
	GetFoodLogByUser(userid string) ([]*FoodLog, error)
	GetFoodLogByUserDate(userid string, date string) (*FoodLog, error)
	DeleteFoodLog(id string) error
	UpdateFoodLog(doc *UpdateFoodLogDto, id string) (*FoodLog, error)
}

func (fs *FoodLogService) CreateFoodLog(foodlog *CreateFoodLogDto) (*FoodLog, error) {
	foodlog.CreatedAt = time.Now()
	result, err := fs.DB.Collection("foodlog").InsertOne(context.Background(), foodlog)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: result.InsertedID}}
	createdRecord := fs.DB.Collection("foodlog").FindOne(context.Background(), filter)
	createdFoodLog := &FoodLog{}
	if err := createdRecord.Decode(createdFoodLog); err != nil {
		return nil, err
	}
	return createdFoodLog, nil
}

func (fs *FoodLogService) GetFoodLog(id string) (*FoodLog, error) {
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
	cursor, err := fs.DB.Collection("foodlog").Find(context.Background(), bson.D{{Key: "userid", Value: userid}})
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
	filter := bson.D{{Key: "userid", Value: userid}, {Key: "date", Value: date}}
	foodlog := &FoodLog{}
	err := fs.DB.Collection("foodlog").FindOne(context.Background(), filter).Decode(foodlog)
	if err != nil {
		return nil, err
	}
	return foodlog, nil
}

func (fs *FoodLogService) DeleteFoodLog(id string) error {
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

func (fs *FoodLogService) UpdateFoodLog(doc *UpdateFoodLogDto, id string) (*FoodLog, error) {
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
	foodlog.UserID = doc.UserID
	foodlog.Date = doc.Date
	foodlog.Meals = doc.Meals
	_, err = fs.DB.Collection("foodlog").UpdateOne(context.Background(), filter, bson.D{{Key: "$set", Value: foodlog}})
	if err != nil {
		return nil, err
	}
	return foodlog, nil
}
