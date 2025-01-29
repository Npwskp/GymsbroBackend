package foodlog_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	foodlog "github.com/Npwskp/GymsbroBackend/api/v1/nutrition/foodLog"
	"github.com/Npwskp/GymsbroBackend/api/v1/nutrition/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock service
type MockFoodLogService struct {
	mock.Mock
}

func (m *MockFoodLogService) AddMealToFoodLog(dto *foodlog.AddMealToFoodLogDto, userid string) (*foodlog.FoodLog, error) {
	args := m.Called(dto, userid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*foodlog.FoodLog), args.Error(1)
}

func (m *MockFoodLogService) GetFoodLog(id string, userid string) (*foodlog.FoodLog, error) {
	args := m.Called(id, userid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*foodlog.FoodLog), args.Error(1)
}

func (m *MockFoodLogService) GetFoodLogByUser(userid string) ([]*foodlog.FoodLog, error) {
	args := m.Called(userid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*foodlog.FoodLog), args.Error(1)
}

func (m *MockFoodLogService) GetFoodLogByUserDate(userid string, date string) (*foodlog.FoodLog, error) {
	args := m.Called(userid, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*foodlog.FoodLog), args.Error(1)
}

func (m *MockFoodLogService) DeleteFoodLog(id string, userid string) error {
	args := m.Called(id, userid)
	return args.Error(0)
}

func (m *MockFoodLogService) UpdateFoodLog(doc *foodlog.UpdateFoodLogDto, id string, userid string) (*foodlog.FoodLog, error) {
	args := m.Called(doc, id, userid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*foodlog.FoodLog), args.Error(1)
}

func (m *MockFoodLogService) CalculateDailyNutrients(date string, userid string) (*foodlog.DailyNutrientResponse, error) {
	args := m.Called(date, userid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*foodlog.DailyNutrientResponse), args.Error(1)
}

func setupTest() (*fiber.App, *MockFoodLogService) {
	app := fiber.New()
	mockService := new(MockFoodLogService)

	// Create a mock JWT token using golang-jwt/jwt/v4
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "test-user-id",
	})

	// Replace the middleware with one that properly sets up the test context
	app.Use(func(c *fiber.Ctx) error {
		// Set the mock token in the context
		c.Locals("user", token)
		return c.Next()
	})

	controller := &foodlog.FoodLogController{
		Instance: app,
		Service:  mockService,
	}
	controller.Handle()
	return app, mockService
}

