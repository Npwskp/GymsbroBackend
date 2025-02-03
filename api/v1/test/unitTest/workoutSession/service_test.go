package workoutSession_test

import (
	"context"
	"testing"
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/workout/workoutSession"
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

func TestStartSession(t *testing.T) {
	db := setupTestDB(t)
	service := &workoutSession.WorkoutSessionService{DB: db}

	t.Run("Successfully start custom session", func(t *testing.T) {
		userId := "test_user"
		createDto := &workoutSession.CreateWorkoutSessionDto{
			Type:  workoutSession.CustomSession,
			Notes: "Test workout session",
		}

		result, err := service.StartSession(createDto, userId)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userId, result.UserID)
		assert.Equal(t, workoutSession.CustomSession, result.Type)
		assert.Equal(t, workoutSession.StatusInProgress, result.Status)
		assert.Equal(t, createDto.Notes, result.Notes)
	})

	t.Run("Fail to start session with ongoing session", func(t *testing.T) {
		userId := "test_user"
		existingSession := &workoutSession.WorkoutSession{
			UserID:    userId,
			Type:      workoutSession.CustomSession,
			StartTime: time.Now(),
			Status:    workoutSession.StatusInProgress,
		}

		_, err := db.Collection("workoutSessions").InsertOne(context.Background(), existingSession)
		assert.NoError(t, err)

		createDto := &workoutSession.CreateWorkoutSessionDto{
			Type:  workoutSession.CustomSession,
			Notes: "Test workout session",
		}

		result, err := service.StartSession(createDto, userId)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "already has an ongoing session")
	})
}

