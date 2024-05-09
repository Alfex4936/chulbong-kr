package handler

import (
	"log"
	"time"

	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/Alfex4936/chulbong-kr/facade"
	"github.com/Alfex4936/chulbong-kr/middleware"
	"github.com/Alfex4936/chulbong-kr/model"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	UserFacadeService *facade.UserFacadeService

	AuthMiddleware *middleware.AuthMiddleware
}

// NewUserHandler creates a new UserHandler with dependencies injected
func NewUserHandler(authMiddleware *middleware.AuthMiddleware, facade *facade.UserFacadeService,
) *UserHandler {
	return &UserHandler{
		UserFacadeService: facade,

		AuthMiddleware: authMiddleware,
	}
}

func RegisterUserRoutes(api fiber.Router, handler *UserHandler, authMiddleware *middleware.AuthMiddleware) {
	userGroup := api.Group("/users")
	{
		userGroup.Use(authMiddleware.Verify)
		userGroup.Get("/me", handler.HandleProfile)
		userGroup.Get("/favorites", handler.HandleGetFavorites)
		userGroup.Get("/reports", handler.HandleGetMyReports)
		userGroup.Patch("/me", handler.HandleUpdateUser)
		userGroup.Delete("/me", handler.HandleDeleteUser)
		userGroup.Delete("/s3/objects", authMiddleware.CheckAdmin, handler.HandleDeleteObjectFromS3)
	}
}

// UpdateUserHandler
func (h *UserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	userData, err := h.UserFacadeService.GetUserFromContext(c)
	if err != nil {
		return err // fiber err
	}

	var updateReq dto.UpdateUserRequest
	if err := c.BodyParser(&updateReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	user, err := h.UserFacadeService.UpdateUserProfile(userData.UserID, &updateReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	h.UserFacadeService.ResetUserFavCache(userData.UserID)

	return c.JSON(user)
}

// DeleteUserHandler deletes the currently authenticated user
func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID not found"})
	}

	log.Printf("[DEBUG][HANDLER] Deleting user %v", userID)

	if err := h.UserFacadeService.DeleteUserWithRelatedData(c.Context(), userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent) // 204 for successful deletion with no content in response
}

func (h *UserHandler) HandleProfile(c *fiber.Ctx) error {
	userData, err := h.UserFacadeService.GetUserFromContext(c)
	if err != nil {
		return err // fiber err
	}

	userProfileKey := h.UserFacadeService.GetUserProfileKey(userData.UserID)

	var cachedUser *model.User
	// Try to get the user profile from the cache first
	cacheErr := h.UserFacadeService.GetUserCache(userProfileKey, &cachedUser)
	if cacheErr == nil && cachedUser != nil {
		// Cache hit, return the cached user
		c.Append("X-Cache", "hit")
		return c.JSON(cachedUser)
	}

	// If the cache doesn't have the user profile, fetch it from the database
	user, err := h.UserFacadeService.GetUserById(userData.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// After fetching from the database
	h.UserFacadeService.SetRedisCache(userProfileKey, user, 15*time.Minute)

	return c.JSON(user)
}

func (h *UserHandler) HandleGetFavorites(c *fiber.Ctx) error {
	userData, err := h.UserFacadeService.GetUserFromContext(c)
	if err != nil {
		return err // fiber err
	}

	userFavKey := h.UserFacadeService.GetUserFavKey(userData.UserID, userData.Username)

	var cachedFav []dto.MarkerSimpleWithDescrption

	// Try to get the user fav from the cache first
	cacheErr := h.UserFacadeService.GetUserCache(userFavKey, &cachedFav)
	if cacheErr == nil && cachedFav != nil {
		// Cache hit, return the cached fav
		return c.JSON(cachedFav)
	}

	favorites, err := h.UserFacadeService.GetAllFavorites(userData.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// After fetching from the database
	h.UserFacadeService.SetRedisCache(userFavKey, favorites, 10*time.Minute)

	return c.JSON(favorites)
}

// GetMyReportsHandler handles requests to get all reports submitted by the logged-in user.
func (h *UserHandler) HandleGetMyReports(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int) // Make sure to handle errors and cases where userID might not be set
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID not found"})
	}

	reports, err := h.UserFacadeService.GetAllReportsByUser(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get reports"})
	}

	if len(reports) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No reports found"})
	}

	return c.JSON(reports)
}

// DeleteObjectFromS3Handler handles requests to delete objects from S3.
func (h *UserHandler) HandleDeleteObjectFromS3(c *fiber.Ctx) error {
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
	if err := h.UserFacadeService.DeleteDataFromS3(requestBody.ObjectURL); err != nil {
		// Determine if the error should be a 404 not found or a 500 internal server error
		if err.Error() == "object not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Object not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete object from S3"})
	}

	// Return a success response
	return c.SendStatus(fiber.StatusNoContent)
}
