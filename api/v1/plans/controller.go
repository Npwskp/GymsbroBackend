package plans

import (
	"github.com/Npwskp/GymsbroBackend/api/v1/function"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type Error error

type PlanController struct {
	Instance fiber.Router
	Service  IPlanService
}

// @Summary		Create a plan
// @Description	Create a plan
// @Tags		plans
// @Accept		json
// @Produce		json
// @Param		plan body CreatePlanDto true "Create Plan"
// @Success		201	{object} Plan
// @Failure		400	{object} Error
// @Router		/plans [post]
func (pc *PlanController) PostPlansHandler(c *fiber.Ctx) error {
	validate := validator.New()
	plan := new(CreatePlanDto)
	if err := c.BodyParser(plan); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	if err := validate.Struct(*plan); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	if !function.CheckDay(plan.DayOfWeek) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Day of week is not valid"})
	}
	if !function.CheckExerciseType(plan.TypeOfPlan) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Type of plan is not valid"})
	}
	createdPlan, err := pc.Service.CreatePlan(plan)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(createdPlan)
}

// @Summary		Get all plans
// @Description	Get all plans
// @Tags		plans
// @Accept		json
// @Produce		json
// @Success		200	{array}	Plan
// @Failure		400	{object} Error
// @Router		/plans [get]
func (pc *PlanController) GetPlansHandler(c *fiber.Ctx) error {
	plans, err := pc.Service.GetAllPlans()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(plans)
}

// @Summary		Get a plan
// @Description	Get a plan
// @Tags		plans
// @Accept		json
// @Produce		json
// @Param		id path	string true "Plan ID"
// @Success		200	{object} Plan
// @Failure		400	{object} Error
// @Router		/plans/{id} [get]
func (pc *PlanController) GetPlanHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	plan, err := pc.Service.GetPlan(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(plan)
}

// @Summary		Get a plan by user
// @Description	Get a plan by user
// @Tags		plans
// @Accept		json
// @Produce		json
// @Param		user_id path	string true "User ID"
// @Success		200	{array} Plan
// @Failure		400	{object} Error
// @Router		/plans/user/{user_id} [get]
func (pc *PlanController) GetPlanByUserHandler(c *fiber.Ctx) error {
	user_id := c.Params("user_id")
	plans, err := pc.Service.GetAllPlanByUser(user_id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(plans)
}

// @Summary		Get a plan by user and day
// @Description	Get a plan by user and day
// @Tags		plans
// @Accept		json
// @Produce		json
// @Param		user_id path	string true "User ID"
// @Param		day path	string true "Day of week"
// @Success		200	{object} Plan
// @Failure		400	{object} Error
// @Router		/plans/user/{user_id}/{day} [get]
func (pc *PlanController) GetPlanByUserDayHandler(c *fiber.Ctx) error {
	user_id := c.Params("user_id")
	day := c.Params("day")
	plan, err := pc.Service.GetPlanByUserDay(user_id, day)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(plan)
}

// @Summary		Delete a plan
// @Description	Delete a plan
// @Tags		plans
// @Accept		json
// @Produce		json
// @Param		id path	string true "Plan ID"
// @Success		204
// @Failure		400	{object} Error
// @Router		/plans/{id} [delete]
func (pc *PlanController) DeletePlanHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := pc.Service.DeletePlan(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary		Delete a plan by user and day
// @Description	Delete a plan by user and day
// @Tags		plans
// @Accept		json
// @Produce		json
// @Param		user_id path	string true "User ID"
// @Param		day path	string true "Day of week"
// @Success		204
// @Failure		400	{object} Error
// @Router		/plans/user/{user_id}/{day} [delete]
func (pc *PlanController) DeletePlanByUserDayHandler(c *fiber.Ctx) error {
	user_id := c.Params("user_id")
	day := c.Params("day")
	if err := pc.Service.DeleteByUserDay(user_id, day); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary		Update a plan
// @Description	Update a plan
// @Tags		plans
// @Accept		json
// @Produce		json
// @Param		id path	string true "Plan ID"
// @Param		plan body UpdatePlanDto true "Update Plan"
// @Success		200	{object} Plan
// @Failure		400	{object} Error
// @Router		/plans/{id} [put]
func (pc *PlanController) UpdatePlanHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	validate := validator.New()
	doc := new(UpdatePlanDto)
	if err := c.BodyParser(&doc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	if err := validate.Struct(*doc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	if !function.CheckExerciseType(doc.TypeOfPlan) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Type of plan is not valid"})
	}
	plan, err := pc.Service.UpdatePlan(doc, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(plan)
}

// @Summary		Update a plan by user and day
// @Description	Update a plan by user and day
// @Tags		plans
// @Accept		json
// @Produce		json
// @Param		user_id path	string true "User ID"
// @Param		day path	string true "Day of week"
// @Param		plan body UpdatePlanDto true "Update Plan"
// @Success		200	{object} Plan
// @Failure		400	{object} Error
// @Router		/plans/{user_id}/{day} [put]
func (pc *PlanController) UpdatePlanByUserDayHandler(c *fiber.Ctx) error {
	user_id := c.Params("user_id")
	day := c.Params("day")
	validate := validator.New()
	doc := new(UpdatePlanDto)
	if err := c.BodyParser(&doc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	if err := validate.Struct(*doc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	if !function.CheckDay(day) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Day of week is not valid"})
	}
	if !function.CheckExerciseType(doc.TypeOfPlan) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Type of plan is not valid"})
	}
	plan, err := pc.Service.UpdatePlanByUserDay(doc, user_id, day)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(plan)
}

func (pc *PlanController) Handle() {
	g := pc.Instance.Group("/plans")
	g.Post("/", pc.PostPlansHandler)
	g.Get("/", pc.GetPlansHandler)
	g.Get("/:id", pc.GetPlanHandler)
	g.Get("/user/:user_id", pc.GetPlanByUserHandler)
	g.Get("/user/:user_id/:day", pc.GetPlanByUserDayHandler)
	g.Delete("/:id", pc.DeletePlanHandler)
	g.Delete("/:user_id/:day", pc.DeletePlanByUserDayHandler)
	g.Put("/:id", pc.UpdatePlanHandler)
	g.Put("/:user_id/:day", pc.UpdatePlanByUserDayHandler)
}
