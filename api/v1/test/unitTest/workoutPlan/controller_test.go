package workoutPlan_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/workout/workoutPlan"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock service
type MockWorkoutPlanService struct {
	mock.Mock
}

func (m *MockWorkoutPlanService) CreatePlanByDaysOfWeek(dto *workoutPlan.CreatePlanByDaysOfWeekDto, userId string) ([]workoutPlan.WorkoutPlan, error) {
	args := m.Called(dto, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]workoutPlan.WorkoutPlan), args.Error(1)
}

func (m *MockWorkoutPlanService) CreatePlanByCyclicWorkout(dto *workoutPlan.CreatePlanByCyclicWorkoutDto, userId string) ([]workoutPlan.WorkoutPlan, error) {
	args := m.Called(dto, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]workoutPlan.WorkoutPlan), args.Error(1)
}

func (m *MockWorkoutPlanService) GetWorkoutPlansByUser(userId string) ([]workoutPlan.WorkoutPlan, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]workoutPlan.WorkoutPlan), args.Error(1)
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
func setupTest() (*fiber.App, *MockWorkoutPlanService) {
	app := fiber.New()
	api := app.Group("/api/v1", testMiddleware())

	mockService := new(MockWorkoutPlanService)
	controller := &workoutPlan.WorkoutPlanController{
		Instance: api,
		Service:  mockService,
	}
	controller.Handle()
	return app, mockService
}

func TestCreatePlanByDaysOfWeekHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Successfully create plan by days of week", func(t *testing.T) {
		createDto := workoutPlan.CreatePlanByDaysOfWeekDto{
			MondayWorkoutID:    primitive.NewObjectID().Hex(),
			TuesdayWorkoutID:   primitive.NewObjectID().Hex(),
			WednesdayWorkoutID: primitive.NewObjectID().Hex(),
			ThursdayWorkoutID:  primitive.NewObjectID().Hex(),
			FridayWorkoutID:    primitive.NewObjectID().Hex(),
			SaturdayWorkoutID:  primitive.NewObjectID().Hex(),
			SundayWorkoutID:    primitive.NewObjectID().Hex(),
			WeeksDuration:      2,
		}

		expectedResponse := []workoutPlan.WorkoutPlan{
			{
				ID:        primitive.NewObjectID(),
				UserID:    "test_user",
				WorkoutID: createDto.MondayWorkoutID,
				Dates:     []time.Time{time.Now(), time.Now().AddDate(0, 0, 7)},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			// More plans for other days...
		}

		mockService.On("CreatePlanByDaysOfWeek",
			mock.MatchedBy(func(dto *workoutPlan.CreatePlanByDaysOfWeekDto) bool {
				return dto.MondayWorkoutID == createDto.MondayWorkoutID &&
					dto.WeeksDuration == createDto.WeeksDuration
			}),
			"test_user",
		).Return(expectedResponse, nil)

		body, _ := json.Marshal(createDto)
		req := httptest.NewRequest("POST", "/api/v1/workoutPlan/byDaysOfWeek", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		var result []workoutPlan.WorkoutPlan
		json.NewDecoder(resp.Body).Decode(&result)
		assert.NotEmpty(t, result)
		assert.Equal(t, expectedResponse[0].UserID, result[0].UserID)
		assert.Equal(t, expectedResponse[0].WorkoutID, result[0].WorkoutID)
	})

	t.Run("Invalid request - missing required fields", func(t *testing.T) {
		createDto := workoutPlan.CreatePlanByDaysOfWeekDto{
			// Missing required fields
			WeeksDuration: 2,
		}

		body, _ := json.Marshal(createDto)
		req := httptest.NewRequest("POST", "/api/v1/workoutPlan/byDaysOfWeek", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

func TestCreatePlanByCyclicWorkoutHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Successfully create cyclic workout plan", func(t *testing.T) {
		workoutIds := []string{
			primitive.NewObjectID().Hex(),
			primitive.NewObjectID().Hex(),
			primitive.NewObjectID().Hex(),
		}

		createDto := workoutPlan.CreatePlanByCyclicWorkoutDto{
			WorkoutIDs:    workoutIds,
			WeeksDuration: 2,
		}

		expectedResponse := []workoutPlan.WorkoutPlan{
			{
				ID:        primitive.NewObjectID(),
				UserID:    "test_user",
				WorkoutID: workoutIds[0],
				Dates:     []time.Time{time.Now(), time.Now().AddDate(0, 0, 3)},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			// More plans for other workouts...
		}

		mockService.On("CreatePlanByCyclicWorkout",
			mock.MatchedBy(func(dto *workoutPlan.CreatePlanByCyclicWorkoutDto) bool {
				return len(dto.WorkoutIDs) == len(createDto.WorkoutIDs) &&
					dto.WeeksDuration == createDto.WeeksDuration
			}),
			"test_user",
		).Return(expectedResponse, nil)

		body, _ := json.Marshal(createDto)
		req := httptest.NewRequest("POST", "/api/v1/workoutPlan/byCyclicWorkout", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		var result []workoutPlan.WorkoutPlan
		json.NewDecoder(resp.Body).Decode(&result)
		assert.NotEmpty(t, result)
		assert.Equal(t, expectedResponse[0].UserID, result[0].UserID)
		assert.Equal(t, expectedResponse[0].WorkoutID, result[0].WorkoutID)
	})

	t.Run("Invalid request - missing required fields", func(t *testing.T) {
		createDto := workoutPlan.CreatePlanByCyclicWorkoutDto{
			// Missing required fields
			WeeksDuration: 2,
		}

		body, _ := json.Marshal(createDto)
		req := httptest.NewRequest("POST", "/api/v1/workoutPlan/byCyclicWorkout", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

func TestGetWorkoutPlansHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Successfully get user workout plans", func(t *testing.T) {
		expectedResponse := []workoutPlan.WorkoutPlan{
			{
				ID:        primitive.NewObjectID(),
				UserID:    "test_user",
				WorkoutID: primitive.NewObjectID().Hex(),
				Dates:     []time.Time{time.Now(), time.Now().AddDate(0, 0, 7)},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        primitive.NewObjectID(),
				UserID:    "test_user",
				WorkoutID: primitive.NewObjectID().Hex(),
				Dates:     []time.Time{time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 8)},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		mockService.On("GetWorkoutPlansByUser", "test_user").Return(expectedResponse, nil)

		req := httptest.NewRequest("GET", "/api/v1/workoutPlan/byUser", nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result []workoutPlan.WorkoutPlan
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, len(expectedResponse), len(result))
		assert.Equal(t, expectedResponse[0].UserID, result[0].UserID)
		assert.Equal(t, expectedResponse[0].WorkoutID, result[0].WorkoutID)
	})

	t.Run("Successfully get empty workout plans", func(t *testing.T) {
		mockService.On("GetWorkoutPlansByUser", "user_with_no_plans").Return([]workoutPlan.WorkoutPlan{}, nil)

		req := httptest.NewRequest("GET", "/api/v1/workoutPlan/byUser", nil)
		req.Header.Set("userid", "user_with_no_plans")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result []workoutPlan.WorkoutPlan
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Empty(t, result)
	})
}
