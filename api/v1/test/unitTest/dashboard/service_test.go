package dashboard_test

import (
	"context"
	"testing"
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/dashboard"
	userFitnessPreferenceEnums "github.com/Npwskp/GymsbroBackend/api/v1/user/enums"
	bodyCompositionLog "github.com/Npwskp/GymsbroBackend/api/v1/userLog/userBodyComposition"
	"github.com/Npwskp/GymsbroBackend/api/v1/workout/exercise"
	"github.com/Npwskp/GymsbroBackend/api/v1/workout/exerciseLog"
	"github.com/Npwskp/GymsbroBackend/api/v1/workout/workoutSession"
	"github.com/stretchr/testify/assert"
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

func TestGetDashboard(t *testing.T) {
	db := setupTestDB(t)
	service := &dashboard.DashboardService{DB: db}

	t.Run("Successfully get dashboard data", func(t *testing.T) {
		userId := primitive.NewObjectID().Hex()
		now := time.Now()
		startDate := now.AddDate(0, 0, -7)
		endDate := now

		// Create test workout sessions
		sessions := []interface{}{
			&workoutSession.WorkoutSession{
				ID:          primitive.NewObjectID(),
				UserID:      userId,
				StartTime:   now.AddDate(0, 0, -1),
				EndTime:     now,
				Status:      workoutSession.StatusCompleted,
				TotalVolume: 1000,
				Duration:    3600,
			},
			&workoutSession.WorkoutSession{
				ID:          primitive.NewObjectID(),
				UserID:      userId,
				StartTime:   now.AddDate(0, 0, -2),
				EndTime:     now.AddDate(0, 0, -1),
				Status:      workoutSession.StatusCompleted,
				TotalVolume: 1200,
				Duration:    4000,
			},
		}

		_, err := db.Collection("workoutSessions").InsertMany(context.Background(), sessions)
		assert.NoError(t, err)

		// Create test exercise logs
		exerciseId := primitive.NewObjectID()
		exerciseLogs := []interface{}{
			&exerciseLog.ExerciseLog{
				ID:          primitive.NewObjectID(),
				UserID:      userId,
				ExerciseID:  exerciseId.Hex(),
				DateTime:    now.AddDate(0, 0, -1),
				TotalVolume: 500,
				CreatedAt:   now.AddDate(0, 0, -1),
			},
			&exerciseLog.ExerciseLog{
				ID:          primitive.NewObjectID(),
				UserID:      userId,
				ExerciseID:  exerciseId.Hex(),
				DateTime:    now.AddDate(0, 0, -2),
				TotalVolume: 600,
				CreatedAt:   now.AddDate(0, 0, -2),
			},
		}

		_, err = db.Collection("exerciseLogs").InsertMany(context.Background(), exerciseLogs)
		assert.NoError(t, err)

		// Create test exercise
		testExercise := &exercise.Exercise{
			ID:   exerciseId,
			Name: "Test Exercise",
		}
		_, err = db.Collection("exercises").InsertOne(context.Background(), testExercise)
		assert.NoError(t, err)

		result, err := service.GetDashboard(userId, startDate, endDate)
		assert.NoError(t, err)
		assert.NotNil(t, result)

		// Verify dashboard data
		assert.Equal(t, 2, result.Analysis.TotalWorkouts)
		assert.Equal(t, 2, result.Analysis.TotalExercises)
		assert.Equal(t, float64(3300), result.Analysis.TotalVolume)
		assert.Equal(t, float64(3800), result.Analysis.AverageWorkoutDuration)
		assert.Equal(t, 8, len(result.FrequencyGraph.Labels))
		assert.Equal(t, 8, len(result.FrequencyGraph.Values))
		assert.Equal(t, 8, len(result.FrequencyGraph.TrendLine))
	})
}
func TestGetRepMax(t *testing.T) {
	db := setupTestDB(t)
	service := &dashboard.DashboardService{DB: db}

	t.Run("Successfully get rep max", func(t *testing.T) {
		userId := "test_user"
		exerciseId := primitive.NewObjectID().Hex()
		now := time.Now()

		// Create test exercise logs
		logs := []interface{}{
			&exerciseLog.ExerciseLog{
				UserID:     userId,
				ExerciseID: exerciseId,
				DateTime:   now.AddDate(0, 0, -1),
				CreatedAt:  now.AddDate(0, 0, -1),
				Sets: []exerciseLog.SetLog{
					{
						Weight: 100,
						Reps:   10,
						Type:   exerciseLog.WorkingSet,
					},
				},
			},
			&exerciseLog.ExerciseLog{
				UserID:     userId,
				ExerciseID: exerciseId,
				DateTime:   now,
				CreatedAt:  now,
				Sets: []exerciseLog.SetLog{
					{
						Weight: 120,
						Reps:   8,
						Type:   exerciseLog.WorkingSet,
					},
				},
			},
		}

		_, err := db.Collection("exerciseLogs").InsertMany(context.Background(), logs)
		assert.NoError(t, err)

		result, err := service.GetRepMax(userId, exerciseId, false)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotZero(t, result.OneRepMax)
		assert.NotZero(t, result.EightRepMax)
		assert.NotZero(t, result.TwelveRepMax)
	})
}
func TestGetBodyCompositionAnalysis(t *testing.T) {
	db := setupTestDB(t)
	service := &dashboard.DashboardService{DB: db}

	t.Run("Successfully get body composition analysis", func(t *testing.T) {
		userId := "test_user"
		now := time.Now()
		startDate := now.AddDate(0, 0, -7)
		endDate := now

		// Create test body composition logs
		logs := []interface{}{
			&bodyCompositionLog.UserBodyCompositionLog{
				ID:     primitive.NewObjectID(),
				UserID: userId,
				Weight: 70.0,
				BodyComposition: userFitnessPreferenceEnums.BodyCompositionInfo{
					BMI:                22.5,
					BodyFatMass:        12.0,
					BodyFatPercentage:  17.0,
					SkeletalMuscleMass: 35.0,
					ExtracellularWater: 20.0,
					ECWRatio:           0.38,
				},
				CreatedAt: now.AddDate(0, 0, -1),
				UpdatedAt: now.AddDate(0, 0, -1),
			},
			&bodyCompositionLog.UserBodyCompositionLog{
				ID:     primitive.NewObjectID(),
				UserID: userId,
				Weight: 71.0,
				BodyComposition: userFitnessPreferenceEnums.BodyCompositionInfo{
					BMI:                22.8,
					BodyFatMass:        12.5,
					BodyFatPercentage:  17.5,
					SkeletalMuscleMass: 35.2,
					ExtracellularWater: 20.2,
					ECWRatio:           0.385,
				},
				CreatedAt: now.AddDate(0, 0, -2),
				UpdatedAt: now.AddDate(0, 0, -2),
			},
		}

		_, err := db.Collection("bodyCompositionLog").InsertMany(context.Background(), logs)
		assert.NoError(t, err)

		result, err := service.GetBodyCompositionAnalysis(userId, startDate, endDate)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.Labels)
		assert.NotEmpty(t, result.Data)
		assert.NotEmpty(t, result.Changes)
	})
}

// TODO: Add test for nutritionSummary and stregthStandard
