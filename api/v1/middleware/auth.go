package middleware

import (
	"fmt"

	"github.com/Npwskp/GymsbroBackend/api/v1/config"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func AuthMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:  []byte(config.GetJWTSecret()),
		TokenLookup: "cookie:jwt",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			fmt.Printf("Auth failed for path: %s\n", c.Path())
			fmt.Printf("Cookies present: %v\n", c.Cookies("jwt"))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			fmt.Printf("Auth succeeded for path: %s\n", c.Path())
			fmt.Printf("JWT Token: %v\n", c.Cookies("jwt"))
			return c.Next()
		},
	})
}
