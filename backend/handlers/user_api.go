package handlers

import (
	"chulbong-kr/dto"
	"chulbong-kr/middlewares"
	"chulbong-kr/models"
	"chulbong-kr/services"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// RegisterUserRoutes sets up the routes for user handling within the application.
func RegisterUserRoutes(api fiber.Router) {
	userGroup := api.Group("/users")
	{
		userGroup.Use(middlewares.AuthMiddleware)
		userGroup.Get("/me", profileHandler)
		userGroup.Get("/favorites", getFavoritesHandler)
		userGroup.Get("/reports", getMyReportsHandler)
		userGroup.Patch("/me", updateUserHandler)
		userGroup.Delete("/me", deleteUserHandler)
		userGroup.Delete("/s3/objects", middlewares.AdminOnly, deleteObjectFromS3Handler)
	}
}

// UpdateUserHandler
func updateUserHandler(c *fiber.Ctx) error {
	userData, err := services.GetUserFromContext(c)
	if err != nil {
		return err // fiber err
	}

	var updateReq dto.UpdateUserRequest
	if err := c.BodyParser(&updateReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	user, err := services.UpdateUserProfile(userData.UserID, &updateReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	userProfileKey := fmt.Sprintf("%s:%d", services.USER_PROFILE_KEY, userData.UserID)
	services.ResetCache(userProfileKey)

	return c.JSON(user)
}

// DeleteUserHandler deletes the currently authenticated user
func deleteUserHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID not found"})
	}

	log.Printf("[DEBUG][HANDLER] Deleting user %v", userID)

	if err := services.DeleteUserWithRelatedData(c.Context(), userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent) // 204 for successful deletion with no content in response
}

func profileHandler(c *fiber.Ctx) error {
	userData, err := services.GetUserFromContext(c)
	if err != nil {
		return err // fiber err
	}

	userProfileKey := fmt.Sprintf("%s:%d", services.USER_PROFILE_KEY, userData.UserID)

	// Try to get the user profile from the cache first
	cachedUser, cacheErr := services.GetCacheEntry[*models.User](userProfileKey)
	if cacheErr == nil && cachedUser != nil {
		// Cache hit, return the cached user
		c.Append("X-Cache", "hit")
		return c.JSON(cachedUser)
	}

	// If the cache doesn't have the user profile, fetch it from the database
	user, err := services.GetUserById(userData.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// After fetching from the database
	services.SetCacheEntry(userProfileKey, user, 15*time.Minute)

	return c.JSON(user)
}

func getFavoritesHandler(c *fiber.Ctx) error {
	userData, err := services.GetUserFromContext(c)
	if err != nil {
		return err // fiber err
	}

	userFavKey := fmt.Sprintf("%s:%d:%s", services.USER_FAV_KEY, userData.UserID, userData.Username)

	// Try to get the user fav from the cache first
	cachedFav, cacheErr := services.GetCacheEntry[[]dto.MarkerSimpleWithDescrption](userFavKey)
	if cacheErr == nil && cachedFav != nil {
		// Cache hit, return the cached fav
		return c.JSON(cachedFav)
	}

	favorites, err := services.GetAllFavorites(userData.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// After fetching from the database
	services.SetCacheEntry(userFavKey, favorites, 10*time.Minute)

	return c.JSON(favorites)
}

// GetMyReportsHandler handles requests to get all reports submitted by the logged-in user.
func getMyReportsHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int) // Make sure to handle errors and cases where userID might not be set
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID not found"})
	}

	reports, err := services.GetAllReportsByUser(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get reports"})
	}

	if len(reports) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No reports found"})
	}

	return c.JSON(reports)
}

// DeleteObjectFromS3Handler handles requests to delete objects from S3.
func deleteObjectFromS3Handler(c *fiber.Ctx) error {
	var requestBody struct {
		ObjectURL string `json:"objectUrl"`
	}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Request body is not valid"})
	}

	// Ensure the object URL is not empty
	if requestBody.ObjectURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Object URL is required"})
	}

	// Call the service function to delete the object from S3
	if err := services.DeleteDataFromS3(requestBody.ObjectURL); err != nil {
		// Determine if the error should be a 404 not found or a 500 internal server error
		if err.Error() == "object not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Object not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete object from S3"})
	}

	// Return a success response
	return c.SendStatus(fiber.StatusNoContent)
}
