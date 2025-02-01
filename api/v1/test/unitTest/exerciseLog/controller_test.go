package exerciseLog_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/workout/exerciseLog"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock service
type MockExerciseLogService struct {
	mock.Mock
}

func (m *MockExerciseLogService) CreateLog(log *exerciseLog.CreateExerciseLogDto, userId string) (*exerciseLog.ExerciseLog, error) {
	args := m.Called(log, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*exerciseLog.ExerciseLog), args.Error(1)
}

func (m *MockExerciseLogService) GetLogsByUser(userId string) ([]*exerciseLog.ExerciseLog, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*exerciseLog.ExerciseLog), args.Error(1)
}

func (m *MockExerciseLogService) GetLogsByExercise(exerciseId string, userId string) ([]*exerciseLog.ExerciseLog, error) {
	args := m.Called(exerciseId, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*exerciseLog.ExerciseLog), args.Error(1)
}

func (m *MockExerciseLogService) GetLogsByDateRange(userId string, startDate, endDate time.Time) ([]*exerciseLog.ExerciseLog, error) {
	args := m.Called(userId, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*exerciseLog.ExerciseLog), args.Error(1)
}

func (m *MockExerciseLogService) UpdateLog(id string, log *exerciseLog.UpdateExerciseLogDto, userId string) (*exerciseLog.ExerciseLog, error) {
	args := m.Called(id, log, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*exerciseLog.ExerciseLog), args.Error(1)
}

func (m *MockExerciseLogService) DeleteLog(id string, userId string) error {
	args := m.Called(id, userId)
	return args.Error(0)
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
func setupTest() (*fiber.App, *MockExerciseLogService) {
	app := fiber.New()
	api := app.Group("/api/v1", testMiddleware())

	mockService := new(MockExerciseLogService)
	controller := &exerciseLog.ExerciseLogController{
		Instance: api,
		Service:  mockService,
	}
	controller.Handle()
	return app, mockService
}

func TestCreateLogHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Valid exercise log creation", func(t *testing.T) {
		now := time.Now().Truncate(time.Second)
		createDto := exerciseLog.CreateExerciseLogDto{
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

		expectedResponse := &exerciseLog.ExerciseLog{
			ID:            primitive.NewObjectID(),
			UserID:        "test_user",
			ExerciseID:    createDto.ExerciseID,
			CompletedSets: 1,
			TotalVolume:   1000,
			Notes:         createDto.Notes,
			DateTime:      now,
			Sets:          createDto.Sets,
		}

		mockService.On("CreateLog",
			mock.MatchedBy(func(dto *exerciseLog.CreateExerciseLogDto) bool {
				return dto.ExerciseID == createDto.ExerciseID &&
					dto.Notes == createDto.Notes &&
					len(dto.Sets) == len(createDto.Sets) &&
					dto.DateTime.Unix() == createDto.DateTime.Unix()
			}),
			"test_user",
		).Return(expectedResponse, nil)

		body, _ := json.Marshal(createDto)
		req := httptest.NewRequest("POST", "/api/v1/exercise-log", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		var result exerciseLog.ExerciseLog
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, expectedResponse.Notes, result.Notes)
	})
}

func TestGetUserLogsHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Get user logs successfully", func(t *testing.T) {
		expectedLogs := []*exerciseLog.ExerciseLog{
			{
				ID:            primitive.NewObjectID(),
				UserID:        "test_user",
				ExerciseID:    primitive.NewObjectID().Hex(),
				CompletedSets: 1,
				TotalVolume:   1000,
				Notes:         "Test log 1",
			},
			{
				ID:            primitive.NewObjectID(),
				UserID:        "test_user",
				ExerciseID:    primitive.NewObjectID().Hex(),
				CompletedSets: 2,
				TotalVolume:   2000,
				Notes:         "Test log 2",
			},
		}

		mockService.On("GetLogsByUser", "test_user").Return(expectedLogs, nil)

		req := httptest.NewRequest("GET", "/api/v1/exercise-log", nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result []*exerciseLog.ExerciseLog
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, len(expectedLogs), len(result))
	})
}

func TestGetExerciseLogsHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Get exercise logs successfully", func(t *testing.T) {
		exerciseId := primitive.NewObjectID().Hex()
		expectedLogs := []*exerciseLog.ExerciseLog{
			{
				ID:            primitive.NewObjectID(),
				UserID:        "test_user",
				ExerciseID:    exerciseId,
				CompletedSets: 1,
				TotalVolume:   1000,
				Notes:         "Test log 1",
			},
		}

		mockService.On("GetLogsByExercise", exerciseId, "test_user").Return(expectedLogs, nil)

		req := httptest.NewRequest("GET", "/api/v1/exercise-log/exercise/"+exerciseId, nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result []*exerciseLog.ExerciseLog
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, len(expectedLogs), len(result))
	})
}

func TestGetLogsByDateRangeHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Get logs by date range successfully", func(t *testing.T) {
		startDate := time.Now().Add(-24 * time.Hour).UTC().Truncate(time.Second)
		endDate := time.Now().UTC().Truncate(time.Second)

		expectedLogs := []*exerciseLog.ExerciseLog{
			{
				ID:            primitive.NewObjectID(),
				UserID:        "test_user",
				ExerciseID:    primitive.NewObjectID().Hex(),
				CompletedSets: 1,
				TotalVolume:   1000,
				Notes:         "Test log 1",
				DateTime:      startDate.Add(12 * time.Hour),
			},
		}

		mockService.On("GetLogsByDateRange",
			"test_user",
			mock.MatchedBy(func(t time.Time) bool {
				// Compare timestamps in UTC to avoid timezone issues
				return t.UTC().Unix() >= startDate.Unix()-1 && t.UTC().Unix() <= startDate.Unix()+1
			}),
			mock.MatchedBy(func(t time.Time) bool {
				return t.UTC().Unix() >= endDate.Unix()-1 && t.UTC().Unix() <= endDate.Unix()+1
			}),
		).Return(expectedLogs, nil)

		url := fmt.Sprintf("/api/v1/exercise-log/range?startDate=%s&endDate=%s",
			url.QueryEscape(startDate.Format("2006-01-02 15:04:05")),
			url.QueryEscape(endDate.Format("2006-01-02 15:04:05")))

		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result []*exerciseLog.ExerciseLog
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, len(expectedLogs), len(result))
	})
}

func TestUpdateLogHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Update log successfully", func(t *testing.T) {
		logId := primitive.NewObjectID().Hex()
		now := time.Now().Truncate(time.Second)
		updateDto := exerciseLog.UpdateExerciseLogDto{
			Sets: []exerciseLog.SetLog{
				{
					Weight:    120,
					Reps:      8,
					SetNumber: 1,
					Type:      exerciseLog.WorkingSet,
				},
			},
			DateTime: now,
			Notes:    "Updated test log",
		}

		expectedResponse := &exerciseLog.ExerciseLog{
			ID:            primitive.NewObjectID(),
			UserID:        "test_user",
			ExerciseID:    primitive.NewObjectID().Hex(),
			CompletedSets: 1,
			TotalVolume:   960,
			Notes:         updateDto.Notes,
			DateTime:      now,
			Sets:          updateDto.Sets,
		}

		mockService.On("UpdateLog",
			logId,
			mock.MatchedBy(func(dto *exerciseLog.UpdateExerciseLogDto) bool {
				return dto.Notes == updateDto.Notes &&
					len(dto.Sets) == len(updateDto.Sets) &&
					dto.DateTime.Unix() == updateDto.DateTime.Unix()
			}),
			"test_user",
		).Return(expectedResponse, nil)

		body, _ := json.Marshal(updateDto)
		req := httptest.NewRequest("PUT", "/api/v1/exercise-log/"+logId, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result exerciseLog.ExerciseLog
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, expectedResponse.Notes, result.Notes)
	})
}

func TestDeleteLogHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Delete log successfully", func(t *testing.T) {
		logId := primitive.NewObjectID().Hex()

		mockService.On("DeleteLog", logId, "test_user").Return(nil)

		req := httptest.NewRequest("DELETE", "/api/v1/exercise-log/"+logId, nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	})
}
