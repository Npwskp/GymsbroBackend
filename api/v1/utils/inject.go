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
	api := app.Group("/api/v1")

	userService := user.UserService{DB: db}
	userController := user.UserController{Instance: api, Service: &userService}
	userController.Handle()

	planService := plans.PlanService{DB: db}
	planController := plans.PlanController{Instance: api, Service: &planService}
	planController.Handle()

	exerciseService := exercise.ExerciseService{DB: db}
	exerciseController := exercise.ExerciseController{Instance: api, Service: &exerciseService}
	exerciseController.Handle()

	ingredientService := ingredient.IngredientService{DB: db}
	ingredientController := ingredient.IngredientController{Instance: api, Service: &ingredientService}
	ingredientController.Handle()

	mealService := meal.MealService{DB: db}
	mealController := meal.MealController{Instance: api, Service: &mealService}
	mealController.Handle()

	foodLogService := foodlog.FoodLogService{DB: db}
	foodLogController := foodlog.FoodLogController{Instance: api, Service: &foodLogService}
	foodLogController.Handle()

	authService := auth.AuthService{DB: db}
	authController := auth.AuthController{Instance: api, Service: &authService}
	authController.Handle()
}
