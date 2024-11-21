package foodlog

import (
	"github.com/Npwskp/GymsbroBackend/api/v1/function"
	"github.com/gofiber/fiber/v2"
)

type Error error

type FoodLogController struct {
	Instance fiber.Router
	Service  IFoodLogService
}

// @Summary		Create a food log
// @Description	Create a food log
// @Tags		foodlog
// @Accept		json
// @Produce		json
// @Param		foodlog body CreateFoodLogDto true "Food log object that needs to be created"
// @Success		201	{object} FoodLog
// @Failure		400	{object} Error
// @Router		/foodlog [post]
func (fc *FoodLogController) CreateFoodLog(c *fiber.Ctx) error {
	dto := new(CreateFoodLogDto)
	userid := function.GetUserIDFromContext(c)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	foodlog, err := fc.Service.CreateFoodLog(dto, userid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(foodlog)
}

// @Summary		Get a food log
// @Description	Get a food log
// @Tags		foodlog
// @Accept		json
// @Produce		json
// @Param		id path	string true "Food log ID"
// @Success		200	{object} FoodLog
// @Failure		400	{object} Error
// @Router		/foodlog/{id} [get]
func (fc *FoodLogController) GetFoodLog(c *fiber.Ctx) error {
	id := c.Params("id")
	userid := function.GetUserIDFromContext(c)
	foodlog, err := fc.Service.GetFoodLog(id, userid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(foodlog)
}

// @Summary		Get a food log by user
// @Description	Get a food log by user
// @Tags		foodlog
// @Accept		json
// @Produce		json
// @Success		200	{object} []FoodLog
// @Failure		400	{object} Error
// @Router		/foodlog/user [get]
func (fc *FoodLogController) GetFoodLogByUser(c *fiber.Ctx) error {
	userid := function.GetUserIDFromContext(c)
	foodlogs, err := fc.Service.GetFoodLogByUser(userid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(foodlogs)
}

// @Summary		Get a food log by user and date
// @Description	Get a food log by user and date
// @Tags		foodlog
// @Accept		json
// @Produce		json
// @Param		date path	string true "Date"
// @Success		200	{object} []FoodLog
// @Failure		400	{object} Error
// @Router		/foodlog/date/{date} [get]
func (fc *FoodLogController) GetFoodLogByUserDate(c *fiber.Ctx) error {
	userid := function.GetUserIDFromContext(c)
	date := c.Params("date")
	foodlog, err := fc.Service.GetFoodLogByUserDate(userid, date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(foodlog)
}

// @Summary		Delete a food log
// @Description	Delete a food log
// @Tags		foodlog
// @Accept		json
// @Produce		json
// @Param		id path	string true "Food log ID"
// @Success		204
// @Failure		400	{object} Error
// @Router		/foodlog/{id} [delete]
func (fc *FoodLogController) DeleteFoodLog(c *fiber.Ctx) error {
	id := c.Params("id")
	userid := function.GetUserIDFromContext(c)
	err := fc.Service.DeleteFoodLog(id, userid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary		Update a food log
// @Description	Update a food log
// @Tags		foodlog
// @Accept		json
// @Produce		json
// @Param		id path	string true "Food log ID"
// @Success		200	{object} FoodLog
// @Failure		400	{object} Error
// @Router		/foodlog/{id} [put]
func (fc *FoodLogController) UpdateFoodLog(c *fiber.Ctx) error {
	id := c.Params("id")
	userid := function.GetUserIDFromContext(c)
	dto := new(UpdateFoodLogDto)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	foodlog, err := fc.Service.UpdateFoodLog(dto, id, userid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(foodlog)
}

func (fc *FoodLogController) Handle() {
	g := fc.Instance.Group("/foodlog")

	g.Post("/", fc.CreateFoodLog)
	g.Get("/:id", fc.GetFoodLog)
	g.Get("/user", fc.GetFoodLogByUser)
	g.Get("/date/:date", fc.GetFoodLogByUserDate)
	g.Delete("/:id", fc.DeleteFoodLog)
	g.Put("/:id", fc.UpdateFoodLog)
}
