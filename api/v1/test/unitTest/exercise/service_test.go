package exercise_test

import (
	"context"
	"testing"
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/workout/exercise"
	exerciseEnums "github.com/Npwskp/GymsbroBackend/api/v1/workout/exercise/enums"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Setup functions
func setupTestDB(t *testing.T) *mongo.Database {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	db := client.Database("testdb_" + primitive.NewObjectID().Hex())

	t.Cleanup(func() {
		if err := db.Drop(context.Background()); err != nil {
			t.Errorf("Failed to drop test database: %v", err)
		}
		if err := client.Disconnect(context.Background()); err != nil {
			t.Errorf("Failed to disconnect from MongoDB: %v", err)
		}
	})

	return db
}

func TestCreateExercise(t *testing.T) {
	db := setupTestDB(t)
	service := &exercise.ExerciseService{DB: db}

	t.Run("Valid exercise creation", func(t *testing.T) {
		dto := &exercise.CreateExerciseDto{
			Name:         "Push-up",
			Equipment:    exerciseEnums.BodyWeight,
			Mechanics:    exerciseEnums.Compound,
			Force:        exerciseEnums.Push,
			Preparation:  []string{"Get into plank position"},
			Execution:    []string{"Lower body", "Push back up"},
			BodyPart:     []exerciseEnums.BodyPart{exerciseEnums.Chest},
			TargetMuscle: []exerciseEnums.TargetMuscle{exerciseEnums.PectoralisMajorSternal},
		}
		userId := "test_user"

		result, err := service.CreateExercise(dto, userId)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userId, result.UserID)
		assert.Equal(t, dto.Name, result.Name)
		assert.Equal(t, dto.Equipment, result.Equipment)
		assert.Equal(t, dto.Mechanics, result.Mechanics)
		assert.Equal(t, dto.Force, result.Force)
	})
}

