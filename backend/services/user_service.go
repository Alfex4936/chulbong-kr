package services

import (
	"chulbong-kr/database"
	"chulbong-kr/dto"
	"chulbong-kr/models"
	"chulbong-kr/util"
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var (
	TOKEN_DURATION         time.Duration
	NAVER_EMAIL_VERIFY_URL = os.Getenv("NAVER_EMAIL_VERIFY_URL")
)

// GetUserById retrieves a user by their email address
func GetUserById(userID int) (*models.User, error) {
	var user models.User

	// Define the query to select the user
	query := `SELECT UserID, Username, Email, PasswordHash, Provider, ProviderID, CreatedAt, UpdatedAt FROM Users WHERE UserID = ?`

	// Execute the query
	err := database.DB.Get(&user, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with userID %d", userID)
		}
		return nil, fmt.Errorf("error fetching user by userID: %w", err)
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by their email address
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	// Define the query to select the user
	query := `SELECT UserID, Username, Email, PasswordHash, Provider, ProviderID, CreatedAt, UpdatedAt FROM Users WHERE Email = ?`

	// Execute the query
	err := database.DB.Get(&user, query, email)
	if err != nil {
		return nil, err
		// if err == sql.ErrNoRows {
		// 	// No user found with the provided email
		// 	return nil, fmt.Errorf("no user found with email %s", email)
		// }
		// // An error occurred during the query execution
		// return nil, fmt.Errorf("error fetching user by email: %w", err)
	}

	return &user, nil
}

// SaveUser creates a new user with hashed password
func SaveUser(signUpReq *dto.SignUpRequest) (*models.User, error) {
	tx, err := database.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	hashedPassword, err := hashPassword(signUpReq.Password)
	if err != nil {
		return nil, err
	}

	userID, err := insertUserWithRetry(tx, signUpReq, hashedPassword)
	if err != nil {
		return nil, err
	}

	newUser, err := fetchNewUser(tx, userID)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("DELETE FROM PasswordTokens WHERE Email = ? AND Verified = TRUE", newUser.Email)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error removing verified token: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return newUser, nil
}

// Login checks if a user exists with the given email and password.
func Login(email, password string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT UserID, Username, Email, PasswordHash, Provider FROM Users WHERE Email = ?`
	err := database.DB.Get(user, query, email)
	if err != nil {
		return nil, err // User not found or db error
	}

	// Check if the user was registered through an external provider
	if user.Provider.Valid && user.Provider.String != "website" {
		// The user did not register through the website's traditional sign-up process
		return nil, fmt.Errorf("external provider login not supported here")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(password)) // heavy
	if err != nil {
		// Password does not match
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

func UpdateUserProfile(userID int, updateReq *dto.UpdateUserRequest) (*models.User, error) {
	tx, err := database.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if updateReq.Username != nil {
		var existingID int
		err = tx.Get(&existingID, "SELECT UserID FROM Users WHERE Username = ?", *updateReq.Username)
		if err == nil || err != sql.ErrNoRows {
			return nil, fmt.Errorf("username %s is already in use", *updateReq.Username)
		}
	}

	if updateReq.Email != nil {
		var existingID int
		err = tx.Get(&existingID, "SELECT UserID FROM Users WHERE Email = ?", *updateReq.Email)
		if err == nil || err != sql.ErrNoRows {
			return nil, fmt.Errorf("email %s is already in use", *updateReq.Email)
		}
	}

	var setParts []string
	var args []any

	if updateReq.Username != nil {
		setParts = append(setParts, "Username = ?")
		args = append(args, *updateReq.Username)
	}

	if updateReq.Email != nil {
		setParts = append(setParts, "Email = ?")
		args = append(args, *updateReq.Email)
	}

	if updateReq.Password != nil {
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(*updateReq.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			return nil, hashErr
		}
		setParts = append(setParts, "PasswordHash = ?")
		args = append(args, string(hashedPassword))
	}

	if len(setParts) > 0 {
		args = append(args, userID)
		query := fmt.Sprintf("UPDATE Users SET %s WHERE UserID = ?", strings.Join(setParts, ", "))
		_, err = tx.Exec(query, args...)
		if err != nil {
			return nil, fmt.Errorf("error updating user: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing update: %w", err)
	}

	// Fetch the updated user details
	updatedUser, err := GetUserById(userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching updated user: %w", err)
	}

	return updatedUser, nil
}

// GetAllReportsByUser retrieves all reports submitted by a specific user from the database.
func GetAllReportsByUser(userID int) ([]dto.MarkerReportResponse, error) {
	const query = `
    SELECT ReportID, MarkerID, UserID, ST_X(Location) AS Latitude, ST_Y(Location) AS Longitude,
           Description, ReportImageURL, CreatedAt
    FROM Reports
    WHERE UserID = ?
    ORDER BY CreatedAt DESC
    `
	reports := make([]dto.MarkerReportResponse, 0)
	if err := database.DB.Select(&reports, query, userID); err != nil {
		return nil, fmt.Errorf("error querying reports: %w", err)
	}

	return reports, nil
}

func GetAllFavorites(userID int) ([]dto.MarkerSimpleWithDescrption, error) {
	favorites := make([]dto.MarkerSimpleWithDescrption, 0)
	const query = `
    SELECT Markers.MarkerID, ST_X(Markers.Location) AS Latitude, ST_Y(Markers.Location) AS Longitude, Markers.Description, Markers.Address
    FROM Favorites
    JOIN Markers ON Favorites.MarkerID = Markers.MarkerID
    WHERE Favorites.UserID = ?
    ORDER BY Markers.CreatedAt DESC` // Order by CreatedAt in descending order

	err := database.DB.Select(&favorites, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching favorites: %w", err)
	}

	return favorites, nil
}

// DeleteUserWithRelatedData
func DeleteUserWithRelatedData(ctx context.Context, userID int) error {
	// Begin a transaction
	tx, err := database.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}

	// Fetch Photo URLs associated with the user
	var photoURLs []string
	fetchPhotosQuery := `SELECT PhotoURL FROM Photos WHERE MarkerID IN (SELECT MarkerID FROM Markers WHERE UserID = ?)`
	if err := tx.SelectContext(ctx, &photoURLs, fetchPhotosQuery, userID); err != nil {
		tx.Rollback()
		return fmt.Errorf("fetching photo URLs: %w", err)
	}

	// Delete each photo from S3
	for _, url := range photoURLs {
		if err := DeleteDataFromS3(url); err != nil {
			tx.Rollback()
			return fmt.Errorf("deleting photo from S3: %w", err)
		}
	}

	// Note: Order matters due to foreign key constraints
	deletionQueries := []string{
		"DELETE FROM OpaqueTokens WHERE UserID = ?",
		"DELETE FROM Comments WHERE UserID = ?",
		"DELETE FROM MarkerDislikes WHERE UserID = ?",
		"DELETE FROM Photos WHERE MarkerID IN (SELECT MarkerID FROM Markers WHERE UserID = ?)",
		"UPDATE Markers SET UserID = NULL WHERE UserID = ?", // Set UserID to NULL for Markers instead of deleting
		"DELETE FROM Users WHERE UserID = ?",
	}

	// Execute each deletion query within the transaction
	for _, query := range deletionQueries {
		if _, err := tx.ExecContext(ctx, query, userID); err != nil {
			tx.Rollback() // Attempt to rollback, but don't override the original error
			return fmt.Errorf("executing deletion query (%s): %w", query, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

func ResetPassword(token string, newPassword string) error {
	// Start a transaction
	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}

	// Ensure the transaction is rolled back if an error occurs
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var userID int
	// Use the transaction (tx) to perform the query
	err = tx.Get(&userID, "SELECT UserID FROM PasswordResetTokens WHERE Token = ? AND ExpiresAt > NOW()", token)
	if err != nil {
		return err // Token not found or expired
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Use the transaction (tx) to update the user's password
	_, err = tx.Exec("UPDATE Users SET PasswordHash = ? WHERE UserID = ?", string(hashedPassword), userID)
	if err != nil {
		return err
	}

	// Use the transaction (tx) to delete the reset token
	_, err = tx.Exec("DELETE FROM PasswordResetTokens WHERE Token = ?", token)
	if err != nil {
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

func GeneratePasswordResetToken(email string) (string, error) {
	user := models.User{}
	err := database.DB.Get(&user, "SELECT UserID FROM Users WHERE Email = ?", email)
	if err != nil {
		return "", err // User not found or db error
	}

	token, err := util.GenerateOpaqueToken()
	if err != nil {
		return "", err
	}

	_, err = database.DB.Exec(`
    INSERT INTO PasswordResetTokens (UserID, Token, ExpiresAt)
    VALUES (?, ?, ?)
    ON DUPLICATE KEY UPDATE Token = VALUES(Token), ExpiresAt = VALUES(ExpiresAt)`,
		user.UserID, token, time.Now().Add(24*time.Hour))
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserFromContext extracts and validates the user information from the Fiber context.
func GetUserFromContext(c *fiber.Ctx) (*dto.UserData, error) {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	username, ok := c.Locals("username").(string)
	if !ok {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username not found"})
	}

	return &dto.UserData{
		UserID:   userID,
		Username: username,
	}, nil
}

// VerifyNaverEmail can check naver email existence before sending
func VerifyNaverEmail(naverAddress string) (bool, error) {
	naverAddress = strings.Split(naverAddress, "@naver.com")[0]
	reqURL := fmt.Sprintf("%s=%s", NAVER_EMAIL_VERIFY_URL, naverAddress)
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return false, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	resp, err := HTTPClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read the body: %v", err)
	}

	// Convert body bytes to a string
	bodyString := string(bodyBytes)

	// Check if the body is non-empty and ends with 'N'
	if len(bodyString) > 0 && bodyString[len(bodyString)-1] == 'N' {
		return true, nil
	} else {
		return false, nil
	}
}

// private
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func generateUsername(signUpReq *dto.SignUpRequest) string {
	if signUpReq.Username != nil && *signUpReq.Username != "" {
		return *signUpReq.Username
	}
	emailParts := strings.Split(signUpReq.Email, "@")
	return emailParts[0]
}

func insertUserWithRetry(tx *sqlx.Tx, signUpReq *dto.SignUpRequest, hashedPassword string) (int64, error) {
	username := generateUsername(signUpReq)
	const maxRetries = 5
	for i := 0; i < maxRetries; i++ {
		res, err := tx.Exec(`INSERT INTO Users (Username, Email, PasswordHash, Provider, ProviderID, Role, CreatedAt, UpdatedAt) VALUES (?, ?, ?, ?, ?, 'user', NOW(), NOW())`,
			username, signUpReq.Email, hashedPassword, signUpReq.Provider, signUpReq.ProviderID)
		if err != nil {
			if strings.Contains(err.Error(), "Duplicate entry") && strings.Contains(err.Error(), "for key 'idx_users_username'") {
				username = fmt.Sprintf("%s-%s", username, util.GenerateRandomString(5))
				continue
			}
			return 0, fmt.Errorf("error registering user: %w", err)
		}
		userID, _ := res.LastInsertId()
		return userID, nil
	}
	return 0, fmt.Errorf("failed to insert user after retries")
}

func fetchNewUser(tx *sqlx.Tx, userID int64) (*models.User, error) {
	var newUser models.User
	query := `SELECT UserID, Username, Email, Provider, ProviderID, Role, CreatedAt, UpdatedAt FROM Users WHERE UserID = ?`
	err := tx.QueryRowx(query, userID).StructScan(&newUser)
	if err != nil {
		return nil, fmt.Errorf("error fetching newly created user: %w", err)
	}
	return &newUser, nil
}
