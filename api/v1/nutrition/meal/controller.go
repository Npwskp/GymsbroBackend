package meal

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type Error error

type MealController struct {
	Instance *fiber.App
	Service  IMealService
}

// @Summary		Create a meal
// @Description	Create a meal
// @Tags		meals
// @Accept		json
// @Produce		json
// @Param		meal body CreateMealDto true "Create Meal"
// @Success		201	{object} Meal
// @Failure		400	{object} Error
// @Router		/meals [post]
func (nc *MealController) CreateMealHandler(c *fiber.Ctx) error {
	meal := new(CreateMealDto)
	validate := validator.New()
	if err := c.BodyParser(meal); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	if err := validate.Struct(*meal); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	createdMeal, err := nc.Service.CreateMeal(meal)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(createdMeal)
}

// @Summary		Get all meals
// @Description	Get all meals
// @Tags		meals
// @Accept		json
// @Produce		json
// @Success		200	{object} []Meal
// @Failure		400	{object} Error
// @Router		/meals [get]
func (nc *MealController) GetMealsHandler(c *fiber.Ctx) error {
	meals, err := nc.Service.GetAllMeals()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(meals)
}

// @Summary		Get a meal
// @Description	Get a meal
// @Tags		meals
// @Accept		json
// @Produce		json
// @Param		id path string true "Meal ID"
// @Success		200	{object} Meal
// @Failure		400	{object} Error
// @Router		/meals/{id} [get]
func (nc *MealController) GetMealHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	meal, err := nc.Service.GetMeal(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(meal)
}

// @Summary		Get meals by user
// @Description	Get meals by user
// @Tags		meals
// @Accept		json
// @Produce		json
// @Param		userid path string true "User ID"
// @Success		200	{object} []Meal
// @Failure		400	{object} Error
// @Router		/meals/user/{userid} [get]
func (nc *MealController) GetMealByUserHandler(c *fiber.Ctx) error {
	userid := c.Params("userid")
	meals, err := nc.Service.GetMealByUser(userid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(meals)
}

// @Summary		Get meals by user and date
// @Description	Get meals by user and date
// @Tags		meals
// @Accept		json
// @Produce		json
// @Param		userid path string true "User ID"
// @Param		start query int true "Start date"
// @Param		end query int true "End date"
// @Success		200	{object} []Meal
// @Failure		400	{object} Error
// @Router		/meals/userdate/{userid} [get]
func (nc *MealController) GetMealByUserDateHandler(c *fiber.Ctx) error {
	userid := c.Params("userid")
	start := c.QueryInt("start")
	end := c.QueryInt("end")
	meals, err := nc.Service.GetMealByUserDate(userid, start, end)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(meals)
}

// @Summary		Delete a meal
// @Description	Delete a meal
// @Tags		meals
// @Accept		json
// @Produce		json
// @Param		id path string true "Meal ID"
// @Success		204
// @Failure		400	{object} Error
// @Router		/meals/{id} [delete]
func (nc *MealController) DeleteMealHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	err := nc.Service.DeleteMeal(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "Meal deleted"})
}

// @Summary		Update a meal
// @Description	Update a meal
// @Tags		meals
// @Accept		json
// @Produce		json
// @Param		id path string true "Meal ID"
// @Param		meal body UpdateMealDto true "Update Meal"
// @Success		200	{object} Meal
// @Failure		400	{object} Error
// @Router		/meals/{id} [put]
func (nc *MealController) UpdateMealHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	validate := validator.New()
	doc := new(UpdateMealDto)
	if err := c.BodyParser(doc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	if err := validate.Struct(*doc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	meal, err := nc.Service.UpdateMeal(doc, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(meal)
}

func (nc *MealController) Handle() {
	g := nc.Instance.Group("/meals")
	g.Post("/", nc.CreateMealHandler)
	g.Get("/", nc.GetMealsHandler)
	g.Get("/:id", nc.GetMealHandler)
	g.Get("/user/:userid", nc.GetMealByUserHandler)
	g.Get("/userdate/:userid", nc.GetMealByUserDateHandler)
	g.Delete("/:id", nc.DeleteMealHandler)
	g.Put("/:id", nc.UpdateMealHandler)
}
