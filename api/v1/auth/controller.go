package auth

import (
	"fmt"
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/config"
	"github.com/Npwskp/GymsbroBackend/api/v1/middleware"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	Instance fiber.Router
	Service  IAuthService
}

type Error error

// @Summary		Login
// @Description	Login
// @Tags		auth
// @Accept		json
// @Produce		json
// @Param		user body LoginDto true "Login"
// @Success		200	{object} ReturnToken
// @Failure		400	{object} Error
// @Router		/auth/login [post]
func (ac *AuthController) PostLoginHandler(c *fiber.Ctx) error {
	validate := validator.New()
	login := new(LoginDto)
	if err := c.BodyParser(login); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	if err := validate.Struct(login); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	token, exp, err := ac.Service.Login(login)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}
	cookie := &fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Unix(exp, 0),
		Secure:   config.CookieSecure,
		HTTPOnly: config.CookieHTTPOnly,
		SameSite: config.CookieSameSite,
		Path:     "/",
	}
	c.Cookie(cookie)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"exp":     exp,
	})
}

// @Summary		Register
// @Description	Register
// @Tags		auth
// @Accept		json
// @Produce		json
// @Param		user body RegisterDto true "Register"
// @Success		201	{object} user.User
// @Failure		400	{object} Error
// @Router		/auth/register [post]
func (ac *AuthController) PostRegisterHandler(c *fiber.Ctx) error {
	validate := validator.New()
	register := new(RegisterDto)
	fmt.Println("Hello")
	if err := c.BodyParser(register); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	fmt.Println(register)
	if err := validate.Struct(register); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	user, err := ac.Service.Register(register)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	fmt.Println(user, err)
	return c.Status(fiber.StatusCreated).JSON(user)
}

// @Summary		Get me
// @Description	Get me
// @Tags		auth
// @Accept		json
// @Produce		json
// @Success		200	{object} user.User
// @Failure		400	{object} Error
// @Router		/auth/me [get]
func (ac *AuthController) GetMeHandler(c *fiber.Ctx) error {
	userClaims, err := middleware.GetCurrentUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	user, err := ac.Service.Me(userClaims.UserID.Hex())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// @Summary		Logout
// @Description	Logout
// @Tags		auth
// @Accept		json
// @Produce		json
// @Success		200	{object} string
// @Failure		400	{object} Error
// @Router		/auth/logout [post]
func (ac *AuthController) PostLogoutHandler(c *fiber.Ctx) error {
	cookie := new(fiber.Cookie)
	cookie.Name = "jwt"
	cookie.Value = ""
	c.Cookie(cookie)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "logout success"})
}

func (ac *AuthController) Handle() {
	g := ac.Instance.Group("/auth")
	g.Post("/login", ac.PostLoginHandler)
	g.Post("/register", ac.PostRegisterHandler)
	g.Get("/me", ac.GetMeHandler)
	g.Post("/logout", ac.PostLogoutHandler)
}
