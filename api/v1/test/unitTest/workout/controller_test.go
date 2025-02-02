package workout_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Npwskp/GymsbroBackend/api/v1/workout/workout"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock service
type MockWorkoutService struct {
	mock.Mock
}

func (m *MockWorkoutService) CreateWorkout(dto *workout.CreateWorkoutDto, userId string) (*workout.Workout, error) {
	args := m.Called(dto, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*workout.Workout), args.Error(1)
}

func (m *MockWorkoutService) GetWorkout(id string) (*workout.Workout, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*workout.Workout), args.Error(1)
}

func (m *MockWorkoutService) GetWorkouts(userId string) ([]*workout.Workout, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*workout.Workout), args.Error(1)
}

func (m *MockWorkoutService) UpdateWorkout(id string, dto *workout.UpdateWorkoutDto, userId string) (*workout.Workout, error) {
	args := m.Called(id, dto, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*workout.Workout), args.Error(1)
}

func (m *MockWorkoutService) DeleteWorkout(id string, userId string) error {
	args := m.Called(id, userId)
	return args.Error(0)
}

func (m *MockWorkoutService) SearchWorkouts(query string, userId string) ([]*workout.Workout, error) {
	args := m.Called(query, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*workout.Workout), args.Error(1)
}

// TestMiddleware sets up the test context with a mock user
func testMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create a mock JWT token claims
		claims := jwt.MapClaims{
			"sub": c.Get("userid", ""), // Get userid from header, default to empty string
		}
		token := &jwt.Token{
			Claims: claims,
		}
		c.Locals("user", token)
		return c.Next()
	}
}

// Test setup helper
func setupTest() (*fiber.App, *MockWorkoutService) {
	app := fiber.New()
	api := app.Group("/api/v1", testMiddleware())

	mockService := new(MockWorkoutService)
	controller := &workout.WorkoutController{
		Instance: api,
		Service:  mockService,
	}
	controller.Handle()
	return app, mockService
}

func TestCreateWorkoutHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Valid workout creation", func(t *testing.T) {
		workoutDto := &workout.CreateWorkoutDto{
			Name:        "Test Workout",
			Description: "Test Description",
			Exercises: []workout.WorkoutExercise{
				{ExerciseID: primitive.NewObjectID().Hex(), Order: 0},
			},
		}

		expectedResponse := &workout.Workout{
			ID:          primitive.NewObjectID(),
			UserID:      "test_user",
			Name:        workoutDto.Name,
			Description: workoutDto.Description,
			Exercises:   workoutDto.Exercises,
		}

		mockService.On("CreateWorkout", workoutDto, "test_user").Return(expectedResponse, nil)

		body, _ := json.Marshal(workoutDto)
		req := httptest.NewRequest("POST", "/api/v1/workout", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		var result workout.Workout
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, expectedResponse.Name, result.Name)
	})
}

func TestGetWorkoutHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Get existing workout", func(t *testing.T) {
		workoutID := primitive.NewObjectID()
		expectedWorkout := &workout.Workout{
			ID:          workoutID,
			UserID:      "test_user",
			Name:        "Test Workout",
			Description: "Test Description",
			Exercises: []workout.WorkoutExercise{
				{ExerciseID: primitive.NewObjectID().Hex(), Order: 0},
			},
		}

		mockService.On("GetWorkout", workoutID.Hex()).Return(expectedWorkout, nil)

		req := httptest.NewRequest("GET", "/api/v1/workout/"+workoutID.Hex(), nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result workout.Workout
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, expectedWorkout.Name, result.Name)
	})
}

func TestGetWorkoutsHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Get all workouts", func(t *testing.T) {
		expectedWorkouts := []*workout.Workout{
			{
				ID:          primitive.NewObjectID(),
				UserID:      "test_user",
				Name:        "Workout 1",
				Description: "Description 1",
			},
			{
				ID:          primitive.NewObjectID(),
				UserID:      "test_user",
				Name:        "Workout 2",
				Description: "Description 2",
			},
		}

		mockService.On("GetWorkouts", "test_user").Return(expectedWorkouts, nil)

		req := httptest.NewRequest("GET", "/api/v1/workout", nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var results []*workout.Workout
		json.NewDecoder(resp.Body).Decode(&results)
		assert.Len(t, results, 2)
	})
}

func TestUpdateWorkoutHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Update existing workout", func(t *testing.T) {
		workoutID := primitive.NewObjectID()
		updateDto := &workout.UpdateWorkoutDto{
			Name:        "Updated Workout",
			Description: "Updated Description",
			Exercises: []workout.WorkoutExercise{
				{ExerciseID: primitive.NewObjectID().Hex(), Order: 1},
			},
		}

		expectedWorkout := &workout.Workout{
			ID:          workoutID,
			UserID:      "test_user",
			Name:        updateDto.Name,
			Description: updateDto.Description,
			Exercises:   updateDto.Exercises,
		}

		mockService.On("UpdateWorkout", workoutID.Hex(), updateDto, "test_user").Return(expectedWorkout, nil)

		body, _ := json.Marshal(updateDto)
		req := httptest.NewRequest("PUT", "/api/v1/workout/"+workoutID.Hex(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result workout.Workout
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, expectedWorkout.Name, result.Name)
	})
}

func TestDeleteWorkoutHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Delete existing workout", func(t *testing.T) {
		workoutID := primitive.NewObjectID()
		mockService.On("DeleteWorkout", workoutID.Hex(), "test_user").Return(nil)

		req := httptest.NewRequest("DELETE", "/api/v1/workout/"+workoutID.Hex(), nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	})
}

func TestSearchWorkoutsHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Search workouts", func(t *testing.T) {
		expectedWorkouts := []*workout.Workout{
			{
				ID:          primitive.NewObjectID(),
				UserID:      "test_user",
				Name:        "Upper Body Workout",
				Description: "Chest and Back",
			},
		}

		mockService.On("SearchWorkouts", "Upper", "test_user").Return(expectedWorkouts, nil)

		req := httptest.NewRequest("GET", "/api/v1/workout/search?query=Upper", nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var results []*workout.Workout
		json.NewDecoder(resp.Body).Decode(&results)
		assert.Len(t, results, 1)
		assert.Contains(t, results[0].Name, "Upper")
	})
}
