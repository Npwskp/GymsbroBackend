package meal_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Npwskp/GymsbroBackend/api/v1/nutrition/meal"
	"github.com/Npwskp/GymsbroBackend/api/v1/nutrition/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock service
type MockMealService struct {
	mock.Mock
}

func (m *MockMealService) CreateMeal(dto *meal.CreateMealDto, userId string) (*meal.Meal, error) {
	args := m.Called(dto, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*meal.Meal), args.Error(1)
}

func (m *MockMealService) CalculateNutrient(body *meal.CalculateNutrientBody, userId string) (*meal.CalculateNutrientResponse, error) {
	args := m.Called(body, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*meal.CalculateNutrientResponse), args.Error(1)
}

func (m *MockMealService) GetMeal(id string, userId string) (*meal.Meal, error) {
	args := m.Called(id, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*meal.Meal), args.Error(1)
}

func (m *MockMealService) GetMealByUser(userId string) ([]*meal.Meal, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*meal.Meal), args.Error(1)
}

func (m *MockMealService) DeleteMeal(id string, userId string) error {
	args := m.Called(id, userId)
	return args.Error(0)
}

func (m *MockMealService) UpdateMeal(doc *meal.UpdateMealDto, id string, userId string) (*meal.Meal, error) {
	args := m.Called(doc, id, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*meal.Meal), args.Error(1)
}

func (m *MockMealService) SearchFilteredMeals(filters meal.SearchFilters) ([]*meal.Meal, error) {
	args := m.Called(filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*meal.Meal), args.Error(1)
}

func setupTest() (*fiber.App, *MockMealService) {
	app := fiber.New()
	mockService := new(MockMealService)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "test-user-id",
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", token)
		return c.Next()
	})

	controller := &meal.MealController{
		Instance: app,
		Service:  mockService,
	}
	controller.Handle()
	return app, mockService
}

func TestCreateMealController(t *testing.T) {
	app, mockService := setupTest()

	dto := &meal.CreateMealDto{
		Name:        "Test Meal",
		Description: "Test Description",
		Category:    "Test Category",
		Calories:    500,
		Nutrients: []types.Nutrient{
			{Name: "Protein", Amount: 20, Unit: "g"},
		},
		Ingredients: []types.Ingredient{
			{IngredientId: primitive.NewObjectID().Hex(), Amount: 100, Unit: "g"},
		},
	}

	expectedResponse := &meal.Meal{
		ID:          primitive.NewObjectID(),
		UserID:      "test-user-id",
		Name:        "Test Meal",
		Description: "Test Description",
		Category:    "Test Category",
		Calories:    500,
		Nutrients:   dto.Nutrients,
		Ingredients: dto.Ingredients,
	}

	mockService.On("CreateMeal", dto, "test-user-id").Return(expectedResponse, nil)

	jsonBody, _ := json.Marshal(dto)
	req := httptest.NewRequest("POST", "/meal", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestGetMealController(t *testing.T) {
	app, mockService := setupTest()

	expectedResponse := &meal.Meal{
		ID:          primitive.NewObjectID(),
		UserID:      "test-user-id",
		Name:        "Test Meal",
		Description: "Test Description",
		Category:    "Test Category",
		Calories:    500,
	}

	mockService.On("GetMeal", expectedResponse.ID.Hex(), "test-user-id").Return(expectedResponse, nil)

	req := httptest.NewRequest("GET", "/meal/"+expectedResponse.ID.Hex(), nil)

	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestDeleteMealController(t *testing.T) {
	app, mockService := setupTest()

	id := primitive.NewObjectID().Hex()
	mockService.On("DeleteMeal", id, "test-user-id").Return(nil)

	req := httptest.NewRequest("DELETE", "/meal/"+id, nil)

	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestUpdateMealController(t *testing.T) {
	app, mockService := setupTest()

	id := primitive.NewObjectID().Hex()
	dto := &meal.UpdateMealDto{
		Description: "Updated Description",
		Calories:    600,
		Category:    "Updated Category",
		Image:       "Updated Image",
		Nutrients:   []types.Nutrient{{Name: "Protein", Amount: 25, Unit: "g"}},
		Ingredients: []types.Ingredient{{IngredientId: primitive.NewObjectID().Hex(), Amount: 150, Unit: "g"}},
	}

	expectedResponse := &meal.Meal{
		ID:          primitive.NewObjectID(),
		UserID:      "test-user-id",
		Description: "Updated Description",
		Calories:    600,
	}

	mockService.On("UpdateMeal", dto, id, "test-user-id").Return(expectedResponse, nil)

	jsonBody, _ := json.Marshal(dto)
	req := httptest.NewRequest("PUT", "/meal/"+id, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestSearchFilteredMealsController(t *testing.T) {
	app, mockService := setupTest()

	filters := meal.SearchFilters{
		Query:       "Test",
		Category:    "Test Category",
		MinCalories: 300,
		MaxCalories: 800,
		UserID:      "test-user-id",
	}

	expectedResponse := []*meal.Meal{
		{
			ID:          primitive.NewObjectID(),
			UserID:      "test-user-id",
			Name:        "Test Meal",
			Description: "Test Description",
			Category:    "Test Category",
			Calories:    500,
		},
	}

	mockService.On("SearchFilteredMeals", mock.MatchedBy(func(f meal.SearchFilters) bool {
		return f.Category == filters.Category && f.UserID == filters.UserID
	})).Return(expectedResponse, nil)

	req := httptest.NewRequest("GET", "/meal/search?q=Test&category=Test%20Category&minCalories=300&maxCalories=800", nil)

	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}
