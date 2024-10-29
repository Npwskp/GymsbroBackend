package utils

import (
	"github.com/Npwskp/GymsbroBackend/api/v1/auth"
	"github.com/Npwskp/GymsbroBackend/api/v1/exercise"
	"github.com/Npwskp/GymsbroBackend/api/v1/middleware"
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

	// Public routes group (no auth required)
	public := api.Group("")
	authService := auth.AuthService{DB: db}
	authController := auth.AuthController{Instance: public, Service: &authService}
	authController.Handle() // Login, Register, etc.

	// Protected routes group (requires auth)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	protected.Use(middleware.ExtractUserContext())

	// All protected controllers
	userService := user.UserService{DB: db}
	userController := user.UserController{Instance: protected, Service: &userService}
	userController.Handle()

	planService := plans.PlanService{DB: db}
	planController := plans.PlanController{Instance: protected, Service: &planService}
	planController.Handle()

	exerciseService := exercise.ExerciseService{DB: db}
	exerciseController := exercise.ExerciseController{Instance: protected, Service: &exerciseService}
	exerciseController.Handle()

	ingredientService := ingredient.IngredientService{DB: db}
	ingredientController := ingredient.IngredientController{Instance: protected, Service: &ingredientService}
	ingredientController.Handle()

	mealService := meal.MealService{DB: db}
	mealController := meal.MealController{Instance: protected, Service: &mealService}
	mealController.Handle()

	foodLogService := foodlog.FoodLogService{DB: db}
	foodLogController := foodlog.FoodLogController{Instance: protected, Service: &foodLogService}
	foodLogController.Handle()
}
