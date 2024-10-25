package utils

import (
	"github.com/Npwskp/GymsbroBackend/api/v1/auth"
	"github.com/Npwskp/GymsbroBackend/api/v1/exercise"
	foodlog "github.com/Npwskp/GymsbroBackend/api/v1/nutrition/foodLog"
	"github.com/Npwskp/GymsbroBackend/api/v1/nutrition/ingredient"
	"github.com/Npwskp/GymsbroBackend/api/v1/nutrition/meal"
	"github.com/Npwskp/GymsbroBackend/api/v1/plans"
	"github.com/Npwskp/GymsbroBackend/api/v1/user"
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

	exerciseService := exercise.ExerciseService{DB: db}
	exerciseController := exercise.ExerciseController{Instance: app, Service: &exerciseService}
	exerciseController.Handle()

	ingredientService := ingredient.IngredientService{DB: db}
	ingredientController := ingredient.IngredientController{Instance: app, Service: &ingredientService}
	ingredientController.Handle()

	mealService := meal.MealService{DB: db}
	mealController := meal.MealController{Instance: app, Service: &mealService}
	mealController.Handle()

	foodLogService := foodlog.FoodLogService{DB: db}
	foodLogController := foodlog.FoodLogController{Instance: app, Service: &foodLogService}
	foodLogController.Handle()

	authService := auth.AuthService{DB: db}
	authController := auth.AuthController{Instance: app, Service: &authService}
	authController.Handle()
}
