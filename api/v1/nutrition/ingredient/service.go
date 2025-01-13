package ingredient

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	minio "github.com/Npwskp/GymsbroBackend/api/v1/storage"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	IngredientImageBucketName = "ingredient-image"
)

type IngredientService struct {
	DB           *mongo.Database
	MinioService minio.MinioService
}

type IIngredientService interface {
	CreateIngredient(ingredient *CreateIngredientDto, userId string) (*Ingredient, error)
	GetIngredient(id string, userId string) (*Ingredient, error)
	GetIngredientByUser(userId string) ([]*Ingredient, error)
	DeleteIngredient(id string, userId string) error
	UpdateIngredient(doc *UpdateIngredientDto, id string, userId string) (*Ingredient, error)
	SearchFilteredIngredients(filters SearchFilters) ([]*Ingredient, error)
	UpdateIngredientImage(c *fiber.Ctx, id string, file io.Reader, filename string, contentType string, userId string) (*Ingredient, error)
}

func (is *IngredientService) CreateIngredient(ingredient *CreateIngredientDto, userId string) (*Ingredient, error) {
	ingredientModel := CreateIngredientModel(ingredient)
	ingredientModel.UserID = userId

	result, err := is.DB.Collection("ingredient").InsertOne(context.Background(), ingredientModel)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: result.InsertedID}}
	createdIngredient := &Ingredient{}
	if err := is.DB.Collection("ingredient").FindOne(context.Background(), filter).Decode(createdIngredient); err != nil {
		return nil, err
	}

	return createdIngredient, nil
}

func (is *IngredientService) GetIngredient(id string, userId string) (*Ingredient, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: objectID}}
	record := is.DB.Collection("ingredient").FindOne(context.Background(), filter)
	ingredient := &Ingredient{}
	if err := record.Decode(ingredient); err != nil {
		return nil, err
	}
	return ingredient, nil
}

func (is *IngredientService) GetIngredientByUser(userId string) ([]*Ingredient, error) {
	filter := bson.D{
		{Key: "userid", Value: userId},
		{Key: "$or", Value: []bson.M{
			{"deleted_at": bson.M{"$exists": false}},
			{"deleted_at": ""},
		}},
	}
	cursor, err := is.DB.Collection("ingredient").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	var ingredients []*Ingredient
	if err := cursor.All(context.Background(), &ingredients); err != nil {
		return nil, err
	}
	return ingredients, nil
}

func (is *IngredientService) DeleteIngredient(id string, userId string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: objectID},
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

	result, err := is.DB.Collection("ingredient").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (is *IngredientService) UpdateIngredient(doc *UpdateIngredientDto, id string, userId string) (*Ingredient, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objectID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: doc.Name},
			{Key: "category", Value: doc.Category},
			{Key: "calories", Value: doc.Calories},
			{Key: "nutrients", Value: doc.Nutrients},
			{Key: "updated_at", Value: time.Now()},
		}},
	}
	_, err = is.DB.Collection("ingredient").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}
	ingredient, err := is.GetIngredient(id, userId)
	if err != nil {
		return nil, err
	}
	return ingredient, nil
}

func (is *IngredientService) SearchFilteredIngredients(filters SearchFilters) ([]*Ingredient, error) {
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

	// Create two separate queries: one for public ingredients and one for user-specific ingredients
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
		})}}
	}

	// Execute both queries with limit
	opts := options.Find().SetLimit(20)

	// Get public ingredients
	publicCursor, err := is.DB.Collection("ingredient").Find(context.Background(), publicQuery, opts)
	if err != nil {
		return nil, fmt.Errorf("error searching public ingredients: %w", err)
	}
	defer publicCursor.Close(context.Background())

	var publicIngredients []*Ingredient
	if err := publicCursor.All(context.Background(), &publicIngredients); err != nil {
		return nil, fmt.Errorf("error decoding public ingredients: %w", err)
	}

	// Get user-specific ingredients if UserID is provided
	var userIngredients []*Ingredient
	if filters.UserID != "" {
		userCursor, err := is.DB.Collection("ingredient").Find(context.Background(), userQuery, opts)
		if err != nil {
			return nil, fmt.Errorf("error searching user ingredients: %w", err)
		}
		defer userCursor.Close(context.Background())

		if err := userCursor.All(context.Background(), &userIngredients); err != nil {
			return nil, fmt.Errorf("error decoding user ingredients: %w", err)
		}
	}

	// Combine both results
	return append(publicIngredients, userIngredients...), nil
}

func (is *IngredientService) UpdateIngredientImage(c *fiber.Ctx, id string, file io.Reader, filename string, contentType string, userId string) (*Ingredient, error) {
	// Get ingredient first to verify existence and get current image URL
	ingredient, err := is.GetIngredient(id, userId)
	if err != nil {
		return nil, err
	}

	oldImageURL := ingredient.Image

	ext := strings.ToLower(filepath.Ext(filename))
	// Generate unique filename
	timestamp := time.Now().UnixNano()
	objectName := fmt.Sprintf("ingredients/%s/image_%d%s", id, timestamp, ext)

	// Upload to MinIO
	err = is.MinioService.UploadFile(c.Context(), file, IngredientImageBucketName, objectName, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image: %v", err)
	}

	// Get the URL of the uploaded file
	url, err := is.MinioService.GetFileURL(c.Context(), IngredientImageBucketName, objectName)
	if err != nil {
		return nil, fmt.Errorf("failed to get image URL: %v", err)
	}

	// Update ingredient's image URL in database
	oid, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: oid}, {Key: "userid", Value: userId}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "image", Value: url},
			{Key: "updated_at", Value: time.Now()},
		}},
	}

	result, err := is.DB.Collection("ingredient").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("no ingredient found for the given ID")
	}

	// Delete old image after successful upload and update
	if oldImageURL != "" {
		baseURL := strings.Split(oldImageURL, "?")[0]
		urlParts := strings.Split(baseURL, is.MinioService.GetFullBucketName(IngredientImageBucketName)+"/")
		if len(urlParts) > 1 {
			oldObjectName := urlParts[1]
			if err := is.MinioService.DeleteFile(c.Context(), IngredientImageBucketName, oldObjectName); err != nil {
				fmt.Printf("Warning: Failed to delete old ingredient image: %v\n", err)
			}
		}
	}

	return is.GetIngredient(id, userId)
}