func TestEndSession(t *testing.T) {
	db := setupTestDB(t)
	service := &workoutSession.WorkoutSessionService{DB: db}

	t.Run("Successfully end session", func(t *testing.T) {
		userId := "test_user"
		session := &workoutSession.WorkoutSession{
			ID:        primitive.NewObjectID(),
			UserID:    userId,
			Type:      workoutSession.CustomSession,
			StartTime: time.Now().Add(-1 * time.Hour),
			Status:    workoutSession.StatusInProgress,
			Exercises: []workoutSession.SessionExercise{
				{
					ExerciseID: primitive.NewObjectID().Hex(),
					Order:      0,
				},
			},
		}

		_, err := db.Collection("workoutSessions").InsertOne(context.Background(), session)
		assert.NoError(t, err)

		result, err := service.EndSession(session.ID.Hex(), userId)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, workoutSession.StatusCompleted, result.Status)
		assert.NotZero(t, result.Duration)
		assert.NotNil(t, result.EndTime)
	})

	t.Run("Fail to end non-existent session", func(t *testing.T) {
		userId := "test_user"
		nonExistentId := primitive.NewObjectID().Hex()

		result, err := service.EndSession(nonExistentId, userId)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestUpdateSession(t *testing.T) {
	db := setupTestDB(t)
	service := &workoutSession.WorkoutSessionService{DB: db}

	t.Run("Successfully update session", func(t *testing.T) {
		userId := "test_user"
		session := &workoutSession.WorkoutSession{
			ID:        primitive.NewObjectID(),
			UserID:    userId,
			Type:      workoutSession.CustomSession,
			StartTime: time.Now(),
			Status:    workoutSession.StatusInProgress,
		}

		_, err := db.Collection("workoutSessions").InsertOne(context.Background(), session)
		assert.NoError(t, err)

		updateDto := &workoutSession.UpdateWorkoutSessionDto{
			Status: workoutSession.StatusInProgress,
			Exercises: []workoutSession.SessionExercise{
				{
					ExerciseID: primitive.NewObjectID().Hex(),
					Order:      0,
				},
			},
			Notes: "Updated notes",
		}

		result, err := service.UpdateSession(session.ID.Hex(), updateDto, userId)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, updateDto.Notes, result.Notes)
		assert.Equal(t, len(updateDto.Exercises), len(result.Exercises))
	})

	t.Run("Error - Non-Existing Session", func(t *testing.T) {
		updateDto := &workoutSession.UpdateWorkoutSessionDto{
			Status: workoutSession.StatusInProgress,
		}
		result, err := service.UpdateSession(primitive.NewObjectID().Hex(), updateDto, "test_user")
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestGetSession(t *testing.T) {
	db := setupTestDB(t)
	service := &workoutSession.WorkoutSessionService{DB: db}

	t.Run("Successfully get session", func(t *testing.T) {
		userId := "test_user"
		session := &workoutSession.WorkoutSession{
			ID:        primitive.NewObjectID(),
			UserID:    userId,
			Type:      workoutSession.CustomSession,
			StartTime: time.Now(),
			Status:    workoutSession.StatusInProgress,
		}

		_, err := db.Collection("workoutSessions").InsertOne(context.Background(), session)
		assert.NoError(t, err)

		result, err := service.GetSession(session.ID.Hex(), userId)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, session.ID, result.ID)
		assert.Equal(t, session.UserID, result.UserID)
	})
}

func TestGetUserSessions(t *testing.T) {
	db := setupTestDB(t)
	service := &workoutSession.WorkoutSessionService{DB: db}

	t.Run("Successfully get user sessions", func(t *testing.T) {
		userId := "test_user"
		sessions := []interface{}{
			&workoutSession.WorkoutSession{
				ID:        primitive.NewObjectID(),
				UserID:    userId,
				Type:      workoutSession.CustomSession,
				StartTime: time.Now(),
				Status:    workoutSession.StatusCompleted,
			},
			&workoutSession.WorkoutSession{
				ID:        primitive.NewObjectID(),
				UserID:    userId,
				Type:      workoutSession.PlannedSession,
				StartTime: time.Now().Add(-24 * time.Hour),
				Status:    workoutSession.StatusCompleted,
			},
		}

		_, err := db.Collection("workoutSessions").InsertMany(context.Background(), sessions)
		assert.NoError(t, err)

		results, err := service.GetUserSessions(userId)
		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.Equal(t, 2, len(results))
	})
}

func TestGetOnGoingSession(t *testing.T) {
	db := setupTestDB(t)
	service := &workoutSession.WorkoutSessionService{DB: db}

	t.Run("Successfully get ongoing session", func(t *testing.T) {
		userId := "test_user"
		session := &workoutSession.WorkoutSession{
			ID:        primitive.NewObjectID(),
			UserID:    userId,
			Type:      workoutSession.CustomSession,
			StartTime: time.Now(),
			Status:    workoutSession.StatusInProgress,
		}

		_, err := db.Collection("workoutSessions").InsertOne(context.Background(), session)
		assert.NoError(t, err)

		result, err := service.GetOnGoingSession(userId)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, session.ID, result.ID)
		assert.Equal(t, workoutSession.StatusInProgress, result.Status)
	})
}

func TestDeleteSession(t *testing.T) {
	db := setupTestDB(t)
	service := &workoutSession.WorkoutSessionService{DB: db}

	t.Run("Successfully delete session", func(t *testing.T) {
		userId := "test_user"
		session := &workoutSession.WorkoutSession{
			ID:        primitive.NewObjectID(),
			UserID:    userId,
			Type:      workoutSession.CustomSession,
			StartTime: time.Now(),
			Status:    workoutSession.StatusCompleted,
		}

		_, err := db.Collection("workoutSessions").InsertOne(context.Background(), session)
		assert.NoError(t, err)

		err = service.DeleteSession(session.ID.Hex(), userId)
		assert.NoError(t, err)

		// Verify deletion
		var deletedSession workoutSession.WorkoutSession
		err = db.Collection("workoutSessions").FindOne(context.Background(), bson.D{{Key: "_id", Value: session.ID}}).Decode(&deletedSession)
		assert.Equal(t, mongo.ErrNoDocuments, err)
	})
}

func TestLogSession(t *testing.T) {
	db := setupTestDB(t)
	service := &workoutSession.WorkoutSessionService{DB: db}

	t.Run("Successfully log session", func(t *testing.T) {
		userId := "test_user"
		startTime := time.Now().Add(-1 * time.Hour)
		endTime := time.Now()

		logDto := &workoutSession.LoggedSessionDto{
			StartTime: startTime,
			EndTime:   endTime,
			Status:    workoutSession.StatusCompleted,
			Exercises: []workoutSession.SessionExercise{
				{
					ExerciseID: primitive.NewObjectID().Hex(),
					Order:      0,
				},
			},
			Notes: "Logged workout session",
		}

		result, err := service.LogSession(logDto, userId)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userId, result.UserID)
		assert.Equal(t, workoutSession.LoggedSession, result.Type)
		assert.Equal(t, workoutSession.StatusCompleted, result.Status)
		assert.Equal(t, startTime.Unix(), result.StartTime.Unix())
		assert.Equal(t, endTime.Unix(), result.EndTime.Unix())
		assert.Equal(t, logDto.Notes, result.Notes)
		assert.Equal(t, len(logDto.Exercises), len(result.Exercises))
	})
}
