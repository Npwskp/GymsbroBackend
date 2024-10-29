package middleware

import (
	"github.com/Npwskp/GymsbroBackend/api/v1/config"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func AuthMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:  []byte(config.GetJWTSecret()),
		TokenLookup: "cookie:jwt",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			return c.Next()
		},
	})
}
