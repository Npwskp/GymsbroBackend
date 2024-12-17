package workout

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WorkoutService struct {
	DB *mongo.Database
}

type IWorkoutService interface {
	CreateWorkout(workout *CreateWorkoutDto, userId string) (*Workout, error)
	GetWorkout(id string, userId string) (*Workout, error)
	GetWorkoutsByUser(userId string) ([]*Workout, error)
	UpdateWorkout(id string, workout *UpdateWorkoutDto, userId string) (*Workout, error)
	DeleteWorkout(id string, userId string) error
}

func (ws *WorkoutService) CreateWorkout(dto *CreateWorkoutDto, userId string) (*Workout, error) {
	workout := &Workout{
		UserID:      userId,
		Name:        dto.Name,
		Description: dto.Description,
		Exercises:   dto.Exercises,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result, err := ws.DB.Collection("workout").InsertOne(context.Background(), workout)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: result.InsertedID}}
	createdWorkout := &Workout{}
	if err := ws.DB.Collection("workout").FindOne(context.Background(), filter).Decode(createdWorkout); err != nil {
		return nil, err
	}

	return createdWorkout, nil
}

func (ws *WorkoutService) GetWorkout(id string, userId string) (*Workout, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "userId", Value: userId},
	}

	workout := &Workout{}
	if err := ws.DB.Collection("workout").FindOne(context.Background(), filter).Decode(workout); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("workout not found")
		}
		return nil, err
	}

	return workout, nil
}

func (ws *WorkoutService) GetWorkoutsByUser(userId string) ([]*Workout, error) {
	filter := bson.D{{Key: "userId", Value: userId}}

	cursor, err := ws.DB.Collection("workout").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	var workouts []*Workout
	if err := cursor.All(context.Background(), &workouts); err != nil {
		return nil, err
	}

	return workouts, nil
}

func (ws *WorkoutService) UpdateWorkout(id string, dto *UpdateWorkoutDto, userId string) (*Workout, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "userId", Value: userId},
	}

	workout := &Workout{}
	if err := ws.DB.Collection("workout").FindOne(context.Background(), filter).Decode(workout); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("workout not found")
		}
		return nil, err
	}

	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "name", Value: dto.Name},
		{Key: "description", Value: dto.Description},
		{Key: "exercises", Value: dto.Exercises},
		{Key: "updatedAt", Value: time.Now()},
	}}}

	if err := ws.DB.Collection("workout").FindOneAndUpdate(
		context.Background(),
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(workout); err != nil {
		return nil, err
	}

	return workout, nil
}

func (ws *WorkoutService) DeleteWorkout(id string, userId string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "userId", Value: userId},
	}

	result, err := ws.DB.Collection("workout").DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("workout not found")
	}

	return nil
}
