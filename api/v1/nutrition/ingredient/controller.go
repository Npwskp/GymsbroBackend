package ingredient

import (
	"github.com/Npwskp/GymsbroBackend/api/v1/function"
	"github.com/gofiber/fiber/v2"
)

type Error error

type IngredientController struct {
	Instance fiber.Router
	Service  IIngredientService
}

// @Summary Create new ingredient
// @Description Create new ingredient
// @Tags ingredient
// @Accept json
// @Produce json
// @Param ingredient body CreateIngredientDto true "Ingredient object that needs to be created"
// @Success 201 {object} Ingredient
// @Failure 400 {object} Error
// @Router /ingredient [post]
func (ic *IngredientController) CreateIngredient(c *fiber.Ctx) error {
	dto := new(CreateIngredientDto)
	userId := function.GetUserIDFromContext(c)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	ingredient, err := ic.Service.CreateIngredient(dto, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(ingredient)
}

// @Summary Get all ingredients
// @Description Get all ingredients
// @Tags ingredient
// @Accept json
// @Produce json
// @Success 200 {object} []Ingredient
// @Failure 400 {object} Error
// @Router /ingredient [get]
func (ic *IngredientController) GetAllIngredients(c *fiber.Ctx) error {
	userId := function.GetUserIDFromContext(c)
	ingredients, err := ic.Service.GetAllIngredients(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(ingredients)
}

// @Summary Get an ingredient
// @Description Get an ingredient
// @Tags ingredient
// @Accept json
// @Produce json
// @Param id path string true "Ingredient ID"
// @Success 200 {object} Ingredient
// @Failure 400 {object} Error
// @Router /ingredient/{id} [get]
func (ic *IngredientController) GetIngredient(c *fiber.Ctx) error {
	id := c.Params("id")
	userId := function.GetUserIDFromContext(c)
	ingredient, err := ic.Service.GetIngredient(id, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(ingredient)
}

// @Summary Get ingredients by user
// @Description Get ingredients by user
// @Tags ingredient
// @Accept json
// @Produce json
// @Param userid path string true "User ID"
// @Success 200 {object} []Ingredient
// @Failure 400 {object} Error
// @Router /ingredient/user [get]
func (ic *IngredientController) GetIngredientByUser(c *fiber.Ctx) error {
	userId := function.GetUserIDFromContext(c)
	ingredients, err := ic.Service.GetIngredientByUser(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(ingredients)
}

// @Summary Delete an ingredient
// @Description Delete an ingredient
// @Tags ingredient
// @Accept json
// @Produce json
// @Param id path string true "Ingredient ID"
// @Success 204 "No Content"
// @Failure 400 {object} Error
// @Router /ingredient/{id} [delete]
func (ic *IngredientController) DeleteIngredient(c *fiber.Ctx) error {
	id := c.Params("id")
	userId := function.GetUserIDFromContext(c)

	err := ic.Service.DeleteIngredient(id, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary Update an ingredient
// @Description Update an ingredient
// @Tags ingredient
// @Accept json
// @Produce json
// @Param id path string true "Ingredient ID"
// @Param ingredient body UpdateIngredientDto true "Ingredient object that needs to be updated"
// @Success 200 {object} Ingredient
// @Failure 400 {object} Error
// @Router /ingredient/{id} [put]
func (ic *IngredientController) UpdateIngredient(c *fiber.Ctx) error {
	id := c.Params("id")
	user := c.Locals("user").(map[string]interface{})
	userId := user["id"].(string)

	var doc UpdateIngredientDto
	if err := c.BodyParser(&doc); err != nil {
		return err
	}

	ingredient, err := ic.Service.UpdateIngredient(&doc, id, userId)
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
