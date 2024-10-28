package function

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckOwnership(db *mongo.Database, id string, userid string, collection string, objType interface{}) (bool, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	filter := bson.D{{Key: "_id", Value: objectID}}
	record := db.Collection("ingredient").FindOne(context.Background(), filter)
	data := &objType
	if err := record.Decode(data); err != nil {
		return false, err
	}

	value := reflect.ValueOf(data)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return false, fmt.Errorf("expected a struct, but got %s", value.Kind())
	}

	fieldvalue := value.FieldByName("UserID")
	if !fieldvalue.IsValid() {
		return false, fmt.Errorf("no such field: UserID in obj")
	}

	if str, ok := fieldvalue.Interface().(string); ok {
		if str == userid {
			return true, nil
		} else {
			return false, errors.New("update failed: user does not own this entity")
		}
	} else {
		return false, errors.New("field is not a string")
	}
}
