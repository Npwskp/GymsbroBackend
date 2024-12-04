package meal

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Npwskp/GymsbroBackend/api/v1/function"
	"github.com/Npwskp/GymsbroBackend/api/v1/nutrition/ingredient"
	"github.com/Npwskp/GymsbroBackend/api/v1/nutrition/types"
)

type MealService struct {
	DB *mongo.Database
}

type IMealService interface {
	CreateMeal(meal *CreateMealDto, userid string) (*Meal, error)
	CalculateNutrient(body *CalculateNutrientBody, userid string) (*CalculateNutrientResponse, error)
	GetAllMeals(userid string) ([]*Meal, error)
	GetMeal(id string, userid string) (*Meal, error)
	GetMealByUser(userid string) ([]*Meal, error)
	DeleteMeal(id string, userid string) error
	UpdateMeal(doc *UpdateMealDto, id string, userid string) (*Meal, error)
	SearchFilteredMeals(filters SearchFilters) ([]*Meal, error)
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
		amountInGrams, err := function.ConvertUnit(ingBody.Amount, ingBody.Unit, "g")
		if err != nil {
			return nil, fmt.Errorf("error converting units for ingredient %s: %w", ingBody.IngredientId, err)
		}

		// Calculate the serving ratio (amount in grams / 100g base)
		servingRatio := amountInGrams / 100.0

		// Calculate calories based on the actual amount in grams
		totalCalories, _ = strconv.ParseFloat(fmt.Sprintf("%.5f", fullIngredient.Calories*servingRatio), 64)

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

func (ns *MealService) GetAllMeals(userid string) ([]*Meal, error) {
	filter := bson.D{{Key: "userid", Value: userid}, {Key: "userid", Value: primitive.Null{}}, {Key: "userid", Value: ""}}
	cursor, err := ns.DB.Collection("meal").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	var Meals []*Meal
	if err := cursor.All(context.Background(), &Meals); err != nil {
		return nil, err
	}
	return Meals, nil
}

func (ns *MealService) GetMeal(id string, userid string) (*Meal, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: oid}, {Key: "userid", Value: userid}}
	meal := &Meal{}
	if err := ns.DB.Collection("meal").FindOne(context.Background(), filter).Decode(meal); err != nil {
		return nil, err
	}
	return meal, nil
}

func (ns *MealService) GetMealByUser(userid string) ([]*Meal, error) {
	filter := bson.D{{Key: "userid", Value: userid}}
	cursor, err := ns.DB.Collection("meal").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var meals []*Meal
	if err := cursor.All(context.Background(), &meals); err != nil {
		return nil, err
	}
	return meals, nil
}

func (ns *MealService) DeleteMeal(id string, userid string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: oid}, {Key: "userid", Value: userid}}
	if _, err := ns.DB.Collection("meal").DeleteOne(context.Background(), filter); err != nil {
		return err
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
			{Key: "carb", Value: doc.Carb},
			{Key: "protein", Value: doc.Protein},
			{Key: "fat", Value: doc.Fat},
			{Key: "calories", Value: doc.Calories},
			{Key: "updated_at", Value: time.Now()},
		}},
	}
	if _, err := ns.DB.Collection("meal").UpdateOne(context.Background(), filter, update); err != nil {
		return nil, err
	}
	return ns.GetMeal(id, userid)
}

func (ns *MealService) SearchFilteredMeals(filters SearchFilters) ([]*Meal, error) {
	filterQuery := bson.D{}
	andConditions := []bson.D{}

	// Add name search if query is provided
	if filters.Query != "" {
		andConditions = append(andConditions, bson.D{
			{Key: "name", Value: bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: filters.Query, Options: "i"}}}},
		})
	}

	// Add user filter
	if filters.UserID != "" {
		andConditions = append(andConditions, bson.D{
			{Key: "$or", Value: []bson.D{
				{{Key: "userid", Value: ""}},
				{{Key: "userid", Value: filters.UserID}},
				{{Key: "userid", Value: primitive.Null{}}},
			}},
		})
	}

	// Add category filter
	if filters.Category != "" {
		andConditions = append(andConditions, bson.D{
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
		andConditions = append(andConditions, bson.D{
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
			andConditions = append(andConditions, bson.D{
				{Key: "$or", Value: nutrientFilters},
			})
		}
	}

	// Combine all conditions
	if len(andConditions) > 0 {
		filterQuery = bson.D{{Key: "$and", Value: andConditions}}
	}

	// Execute query with limit
	opts := options.Find().SetLimit(20)
	cursor, err := ns.DB.Collection("meal").Find(context.Background(), filterQuery, opts)
	if err != nil {
		return nil, fmt.Errorf("error searching meals: %w", err)
	}
	defer cursor.Close(context.Background())

	var meals []*Meal
	if err := cursor.All(context.Background(), &meals); err != nil {
		return nil, fmt.Errorf("error decoding meals: %w", err)
	}

	return meals, nil
}
