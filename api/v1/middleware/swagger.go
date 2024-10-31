package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SwaggerMiddleware() fiber.Handler {
	return swagger.New(swagger.Config{
		Title: "GymsBro API",
	})
}
