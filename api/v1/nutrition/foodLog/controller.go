package foodlog

import (
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/function"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type Error error

type FoodLogController struct {
	Instance fiber.Router
	Service  IFoodLogService
}

// @Summary		Add meal to food log
// @Description	Add meal to food log
// @Tags		foodlog
// @Accept		json
// @Produce		json
// @Param		foodlog body AddMealToFoodLogDto true "Food log object that needs to be created"
// @Success		201	{object} FoodLog
// @Failure		400	{object} Error
// @Failure		500	{object} Error
// @Router		/foodlog [post]
func (fc *FoodLogController) AddMealToFoodLog(c *fiber.Ctx) error {
	dto := new(AddMealToFoodLogDto)
	userid := function.GetUserIDFromContext(c)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	foodlog, err := fc.Service.AddMealToFoodLog(dto, userid)
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
// @Failure		404	{object} Error
// @Failure		400	{object} Error
// @Failure		500	{object} Error
// @Router		/foodlog/{id} [get]
func (fc *FoodLogController) GetFoodLog(c *fiber.Ctx) error {
	id := c.Params("id")
	userid := function.GetUserIDFromContext(c)
	foodlog, err := fc.Service.GetFoodLog(id, userid)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Food log not found",
			})
		}
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
// @Failure		500	{object} Error
// @Router		/foodlog/user [get]
func (fc *FoodLogController) GetFoodLogByUser(c *fiber.Ctx) error {
	userid := function.GetUserIDFromContext(c)
	foodlogs, err := fc.Service.GetFoodLogByUser(userid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if foodlogs == nil {
		foodlogs = []*FoodLog{}
	}
	return c.JSON(foodlogs)
}

// @Summary		Get a food log by user and date
// @Description	Get a food log by user and date
// @Tags		foodlog
// @Accept		json
// @Produce		json
// @Param		date path	string true "Date"
// @Success		200	{object} FoodLog
// @Failure		404	{object} Error
// @Failure		400	{object} Error
// @Router		/foodlog/user/{date} [get]
func (fc *FoodLogController) GetFoodLogByUserDate(c *fiber.Ctx) error {
	userid := function.GetUserIDFromContext(c)
	date := c.Params("date")
	foodlog, err := fc.Service.GetFoodLogByUserDate(userid, date)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Food log not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(foodlog)
}

// @Summary     Calculate daily nutrients
// @Description Calculate total nutrients and calories for a specific date
// @Tags        foodlog
// @Accept      json
// @Produce     json
// @Param       date path string true "Date (YYYY-MM-DD)"
// @Success     200  {object} DailyNutrientResponse
// @Failure     400  {object} Error
// @Failure     500  {object} Error
// @Router      /foodlog/nutrients/{date} [get]
func (fc *FoodLogController) CalculateDailyNutrients(c *fiber.Ctx) error {
	date := c.Params("date")
	userid := function.GetUserIDFromContext(c)

	// Validate date format (YYYY-MM-DD)
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid date format. Please use YYYY-MM-DD",
		})
	}

	nutrients, err := fc.Service.CalculateDailyNutrients(date, userid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(nutrients)
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
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Food log not found",
			})
		}
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
// @Param		foodlog body UpdateFoodLogDto true "Food log object that needs to be updated"
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
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Food log not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(foodlog)
}

func (fc *FoodLogController) Handle() {
	g := fc.Instance.Group("/foodlog")

	g.Post("/", fc.AddMealToFoodLog)
	g.Get("/user", fc.GetFoodLogByUser)
	g.Get("/:id", fc.GetFoodLog)
	g.Get("/user/:date", fc.GetFoodLogByUserDate)
	g.Get("/nutrients/:date", fc.CalculateDailyNutrients)
	g.Delete("/:id", fc.DeleteFoodLog)
	g.Put("/:id", fc.UpdateFoodLog)
}
