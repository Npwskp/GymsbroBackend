package exerciseLog

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ExerciseLogService struct {
	DB *mongo.Database
}

type IExerciseLogService interface {
	CreateLog(log *CreateExerciseLogDto, userId string) (*ExerciseLog, error)
	GetLogsByUser(userId string) ([]*ExerciseLog, error)
	GetLogsByExercise(exerciseId string, userId string) ([]*ExerciseLog, error)
	GetLogsByDateRange(userId string, startDate, endDate time.Time) ([]*ExerciseLog, error)
	UpdateLog(id string, log *UpdateExerciseLogDto, userId string) (*ExerciseLog, error)
	DeleteLog(id string, userId string) error
}

func (s *ExerciseLogService) CreateLog(dto *CreateExerciseLogDto, userId string) (*ExerciseLog, error) {
	// Calculate total volume and completed sets
	var totalVolume float64
	completedSets := len(dto.Sets)

	for _, set := range dto.Sets {
		totalVolume += float64(set.Reps) * set.Weight
	}

	log := &ExerciseLog{
		UserID:        userId,
		ExerciseID:    dto.ExerciseID,
		CompletedSets: completedSets,
		TotalVolume:   totalVolume,
		Notes:         dto.Notes,
		Duration:      0, // This will be updated when the session ends
		DateTime:      time.Now(),
		Sets:          dto.Sets,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	result, err := s.DB.Collection("exerciseLogs").InsertOne(context.Background(), log)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: result.InsertedID}}
	createdLog := &ExerciseLog{}
	if err := s.DB.Collection("exerciseLogs").FindOne(context.Background(), filter).Decode(createdLog); err != nil {
		return nil, err
	}

	return createdLog, nil
}

func (s *ExerciseLogService) GetLogsByUser(userId string) ([]*ExerciseLog, error) {
	filter := bson.D{{Key: "userid", Value: userId}}
	opts := options.Find().SetSort(bson.D{{Key: "datetime", Value: -1}})

	cursor, err := s.DB.Collection("exerciseLogs").Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}

	var logs []*ExerciseLog
	if err := cursor.All(context.Background(), &logs); err != nil {
		return nil, err
	}

	return logs, nil
}

func (s *ExerciseLogService) GetLogsByExercise(exerciseId string, userId string) ([]*ExerciseLog, error) {
	filter := bson.D{
		{Key: "userid", Value: userId},
		{Key: "exerciseId", Value: exerciseId},
	}
	opts := options.Find().SetSort(bson.D{{Key: "datetime", Value: -1}})

	cursor, err := s.DB.Collection("exerciseLogs").Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}

	var logs []*ExerciseLog
	if err := cursor.All(context.Background(), &logs); err != nil {
		return nil, err
	}

	return logs, nil
}

func (s *ExerciseLogService) GetLogsByDateRange(userId string, startDate, endDate time.Time) ([]*ExerciseLog, error) {
	filter := bson.D{
		{Key: "userid", Value: userId},
		{Key: "datetime", Value: bson.D{
			{Key: "$gte", Value: startDate},
			{Key: "$lte", Value: endDate},
		}},
	}
	opts := options.Find().SetSort(bson.D{{Key: "datetime", Value: -1}})

	cursor, err := s.DB.Collection("exerciseLogs").Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}

	var logs []*ExerciseLog
	if err := cursor.All(context.Background(), &logs); err != nil {
		return nil, err
	}

	return logs, nil
}

func (s *ExerciseLogService) UpdateLog(id string, dto *UpdateExerciseLogDto, userId string) (*ExerciseLog, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "userid", Value: userId},
	}

	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "sets", Value: dto.Sets},
		{Key: "notes", Value: dto.Notes},
	}}}

	after := options.After
	opts := options.FindOneAndUpdate().SetReturnDocument(after)

	result := &ExerciseLog{}
	err = s.DB.Collection("exerciseLogs").FindOneAndUpdate(context.Background(), filter, update, opts).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *ExerciseLogService) DeleteLog(id string, userId string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "userid", Value: userId},
	}

	result, err := s.DB.Collection("exerciseLogs").DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
