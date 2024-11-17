package function

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
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

func CheckOwnership(
	db *mongo.Database,
	id string,
	userid string,
	collection string,
	objType interface{},
) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
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
		return err
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
		return fmt.Errorf("no such field: userId in obj")
	}

	if str, ok := userIDField.Interface().(string); ok {
		if str == userid {
			return nil
		} else {
			return errors.New("user does not own this entity")
		}
	}

	return errors.New("userId field is not a string")
}

func GetUserIDFromContext(c *fiber.Ctx) string {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userId := claims["sub"].(string)
	return userId
}