func TestGetExercise(t *testing.T) {
	db := setupTestDB(t)
	service := &exercise.ExerciseService{DB: db}

	testExercise := &exercise.Exercise{
		ID:        primitive.NewObjectID(),
		UserID:    "test_user",
		Name:      "Squat",
		Equipment: exerciseEnums.BodyWeight,
		Mechanics: exerciseEnums.Compound,

		Force:        exerciseEnums.Push,
		Preparation:  []string{"Stand with feet shoulder-width apart"},
		Execution:    []string{"Bend knees", "Lower body", "Stand back up"},
		BodyPart:     []exerciseEnums.BodyPart{exerciseEnums.Thighs},
		TargetMuscle: []exerciseEnums.TargetMuscle{exerciseEnums.Quadriceps},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err := db.Collection("exercises").InsertOne(context.Background(), testExercise)
	assert.NoError(t, err)

	t.Run("Get existing exercise", func(t *testing.T) {
		result, err := service.GetExercise(testExercise.ID.Hex(), testExercise.UserID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, testExercise.ID, result.ID)
		assert.Equal(t, testExercise.Name, result.Name)
	})

	t.Run("Get non-existing exercise", func(t *testing.T) {
		nonExistingID := primitive.NewObjectID().Hex()
		result, err := service.GetExercise(nonExistingID, "test_user")
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestGetAllExercises(t *testing.T) {
	db := setupTestDB(t)
	service := &exercise.ExerciseService{DB: db}
	userId := "test_user"

	testExercises := []interface{}{
		&exercise.Exercise{
			ID:        primitive.NewObjectID(),
			UserID:    userId,
			Name:      "Push-up",
			Equipment: exerciseEnums.BodyWeight,
			Mechanics: exerciseEnums.Compound,
			Force:     exerciseEnums.Push,

			BodyPart:     []exerciseEnums.BodyPart{exerciseEnums.Chest},
			TargetMuscle: []exerciseEnums.TargetMuscle{exerciseEnums.PectoralisMajorSternal},
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		&exercise.Exercise{
			ID:        primitive.NewObjectID(),
			UserID:    userId,
			Name:      "Squat",
			Equipment: exerciseEnums.BodyWeight,
			Mechanics: exerciseEnums.Compound,
			Force:     exerciseEnums.Push,

			BodyPart:     []exerciseEnums.BodyPart{exerciseEnums.Thighs},
			TargetMuscle: []exerciseEnums.TargetMuscle{exerciseEnums.Quadriceps},
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	_, err := db.Collection("exercises").InsertMany(context.Background(), testExercises)
	assert.NoError(t, err)

	t.Run("Get all exercises for user", func(t *testing.T) {
		exercises, err := service.GetAllExercises(userId)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(exercises))
	})
}

func TestDeleteExercise(t *testing.T) {
	db := setupTestDB(t)
	service := &exercise.ExerciseService{DB: db}

	testExercise := &exercise.Exercise{
		ID:        primitive.NewObjectID(),
		UserID:    "test_user",
		Name:      "Push-up",
		Equipment: exerciseEnums.BodyWeight,
		Mechanics: exerciseEnums.Compound,
		Force:     exerciseEnums.Push,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := db.Collection("exercises").InsertOne(context.Background(), testExercise)
	assert.NoError(t, err)

	t.Run("Delete existing exercise", func(t *testing.T) {
		err := service.DeleteExercise(testExercise.ID.Hex(), testExercise.UserID)
		assert.NoError(t, err)

		// Verify deletion (should have deleted_at set)
		var found exercise.Exercise
		err = db.Collection("exercises").FindOne(context.Background(), bson.M{"_id": testExercise.ID}).Decode(&found)
		assert.NoError(t, err)
		assert.False(t, found.DeletedAt.IsZero())
	})

	t.Run("Delete non-existing exercise", func(t *testing.T) {
		nonExistingID := primitive.NewObjectID().Hex()
		err := service.DeleteExercise(nonExistingID, "test_user")
		assert.Error(t, err)
	})
}

func TestUpdateExercise(t *testing.T) {
	db := setupTestDB(t)
	service := &exercise.ExerciseService{DB: db}

	testExercise := &exercise.Exercise{
		ID:        primitive.NewObjectID(),
		UserID:    "test_user",
		Name:      "Push-up",
		Equipment: exerciseEnums.BodyWeight,
		Mechanics: exerciseEnums.Compound,
		Force:     exerciseEnums.Push,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := db.Collection("exercises").InsertOne(context.Background(), testExercise)
	assert.NoError(t, err)

	t.Run("Update existing exercise", func(t *testing.T) {
		updateDto := &exercise.UpdateExerciseDto{
			Name:      "Modified Push-up",
			Equipment: exerciseEnums.Barbell,
		}

		updated, err := service.UpdateExercise(updateDto, testExercise.ID.Hex(), testExercise.UserID)
		assert.NoError(t, err)
		assert.Equal(t, updateDto.Name, updated.Name)
		assert.Equal(t, updateDto.Equipment, updated.Equipment)
		assert.Equal(t, testExercise.Mechanics, updated.Mechanics) // Unchanged field
	})
}

func TestSearchAndFilterExercise(t *testing.T) {
	db := setupTestDB(t)
	service := &exercise.ExerciseService{DB: db}
	userId := "test_user"

	testExercises := []interface{}{
		&exercise.Exercise{
			ID:        primitive.NewObjectID(),
			UserID:    userId,
			Name:      "Push-up",
			Equipment: exerciseEnums.BodyWeight,
			Mechanics: exerciseEnums.Compound,
			Force:     exerciseEnums.Push,

			BodyPart:     []exerciseEnums.BodyPart{exerciseEnums.Chest},
			TargetMuscle: []exerciseEnums.TargetMuscle{exerciseEnums.PectoralisMajorSternal},
		},
		&exercise.Exercise{
			ID:        primitive.NewObjectID(),
			UserID:    userId,
			Name:      "Barbell Squat",
			Equipment: exerciseEnums.Barbell,
			Mechanics: exerciseEnums.Compound,
			Force:     exerciseEnums.Push,

			BodyPart:     []exerciseEnums.BodyPart{exerciseEnums.Thighs},
			TargetMuscle: []exerciseEnums.TargetMuscle{exerciseEnums.Quadriceps},
		},
	}

	_, err := db.Collection("exercises").InsertMany(context.Background(), testExercises)
	assert.NoError(t, err)

	t.Run("Search by equipment", func(t *testing.T) {
		results, err := service.SearchAndFilterExercise(
			[]exerciseEnums.Equipment{exerciseEnums.BodyWeight},
			nil, nil, nil, nil,
			"", userId,
		)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(results))
		assert.Equal(t, "Push-up", results[0].Name)
	})

	t.Run("Search by query", func(t *testing.T) {
		results, err := service.SearchAndFilterExercise(
			nil, nil, nil, nil, nil,
			"Barbell", userId,
		)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(results))
		assert.Equal(t, "Barbell Squat", results[0].Name)
	})
}
