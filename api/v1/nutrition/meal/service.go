package meal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Npwskp/GymsbroBackend/api/v1/nutrition/ingredient"
	"github.com/Npwskp/GymsbroBackend/api/v1/nutrition/types"
	minio "github.com/Npwskp/GymsbroBackend/api/v1/storage"
	"github.com/Npwskp/GymsbroBackend/api/v1/unit"
	"github.com/gofiber/fiber/v2"
)

const (
	MealImageBucketName = "meal-image"
)

type MealService struct {
	DB           *mongo.Database
	MinioService minio.MinioService
}

type IMealService interface {
	CreateMeal(meal *CreateMealDto, userid string) (*Meal, error)
	CalculateNutrient(body *CalculateNutrientBody, userid string) (*CalculateNutrientResponse, error)
	GetMeal(id string, userid string) (*Meal, error)
	GetMealByUser(userid string) ([]*Meal, error)
	DeleteMeal(id string, userid string) error
	UpdateMeal(doc *UpdateMealDto, id string, userid string) (*Meal, error)
	SearchFilteredMeals(filters SearchFilters) ([]*Meal, error)
	UpdateMealImage(c *fiber.Ctx, id string, file io.Reader, filename string, contentType string, userId string) (*Meal, error)
}

func (ns *MealService) CreateMeal(meal *CreateMealDto, userid string) (*Meal, error) {
	mealModel := CreateMealModel(meal)
	mealModel.UserID = userid

	result, err := ns.DB.Collection("meal").InsertOne(context.Background(), mealModel)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: result.InsertedID}}
	createdRecord := ns.DB.Collection("meal").FindOne(context.Background(), filter)
	createdMeal := &Meal{}
	if err := createdRecord.Decode(createdMeal); err != nil {
		return nil, err
	}
	return createdMeal, nil
}

func (ns *MealService) CalculateNutrient(body *CalculateNutrientBody, userid string) (*CalculateNutrientResponse, error) {
	ingredientsBody := body.Ingredients
	ingredientService := &ingredient.IngredientService{DB: ns.DB}

	ingredientIDs := make([]primitive.ObjectID, len(ingredientsBody))
	for i, ing := range ingredientsBody {
		oid, err := primitive.ObjectIDFromHex(ing.IngredientId)
		if err != nil {
			return nil, fmt.Errorf("invalid ingredient ID %s: %w", ing.IngredientId, err)
		}
		ingredientIDs[i] = oid
	}

	// Query all ingredients at once
	filter := bson.M{
		"_id": bson.M{"$in": ingredientIDs},
		"$or": []bson.M{
			{"userid": userid},
			{"userid": ""},
			{"userid": nil},
		},
	}
	var ingredients []ingredient.Ingredient
	cursor, err := ingredientService.DB.Collection("ingredient").Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("error fetching ingredients: %w", err)
	}
	if err := cursor.All(context.Background(), &ingredients); err != nil {
		return nil, fmt.Errorf("error decoding ingredients: %w", err)
	}

	// Convert slice to map for easier lookup
	ingredientMap := make(map[primitive.ObjectID]*ingredient.Ingredient)
	for i := range ingredients {
		ingredientMap[ingredients[i].ID] = &ingredients[i]
	}

	totalNutrients := make(map[string]types.Nutrient)
	totalCalories := 0.0

	for _, ingBody := range ingredientsBody {
		oid, err := primitive.ObjectIDFromHex(ingBody.IngredientId)
		if err != nil {
			return nil, fmt.Errorf("invalid ingredient ID %s: %w", ingBody.IngredientId, err)
		}

		fullIngredient, ok := ingredientMap[oid]
		if !ok {
			return nil, fmt.Errorf("ingredient not found: %s", ingBody.IngredientId)
		}

		// Convert ingredient amount to grams for calculation
		amountInGrams, err := unit.Service.ConvertBetweenUnits(ingBody.Amount, ingBody.Unit, "g")
		if err != nil {
			return nil, fmt.Errorf("error converting units for ingredient %s: %w", ingBody.IngredientId, err)
		}

		// Calculate the serving ratio (amount in grams / 100g base)
		servingRatio := amountInGrams / 100.0

		// Calculate calories based on the actual amount in grams
		totalCalories, _ = strconv.ParseFloat(fmt.Sprintf("%.5f", totalCalories+fullIngredient.Calories*servingRatio), 64)

		// Calculate nutrients based on the actual amount in grams
		if fullIngredient.Nutrients != nil {
			for _, nutrient := range fullIngredient.Nutrients {
				if existing, ok := totalNutrients[nutrient.Name]; ok {
					// Update existing nutrient with rounding
					roundedAmount, _ := strconv.ParseFloat(fmt.Sprintf("%.5f", existing.Amount+nutrient.Amount*servingRatio), 64)
					existing.Amount = roundedAmount
					totalNutrients[nutrient.Name] = existing
				} else {
					// Add new nutrient with rounding
					roundedAmount, _ := strconv.ParseFloat(fmt.Sprintf("%.5f", nutrient.Amount*servingRatio), 64)
					totalNutrients[nutrient.Name] = types.Nutrient{
						Name:   nutrient.Name,
						Amount: roundedAmount,
						Unit:   nutrient.Unit,
					}
				}
			}
		}
	}

	// Convert map directly to []types.Nutrient
	nutrientValues := make([]types.Nutrient, 0, len(totalNutrients))
	for _, nutrient := range totalNutrients {
		nutrientValues = append(nutrientValues, nutrient)
	}

	return &CalculateNutrientResponse{
		Calories:  totalCalories,
		Nutrients: nutrientValues,
	}, nil
}

