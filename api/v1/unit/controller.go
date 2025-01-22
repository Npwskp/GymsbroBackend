package unit

import (
	unitEnums "github.com/Npwskp/GymsbroBackend/api/v1/unit/enums"
	"github.com/gofiber/fiber/v2"
)

type UnitController struct {
	Instance fiber.Router
	Service  IUnitService
}

type Error error

// GetAllUnits returns all available units
// @Summary Get all units
// @Description Get a list of all available units with their information
// @Tags Units
// @Accept json
// @Produce json
// @Success 200 {array} unitEnums.UnitInfo
// @Router /unit [get]
func (c *UnitController) GetAllUnits(ctx *fiber.Ctx) error {
	units := c.Service.GetAllUnits()
	return ctx.JSON(units)
}

// GetUnit returns information about a specific unit
// @Summary Get unit information
// @Description Get detailed information about a specific unit by its symbol
// @Tags Units
// @Accept json
// @Produce json
// @Param symbol path string true "Unit symbol"
// @Success 200 {object} unitEnums.UnitInfo
// @Failure 404 {object} Error
// @Router /unit/{symbol} [get]
func (c *UnitController) GetUnit(ctx *fiber.Ctx) error {
	symbol := ctx.Params("symbol")
	unit, exists := c.Service.GetUnit(symbol)
	if !exists {
		return fiber.NewError(fiber.StatusNotFound, "Unit not found")
	}
	return ctx.JSON(unit)
}

// GetWeightUnits returns all available weight units
// @Summary Get all weight units
// @Description Get a list of all available weight units
// @Tags Units
// @Accept json
// @Produce json
// @Success 200 {array} unitEnums.ExerciseWeightUnit
// @Router /unit/weight [get]
func (c *UnitController) GetWeightUnits(ctx *fiber.Ctx) error {
	units := unitEnums.GetAllExerciseWeightUnit()
	return ctx.JSON(units)
}

// ConvertUnits converts a value between units
// @Summary Convert between units
// @Description Convert a value from one unit to another
// @Tags Units
// @Accept json
// @Produce json
// @Param body body ConversionRequest true "Conversion request"
// @Success 200 {object} ConversionResponse
// @Failure 400 {object} Error
// @Router /unit/convert [post]
func (c *UnitController) ConvertUnits(ctx *fiber.Ctx) error {
	var req ConversionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := c.Service.ConvertBetweenUnits(req.Value, req.FromUnit, req.ToUnit)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.JSON(ConversionResponse{
		Value: result,
		Unit:  req.ToUnit,
	})
}

// Handle sets up all the unit routes
func (c *UnitController) Handle() {
	g := c.Instance.Group("/unit")

	g.Get("/", c.GetAllUnits)
	g.Get("/weight", c.GetWeightUnits)
	g.Get("/:symbol", c.GetUnit)
	g.Post("/convert", c.ConvertUnits)
}
