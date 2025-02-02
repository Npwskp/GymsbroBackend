package workout_test

import (
	"context"
	"testing"
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/workout/workout"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Setup functions
func setupTestDB(t *testing.T) *mongo.Database {
	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Create a test database
	db := client.Database("testdb_" + primitive.NewObjectID().Hex())

	// Clean up function
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

func TestCreateWorkout(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	service := &workout.WorkoutService{DB: db}

	t.Run("Success - Create Valid Workout", func(t *testing.T) {
		dto := &workout.CreateWorkoutDto{
			Name:        "Test Workout",
			Description: "Test Description",
			Exercises: []workout.WorkoutExercise{
				{ExerciseID: primitive.NewObjectID().Hex(), Order: 0},
			},
		}
		userId := "test_user"

		result, err := service.CreateWorkout(dto, userId)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, dto.Name, result.Name)
		assert.Equal(t, dto.Description, result.Description)
		assert.Equal(t, userId, result.UserID)
		assert.NotEmpty(t, result.ID)
		assert.NotZero(t, result.CreatedAt)
		assert.NotZero(t, result.UpdatedAt)
	})

	t.Run("Success - Create Workout with Empty Name (Validation at Controller)", func(t *testing.T) {
		dto := &workout.CreateWorkoutDto{
			Name:        "",
			Description: "Test Description",
			Exercises:   []workout.WorkoutExercise{},
		}
		userId := "test_user"

		result, err := service.CreateWorkout(dto, userId)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, dto.Name, result.Name)
	})
}

