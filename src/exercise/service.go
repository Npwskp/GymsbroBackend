package exercise

import (
	"context"
	"fmt"

	"github.com/Npwskp/GymsbroBackend/src/function"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Exercise struct {
	ID          string   `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Type        []string `json:"type" validate:"required"`
	Muscle      []string `json:"muscle" validate:"required"`
	Image       string   `json:"image" validate:"required"`
}

type ExerciseService struct {
	DB *mongo.Database
}

type IExerciseService interface {
	CreateExercise(exercise *CreateExerciseDto) (*Exercise, error)
	CreateManyExercises(exercises *[]CreateExerciseDto) ([]*Exercise, error)
	GetAllExercises() ([]*Exercise, error)
	GetExercise(id string) (*Exercise, error)
	GetExerciseByType(exerciseType string) ([]*Exercise, error)
	DeleteExercise(id string) error
	UpdateExercise(doc *UpdateExerciseDto, id string) (*Exercise, error)
}

func (es *ExerciseService) CreateExercise(exercise *CreateExerciseDto) (*Exercise, error) {
	if exercise.Type == nil {
		exercise.Type = []string{}
	}
	result, err := es.DB.Collection("exercises").InsertOne(context.Background(), exercise)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: result.InsertedID}}
	createdRecord := es.DB.Collection("exercises").FindOne(context.Background(), filter)
	createdExercise := &Exercise{}
	if err := createdRecord.Decode(createdExercise); err != nil {
		return nil, err
	}
	return createdExercise, nil
}

func (es *ExerciseService) CreateManyExercises(exercises *[]CreateExerciseDto) ([]*Exercise, error) {
	var result []interface{}
	for _, exercise := range *exercises {
		if exercise.Type == nil {
			exercise.Type = []string{}
		}
		if exercise.Muscle == nil {
			exercise.Muscle = []string{}
		}
		result = append(result, exercise)
	}
	if _, err := es.DB.Collection("exercises").InsertMany(context.Background(), result); err != nil {
		return nil, err
	}
	var createdExercises []*Exercise
	for _, exercise := range *exercises {
		filter := bson.D{{Key: "name", Value: exercise.Name}}
		createdRecord := es.DB.Collection("exercises").FindOne(context.Background(), filter)
		createdExercise := &Exercise{}
		if err := createdRecord.Decode(createdExercise); err != nil {
			return nil, err
		}
		createdExercises = append(createdExercises, createdExercise)
	}
	return createdExercises, nil
}

func (es *ExerciseService) GetAllExercises() ([]*Exercise, error) {
	cursor, err := es.DB.Collection("exercises").Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	var exercises []*Exercise
	if err := cursor.All(context.Background(), &exercises); err != nil {
		return nil, err
	}
	return exercises, nil
}

func (es *ExerciseService) GetExercise(id string) (*Exercise, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: oid}}
	exercise := &Exercise{}
	if err := es.DB.Collection("exercises").FindOne(context.Background(), filter).Decode(exercise); err != nil {
		return nil, err
	}
	return exercise, nil
}

func (es *ExerciseService) GetExerciseByType(exerciseType string) ([]*Exercise, error) {
	filter := bson.M{"type": bson.M{"$in": []string{exerciseType}}}
	cursor, err := es.DB.Collection("exercises").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	var exercises []*Exercise
	if err := cursor.All(context.Background(), &exercises); err != nil {
		return nil, err
	}
	return exercises, nil
}

func (es *ExerciseService) DeleteExercise(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: oid}}
	if _, err := es.DB.Collection("exercises").DeleteOne(context.Background(), filter); err != nil {
		return err
	}
	return nil
}

func (es *ExerciseService) UpdateExercise(doc *UpdateExerciseDto, id string) (*Exercise, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: oid}}
	exercise := &Exercise{}
	if err := es.DB.Collection("exercises").FindOne(context.Background(), filter).Decode(exercise); err != nil {
		return nil, err
	}
	fmt.Println(exercise, "(", doc.Image, ")")
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: function.Coalesce(doc.Name, exercise.Name)},
			{Key: "description", Value: function.Coalesce(doc.Description, exercise.Description)},
			{Key: "type", Value: function.Coalesce(doc.Type, exercise.Type)},
			{Key: "muscle", Value: function.Coalesce(doc.Muscle, exercise.Muscle)},
			{Key: "image", Value: function.Coalesce(doc.Image, exercise.Image)},
		}},
	}
	if _, err := es.DB.Collection("exercises").UpdateOne(context.Background(), filter, update); err != nil {
		return nil, err
	}
	return es.GetExercise(id)
}
