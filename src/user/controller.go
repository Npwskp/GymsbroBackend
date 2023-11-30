package user

import (
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	Instance *fiber.App
	Service  IUserService
}

type Error error

type CreateUserDto struct {
	Username      string    `json:"username" validate:"required,min=3,max=20"`
	Email         string    `json:"email" validate:"required"`
	Password      string    `json:"password" validate:"required"`
	Weight        float64   `json:"weight" default:"0"` // default:"0" is not working
	Height        float64   `json:"height" default:"0"` // default:"0" is not working
	Age           int       `json:"age" validate:"required,min=1,max=120"`
	Gender        string    `json:"gender" validate:"required"`
	Neck          float64   `json:"neck" default:"0"`          // default:"0" is not working
	Waist         float64   `json:"waist" default:"0"`         // default:"0" is not working
	Hip           float64   `json:"hip" default:"0"`           // default:"0" is not working
	ActivityLevel int       `json:"activityLevel" default:"0"` // default:"0" is not working
	CreatedAt     time.Time `json:"created_at,omitempty" bson:"created_at,omitempty" set:"omitempty"`
}

type UpadateUsernamePasswordDto struct {
	Username    string `json:"username"`
	Password    string `json:"password" validate:"required"`
	NewPassword string `json:"newPassword"`
}

type UpdateBodyDto struct {
	Weight        float64 `json:"weight"`
	Height        float64 `json:"height"`
	Age           int     `json:"age"`
	Gender        string  `json:"gender"`
	Neck          float64 `json:"neck"`
	Waist         float64 `json:"waist"`
	Hip           float64 `json:"hip"`
	ActivityLevel int     `json:"activityLevel"`
}

// @Summary		Create a user
// @Description	Create a user
// @Tags		users
// @Accept		json
// @Produce		json
// @Param		user body CreateUserDto true "Create User"
// @Success		201	{object} User
// @Failure		400	{object} Error
// @Router		/users [post]
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
// @Router		/users [get]
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
// @Param		id path	string true "User ID"
// @Success		200	{object} User
// @Failure		400	{object} Error
// @Router		/users/{id} [get]
func (uc *UserController) GetUserHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := uc.Service.GetUser(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

// @Summary		Delete a user
// @Description	Delete a user
// @Tags		users
// @Accept		json
// @Produce		json
// @Param		id path	string true "User ID"
// @Success		204
// @Failure		400	{object} Error
// @Router		/users/{id} [delete]
func (uc *UserController) DeleteUserHandler(c *fiber.Ctx) error {
	id := c.Params("id")
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
// @Param		id path string true "User ID"
// @Param		user body UpadateUsernamePasswordDto true "UpdateUsernamePassword User"
// @Success		200	{object} User
// @Failure		400	{object} Error
// @Router		/users/{id}/usepass [patch]
func (uc *UserController) UpdateUsernamePassword(c *fiber.Ctx) error {
	id := c.Params("id")
	validate := validator.New()
	doc := new(UpadateUsernamePasswordDto)
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
// @Param		id path string true "User ID"
// @Param		user body UpdateBodyDto true "UpdateBody User"
// @Success		200	{object} User
// @Failure		400	{object} Error
// @Router		/users/{id}/body [patch]
func (uc *UserController) UpdateBody(c *fiber.Ctx) error {
	id := c.Params("id")
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

func (uc *UserController) Handle() {
	g := uc.Instance.Group("/users")
	g.Post("", uc.PostUsersHandler)
	g.Get("", uc.GetAllUsersHandler)
	g.Get("/:id", uc.GetUserHandler)
	g.Delete("/:id", uc.DeleteUserHandler)
	g.Patch("/:id/usepass", uc.UpdateUsernamePassword)
	g.Patch("/:id/body", uc.UpdateBody)
}
