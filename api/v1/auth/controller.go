package auth

import (
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/config"
	"github.com/Npwskp/GymsbroBackend/api/v1/middleware"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
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
	// Check if there's an existing JWT cookie
	if existingJWT := c.Cookies("jwt"); existingJWT != "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "another session is already active, please logout first",
		})
	}

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
			"error": err.Error(),
		})
	}
	cookie := &fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Unix(exp, 0),
		HTTPOnly: true,
		Secure:   config.CookieSecure,
		SameSite: config.CookieSameSite,
	}
	c.Cookie(cookie)
	return c.Status(fiber.StatusOK).JSON(ReturnToken{
		Token: token,
		Exp:   exp,
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
	if err := c.BodyParser(register); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	if err := validate.Struct(register); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	user, err := ac.Service.Register(register)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
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
	// Get JWT from cookie
	jwtCookie := c.Cookies("jwt")
	if jwtCookie == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "no active session",
		})
	}

	// Parse the token to get user claims
	token, err := jwt.Parse(jwtCookie, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetJWTSecret()), nil
	})
	if err != nil || !token.Valid {
		// Clear invalid cookie anyway
		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    "",
			Expires:  time.Now().Add(-time.Hour),
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Strict",
		})
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	// Clear cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}

func (ac *AuthController) Handle() {
	g := ac.Instance.Group("/auth")
	g.Post("/login", middleware.CheckNotLoggedIn(), ac.PostLoginHandler)
	g.Post("/register", middleware.CheckNotLoggedIn(), ac.PostRegisterHandler)
	g.Get("/me", middleware.AuthMiddleware(), ac.GetMeHandler)
	g.Post("/logout", middleware.AuthMiddleware(), ac.PostLogoutHandler)

	// Add Google OAuth routes
	g.Get("/google/login", ac.GoogleLogin)
	g.Get("/google/callback", ac.GoogleCallback)
}
