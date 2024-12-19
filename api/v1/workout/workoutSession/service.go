package workoutSession

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WorkoutSessionService struct {
	DB *mongo.Database
}

type IWorkoutSessionService interface {
	StartSession(dto *CreateWorkoutSessionDto, userId string) (*WorkoutSession, error)
	EndSession(id string, userId string) (*WorkoutSession, error)
	UpdateSession(id string, dto *UpdateWorkoutSessionDto, userId string) (*WorkoutSession, error)
	LogExercise(sessionId string, exerciseId string, dto *CompleteExerciseDto, userId string) (*WorkoutSession, error)
	GetSession(id string, userId string) (*WorkoutSession, error)
	GetUserSessions(userId string) ([]*WorkoutSession, error)
	DeleteSession(id string, userId string) error
}

func (s *WorkoutSessionService) StartSession(dto *CreateWorkoutSessionDto, userId string) (*WorkoutSession, error) {
	session := &WorkoutSession{
		UserID:      userId,
		WorkoutID:   dto.WorkoutID,
		StartTime:   time.Now(),
		Status:      StatusInProgress,
		TotalVolume: 0,
		Duration:    0,
		Exercises:   []ExerciseEntry{},
		Notes:       dto.Notes,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result, err := s.DB.Collection("workoutSessions").InsertOne(context.Background(), session)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: result.InsertedID}}
	createdSession := &WorkoutSession{}
	if err := s.DB.Collection("workoutSessions").FindOne(context.Background(), filter).Decode(createdSession); err != nil {
		return nil, err
	}

	return createdSession, nil
}

func (s *WorkoutSessionService) EndSession(id string, userId string) (*WorkoutSession, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "userId", Value: userId},
	}

	session := &WorkoutSession{}
	if err := s.DB.Collection("workoutSessions").FindOne(context.Background(), filter).Decode(session); err != nil {
		return nil, err
	}

	if session.Status != StatusInProgress {
		return nil, errors.New("session is not in progress")
	}

	endTime := time.Now()
	duration := int(endTime.Sub(session.StartTime).Seconds())

	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "endTime", Value: endTime},
		{Key: "status", Value: StatusCompleted},
		{Key: "duration", Value: duration},
		{Key: "updatedAt", Value: time.Now()},
	}}}

	after := options.After
	opts := options.FindOneAndUpdate().SetReturnDocument(after)

	result := &WorkoutSession{}
	err = s.DB.Collection("workoutSessions").FindOneAndUpdate(context.Background(), filter, update, opts).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *WorkoutSessionService) UpdateSession(id string, dto *UpdateWorkoutSessionDto, userId string) (*WorkoutSession, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "userId", Value: userId},
	}

	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "status", Value: dto.Status},
		{Key: "exercises", Value: dto.Exercises},
		{Key: "notes", Value: dto.Notes},
		{Key: "updatedAt", Value: time.Now()},
	}}}

	after := options.After
	opts := options.FindOneAndUpdate().SetReturnDocument(after)

	result := &WorkoutSession{}
	err = s.DB.Collection("workoutSessions").FindOneAndUpdate(context.Background(), filter, update, opts).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *WorkoutSessionService) LogExercise(sessionId string, exerciseId string, dto *CompleteExerciseDto, userId string) (*WorkoutSession, error) {
	oid, err := primitive.ObjectIDFromHex(sessionId)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "userId", Value: userId},
		{Key: "status", Value: StatusInProgress},
	}

	session := &WorkoutSession{}
	if err := s.DB.Collection("workoutSessions").FindOne(context.Background(), filter).Decode(session); err != nil {
		return nil, err
	}

	// Update exercise entry or add new one
	exerciseEntry := ExerciseEntry{
		ExerciseID:    exerciseId,
		ExerciseLogID: dto.ExerciseLogID,
		TotalVolume:   dto.TotalVolume,
		StartTime:     time.Now(),
		EndTime:       time.Now(),
	}

	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "exercises", Value: exerciseEntry},
		}},
		{Key: "$inc", Value: bson.D{
			{Key: "totalVolume", Value: dto.TotalVolume},
		}},
		{Key: "$set", Value: bson.D{
			{Key: "updatedAt", Value: time.Now()},
		}},
	}

	after := options.After
	opts := options.FindOneAndUpdate().SetReturnDocument(after)

	result := &WorkoutSession{}
	err = s.DB.Collection("workoutSessions").FindOneAndUpdate(context.Background(), filter, update, opts).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *WorkoutSessionService) GetSession(id string, userId string) (*WorkoutSession, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "userId", Value: userId},
	}

	session := &WorkoutSession{}
	if err := s.DB.Collection("workoutSessions").FindOne(context.Background(), filter).Decode(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *WorkoutSessionService) GetUserSessions(userId string) ([]*WorkoutSession, error) {
	filter := bson.D{{Key: "userId", Value: userId}}
	opts := options.Find().SetSort(bson.D{{Key: "startTime", Value: -1}})

	cursor, err := s.DB.Collection("workoutSessions").Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}

	var sessions []*WorkoutSession
	if err := cursor.All(context.Background(), &sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (s *WorkoutSessionService) DeleteSession(id string, userId string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "userId", Value: userId},
	}

	result, err := s.DB.Collection("workoutSessions").DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
