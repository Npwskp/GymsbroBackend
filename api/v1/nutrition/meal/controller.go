package meal

import (
	"strings"

	"github.com/Npwskp/GymsbroBackend/api/v1/function"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
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

// @Summary		Calculate nutrient
// @Description	Calculate nutrient
// @Tags		meals
// @Accept		json
// @Produce		json
// @Param		body body CalculateNutrientBody true "Calculate Nutrient"
// @Success		200	{object} CalculateNutrientResponse
// @Failure		400	{object} Error
// @Router		/meal/calculate [post]
func (nc *MealController) CalculateNutrientHandler(c *fiber.Ctx) error {
	body := new(CalculateNutrientBody)
	userid := function.GetUserIDFromContext(c)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	nutrients, err := nc.Service.CalculateNutrient(body, userid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(nutrients)
}

// @Summary		Get a meal
// @Description	Get a meal
// @Tags		meals
// @Accept		json
// @Produce		json
// @Param		id path string true "Meal ID"
// @Success		200	{object} Meal
// @Failure		404	{object} Error
// @Failure		400	{object} Error
// @Failure		500	{object} Error
// @Router		/meal/{id} [get]
func (nc *MealController) GetMealHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	userid := function.GetUserIDFromContext(c)
	meal, err := nc.Service.GetMeal(id, userid)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Meal not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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

	return c.Status(fiber.StatusOK).JSON(meals)
}

// @Summary		Delete a meal
// @Description	Delete a meal
// @Tags		meals
// @Accept		json
// @Produce		json
// @Param		id path string true "Meal ID"
// @Success		204
// @Failure		404	{object} Error
// @Failure		400	{object} Error
// @Failure		500	{object} Error
// @Router		/meal/{id} [delete]
func (nc *MealController) DeleteMealHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	userid := function.GetUserIDFromContext(c)
	err := nc.Service.DeleteMeal(id, userid)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Meal not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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
// @Failure		404	{object} Error
// @Failure		400	{object} Error
// @Failure		500	{object} Error
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
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Meal not found",
			})
		}
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
// @Failure 500 {object} Error
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

	return c.JSON(meals)
}

// @Summary Update meal image
// @Description Update a meal's image
// @Tags meals
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Meal ID"
// @Param file formData file true "Image file"
// @Success 200 {object} Meal
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Router /meal/{id}/image [patch]
func (mc *MealController) UpdateMealImageHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	userId := function.GetUserIDFromContext(c)

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	// Check file type
	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File must be an image",
		})
	}

	// Open the file
	fileContent, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process file",
		})
	}
	defer fileContent.Close()

	meal, err := mc.Service.UpdateMealImage(c, id, fileContent, file.Filename, contentType, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(meal)
}

func (nc *MealController) Handle() {
	g := nc.Instance.Group("/meal")

	g.Post("/", nc.CreateMealHandler)
	g.Post("/calculate", nc.CalculateNutrientHandler)
	g.Get("/search", nc.SearchFilteredMealsHandler)
	g.Get("/user", nc.GetMealByUserHandler)
	g.Get("/:id", nc.GetMealHandler)
	g.Delete("/:id", nc.DeleteMealHandler)
	g.Put("/:id", nc.UpdateMealHandler)
	g.Patch("/:id/image", nc.UpdateMealImageHandler)
}
