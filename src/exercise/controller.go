package exercise

import (
	"github.com/Npwskp/GymsbroBackend/src/function"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

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

func (ec *ExerciseController) GetExercisesHandler(c *fiber.Ctx) error {
	exercises, err := ec.Service.GetAllExercises()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(exercises)
}

func (ec *ExerciseController) GetExerciseHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	exercise, err := ec.Service.GetExercise(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(exercise)
}

func (ec *ExerciseController) GetExerciseByTypeHandler(c *fiber.Ctx) error {
	t := c.Params("type")
	exercises, err := ec.Service.GetExerciseByType(t)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(exercises)
}

func (ec *ExerciseController) DeleteExerciseHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := ec.Service.DeleteExercise(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "Exercise deleted"})
}

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
	g.Get("/", ec.GetExercisesHandler)
	g.Get("/:id", ec.GetExerciseHandler)
	g.Get("/type/:type", ec.GetExerciseByTypeHandler)
	g.Delete("/:id", ec.DeleteExerciseHandler)
	g.Put("/:id", ec.UpdateExerciseHandler)
}
