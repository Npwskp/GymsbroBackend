package workout

import (
	"github.com/Npwskp/GymsbroBackend/api/v1/function"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type Error error

type WorkoutController struct {
	Instance fiber.Router
	Service  IWorkoutService
}

// @Summary     Create a workout
// @Description Create a new workout plan
// @Tags        workouts
// @Accept      json
// @Produce     json
// @Param       workout body CreateWorkoutDto true "Create Workout"
// @Success     201 {object} Workout
// @Failure     400 {object} Error
// @Failure     500 {object} Error
// @Router      /workout [post]
func (wc *WorkoutController) CreateWorkoutHandler(c *fiber.Ctx) error {
	validate := validator.New()
	workout := new(CreateWorkoutDto)
	userId := function.GetUserIDFromContext(c)

	if err := c.BodyParser(workout); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	if err := validate.Struct(workout); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	createdWorkout, err := wc.Service.CreateWorkout(workout, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(createdWorkout)
}

// @Summary     Get a workout
// @Description Get a workout by ID
// @Tags        workouts
// @Accept      json
// @Produce     json
// @Param       id path string true "Workout ID"
// @Success     200 {object} Workout
// @Failure     404 {object} Error
// @Failure     500 {object} Error
// @Router      /workout/{id} [get]
func (wc *WorkoutController) GetWorkoutHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	userId := function.GetUserIDFromContext(c)

	workout, err := wc.Service.GetWorkout(id, userId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Workout not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(workout)
}

// @Summary     Get user workouts
// @Description Get all workouts for the current user
// @Tags        workouts
// @Accept      json
// @Produce     json
// @Success     200 {array} Workout
// @Failure     400 {object} Error
// @Failure     500 {object} Error
// @Router      /workout [get]
func (wc *WorkoutController) GetWorkoutsHandler(c *fiber.Ctx) error {
	userId := function.GetUserIDFromContext(c)

	workouts, err := wc.Service.GetWorkoutsByUser(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	if workouts == nil {
		workouts = []*Workout{}
	}

	return c.Status(fiber.StatusOK).JSON(workouts)
}

// @Summary     Update a workout
// @Description Update a workout by ID
// @Tags        workouts
// @Accept      json
// @Produce     json
// @Param       id path string true "Workout ID"
// @Param       workout body UpdateWorkoutDto true "Update Workout"
// @Success     200 {object} Workout
// @Failure     400 {object} Error
// @Failure     404 {object} Error
// @Failure     500 {object} Error
// @Router      /workout/{id} [put]
func (wc *WorkoutController) UpdateWorkoutHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	userId := function.GetUserIDFromContext(c)
	workout := new(UpdateWorkoutDto)

	if err := c.BodyParser(workout); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	updatedWorkout, err := wc.Service.UpdateWorkout(id, workout, userId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Workout not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(updatedWorkout)
}

// @Summary     Delete a workout
// @Description Delete a workout by ID
// @Tags        workouts
// @Accept      json
// @Produce     json
// @Param       id path string true "Workout ID"
// @Success     204 {object} nil
// @Failure     404 {object} Error
// @Failure     500 {object} Error
// @Router      /workout/{id} [delete]
func (wc *WorkoutController) DeleteWorkoutHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	userId := function.GetUserIDFromContext(c)

	err := wc.Service.DeleteWorkout(id, userId)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Workout not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func (wc *WorkoutController) Handle() {
	g := wc.Instance.Group("/workout")

	g.Post("/", wc.CreateWorkoutHandler)
	g.Get("/", wc.GetWorkoutsHandler)
	g.Get("/:id", wc.GetWorkoutHandler)
	g.Put("/:id", wc.UpdateWorkoutHandler)
	g.Delete("/:id", wc.DeleteWorkoutHandler)
}
