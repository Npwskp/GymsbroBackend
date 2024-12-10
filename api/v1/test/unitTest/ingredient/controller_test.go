package ingredient_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Npwskp/GymsbroBackend/api/v1/nutrition/ingredient"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock service
type MockIngredientService struct {
	mock.Mock
}

func (m *MockIngredientService) CreateIngredient(dto *ingredient.CreateIngredientDto, userId string) (*ingredient.Ingredient, error) {
	args := m.Called(dto, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ingredient.Ingredient), args.Error(1)
}

func (m *MockIngredientService) GetAllIngredients(userId string) ([]*ingredient.Ingredient, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ingredient.Ingredient), args.Error(1)
}

func (m *MockIngredientService) GetIngredient(id string, userId string) (*ingredient.Ingredient, error) {
	args := m.Called(id, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ingredient.Ingredient), args.Error(1)
}

func (m *MockIngredientService) GetIngredientByUser(userId string) ([]*ingredient.Ingredient, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ingredient.Ingredient), args.Error(1)
}

func (m *MockIngredientService) DeleteIngredient(id string, userId string) error {
	args := m.Called(id, userId)
	return args.Error(0)
}

func (m *MockIngredientService) UpdateIngredient(doc *ingredient.UpdateIngredientDto, id string, userId string) (*ingredient.Ingredient, error) {
	args := m.Called(doc, id, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ingredient.Ingredient), args.Error(1)
}

func (m *MockIngredientService) SearchFilteredIngredients(filters ingredient.SearchFilters) ([]*ingredient.Ingredient, error) {
	args := m.Called(filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ingredient.Ingredient), args.Error(1)
}

func setupTest() (*fiber.App, *MockIngredientService) {
	app := fiber.New()
	mockService := new(MockIngredientService)

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

	controller := &ingredient.IngredientController{
		Instance: app,
		Service:  mockService,
	}
	controller.Handle()
	return app, mockService
}

func TestCreateIngredientController(t *testing.T) {
	app, mockService := setupTest()

	dto := &ingredient.CreateIngredientDto{
		Name:        "Test Ingredient",
		Description: "Test Description",
		Category:    "Test Category",
		Calories:    100,
	}

	expectedResponse := &ingredient.Ingredient{
		ID:          primitive.NewObjectID(),
		UserID:      "test-user-id",
		Name:        "Test Ingredient",
		Description: "Test Description",
		Category:    "Test Category",
		Calories:    100,
	}

	mockService.On("CreateIngredient", dto, "test-user-id").Return(expectedResponse, nil)

	jsonBody, _ := json.Marshal(dto)
	req := httptest.NewRequest("POST", "/ingredient", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "test-user-id")

	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestGetIngredientController(t *testing.T) {
	app, mockService := setupTest()

	expectedResponse := &ingredient.Ingredient{
		ID:          primitive.NewObjectID(),
		UserID:      "test-user-id",
		Name:        "Test Ingredient",
		Description: "Test Description",
		Category:    "Test Category",
		Calories:    100,
	}

	mockService.On("GetIngredient", expectedResponse.ID.Hex(), "test-user-id").Return(expectedResponse, nil)

	req := httptest.NewRequest("GET", "/ingredient/"+expectedResponse.ID.Hex(), nil)
	req.Header.Set("X-User-ID", "test-user-id")

	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestDeleteIngredientController(t *testing.T) {
	app, mockService := setupTest()

	id := primitive.NewObjectID().Hex()
	mockService.On("DeleteIngredient", id, "test-user-id").Return(nil)

	req := httptest.NewRequest("DELETE", "/ingredient/"+id, nil)
	req.Header.Set("X-User-ID", "test-user-id")

	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestUpdateIngredientController(t *testing.T) {
	app, mockService := setupTest()

	id := primitive.NewObjectID().Hex()
	dto := &ingredient.UpdateIngredientDto{
		Name:        "Updated Ingredient",
		Description: "Updated Description",
		Category:    "Updated Category",
		Calories:    150,
	}

	expectedResponse := &ingredient.Ingredient{
		ID:          primitive.NewObjectID(),
		UserID:      "test-user-id",
		Name:        "Updated Ingredient",
		Description: "Updated Description",
		Category:    "Updated Category",
		Calories:    150,
	}

	mockService.On("UpdateIngredient", dto, id, "test-user-id").Return(expectedResponse, nil)

	jsonBody, _ := json.Marshal(dto)
	req := httptest.NewRequest("PUT", "/ingredient/"+id, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "test-user-id")

	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestSearchFilteredIngredientsController(t *testing.T) {
	app, mockService := setupTest()

	filters := ingredient.SearchFilters{
		Query:    "",
		Category: "Category",
		UserID:   "test-user-id",
	}

	expectedResponse := []*ingredient.Ingredient{
		{
			ID:          primitive.NewObjectID(),
			UserID:      "test-user-id",
			Name:        "Test Ingredient",
			Description: "Test Description",
			Category:    "Test Category",
			Calories:    100,
		},
	}

	mockService.On("SearchFilteredIngredients", filters).Return(expectedResponse, nil)

	req := httptest.NewRequest("GET", "/ingredient/search?query=Test&category=Category&userId=test-user-id", nil)
	req.Header.Set("X-User-ID", "test-user-id")

	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}
