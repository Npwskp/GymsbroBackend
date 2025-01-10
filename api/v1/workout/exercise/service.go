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
	exerciseEnums "github.com/Npwskp/GymsbroBackend/api/v1/workout/exercise/enums"
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
	DeleteExercise(id string, userId string) error
	UpdateExercise(doc *UpdateExerciseDto, id string, userId string) (*Exercise, error)
	UpdateExerciseImage(c *fiber.Ctx, id string, file io.Reader, filename string, contentType string, userId string) (*Exercise, error)
	SearchAndFilterExercise(equipment []exerciseEnums.Equipment, mechanics []exerciseEnums.Mechanics, force []exerciseEnums.Force, bodyPart []exerciseEnums.BodyPart, targetMuscle []exerciseEnums.TargetMuscle, query string, userID string) ([]*Exercise, error)
}

func (es *ExerciseService) CreateExercise(exercise *CreateExerciseDto, userId string) (*Exercise, error) {
	if exercise.BodyPart == nil {
		exercise.BodyPart = []exerciseEnums.BodyPart{}
	}
	if exercise.TargetMuscle == nil {
		exercise.TargetMuscle = []exerciseEnums.TargetMuscle{}
	}

	// Create exercise with userId
	exerciseDoc := &Exercise{
		UserID:       userId,
		Name:         exercise.Name,
		Equipment:    exercise.Equipment,
		Mechanics:    exercise.Mechanics,
		Force:        exercise.Force,
		Preparation:  exercise.Preparation,
		Execution:    exercise.Execution,
		Image:        exercise.Image,
		BodyPart:     exercise.BodyPart,
		TargetMuscle: exercise.TargetMuscle,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
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
		if exercise.BodyPart == nil {
			exercise.BodyPart = []exerciseEnums.BodyPart{}
		}
		if exercise.TargetMuscle == nil {
			exercise.TargetMuscle = []exerciseEnums.TargetMuscle{}
		}

		exerciseDoc := Exercise{
			UserID:       userId,
			Name:         exercise.Name,
			Equipment:    exercise.Equipment,
			Mechanics:    exercise.Mechanics,
			Force:        exercise.Force,
			Preparation:  exercise.Preparation,
			Execution:    exercise.Execution,
			Image:        exercise.Image,
			BodyPart:     exercise.BodyPart,
			TargetMuscle: exercise.TargetMuscle,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		result = append(result, exerciseDoc)
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
		{Key: "$or", Value: []bson.M{
			{"userid": userId},
			{"userid": ""},
			{"userid": nil},
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

	// Allow access to both user-specific exercises and public exercises (empty or nil userID)
	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "$or", Value: []bson.M{
			{"userid": userId},
			{"userid": ""},
			{"userid": nil},
		}},
	}

	exercise := &Exercise{}
	if err := es.DB.Collection("exercises").FindOne(context.Background(), filter).Decode(exercise); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("exercise not found")
		}
		return nil, err
	}
	return exercise, nil
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

	// Create update document
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: function.Coalesce(doc.Name, exercise.Name)},
			{Key: "image", Value: function.Coalesce(doc.Image, exercise.Image)},
			{Key: "equipment", Value: function.Coalesce(doc.Equipment, exercise.Equipment)},
			{Key: "mechanics", Value: function.Coalesce(doc.Mechanics, exercise.Mechanics)},
			{Key: "force", Value: function.Coalesce(doc.Force, exercise.Force)},
			{Key: "preparation", Value: function.Coalesce(doc.Preparation, exercise.Preparation)},
			{Key: "execution", Value: function.Coalesce(doc.Execution, exercise.Execution)},
			{Key: "body_part", Value: function.Coalesce(doc.BodyPart, exercise.BodyPart)},
			{Key: "target_muscle", Value: function.Coalesce(doc.TargetMuscle, exercise.TargetMuscle)},
			{Key: "updated_at", Value: time.Now()},
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

func (es *ExerciseService) SearchAndFilterExercise(
	equipment []exerciseEnums.Equipment,
	mechanics []exerciseEnums.Mechanics,
	force []exerciseEnums.Force,
	bodyPart []exerciseEnums.BodyPart,
	targetMuscle []exerciseEnums.TargetMuscle,
	query string, userID string) ([]*Exercise, error) {

	filter := bson.D{}
	andConditions := []bson.D{}

	// Add user filter
	userFilter := bson.D{
		{Key: "$or", Value: []bson.M{
			{"userid": userID},
			{"userid": ""},
			{"userid": nil},
		}},
	}
	andConditions = append(andConditions, userFilter)

	// Add equipment filter if provided
	if len(equipment) > 0 {
		andConditions = append(andConditions, bson.D{{Key: "equipment", Value: bson.D{{Key: "$in", Value: equipment}}}})
	}

	// Add mechanics filter if provided
	if len(mechanics) > 0 {
		andConditions = append(andConditions, bson.D{{Key: "mechanics", Value: bson.D{{Key: "$in", Value: mechanics}}}})
	}

	// Add force filter if provided
	if len(force) > 0 {
		andConditions = append(andConditions, bson.D{{Key: "force", Value: bson.D{{Key: "$in", Value: force}}}})
	}

	// Add body part filter if provided
	if len(bodyPart) > 0 {
		andConditions = append(andConditions, bson.D{{Key: "body_part", Value: bson.D{{Key: "$in", Value: bodyPart}}}})
	}

	// Add target muscle filter if provided
	if len(targetMuscle) > 0 {
		andConditions = append(andConditions, bson.D{{Key: "target_muscle", Value: bson.D{{Key: "$in", Value: targetMuscle}}}})
	}

	// Add name search if query is provided
	if query != "" {
		andConditions = append(andConditions, bson.D{{Key: "name", Value: bson.D{{Key: "$regex", Value: query}, {Key: "$options", Value: "i"}}}})
	}

	// Combine all conditions with $and
	if len(andConditions) > 0 {
		filter = bson.D{{Key: "$and", Value: andConditions}}
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
