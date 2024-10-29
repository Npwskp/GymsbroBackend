package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserClaims struct {
	UserID   primitive.ObjectID `json:"sub"`
	Username string             `json:"username"`
	Email    string             `json:"email"`
}

func ExtractUserContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user from JWT claims (already validated by AuthMiddleware)
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)

		// Convert string ID to ObjectID
		userID, _ := primitive.ObjectIDFromHex(claims["sub"].(string))

		// Create user claims
		userClaims := &UserClaims{
			UserID:   userID,
			Username: claims["username"].(string),
			Email:    claims["email"].(string),
		}

		// Attach to context
		c.Locals("userClaims", userClaims)

		return c.Next()
	}
}
