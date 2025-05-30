package exercise

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"sort"
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
	FindSimilarExercises(id string, userId string, limit int) ([]*Exercise, error)
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
		DeletedAt:    time.Time{},
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
			DeletedAt:    time.Time{},
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
			{"deleted_at": bson.M{"$exists": false}},
			{"deleted_at": ""},
		}},
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
	exercises := make([]*Exercise, 0)
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
		{Key: "$or", Value: []bson.M{
			{"deleted_at": bson.M{"$exists": false}},
			{"deleted_at": ""},
		}},
	}

	now := time.Now()
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "deleted_at", Value: now},
			{Key: "updated_at", Value: now},
		}},
	}

	result, err := es.DB.Collection("exercises").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("exercise not found or already deleted")
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

	// Update exercise's image URL in database - removed deleted_at filter
	oid, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "userid", Value: userId},
		{Key: "$or", Value: []bson.M{
			{"deleted_at": bson.M{"$exists": false}},
			{"deleted_at": ""},
		}},
	}
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

	// Base filter for non-deleted records and user access
	filter := bson.D{
		{Key: "$or", Value: []bson.M{
			{"deleted_at": bson.M{"$exists": false}},
			{"deleted_at": ""},
		}},
		{Key: "$or", Value: []bson.M{
			{"userid": userID},
			{"userid": ""},
			{"userid": nil},
		}},
	}

	// Add search conditions if query is provided
	if query != "" {
		filter = append(filter, bson.E{
			Key: "name",
			Value: bson.M{
				"$regex":   fmt.Sprintf(".*%s.*", query),
				"$options": "i",
			},
		})
	}

	// Add equipment filter if provided
	if len(equipment) > 0 {
		equipmentStrings := make([]string, len(equipment))
		for i, eq := range equipment {
			equipmentStrings[i] = string(eq)
		}
		filter = append(filter, bson.E{Key: "equipment", Value: bson.M{"$in": equipmentStrings}})
	}

	// Add mechanics filter if provided
	if len(mechanics) > 0 {
		mechanicsStrings := make([]string, len(mechanics))
		for i, m := range mechanics {
			mechanicsStrings[i] = string(m)
		}
		filter = append(filter, bson.E{Key: "mechanics", Value: bson.M{"$in": mechanicsStrings}})
	}

	// Add force filter if provided
	if len(force) > 0 {
		forceStrings := make([]string, len(force))
		for i, f := range force {
			forceStrings[i] = string(f)
		}
		filter = append(filter, bson.E{Key: "force", Value: bson.M{"$in": forceStrings}})
	}

	// Add body part filter if provided
	if len(bodyPart) > 0 {
		bodyPartStrings := make([]string, len(bodyPart))
		for i, bp := range bodyPart {
			bodyPartStrings[i] = string(bp)
		}
		filter = append(filter, bson.E{Key: "body_part", Value: bson.M{"$in": bodyPartStrings}})
	}

	// Add target muscle filter if provided
	if len(targetMuscle) > 0 {
		targetMuscleStrings := make([]string, len(targetMuscle))
		for i, tm := range targetMuscle {
			targetMuscleStrings[i] = string(tm)
		}
		filter = append(filter, bson.E{Key: "target_muscle", Value: bson.M{"$in": targetMuscleStrings}})
	}

	cursor, err := es.DB.Collection("exercises").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	exercises := make([]*Exercise, 0)
	if err := cursor.All(context.Background(), &exercises); err != nil {
		return nil, err
	}
	return exercises, nil
}

type SimilarityScore struct {
	Exercise *Exercise
	Score    float64
}

func (es *ExerciseService) FindSimilarExercises(id string, userId string, limit int) ([]*Exercise, error) {
	// Get the source exercise
	sourceExercise, err := es.GetExercise(id, userId)
	if err != nil {
		return nil, err
	}

	// Get all exercises
	allExercises, err := es.GetAllExercises(userId)
	if err != nil {
		return nil, err
	}

	// Calculate similarity scores
	var scores []SimilarityScore
	for _, exercise := range allExercises {
		// Skip the source exercise itself
		if exercise.ID == sourceExercise.ID {
			continue
		}

		// Calculate name similarity (using case-insensitive comparison)
		nameSimilarity := calculateStringSimilarity(strings.ToLower(sourceExercise.Name), strings.ToLower(exercise.Name))

		// Calculate body part overlap
		bodyPartOverlap := calculateOverlap(sourceExercise.BodyPart, exercise.BodyPart)

		// Calculate target muscle overlap
		targetMuscleOverlap := calculateOverlap(sourceExercise.TargetMuscle, exercise.TargetMuscle)

		// Equipment similarity
		equipmentSimilarity := 0.0
		if sourceExercise.Equipment == exercise.Equipment {
			equipmentSimilarity = 1.0
		}

		// Mechanics similarity
		mechanicsSimilarity := 0.0
		if sourceExercise.Mechanics == exercise.Mechanics {
			mechanicsSimilarity = 1.0
		}

		// Calculate total score (weighted sum)
		totalScore := (nameSimilarity * 0.2) +
			(bodyPartOverlap * 0.3) +
			(targetMuscleOverlap * 0.3) +
			(equipmentSimilarity * 0.1) +
			(mechanicsSimilarity * 0.1)

		scores = append(scores, SimilarityScore{
			Exercise: exercise,
			Score:    totalScore,
		})
	}

	// Sort by similarity score
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	// Get top N similar exercises
	if limit <= 0 {
		limit = 5 // Default limit
	}
	if limit > len(scores) {
		limit = len(scores)
	}

	result := make([]*Exercise, limit)
	for i := 0; i < limit; i++ {
		result[i] = scores[i].Exercise
	}

	return result, nil
}

// calculateStringSimilarity calculates similarity between two strings
// using Levenshtein distance normalized to [0,1]
func calculateStringSimilarity(s1, s2 string) float64 {
	d := levenshteinDistance(s1, s2)
	maxLen := math.Max(float64(len(s1)), float64(len(s2)))
	if maxLen == 0 {
		return 1.0
	}
	return 1.0 - float64(d)/maxLen
}

// levenshteinDistance calculates the Levenshtein distance between two strings
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	matrix := make([][]float64, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]float64, len(s2)+1)
		matrix[i][0] = float64(i)
	}
	for j := range matrix[0] {
		matrix[0][j] = float64(j)
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}
			matrix[i][j] = math.Min(
				math.Min(matrix[i-1][j]+1,
					matrix[i][j-1]+1),
				matrix[i-1][j-1]+float64(cost),
			)
		}
	}

	return int(matrix[len(s1)][len(s2)])
}

// calculateOverlap calculates the overlap ratio between two slices
func calculateOverlap[T comparable](s1, s2 []T) float64 {
	if len(s1) == 0 && len(s2) == 0 {
		return 1.0
	}
	if len(s1) == 0 || len(s2) == 0 {
		return 0.0
	}

	// Count matching elements
	matches := 0
	for _, v1 := range s1 {
		for _, v2 := range s2 {
			if v1 == v2 {
				matches++
				break
			}
		}
	}

	// Calculate Jaccard similarity coefficient
	union := float64(len(s1) + len(s2) - matches)
	if union == 0 {
		return 1.0
	}
	return float64(matches) / union
}
