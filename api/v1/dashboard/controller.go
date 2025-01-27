package dashboard

import (
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/function"
	"github.com/gofiber/fiber/v2"
)

type Error error

type DashboardController struct {
	Instance fiber.Router
	Service  IDashboardService
}

// @Summary     Get workout dashboard
// @Description Get workout frequency graph and analysis
// @Tags        dashboard
// @Accept      json
// @Produce     json
// @Param       startDate query string false "Start date"
// @Param       endDate query string false "End date"
// @Success     200 {object} DashboardResponse
// @Failure     400 {object} Error
// @Router      /dashboard [get]
func (c *DashboardController) GetDashboardHandler(ctx *fiber.Ctx) error {
	userId := function.GetUserIDFromContext(ctx)

	startDate, err := time.Parse("2006-01-02 15:04:05", ctx.Query("startDate"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid start date format. Expected format: YYYY-MM-DD HH:mm:ss",
		})
	}

	endDate, err := time.Parse("2006-01-02 15:04:05", ctx.Query("endDate"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid end date format. Expected format: YYYY-MM-DD HH:mm:ss",
		})
	}

	dashboard, err := c.Service.GetDashboard(userId, startDate, endDate)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.JSON(dashboard)
}

// @Summary     Get user strength standards
// @Description Get user strength standards
// @Tags        dashboard
// @Accept      json
// @Produce     json
// @Success     200 {object} UserStrengthStandards
// @Failure     400 {object} Error
// @Router      /dashboard/strength-standards [get]
func (c *DashboardController) GetUserStrengthStandardsHandler(ctx *fiber.Ctx) error {
	userId := function.GetUserIDFromContext(ctx)

	strengthStandards, err := c.Service.GetUserStrengthStandards(userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.JSON(strengthStandards)
}

// @Summary     Get rep max estimates
// @Description Get estimated rep maxes for a specific exercise
// @Tags        dashboard
// @Accept      json
// @Produce     json
// @Param       exerciseId path string true "Exercise ID"
// @Param       useLatest query boolean false "Use only latest exercise log"
// @Success     200 {object} RepMaxResponse
// @Failure     400 {object} Error
// @Router      /dashboard/rep-max/{exerciseId} [get]
func (c *DashboardController) GetRepMaxHandler(ctx *fiber.Ctx) error {
	userId := function.GetUserIDFromContext(ctx)
	exerciseId := ctx.Params("exerciseId")
	useLatest := ctx.QueryBool("useLatest", false)

	repMax, err := c.Service.GetRepMax(userId, exerciseId, useLatest)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.JSON(repMax)
}

// @Summary     Get nutrition summary
// @Description Get nutrition summary
// @Tags        dashboard
// @Accept      json
// @Produce     json
// @Param       startDate query string false "Start date"
// @Param       endDate query string false "End date"
// @Success     200 {object} NutritionSummaryResponse
// @Failure     400 {object} Error
// @Router      /dashboard/nutrition-summary [get]
func (dc *DashboardController) GetNutritionSummary(c *fiber.Ctx) error {
	userid := function.GetUserIDFromContext(c)
	if userid == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	startDate, err := time.Parse("2006-01-02 15:04:05", c.Query("startDate"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	endDate, err := time.Parse("2006-01-02 15:04:05", c.Query("endDate"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if startDate.After(endDate) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "start date cannot be after end date",
		})
	}

	summary, err := dc.Service.GetNutritionSummary(userid, startDate, endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(summary)
}

func (c *DashboardController) Handle() {
	g := c.Instance.Group("/dashboard")
	g.Get("/", c.GetDashboardHandler)
	g.Get("/nutrition-summary", c.GetNutritionSummary)
	g.Get("/strength-standards", c.GetUserStrengthStandardsHandler)
	g.Get("/rep-max/:exerciseId", c.GetRepMaxHandler)
}
