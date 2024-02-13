package handlers

import (
	"chulbong-kr/models"

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

func CreateUser(user *models.User) error {
	// // Prepare SQL statement
	// stmt, err := database.DB.Prepare("INSERT INTO users(email, hashed_password, created_at, last_login_at, nickname, markers) VALUES(?,?,?,?,?,?)")
	// if err != nil {
	// 	return err
	// }
	// defer stmt.Close()

	// // Execute SQL statement
	// _, err = stmt.Exec(user.Email, user.PasswordHash, time.Now(), time.Now(), user.Username, user.Markers)
	// if err != nil {
	// 	return err
	// }

	// fmt.Println("User added successfully!")
	return nil
}
