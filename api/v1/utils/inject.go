package utils

import (
	"log"

	"github.com/Npwskp/GymsbroBackend/api/v1/auth"
	"github.com/Npwskp/GymsbroBackend/api/v1/dashboard"
	"github.com/Npwskp/GymsbroBackend/api/v1/middleware"
	foodlog "github.com/Npwskp/GymsbroBackend/api/v1/nutrition/foodLog"
	"github.com/Npwskp/GymsbroBackend/api/v1/nutrition/ingredient"
	"github.com/Npwskp/GymsbroBackend/api/v1/nutrition/meal"
	minio "github.com/Npwskp/GymsbroBackend/api/v1/storage"
	"github.com/Npwskp/GymsbroBackend/api/v1/unit"
	"github.com/Npwskp/GymsbroBackend/api/v1/user"
	"github.com/Npwskp/GymsbroBackend/api/v1/workout/exercise"
	"github.com/Npwskp/GymsbroBackend/api/v1/workout/exerciseLog"
	"github.com/Npwskp/GymsbroBackend/api/v1/workout/workout"
	"github.com/Npwskp/GymsbroBackend/api/v1/workout/workoutPlan"
	"github.com/Npwskp/GymsbroBackend/api/v1/workout/workoutSession"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func InjectApp(app *fiber.App, db *mongo.Database) {
	api := app.Group("/api/v1")

	// Initialize MinIO dependencies
	minioDeps, err := minio.InjectDependencies()
	if err != nil {
		log.Fatalf("Failed to inject MinIO dependencies: %v", err)
	}

	// Public routes group (no auth required)
	public := api.Group("")
	authService := auth.AuthService{DB: db}
	authController := auth.AuthController{Instance: public, Service: &authService}
	authController.Handle() // Login, Register, etc.

	// Unit service (public)
	unitService := unit.UnitService{}
	unitController := unit.UnitController{Instance: public, Service: &unitService}
	unitController.Handle()

	// Protected routes group (requires auth)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	protected.Use(middleware.ExtractUserContext())

	// All protected controllers
	userService := user.UserService{DB: db, MinioService: minioDeps.MinioService}
	userController := user.UserController{Instance: protected, Service: &userService}
	userController.Handle()

	exerciseService := exercise.ExerciseService{DB: db, MinioService: minioDeps.MinioService}
	exerciseController := exercise.ExerciseController{Instance: protected, Service: &exerciseService}
	exerciseController.Handle()

	ingredientService := ingredient.IngredientService{DB: db, MinioService: minioDeps.MinioService}
	ingredientController := ingredient.IngredientController{Instance: protected, Service: &ingredientService}
	ingredientController.Handle()

	mealService := meal.MealService{DB: db, MinioService: minioDeps.MinioService, UnitService: &unitService}
	mealController := meal.MealController{Instance: protected, Service: &mealService}
	mealController.Handle()

	foodLogService := foodlog.FoodLogService{DB: db}
	foodLogController := foodlog.FoodLogController{Instance: protected, Service: &foodLogService}
	foodLogController.Handle()

	workoutService := workout.WorkoutService{DB: db}
	workoutController := workout.WorkoutController{Instance: protected, Service: &workoutService}
	workoutController.Handle()

	exerciseLogService := exerciseLog.ExerciseLogService{DB: db}
	exerciseLogController := exerciseLog.ExerciseLogController{Instance: protected, Service: &exerciseLogService}
	exerciseLogController.Handle()

	workoutSessionService := workoutSession.WorkoutSessionService{DB: db}
	workoutSessionController := workoutSession.WorkoutSessionController{Instance: protected, Service: &workoutSessionService}
	workoutSessionController.Handle()

	workoutPlanService := workoutPlan.WorkoutPlanService{DB: db}
	workoutPlanController := workoutPlan.WorkoutPlanController{Instance: protected, Service: &workoutPlanService}
	workoutPlanController.Handle()

	dashboardService := dashboard.DashboardService{DB: db}
	dashboardController := dashboard.DashboardController{Instance: protected, Service: &dashboardService}
	dashboardController.Handle()
}
