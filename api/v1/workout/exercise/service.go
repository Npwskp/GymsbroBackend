package exercise

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/function"
	minio "github.com/Npwskp/GymsbroBackend/api/v1/storage"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ExerciseImageBucketName = "exercise-image"
)

type ExerciseService struct {
	DB           *mongo.Database
	MinioService minio.MinioService
}

type IExerciseService interface {
	CreateExercise(exercise *CreateExerciseDto, userId string) (*Exercise, error)
	CreateManyExercises(exercises *[]CreateExerciseDto, userId string) ([]*Exercise, error)
	GetAllExercises(userId string) ([]*Exercise, error)
	GetExercise(id string, userId string) (*Exercise, error)
	GetExerciseByType(exerciseType string, userId string) ([]*Exercise, error)
	DeleteExercise(id string, userId string) error
	UpdateExercise(doc *UpdateExerciseDto, id string, userId string) (*Exercise, error)
	UpdateExerciseImage(c *fiber.Ctx, id string, file io.Reader, filename string, contentType string, userId string) (*Exercise, error)
}

func (es *ExerciseService) CreateExercise(exercise *CreateExerciseDto, userId string) (*Exercise, error) {
	if exercise.Type == nil {
		exercise.Type = []string{}
	}
	if exercise.Muscle == nil {
		exercise.Muscle = []string{}
	}

	// Create exercise with userId
	exerciseDoc := &Exercise{
		UserID:      userId,
		Name:        exercise.Name,
		Description: exercise.Description,
		Type:        exercise.Type,
		Muscle:      exercise.Muscle,
		Image:       exercise.Image,
	}

	result, err := es.DB.Collection("exercises").InsertOne(context.Background(), exerciseDoc)
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

func (es *ExerciseService) CreateManyExercises(exercises *[]CreateExerciseDto, userId string) ([]*Exercise, error) {
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

func (es *ExerciseService) GetAllExercises(userId string) ([]*Exercise, error) {
	filter := bson.D{
		{"$or", bson.A{
			bson.D{{Key: "userid", Value: userId}},
			bson.D{{Key: "userid", Value: nil}},
			bson.D{{Key: "userid", Value: ""}},
		}},
	}

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

func (es *ExerciseService) GetExercise(id string, userId string) (*Exercise, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "userid", Value: userId},
	}
	exercise := &Exercise{}
	if err := es.DB.Collection("exercises").FindOne(context.Background(), filter).Decode(exercise); err != nil {
		return nil, err
	}
	return exercise, nil
}

func (es *ExerciseService) GetExerciseByType(exerciseType string, userId string) ([]*Exercise, error) {
	filter := bson.M{"type": bson.M{"$in": []string{exerciseType}}, "userid": userId}
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

func (es *ExerciseService) DeleteExercise(id string, userId string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "userid", Value: userId},
	}
	if _, err := es.DB.Collection("exercises").DeleteOne(context.Background(), filter); err != nil {
		return err
	}
	return nil
}

func (es *ExerciseService) UpdateExercise(doc *UpdateExerciseDto, id string, userId string) (*Exercise, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "userid", Value: userId},
	}
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
	return es.GetExercise(id, userId)
}

func (es *ExerciseService) UpdateExerciseImage(c *fiber.Ctx, id string, file io.Reader, filename string, contentType string, userId string) (*Exercise, error) {
	// Get exercise first to verify existence and get current image URL
	exercise, err := es.GetExercise(id, userId)
	if err != nil {
		return nil, err
	}

	oldImageURL := exercise.Image

	ext := strings.ToLower(filepath.Ext(filename))
	// Generate unique filename
	timestamp := time.Now().UnixNano()
	objectName := fmt.Sprintf("exercises/%s/image_%d%s", id, timestamp, ext)

	// Upload to MinIO
	err = es.MinioService.UploadFile(c.Context(), file, ExerciseImageBucketName, objectName, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image: %v", err)
	}

	// Get the URL of the uploaded file
	url, err := es.MinioService.GetFileURL(c.Context(), ExerciseImageBucketName, objectName)
	if err != nil {
		return nil, fmt.Errorf("failed to get image URL: %v", err)
	}

	// Update exercise's image URL in database
	oid, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: oid}, {Key: "userid", Value: userId}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "image", Value: url},
			{Key: "updated_at", Value: time.Now()},
		}},
	}

	result, err := es.DB.Collection("exercises").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("no exercise found for the given ID")
	}

	// Delete old image after successful upload and update
	if oldImageURL != "" {
		baseURL := strings.Split(oldImageURL, "?")[0]
		urlParts := strings.Split(baseURL, es.MinioService.GetFullBucketName(ExerciseImageBucketName)+"/")
		if len(urlParts) > 1 {
			oldObjectName := urlParts[1]
			if err := es.MinioService.DeleteFile(c.Context(), ExerciseImageBucketName, oldObjectName); err != nil {
				fmt.Printf("Warning: Failed to delete old exercise image: %v\n", err)
			}
		}
	}

	return es.GetExercise(id, userId)
}
