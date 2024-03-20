package services

import (
	"fmt"
	"strings"

	"chulbong-kr/database"
)

// LeaveDislike user's dislike for a marker
func LeaveDislike(userID int, markerID int) error {
	_, err := database.DB.Exec(
		"INSERT INTO MarkerDislikes (MarkerID, UserID) VALUES (?, ?) ON DUPLICATE KEY UPDATE DislikedAt=VALUES(DislikedAt)",
		markerID, userID,
	)
	if err != nil {
		return fmt.Errorf("inserting dislike: %w", err)
	}
	return nil
}

// UndoDislike nudo user's dislike for a marker
func UndoDislike(userID int, markerID int) error {
	result, err := database.DB.Exec(
		"DELETE FROM MarkerDislikes WHERE UserID = ? AND MarkerID = ?",
		userID, markerID,
	)
	if err != nil {
		return fmt.Errorf("deleting dislike: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no dislike found to undo")
	}

	return nil
}

// This service function checks if the given user has disliked the specified marker.
func CheckUserDislike(userID, markerID int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM MarkerDislikes WHERE UserID = ? AND MarkerID = ?)"
	err := database.DB.Get(&exists, query, userID, markerID)
	return exists, err
}

func CheckUserFavorite(userID, markerID int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM Favorites WHERE UserID = ? AND MarkerID = ?)"
	err := database.DB.Get(&exists, query, userID, markerID)
	return exists, err
}

// AddFavoriteHandler adds a new favorite marker for the user.
func AddFavorite(userID, markerID int) error {
	// First, count the existing favorites for the user
	var count int
	err := database.DB.QueryRowx("SELECT COUNT(*) FROM Favorites WHERE UserID = ?", userID).Scan(&count)
	if err != nil {
		return err
	}

	// Check if the user already has 10 favorites
	if count >= 10 {
		return fmt.Errorf("maximum number of favorites reached")
	}

	// If not, insert the new favorite
	_, err = database.DB.Exec("INSERT INTO Favorites (UserID, MarkerID) VALUES (?, ?)", userID, markerID)
	if err != nil {
		// Convert error to string and check if it contains the MySQL error code for duplicate entry
		if strings.Contains(err.Error(), "1062") {
			return fmt.Errorf("you have already favorited this marker")
		}
		return fmt.Errorf("failed to add favorite: %w", err)
	}

	// Retrieve the UserID (owner) of the marker
	var ownerUserID int
	err = database.DB.QueryRowx("SELECT UserID FROM Markers WHERE MarkerID = ?", markerID).Scan(&ownerUserID)
	if err != nil {
		return fmt.Errorf("failed to retrieve marker owner: %w", err)
	}

	// TODO: update when frontend updates
	// key := fmt.Sprintf("%d-%d", ownerUserID, markerID)
	// PublishLikeEvent(key)
	return nil
}

func RemoveFavorite(userID, markerID int) error {
	// Delete the specified favorite for the user
	_, err := database.DB.Exec("DELETE FROM Favorites WHERE UserID = ? AND MarkerID = ?", userID, markerID)
	if err != nil {
		return fmt.Errorf("error removing favorite: %w", err)
	}
	return nil
}
