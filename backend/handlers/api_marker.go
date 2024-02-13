package handlers

import (
	"chulbong-kr/models"
	"chulbong-kr/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetExample handler
func GetExample(c *fiber.Ctx) error {
	return c.SendString("GET request example")
}

// PutExample handler
func PutExample(c *fiber.Ctx) error {
	return c.SendString("PUT request example")
}

// DynamicRouteExample handler
func DynamicRouteExample(c *fiber.Ctx) error {
	// Capture string and id from the path
	stringParam := c.Params("string")
	idParam := c.Params("id")

	// Optionally, convert idParam to integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID must be a number",
		})
	}

	return c.JSON(fiber.Map{
		"string": stringParam,
		"id":     id,
	})
}

// QueryParamsExample handler
func QueryParamsExample(c *fiber.Ctx) error {
	// Capture query parameters
	query1 := c.Query("query")
	query2 := c.Query("query2")

	// You can also provide a default value if a query parameter is missing
	query3 := c.Query("query3", "default value")

	return c.JSON(fiber.Map{
		"query1": query1,
		"query2": query2,
		"query3": query3,
	})
}

// CreateMarker handler
func CreateMarker(c *fiber.Ctx) error {
	var marker models.Marker
	if err := c.BodyParser(&marker); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	err := services.CreateMarker(&marker)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(marker)
}

// GetMarker handler
func GetMarker(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	marker, err := services.GetMarker(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Marker not found"})
	}
	return c.JSON(marker)
}

// UpdateMarker updates an existing marker
func UpdateMarker(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	marker, _ := services.GetMarker(id)

	if err := c.BodyParser(marker); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := services.UpdateMarker(marker); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(marker)
}
