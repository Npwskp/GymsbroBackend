package plans

import (
	"github.com/Npwskp/GymsbroBackend/src/function"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type PlanController struct {
	Instance *fiber.App
	Service  IPlanService
}

type CreatePlanDto struct {
	UserID     string   `json:"userId" validate:"required"`
	TypeOfPlan string   `json:"typeOfPlan" validate:"required"`
	DayOfWeek  string   `json:"dayOfWeek" validate:"required"`
	Exercise   []string `json:"exercise" default:"[]"`
}

type UpdatePlanDto struct {
	UserID     string   `json:"userId"`
	TypeOfPlan string   `json:"typeOfPlan"`
	Exercise   []string `json:"exercise"`
}

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

func (pc *PlanController) GetPlansHandler(c *fiber.Ctx) error {
	plans, err := pc.Service.GetAllPlans()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(plans)
}

func (pc *PlanController) GetPlanHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	plan, err := pc.Service.GetPlan(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(plan)
}

func (pc *PlanController) GetPlanByUserHandler(c *fiber.Ctx) error {
	user_id := c.Params("user_id")
	plans, err := pc.Service.GetAllPlanByUser(user_id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(plans)
}

func (pc *PlanController) GetPlanByUserDayHandler(c *fiber.Ctx) error {
	user_id := c.Params("user_id")
	day := c.Params("day")
	plan, err := pc.Service.GetPlanByUserDay(user_id, day)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(plan)
}

func (pc *PlanController) DeletePlanHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := pc.Service.DeletePlan(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (pc *PlanController) DeletePlanByUserDayHandler(c *fiber.Ctx) error {
	user_id := c.Params("user_id")
	day := c.Params("day")
	if err := pc.Service.DeleteByUserDay(user_id, day); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

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
	plan, err := pc.Service.UpdatePlan(doc, id)
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
}
