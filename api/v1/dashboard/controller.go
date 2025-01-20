package dashboard

import (
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
// @Success     200 {object} DashboardResponse
// @Failure     400 {object} Error
// @Router      /dashboard [get]
func (c *DashboardController) GetDashboardHandler(ctx *fiber.Ctx) error {
	userId := function.GetUserIDFromContext(ctx)

	dashboard, err := c.Service.GetDashboard(userId)
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
// @Success     200 {object} RepMaxResponse
// @Failure     400 {object} Error
// @Router      /dashboard/rep-max/{exerciseId} [get]
func (c *DashboardController) GetRepMaxHandler(ctx *fiber.Ctx) error {
	userId := function.GetUserIDFromContext(ctx)
	exerciseId := ctx.Params("exerciseId")

	repMax, err := c.Service.GetRepMax(userId, exerciseId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.JSON(repMax)
}

func (c *DashboardController) Handle() {
	g := c.Instance.Group("/dashboard")
	g.Get("/", c.GetDashboardHandler)
	g.Get("/strength-standards", c.GetUserStrengthStandardsHandler)
	g.Get("/rep-max/:exerciseId", c.GetRepMaxHandler)
}
