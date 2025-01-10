package exercise

import (
	"fmt"
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
// @Failure		401	{object} Error
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
// @Success		201	{array} Exercise
// @Failure		400	{object} Error
// @Failure		401	{object} Error
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
// @Success		200	{array} Exercise
// @Failure		400	{object} Error
// @Failure		401	{object} Error
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
// @Failure		401	{object} Error
// @Failure		404	{object} Error
// @Router		/exercise/{id} [get]
func (ec *ExerciseController) GetExerciseHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	userId := function.GetUserIDFromContext(c)

	exercise, err := ec.Service.GetExercise(id, userId)
	if err != nil {
		if err.Error() == "exercise not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Exercise not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(exercise)
}

// @Summary		Get all equipment
// @Description	Get all equipment
// @Tags		exercises
// @Accept		json
// @Produce		json
// @Success		200	{array} exerciseEnums.Equipment
// @Failure		400	{object} Error
// @Failure		401	{object} Error
// @Router		/exercise/equipment [get]
func (ec *ExerciseController) GetAllEquipmentHandler(c *fiber.Ctx) error {
	equipment := exerciseEnums.GetAllEquipment()
	return c.Status(fiber.StatusOK).JSON(equipment)
}

// @Summary		Get all mechanics
// @Description	Get all mechanics
// @Tags		exercises
// @Accept		json
// @Produce		json
// @Success		200	{array} exerciseEnums.Mechanics
// @Failure		400	{object} Error
// @Failure		401	{object} Error
// @Router		/exercise/mechanics [get]
func (ec *ExerciseController) GetAllMechanicsHandler(c *fiber.Ctx) error {
	mechanics := exerciseEnums.GetAllMechanics()
	return c.Status(fiber.StatusOK).JSON(mechanics)
}

// @Summary		Get all force
// @Description	Get all force
// @Tags		exercises
// @Accept		json
// @Produce		json
// @Success		200	{array} exerciseEnums.Force
// @Failure		400	{object} Error
// @Failure		401	{object} Error
// @Router		/exercise/force [get]
func (ec *ExerciseController) GetAllForceHandler(c *fiber.Ctx) error {
	force := exerciseEnums.GetAllForces()
	return c.Status(fiber.StatusOK).JSON(force)
}

// @Summary		Get all body parts
// @Description	Get all body parts
// @Tags		exercises
// @Accept		json
// @Produce		json
// @Success		200	{array} exerciseEnums.BodyPart
// @Failure		400	{object} Error
// @Failure		401	{object} Error
// @Router		/exercise/bodypart [get]
func (ec *ExerciseController) GetAllBodyPartHandler(c *fiber.Ctx) error {
	bodyPart := exerciseEnums.GetAllBodyParts()
	return c.Status(fiber.StatusOK).JSON(bodyPart)
}

// @Summary		Get all target muscles
// @Description	Get all target muscles
// @Tags		exercises
// @Accept		json
// @Produce		json
// @Success		200	{array} exerciseEnums.TargetMuscle
// @Failure		400	{object} Error
// @Failure		401	{object} Error
// @Router		/exercise/targetmuscle [get]
func (ex *ExerciseController) GetAllTargetMusclesHandler(c *fiber.Ctx) error {
	targetMuscle := exerciseEnums.GetAllTargetMuscles()
	return c.Status(fiber.StatusOK).JSON(targetMuscle)
}

// @Summary		Delete an exercise
// @Description	Delete an exercise
// @Tags		exercises
// @Accept		json
// @Produce		json
// @Param		id path string true "Exercise ID"
// @Success		204	{object} Error
// @Failure		400	{object} Error
// @Failure		401	{object} Error
// @Failure		404	{object} Error
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
// @Failure		401	{object} Error
// @Failure		404	{object} Error
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
	updatedExercise, err := ec.Service.UpdateExercise(exercise, id, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(updatedExercise)
}

// @Summary     Update exercise image
// @Description Update an exercise's image
// @Tags        exercises
// @Accept      multipart/form-data
// @Produce     json
// @Param       id path string true "Exercise ID"
// @Param       file formData file true "Image file"
// @Success     200 {object} Exercise
// @Failure     400 {object} Error
// @Failure     401 {object} Error
// @Failure     404 {object} Error
// @Router      /exercise/{id}/image [patch]
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

// @Summary     Search and filter exercises
// @Description Search exercises by name, types and muscle groups
// @Tags        exercises
// @Accept      json
// @Produce     json
// @Param       query query string false "Search query for exercise name"
// @Param       equipment query string false "Equipment types (comma-separated)"
// @Param       mechanics query string false "Mechanics types (comma-separated)"
// @Param       force query string false "Force types (comma-separated)"
// @Param       body_part query string false "Body parts (comma-separated)"
// @Param       target_muscle query string false "Target muscles (comma-separated)"
// @Success     200 {array} Exercise
// @Failure     400 {object} Error
// @Failure     401 {object} Error
// @Router      /exercise/search [get]
func (ec *ExerciseController) SearchAndFilterExerciseHandler(c *fiber.Ctx) error {
	var filters SearchExerciseFilters
	if err := c.QueryParser(&filters); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	userId := function.GetUserIDFromContext(c)

	var equipmentList []exerciseEnums.Equipment
	var mechanicsList []exerciseEnums.Mechanics
	var forceList []exerciseEnums.Force
	var bodyPartList []exerciseEnums.BodyPart
	var targetMuscleList []exerciseEnums.TargetMuscle

	// Validate and convert equipment
	if filters.Equipment != "" {
		for _, eq := range strings.Split(strings.TrimSpace(filters.Equipment), ",") {
			eq = strings.TrimSpace(eq)
			if !exerciseEnums.IsValidEquipment(eq) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": fmt.Sprintf("Invalid equipment value: %s", eq),
				})
			}
			equipmentList = append(equipmentList, exerciseEnums.Equipment(eq))
		}
	}

	// Validate and convert mechanics
	if filters.Mechanics != "" {
		for _, m := range strings.Split(strings.TrimSpace(filters.Mechanics), ",") {
			m = strings.TrimSpace(m)
			if !exerciseEnums.IsValidMechanics(m) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": fmt.Sprintf("Invalid mechanics value: %s", m),
				})
			}
			mechanicsList = append(mechanicsList, exerciseEnums.Mechanics(m))
		}
	}

	// Validate and convert force
	if filters.Force != "" {
		for _, f := range strings.Split(strings.TrimSpace(filters.Force), ",") {
			f = strings.TrimSpace(f)
			if !exerciseEnums.IsValidForce(f) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": fmt.Sprintf("Invalid force value: %s", f),
				})
			}
			forceList = append(forceList, exerciseEnums.Force(f))
		}
	}

	// Validate and convert body parts
	if filters.BodyPart != "" {
		for _, bp := range strings.Split(strings.TrimSpace(filters.BodyPart), ",") {
			bp = strings.TrimSpace(bp)
			if !exerciseEnums.IsValidBodyPart(bp) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": fmt.Sprintf("Invalid body part value: %s", bp),
				})
			}
			bodyPartList = append(bodyPartList, exerciseEnums.BodyPart(bp))
		}
	}

	// Validate and convert target muscles
	if filters.TargetMuscle != "" {
		for _, tm := range strings.Split(strings.TrimSpace(filters.TargetMuscle), ",") {
			tm = strings.TrimSpace(tm)
			if !exerciseEnums.IsValidTargetMuscle(tm) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": fmt.Sprintf("Invalid target muscle value: %s", tm),
				})
			}
			targetMuscleList = append(targetMuscleList, exerciseEnums.TargetMuscle(tm))
		}
	}

	exercises, err := ec.Service.SearchAndFilterExercise(equipmentList, mechanicsList, forceList, bodyPartList, targetMuscleList, filters.Query, userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if exercises == nil {
		exercises = []*Exercise{}
	}

	return c.JSON(exercises)
}

func (ec *ExerciseController) Handle() {
	g := ec.Instance.Group("/exercise")

	g.Post("/", ec.PostExerciseHandler)
	g.Post("/many", ec.PostManyExerciseHandler)
	g.Get("/", ec.GetExercisesHandler)
	g.Get("/equipment", ec.GetAllEquipmentHandler)
	g.Get("/mechanics", ec.GetAllMechanicsHandler)
	g.Get("/force", ec.GetAllForceHandler)
	g.Get("/bodypart", ec.GetAllBodyPartHandler)
	g.Get("/targetmuscle", ec.GetAllTargetMusclesHandler)
	g.Get("/search", ec.SearchAndFilterExerciseHandler)
	g.Get("/:id", ec.GetExerciseHandler)
	g.Delete("/:id", ec.DeleteExerciseHandler)
	g.Put("/:id", ec.UpdateExerciseHandler)
	g.Patch("/:id/image", ec.UpdateExerciseImageHandler)
}
