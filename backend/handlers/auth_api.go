package handlers

import (
	"chulbong-kr/dto"
	"chulbong-kr/models"
	"chulbong-kr/services"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// DeleteExample handler
func DeleteExample(c *fiber.Ctx) error {
	return c.SendString("DELETE request example")
}

// PostExample handler
func PostExample(c *fiber.Ctx) error {
	user := new(models.User)

	// Parse the body into the struct
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	// Return the user as JSON
	return c.Status(fiber.StatusOK).JSON(user)
}

func SignUpHandler(c *fiber.Ctx) error {
	var signUpReq dto.SignUpRequest
	if err := c.BodyParser(&signUpReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON, wrong sign up form."})
	}

	// Check if the token is verified before proceeding
	verified, err := services.IsTokenVerified(signUpReq.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check verification status"})
	}
	if !verified {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email not verified"})
	}

	signUpReq.Provider = "website"

	user, err := services.SaveUser(&signUpReq)
	if err != nil {
		// Handle the duplicate email error
		if strings.Contains(err.Error(), "already registered") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		// For other errors, return a generic error message
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "An error occurred while creating the user: " + err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func LoginHandler(c *fiber.Ctx) error {
	var request dto.LoginRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := services.Login(request.Email, request.Password)
	if err != nil {
		log.Printf("Error logging in: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	token, err := services.GenerateAndSaveToken(user.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	// Create a response object that includes both the user and the token
	var response dto.LoginResponse
	response.User = user
	response.Token = token

	// Setting the token in a secure cookie
	cookie := services.GenerateLoginCookie(token)
	c.Cookie(&cookie)

	return c.JSON(response)
}

func SendVerificationEmailHandler(c *fiber.Ctx) error {
	userEmail := c.FormValue("email")
	_, err := services.GetUserByEmail(userEmail)
	if err == nil {
		// If GetUserByEmail does not return an error, it means the email is already in use
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already registered"})
	} else if err != sql.ErrNoRows {
		// if db couldn't find a user, then it's valid. other errors are bad.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "An unexpected error occurred"})
	}

	// No matter if it's verified, send again.
	// Check if there's already a verified token for this user
	// verified, err := services.IsTokenVerified(userEmail)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check verification status"})
	// }
	// if verified {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email already verified"})
	// }

	// token, err := services.GenerateAndSaveSignUpToken(userEmail)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	// }

	// Use a goroutine to send the email without blocking
	go func(email string) {
		token, err := services.GenerateAndSaveSignUpToken(email)
		if err != nil {
			fmt.Printf("Failed to generate token: %v\n", err)
			return
		}

		err = services.SendVerificationEmail(email, token)
		if err != nil {
			fmt.Printf("Failed to send verification email: %v\n", err)
			return
		}
	}(userEmail)

	return c.SendStatus(fiber.StatusOK)
}

func ValidateTokenHandler(c *fiber.Ctx) error {
	token := c.FormValue("token")
	email := c.FormValue("email")

	valid, err := services.ValidateToken(token, email)
	if err != nil {
		// If err is not nil, it could be a database error or token not found
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error validating token"})
	}
	if !valid {
		// Handle both not found and expired cases
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	return c.SendStatus(fiber.StatusOK)
}
