package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIndexes(db *mongo.Database) error {
	// Create text index for ingredients
	ingredientIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: "text"},
			{Key: "description", Value: "text"},
			{Key: "category", Value: "text"},
		},
		Options: options.Index().SetDefaultLanguage("english"),
	}
	_, err := db.Collection("ingredient").Indexes().CreateOne(context.Background(), ingredientIndex)
	if err != nil {
		return err
	}

	// Create text index for meals
	mealIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: "text"},
			{Key: "description", Value: "text"},
			{Key: "category", Value: "text"},
		},
	}
	_, err = db.Collection("meal").Indexes().CreateOne(context.Background(), mealIndex)
	if err != nil {
		return err
	}

	return nil
}
