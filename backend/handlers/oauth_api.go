package handlers

import (
	"chulbong-kr/services"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

// GetGoogleAuthHandler generates a handler to redirect to Google OAuth2
func GetGoogleAuthHandler(conf *oauth2.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate a state string for CSRF protection
		state := services.GenerateState()
		c.Cookie(&fiber.Cookie{
			Name:  "oauthstate",
			Value: state,
			// Secure, HttpOnly, etc?
		})

		URL := conf.AuthCodeURL(state)
		return c.Redirect(URL)
	}
}

// GetGoogleCallbackHandler generates a handler for the OAuth2 callback from Google
func GetGoogleCallbackHandler(conf *oauth2.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Validate state
		state := c.Cookies("oauthstate")
		if state != c.Query("state") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "State mismatch"})
		}

		code := c.Query("code")
		token, err := conf.Exchange(c.Context(), code)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to exchange token"})
		}

		profile, err := services.ConvertGoogleToken(token.AccessToken)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve user profile"})
		}

		return c.JSON(profile)
	}
}
