package macronutrientLog

import (
	"context"
	"time"

	userFitnessPreferenceEnums "github.com/Npwskp/GymsbroBackend/api/v1/user/enums"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MacronutrientLogService struct {
	DB *mongo.Database
}

type CreateMacronutrientLogDto struct {
	UserID         string                                    `json:"userid" bson:"userid"`
	Macronutrients userFitnessPreferenceEnums.Macronutrients `json:"macronutrients" bson:"macronutrients"`
}

type GetLogsByDateRangeDto struct {
	UserID    string    `json:"userid"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type IMacronutrientLogService interface {
	CreateMacronutrientLog(dto *CreateMacronutrientLogDto) (*UserMacronutrientLog, error)
	GetLogsByDateRange(dto *GetLogsByDateRangeDto) ([]*UserMacronutrientLog, error)
}

func (mls *MacronutrientLogService) CreateMacronutrientLog(dto *CreateMacronutrientLogDto) (*UserMacronutrientLog, error) {
	collection := mls.DB.Collection("macronutrientLog")

	// Get today's start and end time
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Check if there's already a log for today
	filter := bson.M{
		"userid": dto.UserID,
		"created_at": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}

	var existingLog UserMacronutrientLog
	err := collection.FindOne(context.Background(), filter).Decode(&existingLog)

	if err == nil {
		// Update existing log
		update := bson.M{
			"$set": bson.M{
				"macronutrients": dto.Macronutrients,
				"updated_at":     now,
			},
		}
		_, err = collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return nil, err
		}

		existingLog.Macronutrients = dto.Macronutrients
		existingLog.UpdatedAt = now
		return &existingLog, nil
	}

	if err != mongo.ErrNoDocuments {
		return nil, err
	}

	// Create new log
	macronutrientLog := &UserMacronutrientLog{
		UserID:         dto.UserID,
		Macronutrients: dto.Macronutrients,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	_, err = collection.InsertOne(context.Background(), macronutrientLog)
	if err != nil {
		return nil, err
	}

	return macronutrientLog, nil
}

func (mls *MacronutrientLogService) GetLogsByDateRange(dto *GetLogsByDateRangeDto) ([]*UserMacronutrientLog, error) {
	collection := mls.DB.Collection("macronutrientLog")

	// Create filter for date range and userID
	filter := bson.M{
		"userid": dto.UserID,
		"created_at": bson.M{
			"$gte": dto.StartDate,
			"$lte": dto.EndDate,
		},
	}

	// Set up options for sorting by date (newest first)
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	// Execute the query
	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Decode results
	var logs []*UserMacronutrientLog
	if err = cursor.All(context.Background(), &logs); err != nil {
		return nil, err
	}

	return logs, nil
}
