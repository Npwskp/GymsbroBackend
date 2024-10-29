package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetCurrentUser returns the current user claims from context
func GetCurrentUser(c *fiber.Ctx) (*UserClaims, error) {
	user, ok := c.Locals("userClaims").(*UserClaims)
	if !ok {
		return nil, errors.New("unauthorized")
	}
	return user, nil
}

// GetUserID returns just the userID from context
func GetUserID(c *fiber.Ctx) (primitive.ObjectID, error) {
	user, err := GetCurrentUser(c)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return user.UserID, nil
}
