package ingredient

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type Error error

type IngredientController struct {
	Instance fiber.Router
	Service  IIngredientService
}

type CreateIngredientDto struct {
	UserID    string    `json:"userid" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Image     string    `json:"image" default:"null"`
	Calories  float64   `json:"calories" default:"0"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty" default:"null"`
}

type UpdateIngredientDto struct {
	UserID   string  `json:"userid" validate:"required"`
	Name     string  `json:"name"`
	Image    string  `json:"image"`
	Calories float64 `json:"calories"`
}

func (ic *IngredientController) CreateIngredient(c *fiber.Ctx) error {
	dto := new(CreateIngredientDto)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	ingredient, err := ic.Service.CreateIngredient(dto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(ingredient)
}

func (ic *IngredientController) GetAllIngredients(c *fiber.Ctx) error {
	ingredients, err := ic.Service.GetAllIngredients()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(ingredients)
}

func (ic *IngredientController) GetIngredient(c *fiber.Ctx) error {
	id := c.Params("id")
	ingredient, err := ic.Service.GetIngredient(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(ingredient)
}

func (ic *IngredientController) GetIngredientByUser(c *fiber.Ctx) error {
	userid := c.Params("userid")
	ingredients, err := ic.Service.GetIngredientByUser(userid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(ingredients)
}

func (ic *IngredientController) DeleteIngredient(c *fiber.Ctx) error {
	id := c.Params("id")
	err := ic.Service.DeleteIngredient(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (ic *IngredientController) UpdateIngredient(c *fiber.Ctx) error {
	id := c.Params("id")
	dto := new(UpdateIngredientDto)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	ingredient, err := ic.Service.UpdateIngredient(dto, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(ingredient)
}

func (ic *IngredientController) Handle() {
	g := ic.Instance.Group("/ingredient")
	g.Post("/", ic.CreateIngredient)
	g.Get("/", ic.GetAllIngredients)
	g.Get("/:id", ic.GetIngredient)
	g.Get("/user/:userid", ic.GetIngredientByUser)
	g.Delete("/:id", ic.DeleteIngredient)
	g.Put("/:id", ic.UpdateIngredient)
}
