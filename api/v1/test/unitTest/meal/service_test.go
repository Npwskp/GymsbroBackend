package meal_test

import (
	"context"
	"testing"

	"github.com/Npwskp/GymsbroBackend/api/v1/nutrition/meal"
	"github.com/Npwskp/GymsbroBackend/api/v1/nutrition/types"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Helper functions
func createTestMeal(userid string) *meal.Meal {
	return &meal.Meal{
		ID:          primitive.NewObjectID(),
		UserID:      userid,
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
}

func createTestDTO() *meal.CreateMealDto {
	return &meal.CreateMealDto{
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
}

// Setup function
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

func TestCreateMeal(t *testing.T) {
	db := setupTestDB(t)
	service := &meal.MealService{DB: db}

	t.Run("Valid meal creation", func(t *testing.T) {
		dto := createTestDTO()
		userid := "test_user"

		createdMeal, err := service.CreateMeal(dto, userid)
		assert.NoError(t, err)
		assert.NotNil(t, createdMeal)
		assert.Equal(t, userid, createdMeal.UserID)
		assert.Equal(t, dto.Name, createdMeal.Name)
		assert.Equal(t, dto.Calories, createdMeal.Calories)
	})
}

func TestGetMeal(t *testing.T) {
	db := setupTestDB(t)
	service := &meal.MealService{DB: db}

	testMeal := createTestMeal("test_user")
	_, err := db.Collection("meal").InsertOne(context.Background(), testMeal)
	assert.NoError(t, err)

	t.Run("Get existing meal", func(t *testing.T) {
		meal, err := service.GetMeal(testMeal.ID.Hex(), testMeal.UserID)
		assert.NoError(t, err)
		assert.NotNil(t, meal)
		assert.Equal(t, testMeal.ID, meal.ID)
	})

	t.Run("Get non-existing meal", func(t *testing.T) {
		nonExistingID := primitive.NewObjectID().Hex()
		meal, err := service.GetMeal(nonExistingID, "test_user")
		assert.Error(t, err)
		assert.Nil(t, meal)
	})
}

func TestGetMealByUser(t *testing.T) {
	db := setupTestDB(t)
	service := &meal.MealService{DB: db}

	userid := "test_user"
	testMeals := []interface{}{
		createTestMeal(userid),
		createTestMeal(userid),
	}

	_, err := db.Collection("meal").InsertMany(context.Background(), testMeals)
	assert.NoError(t, err)

	t.Run("Get user meals", func(t *testing.T) {
		meals, err := service.GetMealByUser(userid)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(meals))
	})
}

func TestUpdateMeal(t *testing.T) {
	db := setupTestDB(t)
	service := &meal.MealService{DB: db}

	testMeal := createTestMeal("test_user")
	_, err := db.Collection("meal").InsertOne(context.Background(), testMeal)
	assert.NoError(t, err)

	t.Run("Update existing meal", func(t *testing.T) {
		updateDto := &meal.UpdateMealDto{
			Description: "Updated Description",
			Calories:    600,
			Category:    "Updated Category",
			Image:       "Updated Image",
			Nutrients:   []types.Nutrient{{Name: "Protein", Amount: 25, Unit: "g"}},
			Ingredients: []types.Ingredient{{IngredientId: primitive.NewObjectID().Hex(), Amount: 150, Unit: "g"}},
		}

		updated, err := service.UpdateMeal(updateDto, testMeal.ID.Hex(), testMeal.UserID)
		assert.NoError(t, err)
		assert.Equal(t, updateDto.Calories, updated.Calories)
	})
}

func TestDeleteMeal(t *testing.T) {
	db := setupTestDB(t)
	service := &meal.MealService{DB: db}

	testMeal := createTestMeal("test_user")
	_, err := db.Collection("meal").InsertOne(context.Background(), testMeal)
	assert.NoError(t, err)

	t.Run("Delete existing meal", func(t *testing.T) {
		err := service.DeleteMeal(testMeal.ID.Hex(), testMeal.UserID)
		assert.NoError(t, err)

		// Verify deletion
		var found meal.Meal
		err = db.Collection("meal").FindOne(context.Background(), bson.M{"_id": testMeal.ID}).Decode(&found)
		assert.Error(t, err)
	})
}

func TestSearchFilteredMeals(t *testing.T) {
	db := setupTestDB(t)
	service := &meal.MealService{DB: db}

	testMeals := []interface{}{
		&meal.Meal{
			ID:          primitive.NewObjectID(),
			UserID:      "test_user",
			Name:        "Chicken Rice",
			Description: "Healthy meal",
			Category:    "Main Course",
			Calories:    500,
			Nutrients: []types.Nutrient{
				{Name: "Protein", Amount: 30, Unit: "g"},
			},
		},
		&meal.Meal{
			ID:          primitive.NewObjectID(),
			UserID:      "test_user",
			Name:        "Salad",
			Description: "Fresh vegetables",
			Category:    "Side Dish",
			Calories:    200,
			Nutrients: []types.Nutrient{
				{Name: "Fiber", Amount: 5, Unit: "g"},
			},
		},
	}

	_, err := db.Collection("meal").InsertMany(context.Background(), testMeals)
	assert.NoError(t, err)

	t.Run("Search with category filter", func(t *testing.T) {
		filters := meal.SearchFilters{
			Category: "Main Course",
			UserID:   "test_user",
		}
		results, err := service.SearchFilteredMeals(filters)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(results))
	})

	t.Run("Search with name query", func(t *testing.T) {
		filters := meal.SearchFilters{
			Query:  "Chicken",
			UserID: "test_user",
		}
		results, err := service.SearchFilteredMeals(filters)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(results))
	})

	t.Run("Search with calories range", func(t *testing.T) {
		filters := meal.SearchFilters{
			MinCalories: 400,
			MaxCalories: 600,
			UserID:      "test_user",
		}
		results, err := service.SearchFilteredMeals(filters)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(results))
	})
}

