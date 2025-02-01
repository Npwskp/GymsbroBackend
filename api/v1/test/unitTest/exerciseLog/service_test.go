package exerciseLog_test

import (
	"context"
	"testing"
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/workout/exerciseLog"
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

func TestCreateLog(t *testing.T) {
	db := setupTestDB(t)
	service := &exerciseLog.ExerciseLogService{DB: db}

	t.Run("Successfully create exercise log", func(t *testing.T) {
		userId := "test_user"
		now := time.Now()
		createDto := &exerciseLog.CreateExerciseLogDto{
			ExerciseID: primitive.NewObjectID().Hex(),
			DateTime:   now,
			Sets: []exerciseLog.SetLog{
				{
					Weight:    100,
					Reps:      10,
					SetNumber: 1,
					Type:      exerciseLog.WorkingSet,
				},
			},
			Notes: "Test workout",
		}

		result, err := service.CreateLog(createDto, userId)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userId, result.UserID)
		assert.Equal(t, createDto.ExerciseID, result.ExerciseID)
		assert.Equal(t, 1, result.CompletedSets)
		assert.Equal(t, float64(1000), result.TotalVolume)
		assert.Equal(t, createDto.Notes, result.Notes)
		assert.Equal(t, createDto.Sets, result.Sets)
	})
}

func TestGetLogsByUser(t *testing.T) {
	db := setupTestDB(t)
	service := &exerciseLog.ExerciseLogService{DB: db}

	t.Run("Successfully get user logs", func(t *testing.T) {
		userId := "test_user"
		// Create some test logs first
		testLog := &exerciseLog.ExerciseLog{
			UserID:     userId,
			ExerciseID: primitive.NewObjectID().Hex(),
			Sets: []exerciseLog.SetLog{
				{
					Weight:    100,
					Reps:      10,
					SetNumber: 1,
					Type:      exerciseLog.WorkingSet,
				},
			},
			CompletedSets: 1,
			TotalVolume:   1000,
			DateTime:      time.Now(),
		}

		_, err := db.Collection("exerciseLogs").InsertOne(context.Background(), testLog)
		assert.NoError(t, err)

		results, err := service.GetLogsByUser(userId)
		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.Equal(t, 1, len(results))
		assert.Equal(t, userId, results[0].UserID)
	})
}

func TestGetLogsByExercise(t *testing.T) {
	db := setupTestDB(t)
	service := &exerciseLog.ExerciseLogService{DB: db}

	t.Run("Successfully get exercise logs", func(t *testing.T) {
		userId := "test_user"
		exerciseId := primitive.NewObjectID().Hex()

		// Create test logs
		testLog := &exerciseLog.ExerciseLog{
			UserID:     userId,
			ExerciseID: exerciseId,
			Sets: []exerciseLog.SetLog{
				{
					Weight:    100,
					Reps:      10,
					SetNumber: 1,
					Type:      exerciseLog.WorkingSet,
				},
			},
			CompletedSets: 1,
			TotalVolume:   1000,
			DateTime:      time.Now(),
		}

		_, err := db.Collection("exerciseLogs").InsertOne(context.Background(), testLog)
		assert.NoError(t, err)

		results, err := service.GetLogsByExercise(exerciseId, userId)
		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.Equal(t, 1, len(results))
		assert.Equal(t, exerciseId, results[0].ExerciseID)
	})
}

func TestGetLogsByDateRange(t *testing.T) {
	db := setupTestDB(t)
	service := &exerciseLog.ExerciseLogService{DB: db}

	t.Run("Successfully get logs by date range", func(t *testing.T) {
		userId := "test_user"
		now := time.Now()
		startDate := now.Add(-24 * time.Hour)
		endDate := now.Add(24 * time.Hour)

		// Create test logs
		testLog := &exerciseLog.ExerciseLog{
			UserID:     userId,
			ExerciseID: primitive.NewObjectID().Hex(),
			Sets: []exerciseLog.SetLog{
				{
					Weight:    100,
					Reps:      10,
					SetNumber: 1,
					Type:      exerciseLog.WorkingSet,
				},
			},
			CompletedSets: 1,
			TotalVolume:   1000,
			DateTime:      now,
		}

		_, err := db.Collection("exerciseLogs").InsertOne(context.Background(), testLog)
		assert.NoError(t, err)

		results, err := service.GetLogsByDateRange(userId, startDate, endDate)
		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.Equal(t, 1, len(results))
	})
}

