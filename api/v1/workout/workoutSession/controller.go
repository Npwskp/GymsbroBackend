package workoutSession

import (
	"github.com/Npwskp/GymsbroBackend/api/v1/function"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type Error error

type WorkoutSessionController struct {
	Instance fiber.Router
	Service  IWorkoutSessionService
}

// @Summary     Start workout session
// @Description Start a new workout session
// @Tags        workoutSessions
// @Accept      json
// @Produce     json
// @Param       session body CreateWorkoutSessionDto true "Create Session"
// @Success     201 {object} WorkoutSession
// @Failure     400 {object} Error
// @Router      /workout-session [post]
func (c *WorkoutSessionController) StartSessionHandler(ctx *fiber.Ctx) error {
	userId := function.GetUserIDFromContext(ctx)
	validate := validator.New()
	dto := new(CreateWorkoutSessionDto)

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

	session, err := c.Service.StartSession(dto, userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(session)
}

// @Summary     Log custom session
// @Description Create a new workout session with custom start and end times
// @Tags        workoutSessions
// @Accept      json
// @Produce     json
// @Param       session body LoggedSessionDto true "Logged Session"
// @Success     201 {object} WorkoutSession
// @Failure     400 {object} Error
// @Router      /workout-session/log [post]
func (c *WorkoutSessionController) LogSessionHandler(ctx *fiber.Ctx) error {
	userId := function.GetUserIDFromContext(ctx)
	validate := validator.New()
	dto := new(LoggedSessionDto)

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

	session, err := c.Service.LogSession(dto, userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(session)
}

// @Summary     End workout session
// @Description End an active workout session
// @Tags        workoutSessions
// @Accept      json
// @Produce     json
// @Param       id path string true "Session ID"
// @Success     200 {object} WorkoutSession
// @Failure     400 {object} Error
// @Router      /workout-session/{id}/end [put]
func (c *WorkoutSessionController) EndSessionHandler(ctx *fiber.Ctx) error {
	userId := function.GetUserIDFromContext(ctx)
	sessionId := ctx.Params("id")

	session, err := c.Service.EndSession(sessionId, userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.JSON(session)
}

// @Summary     Update workout session
// @Description Update a workout session
// @Tags        workoutSessions
// @Accept      json
// @Produce     json
// @Param       id path string true "Session ID"
// @Param       session body UpdateWorkoutSessionDto true "Update Session"
// @Success     200 {object} WorkoutSession
// @Failure     400 {object} Error
// @Router      /workout-session/{id} [put]
func (c *WorkoutSessionController) UpdateSessionHandler(ctx *fiber.Ctx) error {
	userId := function.GetUserIDFromContext(ctx)
	sessionId := ctx.Params("id")
	validate := validator.New()
	dto := new(UpdateWorkoutSessionDto)

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

	session, err := c.Service.UpdateSession(sessionId, dto, userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.JSON(session)
}

// @Summary     Get session
// @Description Get a workout session by ID
// @Tags        workoutSessions
// @Accept      json
// @Produce     json
// @Param       id path string true "Session ID"
// @Success     200 {object} WorkoutSession
// @Failure     400 {object} Error
// @Router      /workout-session/{id} [get]
func (c *WorkoutSessionController) GetSessionHandler(ctx *fiber.Ctx) error {
	userId := function.GetUserIDFromContext(ctx)
	sessionId := ctx.Params("id")

	session, err := c.Service.GetSession(sessionId, userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.JSON(session)
}

// @Summary     Get user sessions
// @Description Get all workout sessions for a user
// @Tags        workoutSessions
// @Accept      json
// @Produce     json
// @Success     200 {array} WorkoutSession
// @Failure     400 {object} Error
// @Router      /workout-session [get]
func (c *WorkoutSessionController) GetUserSessionsHandler(ctx *fiber.Ctx) error {
	userId := function.GetUserIDFromContext(ctx)

	sessions, err := c.Service.GetUserSessions(userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.JSON(sessions)
}

// @Summary     Delete session
// @Description Delete a workout session
// @Tags        workoutSessions
// @Accept      json
// @Produce     json
// @Param       id path string true "Session ID"
// @Success     204 {object} nil
// @Failure     400 {object} Error
// @Router      /workout-session/{id} [delete]
func (c *WorkoutSessionController) DeleteSessionHandler(ctx *fiber.Ctx) error {
	userId := function.GetUserIDFromContext(ctx)
	sessionId := ctx.Params("id")

	if err := c.Service.DeleteSession(sessionId, userId); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (c *WorkoutSessionController) Handle() {
	g := c.Instance.Group("/workout-session")

	g.Post("/", c.StartSessionHandler)
	g.Post("/log", c.LogSessionHandler)
	g.Get("/", c.GetUserSessionsHandler)
	g.Get("/:id", c.GetSessionHandler)
	g.Put("/:id", c.UpdateSessionHandler)
	g.Put("/:id/end", c.EndSessionHandler)
	g.Delete("/:id", c.DeleteSessionHandler)
}
