package workoutSession_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/workout/workoutSession"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock service
type MockWorkoutSessionService struct {
	mock.Mock
}

func (m *MockWorkoutSessionService) StartSession(dto *workoutSession.CreateWorkoutSessionDto, userId string) (*workoutSession.WorkoutSession, error) {
	args := m.Called(dto, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*workoutSession.WorkoutSession), args.Error(1)
}

func (m *MockWorkoutSessionService) EndSession(id string, userId string) (*workoutSession.WorkoutSession, error) {
	args := m.Called(id, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*workoutSession.WorkoutSession), args.Error(1)
}

func (m *MockWorkoutSessionService) UpdateSession(id string, dto *workoutSession.UpdateWorkoutSessionDto, userId string) (*workoutSession.WorkoutSession, error) {
	args := m.Called(id, dto, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*workoutSession.WorkoutSession), args.Error(1)
}

func (m *MockWorkoutSessionService) GetSession(id string, userId string) (*workoutSession.WorkoutSession, error) {
	args := m.Called(id, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*workoutSession.WorkoutSession), args.Error(1)
}

func (m *MockWorkoutSessionService) GetUserSessions(userId string) ([]*workoutSession.WorkoutSession, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*workoutSession.WorkoutSession), args.Error(1)
}

func (m *MockWorkoutSessionService) GetOnGoingSession(userId string) (*workoutSession.WorkoutSession, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*workoutSession.WorkoutSession), args.Error(1)
}

func (m *MockWorkoutSessionService) DeleteSession(id string, userId string) error {
	args := m.Called(id, userId)
	return args.Error(0)
}

func (m *MockWorkoutSessionService) LogSession(dto *workoutSession.LoggedSessionDto, userId string) (*workoutSession.WorkoutSession, error) {
	args := m.Called(dto, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*workoutSession.WorkoutSession), args.Error(1)
}

// TestMiddleware sets up the test context with a mock user
func testMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := jwt.MapClaims{
			"sub": c.Get("userid", ""),
		}
		token := &jwt.Token{
			Claims: claims,
		}
		c.Locals("user", token)
		return c.Next()
	}
}

// Test setup helper
func setupTest() (*fiber.App, *MockWorkoutSessionService) {
	app := fiber.New()
	api := app.Group("/api/v1", testMiddleware())

	mockService := new(MockWorkoutSessionService)
	controller := &workoutSession.WorkoutSessionController{
		Instance: api,
		Service:  mockService,
	}
	controller.Handle()
	return app, mockService
}

func TestStartSessionHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Successfully start session", func(t *testing.T) {
		createDto := workoutSession.CreateWorkoutSessionDto{
			Type:  workoutSession.CustomSession,
			Notes: "Test workout session",
		}

		expectedResponse := &workoutSession.WorkoutSession{
			ID:        primitive.NewObjectID(),
			UserID:    "test_user",
			Type:      workoutSession.CustomSession,
			StartTime: time.Now(),
			Status:    workoutSession.StatusInProgress,
			Notes:     createDto.Notes,
		}

		mockService.On("StartSession",
			mock.MatchedBy(func(dto *workoutSession.CreateWorkoutSessionDto) bool {
				return dto.Type == createDto.Type &&
					dto.Notes == createDto.Notes
			}),
			"test_user",
		).Return(expectedResponse, nil)

		body, _ := json.Marshal(createDto)
		req := httptest.NewRequest("POST", "/api/v1/workout-session", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		var result workoutSession.WorkoutSession
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, expectedResponse.Type, result.Type)
		assert.Equal(t, expectedResponse.Notes, result.Notes)
	})
}

func TestEndSessionHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Successfully end session", func(t *testing.T) {
		sessionId := primitive.NewObjectID()
		expectedResponse := &workoutSession.WorkoutSession{
			ID:        sessionId,
			UserID:    "test_user",
			Type:      workoutSession.CustomSession,
			StartTime: time.Now().Add(-1 * time.Hour),
			EndTime:   time.Now(),
			Status:    workoutSession.StatusCompleted,
		}

		mockService.On("EndSession", sessionId.Hex(), "test_user").Return(expectedResponse, nil)

		req := httptest.NewRequest("PUT", "/api/v1/workout-session/"+sessionId.Hex()+"/end", nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result workoutSession.WorkoutSession
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, workoutSession.StatusCompleted, result.Status)
	})
}

func TestUpdateSessionHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Successfully update session", func(t *testing.T) {
		sessionId := primitive.NewObjectID()
		updateDto := workoutSession.UpdateWorkoutSessionDto{
			Status: workoutSession.StatusInProgress,
			Exercises: []workoutSession.SessionExercise{
				{
					ExerciseID: primitive.NewObjectID().Hex(),
					Order:      0,
				},
			},
			Notes: "Updated notes",
		}

		expectedResponse := &workoutSession.WorkoutSession{
			ID:        sessionId,
			UserID:    "test_user",
			Type:      workoutSession.CustomSession,
			Status:    updateDto.Status,
			Exercises: updateDto.Exercises,
			Notes:     updateDto.Notes,
		}

		mockService.On("UpdateSession", sessionId.Hex(),
			mock.MatchedBy(func(dto *workoutSession.UpdateWorkoutSessionDto) bool {
				return dto.Status == updateDto.Status &&
					dto.Notes == updateDto.Notes &&
					len(dto.Exercises) == len(updateDto.Exercises)
			}),
			"test_user",
		).Return(expectedResponse, nil)

		body, _ := json.Marshal(updateDto)
		req := httptest.NewRequest("PUT", "/api/v1/workout-session/"+sessionId.Hex(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result workoutSession.WorkoutSession
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, expectedResponse.Notes, result.Notes)
		assert.Equal(t, len(expectedResponse.Exercises), len(result.Exercises))
	})
}

func TestGetSessionHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Successfully get session", func(t *testing.T) {
		sessionId := primitive.NewObjectID()
		expectedResponse := &workoutSession.WorkoutSession{
			ID:        sessionId,
			UserID:    "test_user",
			Type:      workoutSession.CustomSession,
			StartTime: time.Now(),
			Status:    workoutSession.StatusInProgress,
		}

		mockService.On("GetSession", sessionId.Hex(), "test_user").Return(expectedResponse, nil)

		req := httptest.NewRequest("GET", "/api/v1/workout-session/"+sessionId.Hex(), nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result workoutSession.WorkoutSession
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, expectedResponse.ID, result.ID)
		assert.Equal(t, expectedResponse.UserID, result.UserID)
	})
}

func TestGetUserSessionsHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Successfully get user sessions", func(t *testing.T) {
		expectedResponse := []*workoutSession.WorkoutSession{
			{
				ID:        primitive.NewObjectID(),
				UserID:    "test_user",
				Type:      workoutSession.CustomSession,
				StartTime: time.Now(),
				Status:    workoutSession.StatusCompleted,
			},
			{
				ID:        primitive.NewObjectID(),
				UserID:    "test_user",
				Type:      workoutSession.PlannedSession,
				StartTime: time.Now().Add(-24 * time.Hour),
				Status:    workoutSession.StatusCompleted,
			},
		}

		mockService.On("GetUserSessions", "test_user").Return(expectedResponse, nil)

		req := httptest.NewRequest("GET", "/api/v1/workout-session", nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result []*workoutSession.WorkoutSession
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, len(expectedResponse), len(result))
	})
}

func TestGetOnGoingSessionHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Successfully get ongoing session", func(t *testing.T) {
		expectedResponse := &workoutSession.WorkoutSession{
			ID:        primitive.NewObjectID(),
			UserID:    "test_user",
			Type:      workoutSession.CustomSession,
			StartTime: time.Now(),
			Status:    workoutSession.StatusInProgress,
		}

		mockService.On("GetOnGoingSession", "test_user").Return(expectedResponse, nil)

		req := httptest.NewRequest("GET", "/api/v1/workout-session/ongoing", nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result workoutSession.WorkoutSession
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, expectedResponse.ID, result.ID)
		assert.Equal(t, workoutSession.StatusInProgress, result.Status)
	})
}

func TestDeleteSessionHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Successfully delete session", func(t *testing.T) {
		sessionId := primitive.NewObjectID()
		mockService.On("DeleteSession", sessionId.Hex(), "test_user").Return(nil)

		req := httptest.NewRequest("DELETE", "/api/v1/workout-session/"+sessionId.Hex(), nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	})
}

func TestLogSessionHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Successfully log session", func(t *testing.T) {
		startTime := time.Now().Add(-1 * time.Hour)
		endTime := time.Now()
		logDto := workoutSession.LoggedSessionDto{
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

		expectedResponse := &workoutSession.WorkoutSession{
			ID:        primitive.NewObjectID(),
			UserID:    "test_user",
			Type:      workoutSession.LoggedSession,
			StartTime: startTime,
			EndTime:   endTime,
			Status:    workoutSession.StatusCompleted,
			Exercises: logDto.Exercises,
			Notes:     logDto.Notes,
		}

		mockService.On("LogSession",
			mock.MatchedBy(func(dto *workoutSession.LoggedSessionDto) bool {
				return dto.StartTime.Unix() == logDto.StartTime.Unix() &&
					dto.EndTime.Unix() == logDto.EndTime.Unix() &&
					dto.Status == logDto.Status &&
					len(dto.Exercises) == len(logDto.Exercises) &&
					dto.Notes == logDto.Notes
			}),
			"test_user",
		).Return(expectedResponse, nil)

		body, _ := json.Marshal(logDto)
		req := httptest.NewRequest("POST", "/api/v1/workout-session/log", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		var result workoutSession.WorkoutSession
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, expectedResponse.Type, result.Type)
		assert.Equal(t, expectedResponse.Status, result.Status)
		assert.Equal(t, expectedResponse.Notes, result.Notes)
		assert.Equal(t, len(expectedResponse.Exercises), len(result.Exercises))
	})
}
