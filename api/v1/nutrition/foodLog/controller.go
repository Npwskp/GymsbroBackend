package foodlog

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type Error error

type FoodLogController struct {
	Instance *fiber.App
	Service  IFoodLogService
}

type CreateFoodLogDto struct {
	UserID    string    `json:"userid" validate:"required"`
	Date      string    `json:"date" validate:"required"`
	Meals     []string  `json:"meals" default:"null"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty" default:"null"`
}

type UpdateFoodLogDto struct {
	UserID string   `json:"userid" validate:"required"`
	Date   string   `json:"date" validate:"required"`
	Meals  []string `json:"meals" default:"null"`
}

func (fc *FoodLogController) CreateFoodLog(c *fiber.Ctx) error {
	dto := new(CreateFoodLogDto)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	foodlog, err := fc.Service.CreateFoodLog(dto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(foodlog)
}

func (fc *FoodLogController) GetAllFoodLogs(c *fiber.Ctx) error {
	foodlogs, err := fc.Service.GetAllFoodLogs()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(foodlogs)
}

func (fc *FoodLogController) GetFoodLog(c *fiber.Ctx) error {
	id := c.Params("id")
	foodlog, err := fc.Service.GetFoodLog(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(foodlog)
}

func (fc *FoodLogController) GetFoodLogByUser(c *fiber.Ctx) error {
	userid := c.Params("userid")
	foodlogs, err := fc.Service.GetFoodLogByUser(userid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(foodlogs)
}

func (fc *FoodLogController) GetFoodLogByUserDate(c *fiber.Ctx) error {
	userid := c.Params("userid")
	date := c.Params("date")
	foodlog, err := fc.Service.GetFoodLogByUserDate(userid, date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(foodlog)
}

func (fc *FoodLogController) DeleteFoodLog(c *fiber.Ctx) error {
	id := c.Params("id")
	err := fc.Service.DeleteFoodLog(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (fc *FoodLogController) UpdateFoodLog(c *fiber.Ctx) error {
	id := c.Params("id")
	dto := new(UpdateFoodLogDto)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	foodlog, err := fc.Service.UpdateFoodLog(dto, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(foodlog)
}

func (fc *FoodLogController) Handle() {
	g := fc.Instance.Group("/foodlog")
	g.Post("/", fc.CreateFoodLog)
	g.Get("/", fc.GetAllFoodLogs)
	g.Get("/:id", fc.GetFoodLog)
	g.Get("/user/:userid", fc.GetFoodLogByUser)
	g.Get("/userdate/:userid", fc.GetFoodLogByUserDate)
	g.Delete("/:id", fc.DeleteFoodLog)
	g.Put("/:id", fc.UpdateFoodLog)
}
