package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/npwskp/GymsbroBackend/src/user"
	"go.mongodb.org/mongo-driver/mongo"
)

func InjectApp(app *fiber.App, db *mongo.Database) {
	userService := user.UserService{DB: db}
	userController := user.UserController{Instance: app, Service: &userService}
	userController.Handle()
}
