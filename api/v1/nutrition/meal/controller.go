package meal

import (
	"github.com/Npwskp/GymsbroBackend/api/v1/function"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type Error error

type MealController struct {
	Instance fiber.Router
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
// @Router		/meal [post]
func (nc *MealController) CreateMealHandler(c *fiber.Ctx) error {
	meal := new(CreateMealDto)
	userid := function.GetUserIDFromContext(c)
	validate := validator.New()
	if err := c.BodyParser(meal); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	if err := validate.Struct(*meal); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	createdMeal, err := nc.Service.CreateMeal(meal, userid)
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
// @Router		/meal [get]
func (nc *MealController) GetMealsHandler(c *fiber.Ctx) error {
	userid := function.GetUserIDFromContext(c)
	meals, err := nc.Service.GetAllMeals(userid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	// Initialize empty slice if meals is nil
	if meals == nil {
		meals = []*Meal{}
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
// @Router		/meal/{id} [get]
func (nc *MealController) GetMealHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	userid := function.GetUserIDFromContext(c)
	meal, err := nc.Service.GetMeal(id, userid)
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
// @Success		200	{object} []Meal
// @Failure		400	{object} Error
// @Router		/meal/user [get]
func (nc *MealController) GetMealByUserHandler(c *fiber.Ctx) error {
	userid := function.GetUserIDFromContext(c)
	meals, err := nc.Service.GetMealByUser(userid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Initialize empty slice if meals is nil
	if meals == nil {
		meals = []*Meal{}
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
// @Router		/meal/{id} [delete]
func (nc *MealController) DeleteMealHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	userid := function.GetUserIDFromContext(c)
	err := nc.Service.DeleteMeal(id, userid)
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
// @Router		/meal/{id} [put]
func (nc *MealController) UpdateMealHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	userid := function.GetUserIDFromContext(c)
	validate := validator.New()
	doc := new(UpdateMealDto)
	if err := c.BodyParser(doc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	if err := validate.Struct(*doc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	meal, err := nc.Service.UpdateMeal(doc, id, userid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(meal)
}

// @Summary Search and filter meals
// @Description Search meals with optional filters
// @Tags meals
// @Accept json
// @Produce json
// @Param q query string false "Search query"
// @Param category query string false "Category filter"
// @Param minCalories query number false "Minimum calories"
// @Param maxCalories query number false "Maximum calories"
// @Param nutrients query string false "Nutrients filter (comma-separated)"
// @Success 200 {array} Meal
// @Failure 400 {object} Error
// @Router /meal/search [get]
func (mc *MealController) SearchFilteredMealsHandler(c *fiber.Ctx) error {
	filters := SearchFilters{
		Query:       c.Query("q"),
		Category:    c.Query("category"),
		MinCalories: c.QueryInt("minCalories", 0),
		MaxCalories: c.QueryInt("maxCalories", 0),
		Nutrients:   c.Query("nutrients"),
		UserID:      function.GetUserIDFromContext(c),
	}

	meals, err := mc.Service.SearchFilteredMeals(filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Initialize empty slice if meals is nil
	if meals == nil {
		meals = []*Meal{}
	}

	return c.JSON(meals)
}

func (nc *MealController) Handle() {
	g := nc.Instance.Group("/meal")

	g.Post("/", nc.CreateMealHandler)
	g.Get("/", nc.GetMealsHandler)
	g.Get("/search", nc.SearchFilteredMealsHandler)
	g.Get("/user", nc.GetMealByUserHandler)
	g.Get("/:id", nc.GetMealHandler)
	g.Delete("/:id", nc.DeleteMealHandler)
	g.Put("/:id", nc.UpdateMealHandler)
}