func TestGetWorkout(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	service := &workout.WorkoutService{DB: db}

	// Create a test workout for reuse
	testWorkout := &workout.Workout{
		ID:          primitive.NewObjectID(),
		UserID:      "test_user",
		Name:        "Test Workout",
		Description: "Test Description",
		Exercises: []workout.WorkoutExercise{
			{ExerciseID: primitive.NewObjectID().Hex(), Order: 0},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := db.Collection("workout").InsertOne(context.Background(), testWorkout)
	assert.NoError(t, err)

	t.Run("Success - Get Existing Workout", func(t *testing.T) {
		result, err := service.GetWorkout(testWorkout.ID.Hex())
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, testWorkout.ID, result.ID)
		assert.Equal(t, testWorkout.Name, result.Name)
	})

	t.Run("Error - Non-Existing Workout", func(t *testing.T) {
		result, err := service.GetWorkout(primitive.NewObjectID().Hex())
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Error - Invalid ID Format", func(t *testing.T) {
		result, err := service.GetWorkout("invalid-id")
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestGetWorkouts(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	service := &workout.WorkoutService{DB: db}

	t.Run("Success - Get Multiple Workouts", func(t *testing.T) {
		// Create test workouts
		workouts := []interface{}{
			&workout.Workout{
				ID:          primitive.NewObjectID(),
				UserID:      "test_user",
				Name:        "Workout 1",
				Description: "Description 1",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			&workout.Workout{
				ID:          primitive.NewObjectID(),
				UserID:      "test_user",
				Name:        "Workout 2",
				Description: "Description 2",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		_, err := db.Collection("workout").InsertMany(context.Background(), workouts)
		assert.NoError(t, err)

		results, err := service.GetWorkouts("test_user")
		assert.NoError(t, err)
		assert.Len(t, results, 2)
	})

	t.Run("Success - Empty Result for New User", func(t *testing.T) {
		results, err := service.GetWorkouts("new_user")
		assert.NoError(t, err)
		assert.Empty(t, results)
	})
}

func TestUpdateWorkout(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	service := &workout.WorkoutService{DB: db}

	// Create a test workout
	testWorkout := &workout.Workout{
		ID:          primitive.NewObjectID(),
		UserID:      "test_user",
		Name:        "Test Workout",
		Description: "Test Description",
		Exercises: []workout.WorkoutExercise{
			{ExerciseID: primitive.NewObjectID().Hex(), Order: 0},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := db.Collection("workout").InsertOne(context.Background(), testWorkout)
	assert.NoError(t, err)

	t.Run("Success - Update Existing Workout", func(t *testing.T) {
		updateDto := &workout.UpdateWorkoutDto{
			Name:        "Updated Workout",
			Description: "Updated Description",
			Exercises: []workout.WorkoutExercise{
				{ExerciseID: primitive.NewObjectID().Hex(), Order: 1},
			},
		}

		result, err := service.UpdateWorkout(testWorkout.ID.Hex(), updateDto, "test_user")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, updateDto.Name, result.Name)
		assert.Equal(t, updateDto.Description, result.Description)
	})

	t.Run("Error - Non-Existing Workout", func(t *testing.T) {
		updateDto := &workout.UpdateWorkoutDto{
			Name:        "Updated Workout",
			Description: "Updated Description",
		}
		result, err := service.UpdateWorkout(primitive.NewObjectID().Hex(), updateDto, "test_user")
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Error - Wrong User ID", func(t *testing.T) {
		updateDto := &workout.UpdateWorkoutDto{
			Name:        "Updated Workout",
			Description: "Updated Description",
		}
		result, err := service.UpdateWorkout(testWorkout.ID.Hex(), updateDto, "wrong_user")
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Error - Invalid Workout ID Format", func(t *testing.T) {
		updateDto := &workout.UpdateWorkoutDto{
			Name:        "Updated Workout",
			Description: "Updated Description",
		}
		result, err := service.UpdateWorkout("invalid-id", updateDto, "test_user")
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestDeleteWorkout(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	service := &workout.WorkoutService{DB: db}

	t.Run("Success - Delete Existing Workout", func(t *testing.T) {
		// Create a test workout
		testWorkout := &workout.Workout{
			ID:          primitive.NewObjectID(),
			UserID:      "test_user",
			Name:        "Test Workout",
			Description: "Test Description",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		_, err := db.Collection("workout").InsertOne(context.Background(), testWorkout)
		assert.NoError(t, err)

		err = service.DeleteWorkout(testWorkout.ID.Hex(), "test_user")
		assert.NoError(t, err)

		// Verify workout is deleted
		var deletedWorkout workout.Workout
		err = db.Collection("workout").FindOne(context.Background(), bson.M{"_id": testWorkout.ID}).Decode(&deletedWorkout)
		assert.Error(t, err)
		assert.Equal(t, mongo.ErrNoDocuments, err)
	})

	t.Run("Error - Non-Existing Workout", func(t *testing.T) {
		err := service.DeleteWorkout(primitive.NewObjectID().Hex(), "test_user")
		assert.Error(t, err)
	})

	t.Run("Error - Wrong User ID", func(t *testing.T) {
		testWorkout := &workout.Workout{
			ID:          primitive.NewObjectID(),
			UserID:      "test_user",
			Name:        "Test Workout",
			Description: "Test Description",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		_, err := db.Collection("workout").InsertOne(context.Background(), testWorkout)
		assert.NoError(t, err)

		err = service.DeleteWorkout(testWorkout.ID.Hex(), "wrong_user")
		assert.Error(t, err)
	})

	t.Run("Error - Invalid Workout ID Format", func(t *testing.T) {
		err := service.DeleteWorkout("invalid-id", "test_user")
		assert.Error(t, err)
	})
}

func TestSearchWorkouts(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	service := &workout.WorkoutService{DB: db}

	// Create test workouts
	workouts := []interface{}{
		&workout.Workout{
			ID:          primitive.NewObjectID(),
			UserID:      "test_user",
			Name:        "Upper Body Workout",
			Description: "Chest and Back",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		&workout.Workout{
			ID:          primitive.NewObjectID(),
			UserID:      "test_user",
			Name:        "Lower Body Workout",
			Description: "Legs and Core",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	_, err := db.Collection("workout").InsertMany(context.Background(), workouts)
	assert.NoError(t, err)

	t.Run("Success - Search by Name", func(t *testing.T) {
		results, err := service.SearchWorkouts("Upper", "test_user")
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Contains(t, results[0].Name, "Upper")
	})

	t.Run("Success - Search by Description", func(t *testing.T) {
		results, err := service.SearchWorkouts("Legs", "test_user")
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Contains(t, results[0].Description, "Legs")
	})

	t.Run("Success - No Results Found", func(t *testing.T) {
		results, err := service.SearchWorkouts("NonExistent", "test_user")
		assert.NoError(t, err)
		assert.Len(t, results, 0)
	})

	t.Run("Success - Empty Query Returns All Workouts", func(t *testing.T) {
		results, err := service.SearchWorkouts("", "test_user")
		assert.NoError(t, err)
		assert.Len(t, results, 2)
	})
}
