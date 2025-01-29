package meal_test

import (
	"context"
	"testing"
	"time"

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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: time.Time{}, // Initialize with zero time for non-deleted meals
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
		assert.True(t, createdMeal.DeletedAt.IsZero()) // Verify meal is not deleted when created
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
		assert.True(t, meal.DeletedAt.IsZero()) // Verify meal is not deleted
	})

	t.Run("Get non-existing meal", func(t *testing.T) {
		nonExistingID := primitive.NewObjectID().Hex()
		meal, err := service.GetMeal(nonExistingID, "test_user")
		assert.Error(t, err)
		assert.Nil(t, meal)
	})

	t.Run("Get soft-deleted meal", func(t *testing.T) {
		// First soft delete the meal
		err := service.DeleteMeal(testMeal.ID.Hex(), testMeal.UserID)
		assert.NoError(t, err)

		// Try to get the soft-deleted meal
		meal, err := service.GetMeal(testMeal.ID.Hex(), testMeal.UserID)
		assert.NoError(t, err)
		assert.NotNil(t, meal)
		assert.False(t, meal.DeletedAt.IsZero()) // Verify DeletedAt is set
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

	t.Run("Get user meals excluding deleted", func(t *testing.T) {
		meals, err := service.GetMealByUser(userid)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(meals))

		// Soft delete one meal
		err = service.DeleteMeal(meals[0].ID.Hex(), userid)
		assert.NoError(t, err)

		// Get meals again - should only return non-deleted meal
		mealsAfterDelete, err := service.GetMealByUser(userid)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(mealsAfterDelete))
		assert.True(t, mealsAfterDelete[0].DeletedAt.IsZero())
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
		assert.True(t, updated.DeletedAt.IsZero()) // Verify update doesn't affect DeletedAt
	})

	t.Run("Update soft-deleted meal should fail", func(t *testing.T) {
		// First soft delete the meal
		err := service.DeleteMeal(testMeal.ID.Hex(), testMeal.UserID)
		assert.NoError(t, err)

		updateDto := &meal.UpdateMealDto{
			Description: "Should Not Update",
			Calories:    700,
		}

		// Attempt to update deleted meal
		_, err = service.UpdateMeal(updateDto, testMeal.ID.Hex(), testMeal.UserID)
		assert.Error(t, err)
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

		// Verify soft deletion
		var found meal.Meal
		err = db.Collection("meal").FindOne(context.Background(), bson.M{"_id": testMeal.ID}).Decode(&found)
		assert.NoError(t, err)
		assert.NotEmpty(t, found.DeletedAt) // Check that DeletedAt is set
	})

	t.Run("Get deleted meal should succeed with DeletedAt set", func(t *testing.T) {
		foundMeal, err := service.GetMeal(testMeal.ID.Hex(), testMeal.UserID)
		assert.NoError(t, err)
		assert.NotNil(t, foundMeal)
		assert.NotEmpty(t, foundMeal.DeletedAt) // Verify DeletedAt is set
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
			DeletedAt:   time.Time{}, // Not deleted
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		&meal.Meal{
			ID:          primitive.NewObjectID(),
			UserID:      "test_user",
			Name:        "Deleted Meal",
			Description: "Should not appear",
			Category:    "Main Course",
			Calories:    500,
			DeletedAt:   time.Now(), // Deleted
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	_, err := db.Collection("meal").InsertMany(context.Background(), testMeals)
	assert.NoError(t, err)

	t.Run("Search should not return deleted meals", func(t *testing.T) {
		filters := meal.SearchFilters{
			Category: "Main Course",
			UserID:   "test_user",
		}
		results, err := service.SearchFilteredMeals(filters)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(results))
		assert.Equal(t, "Chicken Rice", results[0].Name)
		assert.True(t, results[0].DeletedAt.IsZero()) // Verify returned meal is not deleted
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
		assert.True(t, meal.DeletedAt.IsZero()) // Verify meal is not deleted
	})

	t.Run("Get public meal with null UserID", func(t *testing.T) {
		meal, err := service.GetMeal(nullUserMeal.ID.Hex(), "random_user_id")
		assert.NoError(t, err)
		assert.NotNil(t, meal)
		assert.Equal(t, nullUserMeal.ID, meal.ID)
		assert.True(t, meal.DeletedAt.IsZero()) // Verify meal is not deleted
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
