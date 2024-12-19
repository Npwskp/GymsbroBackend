package exerciseLog

import (
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/function"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type Error error

type ExerciseLogController struct {
	Instance fiber.Router
	Service  IExerciseLogService
}

// @Summary     Create exercise log
// @Description Log a completed exercise
// @Tags        exerciseLogs
// @Accept      json
// @Produce     json
// @Param       log body CreateExerciseLogDto true "Exercise Log"
// @Success     201 {object} ExerciseLog
// @Failure     400 {object} Error
// @Router      /workout/log [post]
func (c *ExerciseLogController) CreateLogHandler(ctx *fiber.Ctx) error {
	userId := function.GetUserIDFromContext(ctx)
	validate := validator.New()
	dto := new(CreateExerciseLogDto)

	if err := ctx.BodyParser(dto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(dto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	log, err := c.Service.CreateLog(dto, userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(log)
}

// @Summary     Get user logs
// @Description Get all exercise logs for a user
// @Tags        exerciseLogs
// @Accept      json
// @Produce     json
// @Success     200 {array} ExerciseLog
// @Failure     400 {object} Error
// @Router      /workout/log [get]
func (c *ExerciseLogController) GetUserLogsHandler(ctx *fiber.Ctx) error {
	userId := function.GetUserIDFromContext(ctx)
	logs, err := c.Service.GetLogsByUser(userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.JSON(logs)
}

// @Summary     Get exercise logs
// @Description Get logs for a specific exercise
// @Tags        exerciseLogs
// @Accept      json
// @Produce     json
// @Param       exerciseId path string true "Exercise ID"
// @Success     200 {array} ExerciseLog
// @Failure     400 {object} Error
// @Router      /workout/log/exercise/{exerciseId} [get]
func (c *ExerciseLogController) GetExerciseLogsHandler(ctx *fiber.Ctx) error {
	userId := function.GetUserIDFromContext(ctx)
	exerciseId := ctx.Params("exerciseId")

	logs, err := c.Service.GetLogsByExercise(exerciseId, userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.JSON(logs)
}

// @Summary     Get logs by date range
// @Description Get exercise logs within a date range
// @Tags        exerciseLogs
// @Accept      json
// @Produce     json
// @Param       startDate query string true "Start Date (YYYY-MM-DD)"
// @Param       endDate query string true "End Date (YYYY-MM-DD)"
// @Success     200 {array} ExerciseLog
// @Failure     400 {object} Error
// @Router      /workout/log/range [get]
func (c *ExerciseLogController) GetLogsByDateRangeHandler(ctx *fiber.Ctx) error {
	userId := function.GetUserIDFromContext(ctx)
	startDate, err := time.Parse("2006-01-02", ctx.Query("startDate"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid start date format",
		})
	}

	endDate, err := time.Parse("2006-01-02", ctx.Query("endDate"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid end date format",
		})
	}

	logs, err := c.Service.GetLogsByDateRange(userId, startDate, endDate)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.JSON(logs)
}

// @Summary     Update exercise log
// @Description Update an existing exercise log
// @Tags        exerciseLogs
// @Accept      json
// @Produce     json
// @Param       id path string true "Log ID"
// @Param       log body UpdateExerciseLogDto true "Updated Log"
// @Success     200 {object} ExerciseLog
// @Failure     400 {object} Error
// @Router      /workout/log/{id} [put]
func (c *ExerciseLogController) UpdateLogHandler(ctx *fiber.Ctx) error {
	userId := function.GetUserIDFromContext(ctx)
	logId := ctx.Params("id")
	validate := validator.New()
	dto := new(UpdateExerciseLogDto)

	if err := ctx.BodyParser(dto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(dto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	log, err := c.Service.UpdateLog(logId, dto, userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.JSON(log)
}

// @Summary     Delete exercise log
// @Description Delete an exercise log
// @Tags        exerciseLogs
// @Accept      json
// @Produce     json
// @Param       id path string true "Log ID"
// @Success     204 {object} nil
// @Failure     400 {object} Error
// @Router      /workout/log/{id} [delete]
func (c *ExerciseLogController) DeleteLogHandler(ctx *fiber.Ctx) error {
	userId := function.GetUserIDFromContext(ctx)
	logId := ctx.Params("id")

	if err := c.Service.DeleteLog(logId, userId); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (c *ExerciseLogController) Handle() {
	g := c.Instance.Group("/log")

	g.Post("/", c.CreateLogHandler)
	g.Get("/", c.GetUserLogsHandler)
	g.Get("/exercise/:exerciseId", c.GetExerciseLogsHandler)
	g.Get("/range", c.GetLogsByDateRangeHandler)
	g.Put("/:id", c.UpdateLogHandler)
	g.Delete("/:id", c.DeleteLogHandler)
}
