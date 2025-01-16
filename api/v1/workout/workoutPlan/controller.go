package workoutPlan

import (
	"github.com/Npwskp/GymsbroBackend/api/v1/function"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type Error error

type WorkoutPlanController struct {
	Instance fiber.Router
	Service  IWorkoutPlanService
}

// @Summary		Create workout plan by days of week
// @Description	Create a workout plan with specific workouts for each day of the week
// @Tags		workoutPlan
// @Accept		json
// @Produce		json
// @Param		plan body CreatePlanByDaysOfWeekDto true "Create Workout Plan by Days"
// @Success		201	{array} WorkoutPlan
// @Failure		400	{object} Error
// @Failure		401	{object} Error
// @Router		/workoutPlan/byDaysOfWeek [post]
func (wp *WorkoutPlanController) CreatePlanByDaysOfWeek(c *fiber.Ctx) error {
	validate := validator.New()
	dto := new(CreatePlanByDaysOfWeekDto)
	userId := function.GetUserIDFromContext(c)

	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validate.Struct(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	workoutPlans, err := wp.Service.CreatePlanByDaysOfWeek(dto, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(workoutPlans)
}

// @Summary		Create workout plan by cyclic workouts
// @Description	Create a workout plan by cycling through a list of workouts
// @Tags		workoutPlan
// @Accept		json
// @Produce		json
// @Param		plan body CreatePlanByCyclicWorkoutDto true "Create Workout Plan by Cycle"
// @Success		201	{array} WorkoutPlan
// @Failure		400	{object} Error
// @Failure		401	{object} Error
// @Router		/workoutPlan/byCyclicWorkout [post]
func (wp *WorkoutPlanController) CreatePlanByCyclicWorkout(c *fiber.Ctx) error {
	validate := validator.New()
	dto := new(CreatePlanByCyclicWorkoutDto)
	userId := function.GetUserIDFromContext(c)

	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validate.Struct(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	workoutPlans, err := wp.Service.CreatePlanByCyclicWorkout(dto, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(workoutPlans)
}

// @Summary		Get all workout plans
// @Description	Get all workout plans for the authenticated user
// @Tags		workoutPlan
// @Accept		json
// @Produce		json
// @Success		200	{array} WorkoutPlan
// @Failure		401	{object} Error
// @Router		/workoutPlan/byUser [get]
func (wp *WorkoutPlanController) GetWorkoutPlans(c *fiber.Ctx) error {
	userId := function.GetUserIDFromContext(c)
	workoutPlans, err := wp.Service.GetWorkoutPlansByUser(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(workoutPlans)
}

func (wp *WorkoutPlanController) Handle() {
	g := wp.Instance.Group("/workoutPlan")

	g.Post("/byDaysOfWeek", wp.CreatePlanByDaysOfWeek)
	g.Post("/byCyclicWorkout", wp.CreatePlanByCyclicWorkout)
	g.Get("/byUser", wp.GetWorkoutPlans)
}
