package nutrition

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type NutritionController struct {
	Instance *fiber.App
	Service  INutritionService
}

type CreateNutritionDto struct {
	UserID    string    `json:"userid" validate:"required"`
	Carb      string    `json:"carb"`
	Protein   string    `json:"protein"`
	Fat       string    `json:"fat"`
	Calories  string    `json:"calories"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type UpdateNutritionDto struct {
	Carb     string `json:"carb"`
	Protein  string `json:"protein"`
	Fat      string `json:"fat"`
	Calories string `json:"calories"`
}

func (nc *NutritionController) PostNutritionHandler(c *fiber.Ctx) error {
	return nil
}

func (nc *NutritionController) GetNutritionsHandler(c *fiber.Ctx) error {
	return nil
}

func (nc *NutritionController) GetNutritionHandler(c *fiber.Ctx) error {
	return nil
}

func (nc *NutritionController) DeleteNutritionHandler(c *fiber.Ctx) error {
	return nil
}

func (nc *NutritionController) UpdateNutritionHandler(c *fiber.Ctx) error {
	return nil
}

func (nc *NutritionController) Handle() {
	g := nc.Instance.Group("/nutrition")
	g.Post("/", nc.PostNutritionHandler)
	g.Get("/", nc.GetNutritionsHandler)
	g.Get("/:id", nc.GetNutritionHandler)
	g.Delete("/:id", nc.DeleteNutritionHandler)
	g.Put("/:id", nc.UpdateNutritionHandler)
}
