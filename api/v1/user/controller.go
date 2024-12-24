package user

import (
	"github.com/Npwskp/GymsbroBackend/api/v1/function"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	Instance fiber.Router
	Service  IUserService
}

type Error error

// @Summary		Create a user
// @Description	Create a user
// @Tags		users
// @Accept		json
// @Produce		json
// @Param		user body CreateUserDto true "Create User"
// @Success		201	{object} User
// @Failure		400	{object} Error
// @Router		/user [post]
func (uc *UserController) PostUsersHandler(c *fiber.Ctx) error {
	validate := validator.New()
	user := new(CreateUserDto)
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(*user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	res, err := uc.Service.CreateUser(user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

// @Summary		Get all users
// @Description	Get all users
// @Tags		users
// @Accept		json
// @Produce		json
// @Success		200	{array}	User
// @Failure		400	{object} Error
// @Router		/user [get]
func (uc *UserController) GetAllUsersHandler(c *fiber.Ctx) error {
	users, err := uc.Service.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(users)
}

// @Summary		Get a user
// @Description	Get a user
// @Tags		users
// @Accept		json
// @Produce		json
// @Success		200	{object} User
// @Failure		400	{object} Error
// @Router		/user/me [get]
func (uc *UserController) GetUserHandler(c *fiber.Ctx) error {
	id := function.GetUserIDFromContext(c)
	user, err := uc.Service.GetUser(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

// @Summary		Get a user energy consume plan
// @Description	Get a user energy consume plan
// @Tags		users
// @Accept		json
// @Produce		json
// @Success		200	{object} function.EnergyConsumptionPlan
// @Failure		400	{object} Error
// @Router		/user/energyplan [get]
func (uc *UserController) GetUserEnergyConsumePlanHandler(c *fiber.Ctx) error {
	id := function.GetUserIDFromContext(c)
	plan, err := uc.Service.GetUserEnergyConsumePlan(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(plan)
}

// @Summary		Delete a user
// @Description	Delete a user
// @Tags		users
// @Accept		json
// @Produce		json
// @Success		204
// @Failure		400	{object} Error
// @Router		/user/me [delete]
func (uc *UserController) DeleteUserHandler(c *fiber.Ctx) error {
	id := function.GetUserIDFromContext(c)
	err := uc.Service.DeleteUser(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary		Update a user username and password
// @Description	Update a user username and password
// @Tags		users
// @Accept		json
// @Produce		json
// @Param		user body UpdateUsernamePasswordDto true "UpdateUsernamePassword User"
// @Success		200	{object} User
// @Failure		400	{object} Error
// @Router		/user/usepass [patch]
func (uc *UserController) UpdateUsernamePassword(c *fiber.Ctx) error {
	id := function.GetUserIDFromContext(c)
	validate := validator.New()
	doc := new(UpdateUsernamePasswordDto)
	if err := c.BodyParser(&doc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(*doc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user, err := uc.Service.UpdateUsernamePassword(doc, id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

// @Summary		Update a user body
// @Description	Update a user body
// @Tags		users
// @Accept		json
// @Produce		json
// @Param		user body UpdateBodyDto true "UpdateBody User"
// @Success		200	{object} User
// @Failure		400	{object} Error
// @Router		/user/body [patch]
func (uc *UserController) UpdateBody(c *fiber.Ctx) error {
	id := function.GetUserIDFromContext(c)
	validate := validator.New()
	doc := new(UpdateBodyDto)
	if err := c.BodyParser(&doc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(*doc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user, err := uc.Service.UpdateBody(doc, id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

// @Summary     Update first login status
// @Description Mark user as not first time login
// @Tags        users
// @Accept      json
// @Produce     json
// @Success     204
// @Failure     400 {object} Error
// @Router      /user/first-login [put]
func (uc *UserController) UpdateFirstLoginStatus(c *fiber.Ctx) error {
	id := function.GetUserIDFromContext(c)
	err := uc.Service.UpdateFirstLoginStatus(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (uc *UserController) Handle() {
	g := uc.Instance.Group("/user")

	g.Post("", uc.PostUsersHandler)
	g.Get("", uc.GetAllUsersHandler)
	g.Get("/me", uc.GetUserHandler)
	g.Get("/energyplan", uc.GetUserEnergyConsumePlanHandler)
	g.Delete("/me", uc.DeleteUserHandler)
	g.Patch("/body", uc.UpdateBody)
	g.Patch("/usepass", uc.UpdateUsernamePassword)
	g.Patch("/first-login", uc.UpdateFirstLoginStatus)
}
