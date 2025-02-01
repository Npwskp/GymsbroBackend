package exercise_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/Npwskp/GymsbroBackend/api/v1/workout/exercise"
	exerciseEnums "github.com/Npwskp/GymsbroBackend/api/v1/workout/exercise/enums"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock service
type MockExerciseService struct {
	mock.Mock
}

func (m *MockExerciseService) CreateManyExercises(exercises *[]exercise.CreateExerciseDto, userId string) ([]*exercise.Exercise, error) {
	args := m.Called(exercises, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*exercise.Exercise), args.Error(1)
}

func (m *MockExerciseService) GetAllExercises(userId string) ([]*exercise.Exercise, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*exercise.Exercise), args.Error(1)
}

func (m *MockExerciseService) GetExercise(id string, userId string) (*exercise.Exercise, error) {
	args := m.Called(id, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*exercise.Exercise), args.Error(1)
}

func (m *MockExerciseService) DeleteExercise(id string, userId string) error {
	args := m.Called(id, userId)
	return args.Error(0)
}

func (m *MockExerciseService) UpdateExercise(doc *exercise.UpdateExerciseDto, id string, userId string) (*exercise.Exercise, error) {
	args := m.Called(doc, id, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*exercise.Exercise), args.Error(1)
}

func (m *MockExerciseService) UpdateExerciseImage(c *fiber.Ctx, id string, file io.Reader, filename string, contentType string, userId string) (*exercise.Exercise, error) {
	args := m.Called(c, id, file, filename, contentType, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*exercise.Exercise), args.Error(1)
}

func (m *MockExerciseService) SearchAndFilterExercise(equipment []exerciseEnums.Equipment, mechanics []exerciseEnums.Mechanics, force []exerciseEnums.Force, bodyPart []exerciseEnums.BodyPart, targetMuscle []exerciseEnums.TargetMuscle, query string, userID string) ([]*exercise.Exercise, error) {
	args := m.Called(equipment, mechanics, force, bodyPart, targetMuscle, query, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*exercise.Exercise), args.Error(1)
}

func (m *MockExerciseService) CreateExercise(doc *exercise.CreateExerciseDto, userId string) (*exercise.Exercise, error) {
	args := m.Called(doc, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*exercise.Exercise), args.Error(1)
}

func (m *MockExerciseService) FindSimilarExercises(id string, userId string, limit int) ([]*exercise.Exercise, error) {
	args := m.Called(id, userId, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*exercise.Exercise), args.Error(1)
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
func setupTest() (*fiber.App, *MockExerciseService) {
	app := fiber.New()
	api := app.Group("/api/v1", testMiddleware())

	mockService := new(MockExerciseService)
	controller := &exercise.ExerciseController{
		Instance: api,
		Service:  mockService,
	}
	controller.Handle()
	return app, mockService
}

func TestPostExerciseHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Valid exercise creation", func(t *testing.T) {
		exerciseDto := exercise.CreateExerciseDto{
			Name:         "Push-up",
			Equipment:    exerciseEnums.BodyWeight,
			Mechanics:    exerciseEnums.Compound,
			Force:        exerciseEnums.Push,
			Preparation:  []string{"Get into plank position"},
			Execution:    []string{"Lower body", "Push up"},
			BodyPart:     []exerciseEnums.BodyPart{exerciseEnums.Chest},
			TargetMuscle: []exerciseEnums.TargetMuscle{exerciseEnums.PectoralisMajorSternal},
		}

		expectedResponse := &exercise.Exercise{
			ID:           primitive.NewObjectID(),
			UserID:       "test_user",
			Name:         exerciseDto.Name,
			Equipment:    exerciseDto.Equipment,
			Mechanics:    exerciseDto.Mechanics,
			Force:        exerciseDto.Force,
			Preparation:  exerciseDto.Preparation,
			Execution:    exerciseDto.Execution,
			BodyPart:     exerciseDto.BodyPart,
			TargetMuscle: exerciseDto.TargetMuscle,
		}

		mockService.On("CreateExercise", &exerciseDto, "test_user").Return(expectedResponse, nil)

		body, _ := json.Marshal(exerciseDto)
		req := httptest.NewRequest("POST", "/api/v1/exercise", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		var result exercise.Exercise
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, expectedResponse.Name, result.Name)
	})
}

func TestGetExerciseHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Get existing exercise", func(t *testing.T) {
		exerciseID := primitive.NewObjectID()
		expectedExercise := &exercise.Exercise{
			ID:           exerciseID,
			UserID:       "test_user",
			Name:         "Push-up",
			Equipment:    exerciseEnums.BodyWeight,
			Mechanics:    exerciseEnums.Compound,
			Force:        exerciseEnums.Push,
			BodyPart:     []exerciseEnums.BodyPart{exerciseEnums.Chest},
			TargetMuscle: []exerciseEnums.TargetMuscle{exerciseEnums.PectoralisMajorSternal},
		}

		mockService.On("GetExercise", exerciseID.Hex(), "test_user").Return(expectedExercise, nil)

		req := httptest.NewRequest("GET", "/api/v1/exercise/"+exerciseID.Hex(), nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result exercise.Exercise
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, expectedExercise.Name, result.Name)
	})
}

func TestGetAllExercisesHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Get all exercises", func(t *testing.T) {
		expectedExercises := []*exercise.Exercise{
			{
				ID:        primitive.NewObjectID(),
				UserID:    "test_user",
				Name:      "Push-up",
				Equipment: exerciseEnums.BodyWeight,
				Mechanics: exerciseEnums.Compound,

				Force:        exerciseEnums.Push,
				BodyPart:     []exerciseEnums.BodyPart{exerciseEnums.Chest},
				TargetMuscle: []exerciseEnums.TargetMuscle{exerciseEnums.PectoralisMajorSternal},
			},
			{
				ID:           primitive.NewObjectID(),
				UserID:       "test_user",
				Name:         "Squat",
				Equipment:    exerciseEnums.BodyWeight,
				Mechanics:    exerciseEnums.Compound,
				Force:        exerciseEnums.Push,
				BodyPart:     []exerciseEnums.BodyPart{exerciseEnums.Thighs},
				TargetMuscle: []exerciseEnums.TargetMuscle{exerciseEnums.Quadriceps},
			},
		}

		mockService.On("GetAllExercises", "test_user").Return(expectedExercises, nil)

		req := httptest.NewRequest("GET", "/api/v1/exercise", nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result []*exercise.Exercise
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, len(expectedExercises), len(result))
	})
}

func TestDeleteExerciseHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Delete existing exercise", func(t *testing.T) {
		exerciseID := primitive.NewObjectID()
		mockService.On("DeleteExercise", exerciseID.Hex(), "test_user").Return(nil)

		req := httptest.NewRequest("DELETE", "/api/v1/exercise/"+exerciseID.Hex(), nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	})
}

func TestUpdateExerciseHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Update existing exercise", func(t *testing.T) {
		exerciseID := primitive.NewObjectID()
		updateDto := exercise.UpdateExerciseDto{
			Name:      "Modified Push-up",
			Equipment: exerciseEnums.Barbell,
		}

		expectedResponse := &exercise.Exercise{
			ID:        exerciseID,
			UserID:    "test_user",
			Name:      updateDto.Name,
			Equipment: updateDto.Equipment,
		}

		mockService.On("UpdateExercise", &updateDto, exerciseID.Hex(), "test_user").Return(expectedResponse, nil)

		body, _ := json.Marshal(updateDto)
		req := httptest.NewRequest("PUT", "/api/v1/exercise/"+exerciseID.Hex(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result exercise.Exercise
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, expectedResponse.Name, result.Name)
		assert.Equal(t, expectedResponse.Equipment, result.Equipment)

	})
}

func TestSearchAndFilterExerciseHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Search and filter exercises", func(t *testing.T) {
		expectedExercises := []*exercise.Exercise{
			{
				ID:           primitive.NewObjectID(),
				UserID:       "test_user",
				Name:         "Push-up",
				Equipment:    exerciseEnums.BodyWeight,
				Mechanics:    exerciseEnums.Compound,
				Force:        exerciseEnums.Push,
				BodyPart:     []exerciseEnums.BodyPart{exerciseEnums.Chest},
				TargetMuscle: []exerciseEnums.TargetMuscle{exerciseEnums.PectoralisMajorSternal},
			},
		}

		mockService.On("SearchAndFilterExercise",
			[]exerciseEnums.Equipment{exerciseEnums.BodyWeight},
			[]exerciseEnums.Mechanics(nil),
			[]exerciseEnums.Force(nil),
			[]exerciseEnums.BodyPart(nil),
			[]exerciseEnums.TargetMuscle(nil),
			"",
			"test_user",
		).Return(expectedExercises, nil)

		req := httptest.NewRequest("GET", "/api/v1/exercise/search?equipment=Body+Weight", nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result []*exercise.Exercise
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, len(expectedExercises), len(result))
		assert.Equal(t, expectedExercises[0].Name, result[0].Name)
	})
}

func TestGetAllEquipmentHandler(t *testing.T) {
	app, _ := setupTest()

	req := httptest.NewRequest("GET", "/api/v1/exercise/equipment", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var equipment []exerciseEnums.Equipment
	json.NewDecoder(resp.Body).Decode(&equipment)

	// Verify all defined equipment constants are present
	expectedEquipment := exerciseEnums.GetAllEquipment()
	assert.ElementsMatch(t, expectedEquipment, equipment)
}

func TestGetAllMechanicsHandler(t *testing.T) {
	app, _ := setupTest()

	req := httptest.NewRequest("GET", "/api/v1/exercise/mechanics", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var mechanics []exerciseEnums.Mechanics
	json.NewDecoder(resp.Body).Decode(&mechanics)

	expectedMechanics := exerciseEnums.GetAllMechanics()
	assert.ElementsMatch(t, expectedMechanics, mechanics)
}

func TestGetAllForceHandler(t *testing.T) {
	app, _ := setupTest()

	req := httptest.NewRequest("GET", "/api/v1/exercise/force", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var forces []exerciseEnums.Force
	json.NewDecoder(resp.Body).Decode(&forces)

	expectedForces := exerciseEnums.GetAllForces()
	assert.ElementsMatch(t, expectedForces, forces)
}

func TestGetAllBodyPartHandler(t *testing.T) {
	app, _ := setupTest()

	req := httptest.NewRequest("GET", "/api/v1/exercise/bodypart", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var bodyParts []exerciseEnums.BodyPart
	json.NewDecoder(resp.Body).Decode(&bodyParts)

	// Get all defined body parts from the enums package
	expectedBodyParts := exerciseEnums.GetAllBodyParts()
	assert.ElementsMatch(t, expectedBodyParts, bodyParts)
}

func TestGetAllTargetMusclesHandler(t *testing.T) {
	app, _ := setupTest()

	req := httptest.NewRequest("GET", "/api/v1/exercise/targetmuscle", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var targetMuscles []exerciseEnums.TargetMuscle
	json.NewDecoder(resp.Body).Decode(&targetMuscles)

	// Get all defined target muscles from the enums package
	expectedTargetMuscles := exerciseEnums.GetAllTargetMuscles()
	assert.ElementsMatch(t, expectedTargetMuscles, targetMuscles)
}