func (ns *MealService) GetMeal(id string, userid string) (*Meal, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid meal ID format: %w", err)
	}

	// Allow access to both public meals and user-specific meals
	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "$or", Value: []bson.M{
			{"userid": userid}, // User's own meals
			{"userid": ""},     // Public meals
			{"userid": nil},    // Public meals (null userid)
		}},
	}

	meal := &Meal{}
	err = ns.DB.Collection("meal").FindOne(context.Background(), filter).Decode(meal)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("meal not found with ID: %s", id)
		}
		return nil, fmt.Errorf("error retrieving meal: %w", err)
	}

	return meal, nil
}

func (ns *MealService) GetMealByUser(userId string) ([]*Meal, error) {
	filter := bson.D{
		{Key: "userid", Value: userId},
		{Key: "$or", Value: []bson.M{
			{"deleted_at": bson.M{"$exists": false}},
			{"deleted_at": ""},
		}},
	}

	cursor, err := ns.DB.Collection("meal").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	meals := make([]*Meal, 0)
	if err := cursor.All(context.Background(), &meals); err != nil {
		return nil, err
	}
	return meals, nil
}

func (ns *MealService) DeleteMeal(id string, userId string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "userid", Value: userId},
		{Key: "$or", Value: []bson.M{
			{"deleted_at": bson.M{"$exists": false}},
			{"deleted_at": ""},
		}},
	}

	now := time.Now()
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "deleted_at", Value: now},
			{Key: "updated_at", Value: now},
		}},
	}

	result, err := ns.DB.Collection("meal").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (ns *MealService) UpdateMeal(doc *UpdateMealDto, id string, userid string) (*Meal, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: oid}, {Key: "userid", Value: userid}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "description", Value: doc.Description},
			{Key: "category", Value: doc.Category},
			{Key: "image", Value: doc.Image},
			{Key: "calories", Value: doc.Calories},
			{Key: "nutrients", Value: doc.Nutrients},
			{Key: "ingredients", Value: doc.Ingredients},
			{Key: "updated_at", Value: time.Now()},
		}},
	}
	if _, err := ns.DB.Collection("meal").UpdateOne(context.Background(), filter, update); err != nil {
		return nil, err
	}
	return ns.GetMeal(id, userid)
}