func TestUpdateLog(t *testing.T) {
	db := setupTestDB(t)
	service := &exerciseLog.ExerciseLogService{DB: db}

	t.Run("Successfully update log", func(t *testing.T) {
		userId := "test_user"
		// Create a test log first
		testLog := &exerciseLog.ExerciseLog{
			ID:         primitive.NewObjectID(),
			UserID:     userId,
			ExerciseID: primitive.NewObjectID().Hex(),
			Sets: []exerciseLog.SetLog{
				{
					Weight:    100,
					Reps:      10,
					SetNumber: 1,
					Type:      exerciseLog.WorkingSet,
				},
			},
			CompletedSets: 1,
			TotalVolume:   1000,
			DateTime:      time.Now(),
		}

		_, err := db.Collection("exerciseLogs").InsertOne(context.Background(), testLog)
		assert.NoError(t, err)

		updateDto := &exerciseLog.UpdateExerciseLogDto{
			Sets: []exerciseLog.SetLog{
				{
					Weight:    120,
					Reps:      8,
					SetNumber: 1,
					Type:      exerciseLog.WorkingSet,
				},
			},
			DateTime: time.Now(),
			Notes:    "Updated test log",
		}

		result, err := service.UpdateLog(testLog.ID.Hex(), updateDto, userId)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, updateDto.Notes, result.Notes)
		assert.Equal(t, updateDto.Sets, result.Sets)
	})

	t.Run("Fail to update non-existent log", func(t *testing.T) {
		userId := "test_user"
		logId := primitive.NewObjectID().Hex()
		updateDto := &exerciseLog.UpdateExerciseLogDto{
			Sets: []exerciseLog.SetLog{
				{
					Weight:    120,
					Reps:      8,
					SetNumber: 1,
					Type:      exerciseLog.WorkingSet,
				},
			},
			DateTime: time.Now(),
			Notes:    "Updated test log",
		}

		result, err := service.UpdateLog(logId, updateDto, userId)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestDeleteLog(t *testing.T) {
	db := setupTestDB(t)
	service := &exerciseLog.ExerciseLogService{DB: db}

	t.Run("Successfully delete log", func(t *testing.T) {
		userId := "test_user"
		// Create a test log first
		testLog := &exerciseLog.ExerciseLog{
			ID:         primitive.NewObjectID(),
			UserID:     userId,
			ExerciseID: primitive.NewObjectID().Hex(),
			Sets: []exerciseLog.SetLog{
				{
					Weight:    100,
					Reps:      10,
					SetNumber: 1,
					Type:      exerciseLog.WorkingSet,
				},
			},
			CompletedSets: 1,
			TotalVolume:   1000,
			DateTime:      time.Now(),
		}

		_, err := db.Collection("exerciseLogs").InsertOne(context.Background(), testLog)
		assert.NoError(t, err)

		err = service.DeleteLog(testLog.ID.Hex(), userId)
		assert.NoError(t, err)

		// Verify deletion
		var found exerciseLog.ExerciseLog
		err = db.Collection("exerciseLogs").FindOne(context.Background(), bson.M{"_id": testLog.ID}).Decode(&found)
		assert.Error(t, err)
		assert.Equal(t, mongo.ErrNoDocuments, err)
	})

	t.Run("Fail to delete non-existent log", func(t *testing.T) {
		userId := "test_user"
		logId := primitive.NewObjectID().Hex()

		err := service.DeleteLog(logId, userId)
		assert.Error(t, err)
		assert.Equal(t, mongo.ErrNoDocuments, err)
	})
}
