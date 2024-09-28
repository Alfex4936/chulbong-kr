package handler

import (
	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/Alfex4936/chulbong-kr/facade"
	"github.com/Alfex4936/chulbong-kr/middleware"
	"github.com/Alfex4936/chulbong-kr/service"
	sonic "github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	UserFacadeService *facade.UserFacadeService

	AuthMiddleware *middleware.AuthMiddleware

	CacheService *service.MarkerCacheService
}

// NewUserHandler creates a new UserHandler with dependencies injected
func NewUserHandler(authMiddleware *middleware.AuthMiddleware, facade *facade.UserFacadeService, c *service.MarkerCacheService,
) *UserHandler {
	return &UserHandler{
		UserFacadeService: facade,

		AuthMiddleware: authMiddleware,

		CacheService: c,
	}
}

func RegisterUserRoutes(api fiber.Router, handler *UserHandler, authMiddleware *middleware.AuthMiddleware) {

	// Route to serve the gallery
	api.Get("/gallery", handler.HandleGalleryView)
	api.Get("/login-home", handler.HandleServerLogin)

	userGroup := api.Group("/users")
	{
		userGroup.Use(authMiddleware.Verify)
		userGroup.Get("/me", authMiddleware.VerifySoft, handler.HandleProfile)
		userGroup.Get("/favorites", handler.HandleGetFavorites)
		userGroup.Get("/reports", handler.HandleGetMyReports)                          // getting reports that I made
		userGroup.Get("/reports/for-my-markers", handler.HandleGetReportsForMyMarkers) // getting reports for my markers
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update user profile"})
	}

	// Marshal the updated user profile to byte array
	// userProfileData, err := sonic.Marshal(user)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to encode user profile"})
	// }

	// Update the user profile cache
	// h.CacheService.SetUserProfileCache(userData.UserID, userProfileData)

	go h.CacheService.ResetUserProfileCache(userData.UserID)

	// TODO: reset favorite markers cache if profile update affects them
	// h.UserFacadeService.ResetUserFavCache(userData.UserID)

	return c.JSON(user)
}

// DeleteUserHandler deletes the currently authenticated user
func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok || userID == 1 { // Prevent deletion of the admin user
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID not found"})
	}

	// log.Printf("[DEBUG][HANDLER] Deleting user %v", userID)

	if err := h.UserFacadeService.DeleteUserWithRelatedData(c.Context(), userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	return c.SendStatus(fiber.StatusNoContent) // 204 for successful deletion with no content in response
}

func (h *UserHandler) HandleProfile(c *fiber.Ctx) error {
	userData, err := h.UserFacadeService.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to get user from context"})
	}

	c.Append("Content-Type", "application/json")

	chulbong, _ := c.Locals("chulbong").(bool)

	// Try to get the user profile as byte data from the cache first
	cachedData, cacheErr := h.CacheService.GetUserProfileCache(userData.UserID)
	if cacheErr == nil && len(cachedData) > 0 {
		// Cache hit, return the cached data as JSON directly
		c.Append("X-Cache", "hit")
		return c.Send(cachedData)
	}

	// If the cache doesn't have the user profile, fetch it from the database
	user, err := h.UserFacadeService.GetUserById(userData.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user"})
	}

	contributions, creations, err := h.UserFacadeService.UserService.GetUserStatistics(userData.UserID)
	if err == nil {
		user.ReportCount = contributions
		user.MarkerCount = creations
	}

	contributions, level, err := h.UserFacadeService.UserService.GetUserContributionScores(userData.UserID)
	if err == nil {
		user.ContributionCount = contributions
		user.ContributionLevel = level
	}

	// Check adminship
	if chulbong {
		user.Chulbong = true
	}

	// Marshal the user profile directly into byte data
	userProfileData, err := sonic.Marshal(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to marshal user profile"})
	}

	// Cache the marshalled user profile for future requests
	go h.CacheService.SetUserProfileCache(userData.UserID, userProfileData)

	// Return the user profile
	return c.Send(userProfileData)
}

func (h *UserHandler) HandleGetFavorites(c *fiber.Ctx) error {
	userData, err := h.UserFacadeService.GetUserFromContext(c)
	if err != nil {
		return err // fiber err
	}

	// Try to get the user fav from the cache first
	cachedFav, cacheErr := h.CacheService.GetUserFavorites(userData.UserID)
	if cacheErr == nil && cachedFav != nil {
		// Cache hit, return the cached fav
		c.Append("X-Cache", "hit")
		return c.JSON(cachedFav)
	}

	favorites, err := h.UserFacadeService.GetAllFavorites(userData.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// After fetching from the database
	go h.CacheService.AddFavoritesToCache(userData.UserID, favorites)

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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get reports:" + err.Error()})
	}

	if len(reports) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No reports found"})
	}

	return c.JSON(reports)
}

// HandleGetReportsForMyMarkers handles requests to get all reports for my markers
func (h *UserHandler) HandleGetReportsForMyMarkers(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int) // Make sure to handle errors and cases where userID might not be set
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID not found"})
	}

	reports, err := h.UserFacadeService.GetAllReportsForMyMarkersByUser(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get reports: " + err.Error()})
	}

	if reports.TotalReports == 0 {
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

func (h *UserHandler) HandleGalleryView(c *fiber.Ctx) error {
	photos, err := h.UserFacadeService.S3Service.ListAllObjectsInS3()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Render("index", fiber.Map{
		"Title":  "Photo Gallery",
		"Photos": photos,
	})
}

func (h *UserHandler) HandleServerLogin(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{})
}