func TestGetPublicMeal(t *testing.T) {
	db := setupTestDB(t)
	service := &meal.MealService{DB: db}

	// Create a public meal (empty UserID)
	publicMeal := createTestMeal("")
	_, err := db.Collection("meal").InsertOne(context.Background(), publicMeal)
	assert.NoError(t, err)

	// Create a public meal (null UserID)
	nullUserMeal := createTestMeal("")
	nullUserMeal.UserID = "" // This will be stored as null in MongoDB
	_, err = db.Collection("meal").InsertOne(context.Background(), nullUserMeal)
	assert.NoError(t, err)

	t.Run("Get public meal with empty UserID", func(t *testing.T) {
		// Any user should be able to access public meal
		meal, err := service.GetMeal(publicMeal.ID.Hex(), "random_user_id")
		assert.NoError(t, err)
		assert.NotNil(t, meal)
		assert.Equal(t, publicMeal.ID, meal.ID)
	})

	t.Run("Get public meal with null UserID", func(t *testing.T) {
		meal, err := service.GetMeal(nullUserMeal.ID.Hex(), "random_user_id")
		assert.NoError(t, err)
		assert.NotNil(t, meal)
		assert.Equal(t, nullUserMeal.ID, meal.ID)
	})
}

func TestSearchPublicAndUserMeals(t *testing.T) {
	db := setupTestDB(t)
	service := &meal.MealService{DB: db}

	testMeals := []interface{}{
		&meal.Meal{
			ID:       primitive.NewObjectID(),
			UserID:   "test_user",
			Name:     "User Chicken Rice",
			Category: "Main Course",
			Calories: 500,
		},
		&meal.Meal{
			ID:       primitive.NewObjectID(),
			UserID:   "", // Public meal
			Name:     "Public Chicken Rice",
			Category: "Main Course",
			Calories: 500,
		},
		&meal.Meal{
			ID:       primitive.NewObjectID(),
			UserID:   "", // Will be stored as null
			Name:     "Instant Chicken Rice",
			Category: "Main Course",
			Calories: 500,
		},
	}

	_, err := db.Collection("meal").InsertMany(context.Background(), testMeals)
	assert.NoError(t, err)

	t.Run("Search should return both public and user meals", func(t *testing.T) {
		filters := meal.SearchFilters{
			Query:  "Chicken",
			UserID: "test_user",
		}
		results, err := service.SearchFilteredMeals(filters)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(results)) // Should return all three meals
	})

	t.Run("Search public meals only", func(t *testing.T) {
		filters := meal.SearchFilters{
			Query:  "Public",
			UserID: "different_user", // Different user should still see public meals
		}
		results, err := service.SearchFilteredMeals(filters)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(results))
		assert.Equal(t, "Public Chicken Rice", results[0].Name)
	})

	t.Run("Search instant meals only", func(t *testing.T) {
		filters := meal.SearchFilters{
			Query:  "Instant",
			UserID: "different_user",
		}
		results, err := service.SearchFilteredMeals(filters)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(results))
		assert.Equal(t, "Instant Chicken Rice", results[0].Name)
	})
}
