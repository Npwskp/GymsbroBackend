package exercise

import (
	"github.com/Npwskp/GymsbroBackend/api/v1/function"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type Error error

type ExerciseController struct {
	Instance *fiber.App
	Service  IExerciseService
}

type CreateExerciseDto struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Type        []string `json:"type" validate:"required"`
	Muscle      []string `json:"muscle" validate:"required"`
	Image       string   `json:"image" validate:"required"`
}

type UpdateExerciseDto struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        []string `json:"type"`
	Muscle      []string `json:"muscle"`
	Image       string   `json:"image"`
}

// @Summary		Create an exercise
// @Description	Create an exercise
// @Tags		exercises
// @Accept		json
// @Produce		json
// @Param		exercise body CreateExerciseDto true "Create Exercise"
// @Success		201	{object} Exercise
// @Failure		400	{object} Error
// @Router		/exercises [post]
func (ec *ExerciseController) PostExerciseHandler(c *fiber.Ctx) error {
	validate := validator.New()
	exercise := new(CreateExerciseDto)
	if err := c.BodyParser(exercise); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	if err := validate.Struct(*exercise); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	for _, t := range exercise.Type {
		if !function.CheckExerciseType(t) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Type of exercise is not valid"})
		}
	}
	for _, m := range exercise.Muscle {
		if !function.CheckMuscleGroup(m) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Muscle is not valid"})
		}
	}
	createdExercise, err := ec.Service.CreateExercise(exercise)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(createdExercise)
}

// @Summary		Create many exercises
// @Description	Create many exercises
// @Tags		exercises
// @Accept		json
// @Produce		json
// @Param		exercises body []CreateExerciseDto true "Create Exercises"
// @Success		201	{object} []Exercise
// @Failure		400	{object} Error
// @Router		/exercises/many [post]
func (ec *ExerciseController) PostManyExerciseHandler(c *fiber.Ctx) error {
	validate := validator.New()
	exercises := new([]CreateExerciseDto)
	if err := c.BodyParser(exercises); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	for _, exercise := range *exercises {
		if err := validate.Struct(exercise); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		for _, t := range exercise.Type {
			if !function.CheckExerciseType(t) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Type of exercise is not valid"})
			}
		}
		for _, m := range exercise.Muscle {
			if !function.CheckMuscleGroup(m) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Muscle is not valid"})
			}
		}
	}
	createdExercises, err := ec.Service.CreateManyExercises(exercises)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(createdExercises)
}

// @Summary		Get all exercises
// @Description	Get all exercises
// @Tags		exercises
// @Accept		json
// @Produce		json
// @Success		200	{object} []Exercise
// @Failure		400	{object} Error
// @Router		/exercises [get]
func (ec *ExerciseController) GetExercisesHandler(c *fiber.Ctx) error {
	exercises, err := ec.Service.GetAllExercises()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(exercises)
}

// @Summary		Get an exercise
// @Description	Get an exercise
// @Tags		exercises
// @Accept		json
// @Produce		json
// @Param		id path string true "Exercise ID"
// @Success		200	{object} Exercise
// @Failure		400	{object} Error
// @Router		/exercises/{id} [get]
func (ec *ExerciseController) GetExerciseHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	exercise, err := ec.Service.GetExercise(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(exercise)
}

// @Summary		Get exercises by type
// @Description	Get exercises by type
// @Tags		exercises
// @Accept		json
// @Produce		json
// @Param		type path string true "Exercise Type"
// @Success		200	{object} []Exercise
// @Failure		400	{object} Error
// @Router		/exercises/type/{type} [get]
func (ec *ExerciseController) GetExerciseByTypeHandler(c *fiber.Ctx) error {
	t := c.Params("type")
	exercises, err := ec.Service.GetExerciseByType(t)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(exercises)
}

// @Summary		Delete an exercise
// @Description	Delete an exercise
// @Tags		exercises
// @Accept		json
// @Produce		json
// @Param		id path string true "Exercise ID"
// @Success		204	{object} Error
// @Failure		400	{object} Error
// @Router		/exercises/{id} [delete]
func (ec *ExerciseController) DeleteExerciseHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := ec.Service.DeleteExercise(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "Exercise deleted"})
}

// @Summary		Update an exercise
// @Description	Update an exercise
// @Tags		exercises
// @Accept		json
// @Produce		json
// @Param		id path string true "Exercise ID"
// @Param		exercise body UpdateExerciseDto true "Update Exercise"
// @Success		200	{object} Exercise
// @Failure		400	{object} Error
// @Router		/exercises/{id} [put]
func (ec *ExerciseController) UpdateExerciseHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	validate := validator.New()
	exercise := new(UpdateExerciseDto)
	if err := c.BodyParser(exercise); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	if err := validate.Struct(*exercise); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	for _, t := range exercise.Type {
		if !function.CheckExerciseType(t) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Type of exercise is not valid"})
		}
	}
	for _, m := range exercise.Muscle {
		if !function.CheckMuscleGroup(m) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Muscle is not valid"})
		}
	}
	updatedExercise, err := ec.Service.UpdateExercise(exercise, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(updatedExercise)
}

func (ec *ExerciseController) Handle() {
	g := ec.Instance.Group("/exercises")
	g.Post("/", ec.PostExerciseHandler)
	g.Post("/many", ec.PostManyExerciseHandler)
	g.Get("/", ec.GetExercisesHandler)
	g.Get("/:id", ec.GetExerciseHandler)
	g.Get("/type/:type", ec.GetExerciseByTypeHandler)
	g.Delete("/:id", ec.DeleteExerciseHandler)
	g.Put("/:id", ec.UpdateExerciseHandler)
}
