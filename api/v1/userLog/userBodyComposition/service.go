package bodyCompositionLog

import (
	"context"
	"time"

	userFitnessPreferenceEnums "github.com/Npwskp/GymsbroBackend/api/v1/user/enums"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BodyCompositionLogService struct {
	DB *mongo.Database
}

type CreateBodyCompositionLogDto struct {
	UserID          string                                         `json:"userid" bson:"userid"`
	Weight          float64                                        `json:"weight" default:"0"`
	BodyComposition userFitnessPreferenceEnums.BodyCompositionInfo `json:"body_composition" bson:"body_composition"`
}

type GetLogsByDateRangeDto struct {
	UserID    string    `json:"userid"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type IBodyCompositionLogService interface {
	CreateBodyCompositionLog(dto *CreateBodyCompositionLogDto) (*UserBodyCompositionLog, error)
	GetLogsByDateRange(dto *GetLogsByDateRangeDto) ([]*UserBodyCompositionLog, error)
}

func (bcs *BodyCompositionLogService) CreateBodyCompositionLog(dto *CreateBodyCompositionLogDto) (*UserBodyCompositionLog, error) {
	collection := bcs.DB.Collection("bodyCompositionLog")

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

	var existingLog UserBodyCompositionLog
	err := collection.FindOne(context.Background(), filter).Decode(&existingLog)

	if err == nil {
		// Update existing log
		update := bson.M{
			"$set": bson.M{
				"weight":           dto.Weight,
				"body_composition": dto.BodyComposition,
				"updated_at":       now,
			},
		}
		_, err = collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return nil, err
		}

		existingLog.Weight = dto.Weight
		existingLog.BodyComposition = dto.BodyComposition
		existingLog.UpdatedAt = now
		return &existingLog, nil
	}

	if err != mongo.ErrNoDocuments {
		return nil, err
	}

	// Create new log
	bodyCompositionLog := &UserBodyCompositionLog{
		UserID:          dto.UserID,
		Weight:          dto.Weight,
		BodyComposition: dto.BodyComposition,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	_, err = collection.InsertOne(context.Background(), bodyCompositionLog)
	if err != nil {
		return nil, err
	}

	return bodyCompositionLog, nil
}

func (bcs *BodyCompositionLogService) GetLogsByDateRange(dto *GetLogsByDateRangeDto) ([]*UserBodyCompositionLog, error) {
	collection := bcs.DB.Collection("bodyCompositionLog")

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
	var logs []*UserBodyCompositionLog
	if err = cursor.All(context.Background(), &logs); err != nil {
		return nil, err
	}

	return logs, nil
}
