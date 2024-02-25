package handlers

import (
	"chulbong-kr/dto"
	"chulbong-kr/models"
	"chulbong-kr/services"
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
