package foodlog_test

import (
	"context"
	"testing"
	"time"

	foodlog "github.com/Npwskp/GymsbroBackend/api/v1/nutrition/foodLog"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Common test data
func createTestFoodLog(userid string) *foodlog.FoodLog {
	return &foodlog.FoodLog{
		ID:     primitive.NewObjectID(),
		UserID: userid,
		Date:   time.Now().Format("2006-01-02"),
		Meals:  []string{"meal1"},
	}
}

func createTestDTO() *foodlog.AddMealToFoodLogDto {
	return &foodlog.AddMealToFoodLogDto{
		Date:  time.Now().Format("2006-01-02"),
		Meals: []string{"meal1", "meal2"},
	}
}

// Setup functions
func setupTestDB(t *testing.T) *mongo.Database {
	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Create a test database
	db := client.Database("testdb_" + primitive.NewObjectID().Hex())

	// Clean up function
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

func setupTestApp() (*fiber.App, *MockFoodLogService) {
	app := fiber.New()
	mockService := new(MockFoodLogService)
	controller := &foodlog.FoodLogController{
		Instance: app,
		Service:  mockService,
	}
	controller.Handle()
	return app, mockService
}

func TestCreateFoodLog(t *testing.T) {
	db := setupTestDB(t)
	service := &foodlog.FoodLogService{DB: db}

	// Test case 1: Create food log with valid data
	t.Run("Valid food log creation", func(t *testing.T) {
		dto := &foodlog.AddMealToFoodLogDto{
			Date:  time.Now().Format("2006-01-02"),
			Meals: []string{"meal1", "meal2"},
		}
		userid := "test_user"

		foodLog, err := service.AddMealToFoodLog(dto, userid)
		assert.NoError(t, err)
		assert.NotNil(t, foodLog)
		assert.Equal(t, userid, foodLog.UserID)
		assert.Equal(t, len(dto.Meals), len(foodLog.Meals))
	})
}

func TestGetFoodLog(t *testing.T) {
	db := setupTestDB(t)
	service := &foodlog.FoodLogService{DB: db}

	// Create a test food log
	testFoodLog := &foodlog.FoodLog{
		ID:     primitive.NewObjectID(),
		UserID: "test_user",
		Date:   time.Now().Format("2006-01-02"),
		Meals:  []string{"meal1"},
	}

	_, err := db.Collection("foodlog").InsertOne(context.Background(), testFoodLog)
	assert.NoError(t, err)

	t.Run("Get existing food log", func(t *testing.T) {
		foodLog, err := service.GetFoodLog(testFoodLog.ID.Hex(), testFoodLog.UserID)
		assert.NoError(t, err)
		assert.NotNil(t, foodLog)
		assert.Equal(t, testFoodLog.ID, foodLog.ID)
	})

	t.Run("Get non-existing food log", func(t *testing.T) {
		nonExistingID := primitive.NewObjectID().Hex()
		foodLog, err := service.GetFoodLog(nonExistingID, "test_user")
		assert.Error(t, err)
		assert.Nil(t, foodLog)
	})
}

func TestGetFoodLogByUser(t *testing.T) {
	db := setupTestDB(t)
	service := &foodlog.FoodLogService{DB: db}

	userid := "test_user"
	// Insert test food logs
	testFoodLogs := []interface{}{
		&foodlog.FoodLog{
			ID:     primitive.NewObjectID(),
			UserID: userid,
			Date:   time.Now().Format("2006-01-02"),
			Meals:  []string{"meal1"},
		},
		&foodlog.FoodLog{
			ID:     primitive.NewObjectID(),
			UserID: userid,
			Date:   time.Now().Format("2006-01-02"),
			Meals:  []string{"meal2"},
		},
	}

	_, err := db.Collection("foodlog").InsertMany(context.Background(), testFoodLogs)
	assert.NoError(t, err)

	t.Run("Get user food logs", func(t *testing.T) {
		foodLogs, err := service.GetFoodLogByUser(userid)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(foodLogs))
	})
}

func TestGetFoodLogByUserDate(t *testing.T) {
	db := setupTestDB(t)
	service := &foodlog.FoodLogService{DB: db}

	userid := "test_user"
	// Use specific time for testing to avoid timezone issues
	now := time.Now().Truncate(24 * time.Hour) // truncate to start of day
	yesterday := now.AddDate(0, 0, -1)

	// Insert test food logs
	testFoodLogs := []interface{}{
		&foodlog.FoodLog{
			ID:     primitive.NewObjectID(),
			UserID: userid,
			Date:   now.Format("2006-01-02"), // use truncated time
			Meals:  []string{"meal1"},
		},
		&foodlog.FoodLog{
			ID:     primitive.NewObjectID(),
			UserID: userid,
			Date:   yesterday.Format("2006-01-02"),
			Meals:  []string{"meal2"},
		},
	}

	_, err := db.Collection("foodlog").InsertMany(context.Background(), testFoodLogs)
	assert.NoError(t, err)

	t.Run("Get today's food logs", func(t *testing.T) {
		dateStr := now.Format("2006-01-02") // use same truncated time
		foodLog, err := service.GetFoodLogByUserDate(userid, dateStr)
		assert.NoError(t, err)
		assert.NotNil(t, foodLog)
		assert.Equal(t, userid, foodLog.UserID)
		assert.Equal(t, []string{"meal1"}, foodLog.Meals)
	})
}

func TestDeleteFoodLog(t *testing.T) {
	db := setupTestDB(t)
	service := &foodlog.FoodLogService{DB: db}

	// Create a test food log
	testFoodLog := &foodlog.FoodLog{
		ID:     primitive.NewObjectID(),
		UserID: "test_user",
		Date:   time.Now().Format("2006-01-02"),
		Meals:  []string{"meal1"},
	}

	_, err := db.Collection("foodlog").InsertOne(context.Background(), testFoodLog)
	assert.NoError(t, err)

	t.Run("Delete existing food log", func(t *testing.T) {
		err := service.DeleteFoodLog(testFoodLog.ID.Hex(), testFoodLog.UserID)
		assert.NoError(t, err)

		// Verify deletion
		var found foodlog.FoodLog
		err = db.Collection("foodlog").FindOne(context.Background(), bson.M{"_id": testFoodLog.ID}).Decode(&found)
		assert.Error(t, err)
	})
}