func (ns *MealService) SearchFilteredMeals(filters SearchFilters) ([]*Meal, error) {
	// Add not-deleted condition to base conditions
	baseConditions := []bson.D{
		{{Key: "$or", Value: []bson.M{
			{"deleted_at": bson.M{"$exists": false}},
			{"deleted_at": ""},
		}}},
	}

	// Add name search if query is provided
	if filters.Query != "" {
		baseConditions = append(baseConditions, bson.D{
			{Key: "name", Value: bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: filters.Query, Options: "i"}}}},
		})
	}

	// Add category filter
	if filters.Category != "" {
		baseConditions = append(baseConditions, bson.D{
			{Key: "category", Value: filters.Category},
		})
	}

	// Add calories range filter
	if filters.MinCalories > 0 || filters.MaxCalories > 0 {
		caloriesFilter := bson.D{}
		if filters.MinCalories > 0 {
			caloriesFilter = append(caloriesFilter, bson.E{Key: "$gte", Value: filters.MinCalories})
		}
		if filters.MaxCalories > 0 {
			caloriesFilter = append(caloriesFilter, bson.E{Key: "$lte", Value: filters.MaxCalories})
		}
		baseConditions = append(baseConditions, bson.D{
			{Key: "calories", Value: caloriesFilter},
		})
	}

	// Add nutrients filter
	if filters.Nutrients != "" {
		nutrients := strings.Split(filters.Nutrients, ",")
		nutrientFilters := bson.A{}
		for _, nutrient := range nutrients {
			nutrient = strings.TrimSpace(nutrient)
			if nutrient != "" {
				nutrientFilters = append(nutrientFilters, bson.M{
					"nutrients.name": bson.M{
						"$regex":   nutrient,
						"$options": "i",
					},
				})
			}
		}
		if len(nutrientFilters) > 0 {
			baseConditions = append(baseConditions, bson.D{
				{Key: "$or", Value: nutrientFilters},
			})
		}
	}

	// Create two separate queries: one for public meals and one for user-specific meals
	publicQuery := bson.D{{Key: "$and", Value: append(baseConditions, bson.D{
		{Key: "$or", Value: []bson.D{
			{{Key: "userid", Value: ""}},
			{{Key: "userid", Value: primitive.Null{}}},
		}},
	})}}

	userQuery := bson.D{}
	if filters.UserID != "" {
		userQuery = bson.D{{Key: "$and", Value: append(baseConditions, bson.D{
			{Key: "userid", Value: filters.UserID},
			{Key: "is_quick_add", Value: false},
		})}}
	}

	// Execute both queries with limit
	opts := options.Find().SetLimit(20)

	// Get public meals
	publicCursor, err := ns.DB.Collection("meal").Find(context.Background(), publicQuery, opts)
	if err != nil {
		return nil, fmt.Errorf("error searching public meals: %w", err)
	}
	defer publicCursor.Close(context.Background())

	publicMeals := make([]*Meal, 0)
	if err := publicCursor.All(context.Background(), &publicMeals); err != nil {
		return nil, fmt.Errorf("error decoding public meals: %w", err)
	}

	// Get user-specific meals if UserID is provided
	userMeals := make([]*Meal, 0)
	if filters.UserID != "" {
		userCursor, err := ns.DB.Collection("meal").Find(context.Background(), userQuery, opts)
		if err != nil {
			return nil, fmt.Errorf("error searching user meals: %w", err)
		}
		defer userCursor.Close(context.Background())

		if err := userCursor.All(context.Background(), &userMeals); err != nil {
			return nil, fmt.Errorf("error decoding user meals: %w", err)
		}
	}

	// Combine both results
	return append(publicMeals, userMeals...), nil
}

func (ns *MealService) UpdateMealImage(c *fiber.Ctx, id string, file io.Reader, filename string, contentType string, userId string) (*Meal, error) {
	// Get meal first to verify existence and get current image URL
	meal, err := ns.GetMeal(id, userId)
	if err != nil {
		return nil, err
	}

	oldImageURL := meal.Image

	ext := strings.ToLower(filepath.Ext(filename))
	// Generate unique filename
	timestamp := time.Now().UnixNano()
	objectName := fmt.Sprintf("meals/%s/image_%d%s", id, timestamp, ext)

	// Upload to MinIO
	err = ns.MinioService.UploadFile(c.Context(), file, MealImageBucketName, objectName, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image: %v", err)
	}

	// Get the URL of the uploaded file
	url, err := ns.MinioService.GetFileURL(c.Context(), MealImageBucketName, objectName)
	if err != nil {
		return nil, fmt.Errorf("failed to get image URL: %v", err)
	}

	// Update meal's image URL in database
	oid, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: oid}, {Key: "userid", Value: userId}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "image", Value: url},
			{Key: "updated_at", Value: time.Now()},
		}},
	}

	result, err := ns.DB.Collection("meal").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("no meal found for the given ID")
	}

	// Delete old image after successful upload and update
	if oldImageURL != "" {
		baseURL := strings.Split(oldImageURL, "?")[0]
		urlParts := strings.Split(baseURL, ns.MinioService.GetFullBucketName(MealImageBucketName)+"/")
		if len(urlParts) > 1 {
			oldObjectName := urlParts[1]
			if err := ns.MinioService.DeleteFile(c.Context(), MealImageBucketName, oldObjectName); err != nil {
				fmt.Printf("Warning: Failed to delete old meal image: %v\n", err)
			}
		}
	}

	return ns.GetMeal(id, userId)
}
