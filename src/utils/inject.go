package utils

import (
	"github.com/Npwskp/GymsbroBackend/src/plans"
	"github.com/Npwskp/GymsbroBackend/src/user"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func InjectApp(app *fiber.App, db *mongo.Database) {
	userService := user.UserService{DB: db}
	userController := user.UserController{Instance: app, Service: &userService}
	userController.Handle()

	planService := plans.PlanService{DB: db}
	planController := plans.PlanController{Instance: app, Service: &planService}
	planController.Handle()
}
