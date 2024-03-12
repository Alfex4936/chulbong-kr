package handlers

import (
	"chulbong-kr/dto"
	"chulbong-kr/models"
	"chulbong-kr/services"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// UpdateUserHandler
func UpdateUserHandler(c *fiber.Ctx) error {
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

	userProfileKey := fmt.Sprintf("%s:%d:%s", services.USER_PROFILE_KEY, userData.UserID, userData.Username)
	services.ResetCache(userProfileKey)

	return c.JSON(user)
}

// DeleteUserHandler deletes the currently authenticated user
func DeleteUserHandler(c *fiber.Ctx) error {
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

func ProfileHandler(c *fiber.Ctx) error {
	userData, err := services.GetUserFromContext(c)
	if err != nil {
		return err // fiber err
	}

	userProfileKey := fmt.Sprintf("%s:%d:%s", services.USER_PROFILE_KEY, userData.UserID, userData.Username)

	// Try to get the user profile from the cache first
	cachedUser, cacheErr := services.GetCacheEntry[*models.User](userProfileKey)
	if cacheErr == nil && cachedUser != nil {
		// Cache hit, return the cached user
		return c.JSON(cachedUser)
	}

	// If the cache doesn't have the user profile, fetch it from the database
	user, err := services.GetUserById(userData.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// After fetching from the database
	services.SetCacheEntry(userProfileKey, user, 10*time.Minute)

	return c.JSON(user)
}

func GetFavoritesHandler(c *fiber.Ctx) error {
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
