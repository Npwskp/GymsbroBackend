package function

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getConcreteValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	return v
}

func CheckOwnership(db *mongo.Database, id string, userid string, collection string, objType interface{}) (bool, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	filter := bson.D{{Key: "_id", Value: objectID}}
	record := db.Collection(collection).FindOne(context.Background(), filter)

	// Create a new pointer to the same type as objType
	t := reflect.TypeOf(objType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	value := reflect.New(t).Interface()

	if err := record.Decode(value); err != nil {
		return false, err
	}

	// Get the concrete value
	v := reflect.ValueOf(value)
	v = getConcreteValue(v)

	// Look for UserID field (case insensitive)
	var userIDField reflect.Value
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		if strings.EqualFold(field.Name, "userid") {
			userIDField = v.Field(i)
			break
		}
	}

	if !userIDField.IsValid() {
		return false, fmt.Errorf("no such field: UserID in obj")
	}

	if str, ok := userIDField.Interface().(string); ok {
		return str == userid, nil
	}

	return false, errors.New("UserID field is not a string")
}
