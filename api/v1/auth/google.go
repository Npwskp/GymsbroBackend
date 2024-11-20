package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/config"
	"github.com/Npwskp/GymsbroBackend/api/v1/user"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig *oauth2.Config

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	Picture       string `json:"picture"`
}

func InitGoogleOAuth() {
	googleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func (ac *AuthController) GoogleLogin(c *fiber.Ctx) error {
	url := googleOauthConfig.AuthCodeURL("state")
	return c.Redirect(url, fiber.StatusTemporaryRedirect)
}

func (ac *AuthController) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(c.Context(), code)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Failed to exchange token",
		})
	}

	client := googleOauthConfig.Client(c.Context(), token)
	userInfo, err := getUserInfo(client)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user info",
		})
	}

	// Check if user exists
	userService := user.UserService{DB: ac.Service.(*AuthService).DB}
	existingUser, err := userService.GetUserByEmail(userInfo.Email)

	if err != nil {
		// Create new user if doesn't exist
		register := &RegisterDto{
			Username:      userInfo.Name,
			Email:         userInfo.Email,
			Password:      generateRandomPassword(32),
			OAuthProvider: "google",
			OAuthID:       userInfo.ID,
			Picture:       userInfo.Picture,
			// Optional fields can be set later
			Age:    0,
			Gender: "",
		}

		existingUser, err = ac.Service.Register(register)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create user",
			})
		}
	}

	// Generate JWT token
	jwtToken, exp, err := createJWTToken(existingUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Set cookie
	cookie := &fiber.Cookie{
		Name:     "jwt",
		Value:    jwtToken,
		Expires:  time.Unix(exp, 0),
		HTTPOnly: true,
		Secure:   config.CookieSecure,
		SameSite: config.CookieSameSite,
	}
	c.Cookie(cookie)

	// Redirect to frontend with token
	frontendURL := os.Getenv("FRONTEND_URL")
	return c.Redirect(fmt.Sprintf("%s/auth/callback?token=%s", frontendURL, jwtToken))
}

func getUserInfo(client *http.Client) (*GoogleUser, error) {
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var user GoogleUser
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#$%^&*"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}