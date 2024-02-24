package handlers

import (
	"chulbong-kr/dto"
	"chulbong-kr/services"
	"fmt"
	"os"

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

		// Prepare the SignUpRequest for the OAuth2 user
		signUpReq := dto.SignUpRequest{
			Email:      profile.Email,
			Username:   &profile.Name, // Assuming profile.Name is the user's name from Google
			Provider:   "google",
			ProviderID: profile.SUB, // The unique ID from Google for the user
		}

		// Call the SaveUser function to create or find the OAuth2 user
		user, err := services.SaveUser(&signUpReq)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create or find OAuth user: %v", err)})
		}

		loginToken, err := services.GenerateAndSaveToken(user.UserID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate login token"})
		}

		// Respond with the user information or a success message
		clientAddr := fmt.Sprintf("%s/%s=%s", os.Getenv("CLIENT_ADRR"), os.Getenv("CLIENT_REDIRECT_ENDPOINT"), loginToken)

		// Setting the token in a secure cookie
		cookie := services.GenerateCookie(loginToken)
		c.Cookie(&cookie)
		return c.Redirect(clientAddr)
	}
}
