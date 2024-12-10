package ingredient_test

import (
	"context"
	"testing"

	ingredient "github.com/Npwskp/GymsbroBackend/api/v1/nutrition/ingredient"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Helper functions
func createTestIngredient(userid string) *ingredient.Ingredient {
	return &ingredient.Ingredient{
		ID:          primitive.NewObjectID(),
		UserID:      userid,
		Name:        "Test Ingredient",
		Description: "Test Description",
		Category:    "Test Category",
		Calories:    100,
	}
}

func createTestDTO() *ingredient.CreateIngredientDto {
	return &ingredient.CreateIngredientDto{
		Name:        "Test Ingredient",
		Description: "Test Description",
		Category:    "Test Category",
		Calories:    100,
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

func TestCreateIngredient(t *testing.T) {
	db := setupTestDB(t)
	service := &ingredient.IngredientService{DB: db}

	t.Run("Valid ingredient creation", func(t *testing.T) {
		dto := createTestDTO()
		userid := "test_user"

		ing, err := service.CreateIngredient(dto, userid)
		assert.NoError(t, err)
		assert.NotNil(t, ing)
		assert.Equal(t, userid, ing.UserID)
		assert.Equal(t, dto.Name, ing.Name)
		assert.Equal(t, dto.Calories, ing.Calories)
	})
}

func TestGetIngredient(t *testing.T) {
	db := setupTestDB(t)
	service := &ingredient.IngredientService{DB: db}

	testIng := createTestIngredient("test_user")
	_, err := db.Collection("ingredient").InsertOne(context.Background(), testIng)
	assert.NoError(t, err)

	t.Run("Get existing ingredient", func(t *testing.T) {
		ing, err := service.GetIngredient(testIng.ID.Hex(), testIng.UserID)
		assert.NoError(t, err)
		assert.NotNil(t, ing)
		assert.Equal(t, testIng.ID, ing.ID)
	})

	t.Run("Get non-existing ingredient", func(t *testing.T) {
		nonExistingID := primitive.NewObjectID().Hex()
		ing, err := service.GetIngredient(nonExistingID, "test_user")
		assert.Error(t, err)
		assert.Nil(t, ing)
	})
}

func TestGetIngredientByUser(t *testing.T) {
	db := setupTestDB(t)
	service := &ingredient.IngredientService{DB: db}

	userid := "test_user"
	testIngredients := []interface{}{
		createTestIngredient(userid),
		createTestIngredient(userid),
	}

	_, err := db.Collection("ingredient").InsertMany(context.Background(), testIngredients)
	assert.NoError(t, err)

	t.Run("Get user ingredients", func(t *testing.T) {
		ingredients, err := service.GetIngredientByUser(userid)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(ingredients))
	})
}

func TestUpdateIngredient(t *testing.T) {
	db := setupTestDB(t)
	service := &ingredient.IngredientService{DB: db}

	testIng := createTestIngredient("test_user")
	_, err := db.Collection("ingredient").InsertOne(context.Background(), testIng)
	assert.NoError(t, err)

	t.Run("Update existing ingredient", func(t *testing.T) {
		updateDto := &ingredient.UpdateIngredientDto{
			Name:        "Updated Name",
			Description: "Updated Description",
			Category:    "Updated Category",
			Calories:    200,
		}

		updated, err := service.UpdateIngredient(updateDto, testIng.ID.Hex(), testIng.UserID)
		assert.NoError(t, err)
		assert.Equal(t, updateDto.Name, updated.Name)
		assert.Equal(t, updateDto.Calories, updated.Calories)
	})
}

func TestDeleteIngredient(t *testing.T) {
	db := setupTestDB(t)
	service := &ingredient.IngredientService{DB: db}

	testIng := createTestIngredient("test_user")
	_, err := db.Collection("ingredient").InsertOne(context.Background(), testIng)
	assert.NoError(t, err)

	t.Run("Delete existing ingredient", func(t *testing.T) {
		err := service.DeleteIngredient(testIng.ID.Hex(), testIng.UserID)
		assert.NoError(t, err)

		// Verify deletion
		var found ingredient.Ingredient
		err = db.Collection("ingredient").FindOne(context.Background(), bson.M{"_id": testIng.ID}).Decode(&found)
		assert.Error(t, err)
	})
}

func TestSearchFilteredIngredients(t *testing.T) {
	db := setupTestDB(t)
	service := &ingredient.IngredientService{DB: db}

	testIngredients := []interface{}{
		&ingredient.Ingredient{
			ID:          primitive.NewObjectID(),
			UserID:      "test_user",
			Name:        "Apple",
			Description: "Fresh fruit",
			Category:    "Fruits",
			Calories:    52,
		},
		&ingredient.Ingredient{
			ID:          primitive.NewObjectID(),
			UserID:      "test_user",
			Name:        "Banana",
			Description: "Yellow fruit",
			Category:    "Fruits",
			Calories:    89,
		},
	}

	_, err := db.Collection("ingredient").InsertMany(context.Background(), testIngredients)
	assert.NoError(t, err)

	t.Run("Search with category filter", func(t *testing.T) {
		filters := ingredient.SearchFilters{
			Category: "Fruits",
			UserID:   "test_user",
		}
		results, err := service.SearchFilteredIngredients(filters)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(results))
	})

	t.Run("Search with name query", func(t *testing.T) {
		filters := ingredient.SearchFilters{
			Query:  "Apple",
			UserID: "test_user",
		}
		results, err := service.SearchFilteredIngredients(filters)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(results))
	})
}