func TestAddMealToFoodLogController(t *testing.T) {
	app, mockService := setupTest()

	dto := &foodlog.AddMealToFoodLogDto{
		Date:  "2024-03-20",
		Meals: []string{"Breakfast", "Lunch"},
	}

	expectedResponse := &foodlog.FoodLog{
		ID:     primitive.NewObjectID(),
		UserID: "test-user-id",
		Date:   "2024-03-20",
		Meals:  []string{"Breakfast", "Lunch"},
	}

	mockService.On("AddMealToFoodLog", mock.AnythingOfType("*foodlog.AddMealToFoodLogDto"), "test-user-id").
		Return(expectedResponse, nil)

	jsonBody, _ := json.Marshal(dto)
	req := httptest.NewRequest("POST", "/foodlog", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "test-user-id")

	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestGetFoodLogController(t *testing.T) {
	app, mockService := setupTest()

	expectedResponse := &foodlog.FoodLog{
		ID:     primitive.NewObjectID(),
		UserID: "test-user-id",
		Date:   "2024-03-20",
		Meals:  []string{"Breakfast", "Lunch"},
	}

	mockService.On("GetFoodLog", expectedResponse.ID.Hex(), "test-user-id").Return(expectedResponse, nil)

	req := httptest.NewRequest("GET", "/foodlog/"+expectedResponse.ID.Hex(), nil)
	req.Header.Set("X-User-ID", "test-user-id")

	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestGetFoodLogByUserController(t *testing.T) {
	app, mockService := setupTest()

	expectedResponse := []*foodlog.FoodLog{
		{
			ID:     primitive.NewObjectID(),
			UserID: "test-user-id",
			Date:   "2024-03-20",
			Meals:  []string{"Breakfast"},
		},
		{
			ID:     primitive.NewObjectID(),
			UserID: "test-user-id",
			Date:   "2024-03-21",
			Meals:  []string{"Lunch"},
		},
	}

	mockService.On("GetFoodLogByUser", "test-user-id").Return(expectedResponse, nil)

	req := httptest.NewRequest("GET", "/foodlog/user", nil)
	req.Header.Set("X-User-ID", "test-user-id")

	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestGetFoodLogByUserDateController(t *testing.T) {
	app, mockService := setupTest()

	expectedResponse := &foodlog.FoodLog{
		ID:     primitive.NewObjectID(),
		UserID: "test-user-id",
		Date:   "2024-03-20",
		Meals:  []string{"Breakfast", "Lunch"},
	}

	mockService.On("GetFoodLogByUserDate", "test-user-id", "2024-03-20").Return(expectedResponse, nil)

	req := httptest.NewRequest("GET", "/foodlog/user/2024-03-20", nil)
	req.Header.Set("X-User-ID", "test-user-id")

	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestDeleteFoodLogController(t *testing.T) {
	app, mockService := setupTest()

	id := primitive.NewObjectID().Hex()
	mockService.On("DeleteFoodLog", id, "test-user-id").Return(nil)

	req := httptest.NewRequest("DELETE", "/foodlog/"+id, nil)
	req.Header.Set("X-User-ID", "test-user-id")

	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestUpdateFoodLogController(t *testing.T) {
	app, mockService := setupTest()

	id := primitive.NewObjectID().Hex()
	dto := &foodlog.UpdateFoodLogDto{
		Date:  "2024-03-20",
		Meals: []string{"Updated Breakfast", "Updated Lunch"},
	}

	expectedResponse := &foodlog.FoodLog{
		ID:     primitive.NewObjectID(),
		UserID: "test-user-id",
		Date:   "2024-03-20",
		Meals:  []string{"Updated Breakfast", "Updated Lunch"},
	}

	mockService.On("UpdateFoodLog", dto, id, "test-user-id").Return(expectedResponse, nil)

	jsonBody, _ := json.Marshal(dto)
	req := httptest.NewRequest("PUT", "/foodlog/"+id, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "test-user-id")

	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestCalculateDailyNutrientsController(t *testing.T) {
	app, mockService := setupTest()

	// Test case 1: Successful calculation
	t.Run("Success", func(t *testing.T) {
		expectedResponse := &foodlog.DailyNutrientResponse{
			Date:     "2024-03-20",
			Calories: 2000.0,
			Nutrients: []types.Nutrient{
				{
					Name:   "Protein",
					Amount: 150.0,
					Unit:   "g",
				},
				{
					Name:   "Carbohydrates",
					Amount: 250.0,
					Unit:   "g",
				},
			},
		}

		mockService.On("CalculateDailyNutrients", "2024-03-20", "test-user-id").
			Return(expectedResponse, nil).Once()

		req := httptest.NewRequest("GET", "/foodlog/nutrients/2024-03-20", nil)
		req.Header.Set("X-User-ID", "test-user-id")

		resp, err := app.Test(req)

		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result foodlog.DailyNutrientResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.Nil(t, err)
		assert.Equal(t, expectedResponse.Date, result.Date)
		assert.Equal(t, expectedResponse.Calories, result.Calories)
		assert.Len(t, result.Nutrients, len(expectedResponse.Nutrients))
	})

	// Test case 2: Invalid date format
	t.Run("Invalid Date Format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/foodlog/nutrients/invalid-date", nil)
		req.Header.Set("X-User-ID", "test-user-id")

		resp, err := app.Test(req)

		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var errorResponse map[string]string
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		assert.Nil(t, err)
		assert.Contains(t, errorResponse["error"], "Invalid date format")
	})

	mockService.AssertExpectations(t)
}
