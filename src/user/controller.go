package user

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	Instance *fiber.App
	Service  IUserService
}

type CreateUserDto struct {
	Username      string  `json:"username" validate:"required,min=3,max=20"`
	Password      string  `json:"password" validate:"required"`
	Weight        float64 `json:"weight" default:"0"` // default:"0" is not working
	Height        float64 `json:"height" default:"0"` // default:"0" is not working
	Age           int     `json:"age" validate:"required,min=1,max=120"`
	Gender        string  `json:"gender" validate:"required"`
	Neck          float64 `json:"neck" default:"0"`          // default:"0" is not working
	Waist         float64 `json:"waist" default:"0"`         // default:"0" is not working
	Hip           float64 `json:"hip" default:"0"`           // default:"0" is not working
	ActivityLevel int     `json:"activityLevel" default:"0"` // default:"0" is not working
}

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

func (uc *UserController) GetAllUsersHandler(c *fiber.Ctx) error {
	users, err := uc.Service.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(users)
}

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

func (uc *UserController) Handle() {
	g := uc.Instance.Group("/users")
	g.Post("", uc.PostUsersHandler)
	g.Get("", uc.GetAllUsersHandler)
	g.Get("/:id", uc.GetUserHandler)
}
