package nutrition

import (
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type Error error

type NutritionController struct {
	Instance *fiber.App
	Service  INutritionService
}

type CreateNutritionDto struct {
	UserID    string    `json:"userid" validate:"required"`
	Carb      float64   `json:"carb" default:"0"`
	Protein   float64   `json:"protein" default:"0"`
	Fat       float64   `json:"fat" default:"0"`
	Calories  float64   `json:"calories" default:"0"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type UpdateNutritionDto struct {
	Carb     float64 `json:"carb"`
	Protein  float64 `json:"protein"`
	Fat      float64 `json:"fat"`
	Calories float64 `json:"calories"`
}

// @Summary		Create a nutrition
// @Description	Create a nutrition
// @Tags		nutritions
// @Accept		json
// @Produce		json
// @Param		nutrition body CreateNutritionDto true "Create Nutrition"
// @Success		201	{object} Nutrition
// @Failure		400	{object} Error
// @Router		/nutritions [post]
func (nc *NutritionController) PostNutritionHandler(c *fiber.Ctx) error {
	nutrition := new(CreateNutritionDto)
	validate := validator.New()
	if err := c.BodyParser(nutrition); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	if err := validate.Struct(*nutrition); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	createdNutrition, err := nc.Service.CreateNutrition(nutrition)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(createdNutrition)
}

// @Summary		Get all nutritions
// @Description	Get all nutritions
// @Tags		nutritions
// @Accept		json
// @Produce		json
// @Success		200	{object} []Nutrition
// @Failure		400	{object} Error
// @Router		/nutritions [get]
func (nc *NutritionController) GetNutritionsHandler(c *fiber.Ctx) error {
	nutritions, err := nc.Service.GetAllNutritions()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(nutritions)
}

// @Summary		Get a nutrition
// @Description	Get a nutrition
// @Tags		nutritions
// @Accept		json
// @Produce		json
// @Param		id path string true "Nutrition ID"
// @Success		200	{object} Nutrition
// @Failure		400	{object} Error
// @Router		/nutritions/{id} [get]
func (nc *NutritionController) GetNutritionHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	nutrition, err := nc.Service.GetNutrition(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(nutrition)
}

// @Summary		Get nutritions by user
// @Description	Get nutritions by user
// @Tags		nutritions
// @Accept		json
// @Produce		json
// @Param		userid path string true "User ID"
// @Success		200	{object} []Nutrition
// @Failure		400	{object} Error
// @Router		/nutritions/user/{userid} [get]
func (nc *NutritionController) GetNutritionByUserHandler(c *fiber.Ctx) error {
	userid := c.Params("userid")
	nutritions, err := nc.Service.GetNutritionByUser(userid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(nutritions)
}

// @Summary		Delete a nutrition
// @Description	Delete a nutrition
// @Tags		nutritions
// @Accept		json
// @Produce		json
// @Param		id path string true "Nutrition ID"
// @Success		204
// @Failure		400	{object} Error
// @Router		/nutritions/{id} [delete]
func (nc *NutritionController) DeleteNutritionHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	err := nc.Service.DeleteNutrition(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "Nutrition deleted"})
}

// @Summary		Update a nutrition
// @Description	Update a nutrition
// @Tags		nutritions
// @Accept		json
// @Produce		json
// @Param		id path string true "Nutrition ID"
// @Param		nutrition body UpdateNutritionDto true "Update Nutrition"
// @Success		200	{object} Nutrition
// @Failure		400	{object} Error
// @Router		/nutritions/{id} [put]
func (nc *NutritionController) UpdateNutritionHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	validate := validator.New()
	doc := new(UpdateNutritionDto)
	if err := c.BodyParser(doc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	if err := validate.Struct(*doc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	nutrition, err := nc.Service.UpdateNutrition(doc, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(nutrition)
}

func (nc *NutritionController) Handle() {
	g := nc.Instance.Group("/nutritions")
	g.Post("/", nc.PostNutritionHandler)
	g.Get("/", nc.GetNutritionsHandler)
	g.Get("/:id", nc.GetNutritionHandler)
	g.Get("/user/:userid", nc.GetNutritionByUserHandler)
	g.Delete("/:id", nc.DeleteNutritionHandler)
	g.Put("/:id", nc.UpdateNutritionHandler)
}
