package workoutPlan_test

import (
	"context"
	"testing"
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/workout/workoutPlan"
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

func TestCreatePlanByDaysOfWeek(t *testing.T) {
	db := setupTestDB(t)
	service := &workoutPlan.WorkoutPlanService{DB: db}

	t.Run("Successfully create plan by days of week", func(t *testing.T) {
		userId := "test_user"
		createDto := &workoutPlan.CreatePlanByDaysOfWeekDto{
			MondayWorkoutID:    primitive.NewObjectID().Hex(),
			TuesdayWorkoutID:   primitive.NewObjectID().Hex(),
			WednesdayWorkoutID: primitive.NewObjectID().Hex(),
			ThursdayWorkoutID:  primitive.NewObjectID().Hex(),
			FridayWorkoutID:    primitive.NewObjectID().Hex(),
			SaturdayWorkoutID:  primitive.NewObjectID().Hex(),
			SundayWorkoutID:    primitive.NewObjectID().Hex(),
			WeeksDuration:      2,
		}

		results, err := service.CreatePlanByDaysOfWeek(createDto, userId)
		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.Equal(t, 7, len(results)) // One plan for each day of the week

		// Verify each plan
		for _, plan := range results {
			assert.Equal(t, userId, plan.UserID)
			assert.NotEmpty(t, plan.WorkoutID)
			assert.Equal(t, 2, len(plan.Dates)) // 2 dates for each workout (2 weeks)
			assert.False(t, plan.CreatedAt.IsZero())
			assert.False(t, plan.UpdatedAt.IsZero())
		}
	})

	t.Run("Successfully update existing plans", func(t *testing.T) {
		userId := "test_user"
		workoutId := primitive.NewObjectID().Hex()
		existingPlan := workoutPlan.WorkoutPlan{
			ID:        primitive.NewObjectID(),
			UserID:    userId,
			WorkoutID: workoutId,
			Dates:     []time.Time{time.Now().Add(24 * time.Hour)},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err := db.Collection("workoutPlan").InsertOne(context.Background(), existingPlan)
		assert.NoError(t, err)

		createDto := &workoutPlan.CreatePlanByDaysOfWeekDto{
			MondayWorkoutID:    workoutId,
			TuesdayWorkoutID:   primitive.NewObjectID().Hex(),
			WednesdayWorkoutID: primitive.NewObjectID().Hex(),
			ThursdayWorkoutID:  primitive.NewObjectID().Hex(),
			FridayWorkoutID:    primitive.NewObjectID().Hex(),
			SaturdayWorkoutID:  primitive.NewObjectID().Hex(),
			SundayWorkoutID:    primitive.NewObjectID().Hex(),
			WeeksDuration:      1,
		}

		results, err := service.CreatePlanByDaysOfWeek(createDto, userId)
		assert.NoError(t, err)
		assert.NotNil(t, results)

		// Find the updated Monday plan
		var mondayPlan workoutPlan.WorkoutPlan
		for _, plan := range results {
			if plan.WorkoutID == workoutId {
				mondayPlan = plan
				break
			}
		}

		assert.Equal(t, existingPlan.ID, mondayPlan.ID)
		assert.Greater(t, len(mondayPlan.Dates), len(existingPlan.Dates))
	})
}

func TestCreatePlanByCyclicWorkout(t *testing.T) {
	db := setupTestDB(t)
	service := &workoutPlan.WorkoutPlanService{DB: db}

	t.Run("Successfully create cyclic workout plan", func(t *testing.T) {
		userId := "test_user"
		workoutIds := []string{
			primitive.NewObjectID().Hex(),
			primitive.NewObjectID().Hex(),
			primitive.NewObjectID().Hex(),
		}

		createDto := &workoutPlan.CreatePlanByCyclicWorkoutDto{
			WorkoutIDs:    workoutIds,
			WeeksDuration: 2,
		}

		results, err := service.CreatePlanByCyclicWorkout(createDto, userId)
		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.Equal(t, len(workoutIds), len(results))

		// Verify each plan
		for _, plan := range results {
			assert.Equal(t, userId, plan.UserID)
			assert.Contains(t, workoutIds, plan.WorkoutID)
			assert.NotEmpty(t, plan.Dates)
			assert.False(t, plan.CreatedAt.IsZero())
			assert.False(t, plan.UpdatedAt.IsZero())
		}
	})

	t.Run("Successfully update existing cyclic plans", func(t *testing.T) {
		userId := "test_user"
		workoutId := primitive.NewObjectID().Hex()
		existingPlan := workoutPlan.WorkoutPlan{
			ID:        primitive.NewObjectID(),
			UserID:    userId,
			WorkoutID: workoutId,
			Dates:     []time.Time{time.Now().Add(24 * time.Hour)},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err := db.Collection("workoutPlan").InsertOne(context.Background(), existingPlan)
		assert.NoError(t, err)

		createDto := &workoutPlan.CreatePlanByCyclicWorkoutDto{
			WorkoutIDs:    []string{workoutId, primitive.NewObjectID().Hex()},
			WeeksDuration: 1,
		}

		results, err := service.CreatePlanByCyclicWorkout(createDto, userId)
		assert.NoError(t, err)
		assert.NotNil(t, results)

		// Find the updated plan
		var updatedPlan workoutPlan.WorkoutPlan
		for _, plan := range results {
			if plan.WorkoutID == workoutId {
				updatedPlan = plan
				break
			}
		}

		assert.Equal(t, existingPlan.ID, updatedPlan.ID)
		assert.Greater(t, len(updatedPlan.Dates), len(existingPlan.Dates))
	})
}

func TestGetWorkoutPlansByUser(t *testing.T) {
	db := setupTestDB(t)
	service := &workoutPlan.WorkoutPlanService{DB: db}

	t.Run("Successfully get user workout plans", func(t *testing.T) {
		userId := "test_user"
		plans := []interface{}{
			workoutPlan.WorkoutPlan{
				ID:        primitive.NewObjectID(),
				UserID:    userId,
				WorkoutID: primitive.NewObjectID().Hex(),
				Dates:     []time.Time{time.Now().Add(24 * time.Hour)},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			workoutPlan.WorkoutPlan{
				ID:        primitive.NewObjectID(),
				UserID:    userId,
				WorkoutID: primitive.NewObjectID().Hex(),
				Dates:     []time.Time{time.Now().Add(48 * time.Hour)},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		_, err := db.Collection("workoutPlan").InsertMany(context.Background(), plans)
		assert.NoError(t, err)

		results, err := service.GetWorkoutPlansByUser(userId)
		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.Equal(t, 2, len(results))

		// Verify each plan belongs to the user
		for _, plan := range results {
			assert.Equal(t, userId, plan.UserID)
			assert.NotEmpty(t, plan.WorkoutID)
			assert.NotEmpty(t, plan.Dates)
		}
	})

	t.Run("Return empty list for user with no plans", func(t *testing.T) {
		userId := "user_with_no_plans"
		results, err := service.GetWorkoutPlansByUser(userId)
		assert.NoError(t, err)
		assert.Empty(t, results)
	})
}
