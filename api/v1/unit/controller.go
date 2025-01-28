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

// GetAllScaleUnits returns all available units
// @Summary Get all units
// @Description Get a list of all available units with their information
// @Tags units
// @Accept json
// @Produce json
// @Param unitType query string true "Unit type"
// @Success 200 {array} unitEnums.ScaleUnitInfo
// @Router /unit [get]
func (uc *UnitController) GetAllScaleUnits(ctx *fiber.Ctx) error {
	uType := unitType(ctx.Query("unitType"))
	units := uc.Service.GetAllUnits(uType)
	return ctx.JSON(units)
}

// GetScaleUnit returns information about a specific unit
// @Summary Get unit information
// @Description Get detailed information about a specific unit by its symbol
// @Tags units
// @Accept json
// @Produce json
// @Param symbol path string true "Unit symbol"
// @Param unitType query string true "Unit type"
// @Success 200 {object} unitEnums.ScaleUnitInfo
// @Failure 404 {object} Error
// @Router /unit/{symbol} [get]
func (uc *UnitController) GetScaleUnit(ctx *fiber.Ctx) error {
	symbol := ctx.Params("symbol")
	uType := unitType(ctx.Query("unitType"))
	unit, exists := uc.Service.GetUnit(symbol, uType)
	if !exists {
		return fiber.NewError(fiber.StatusNotFound, "Unit not found")
	}
	return ctx.JSON(unit)
}

// GetWeightUnits returns all available weight units
// @Summary Get all weight units
// @Description Get a list of all available weight units
// @Tags units
// @Accept json
// @Produce json
// @Success 200 {array} unitEnums.ExerciseWeightUnit
// @Router /unit/weight [get]
func (uc *UnitController) GetWeightUnits(ctx *fiber.Ctx) error {
	units := unitEnums.GetAllExerciseWeightUnit()
	return ctx.JSON(units)
}

// GetBodyPartMeasureUnits returns all available body part measure units
// @Summary Get all body part measure units
// @Description Get a list of all available body part measure units
// @Tags units
// @Accept json
// @Produce json
// @Success 200 {array} unitEnums.BodyPartMeasureUnit
// @Router /unit/bodypart [get]
func (uc *UnitController) GetBodyPartMeasureUnits(ctx *fiber.Ctx) error {
	units := unitEnums.GetAllBodyPartMeasureUnit()
	return ctx.JSON(units)
}

// ConvertUnits converts a value between units
// @Summary Convert between units
// @Description Convert a value from one unit to another
// @Tags units
// @Accept json
// @Produce json
// @Param unitType query string true "Unit type"
// @Param body body ConversionRequest true "Conversion request"
// @Success 200 {object} ConversionResponse
// @Failure 400 {object} Error
// @Router /unit/convert [post]
func (uc *UnitController) ConvertUnits(ctx *fiber.Ctx) error {
	var req ConversionRequest
	uType := unitType(ctx.Query("unitType"))
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := uc.Service.ConvertUnits(req.Value, req.FromUnit, req.ToUnit, uType)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.JSON(ConversionResponse{
		Value: result,
		Unit:  req.ToUnit,
	})
}

// Handle sets up all the unit routes
func (uc *UnitController) Handle() {
	g := uc.Instance.Group("/unit")

	g.Get("/", uc.GetAllScaleUnits)
	g.Get("/weight", uc.GetWeightUnits)
	g.Get("/bodypart", uc.GetBodyPartMeasureUnits)
	g.Get("/:symbol", uc.GetScaleUnit)
	g.Post("/convert", uc.ConvertUnits)
}
