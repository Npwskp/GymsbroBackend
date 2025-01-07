package exercise

import (
	"strings"

	"github.com/Npwskp/GymsbroBackend/api/v1/function"
	exerciseEnums "github.com/Npwskp/GymsbroBackend/api/v1/workout/exercise/enums"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type Error error

type ExerciseController struct {
	Instance fiber.Router
	Service  IExerciseService
}

// @Summary		Create an exercise
// @Description	Create an exercise
// @Tags		exercises
// @Accept		json
// @Produce		json
// @Param		exercise body CreateExerciseDto true "Create Exercise"
// @Success		201	{object} Exercise
// @Failure		400	{object} Error
// @Router		/exercise [post]
func (ec *ExerciseController) PostExerciseHandler(c *fiber.Ctx) error {
	validate := validator.New()
	exercise := new(CreateExerciseDto)
	userId := function.GetUserIDFromContext(c)

	if err := c.BodyParser(exercise); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	if err := validate.Struct(*exercise); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	for _, t := range exercise.Type {
		if _, err := exerciseEnums.ParseExerciseType(t); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Type of exercise is not valid"})
		}
	}
	for _, m := range exercise.Muscle {
		if _, err := exerciseEnums.ParseMuscleGroup(m); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Muscle is not valid"})
		}
	}
	createdExercise, err := ec.Service.CreateExercise(exercise, userId)
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
// @Router		/exercise/many [post]
func (ec *ExerciseController) PostManyExerciseHandler(c *fiber.Ctx) error {
	validate := validator.New()
	exercises := new([]CreateExerciseDto)
	userId := function.GetUserIDFromContext(c)

	if err := c.BodyParser(exercises); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	for _, exercise := range *exercises {
		if err := validate.Struct(exercise); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		for _, t := range exercise.Type {
			if _, err := exerciseEnums.ParseExerciseType(t); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Type of exercise is not valid"})
			}
		}
		for _, m := range exercise.Muscle {
			if _, err := exerciseEnums.ParseMuscleGroup(m); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Muscle is not valid"})
			}
		}
	}
	createdExercises, err := ec.Service.CreateManyExercises(exercises, userId)
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
// @Router		/exercise [get]
func (ec *ExerciseController) GetExercisesHandler(c *fiber.Ctx) error {
	userId := function.GetUserIDFromContext(c)
	exercises, err := ec.Service.GetAllExercises(userId)
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
// @Router		/exercise/{id} [get]
func (ec *ExerciseController) GetExerciseHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	userId := function.GetUserIDFromContext(c)
	exercise, err := ec.Service.GetExercise(id, userId)
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
// @Router		/exercise/type/{type} [get]
func (ec *ExerciseController) GetExerciseByTypeHandler(c *fiber.Ctx) error {
	t := c.Params("type")
	userId := function.GetUserIDFromContext(c)
	exercises, err := ec.Service.GetExerciseByType(t, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(exercises)
}

// @Summary		Get all exercise types
// @Description	Get all exercise types
// @Tags		exercises
// @Accept		json
// @Produce		json
// @Success		200	{array}	exerciseEnums.ExerciseType
// @Failure		400	{object} Error
// @Router		/exercise/types [get]
func (ec *ExerciseController) GetAllExerciseTypesHandler(c *fiber.Ctx) error {
	types := exerciseEnums.GetAllExerciseTypes()
	return c.Status(fiber.StatusOK).JSON(types)
}

// @Summary		Get all muscle groups
// @Description	Get all muscle groups
// @Tags		exercises
// @Accept		json
// @Produce		json
// @Success		200	{array}	exerciseEnums.MuscleGroup
// @Failure		400	{object} Error
// @Router		/exercise/muscles [get]
func (ec *ExerciseController) GetAllMuscleGroupsHandler(c *fiber.Ctx) error {
	muscles := exerciseEnums.GetAllMuscleGroups()
	return c.Status(fiber.StatusOK).JSON(muscles)
}

// @Summary		Delete an exercise
// @Description	Delete an exercise
// @Tags		exercises
// @Accept		json
// @Produce		json
// @Param		id path string true "Exercise ID"
// @Success		204	{object} Error
// @Failure		400	{object} Error
// @Router		/exercise/{id} [delete]
func (ec *ExerciseController) DeleteExerciseHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	userId := function.GetUserIDFromContext(c)
	if err := ec.Service.DeleteExercise(id, userId); err != nil {
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
// @Router		/exercise/{id} [put]
func (ec *ExerciseController) UpdateExerciseHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	userId := function.GetUserIDFromContext(c)
	validate := validator.New()
	exercise := new(UpdateExerciseDto)
	if err := c.BodyParser(exercise); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	if err := validate.Struct(*exercise); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	for _, t := range exercise.Type {
		if _, err := exerciseEnums.ParseExerciseType(t); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Type of exercise is not valid"})
		}
	}
	for _, m := range exercise.Muscle {
		if _, err := exerciseEnums.ParseMuscleGroup(m); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Muscle is not valid"})
		}
	}
	updatedExercise, err := ec.Service.UpdateExercise(exercise, id, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(updatedExercise)
}

// @Summary Update exercise image
// @Description Update an exercise's image
// @Tags exercises
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Exercise ID"
// @Param file formData file true "Image file"
// @Success 200 {object} Exercise
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Router /exercise/{id}/image [put]
func (ec *ExerciseController) UpdateExerciseImageHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	userId := function.GetUserIDFromContext(c)

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	// Check file type
	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File must be an image",
		})
	}

	// Open the file
	fileContent, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process file",
		})
	}
	defer fileContent.Close()

	exercise, err := ec.Service.UpdateExerciseImage(c, id, fileContent, file.Filename, contentType, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(exercise)
}

func (ec *ExerciseController) Handle() {
	g := ec.Instance.Group("/exercise")

	g.Post("/", ec.PostExerciseHandler)
	g.Post("/many", ec.PostManyExerciseHandler)
	g.Get("/", ec.GetExercisesHandler)
	g.Get("/types", ec.GetAllExerciseTypesHandler)
	g.Get("/muscles", ec.GetAllMuscleGroupsHandler)
	g.Get("/:id", ec.GetExerciseHandler)
	g.Get("/type/:type", ec.GetExerciseByTypeHandler)
	g.Delete("/:id", ec.DeleteExerciseHandler)
	g.Put("/:id", ec.UpdateExerciseHandler)
	g.Put("/:id/image", ec.UpdateExerciseImageHandler)
}
