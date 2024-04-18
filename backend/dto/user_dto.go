package dto

import "chulbong-kr/models"

type SignUpRequest struct {
	Username   *string `json:"username,omitempty"` // Optional for traditional, used if provided for OAuth2
	Email      string  `json:"email"`
	Password   string  `json:"password,omitempty"`   // Optional, not used for OAuth2 sign-ups
	Provider   string  `json:"provider,omitempty"`   // e.g., "google", "kakao", empty for traditional sign-ups
	ProviderID string  `json:"providerId,omitempty"` // Unique ID from the provider, empty for traditional sign-ups
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

type UpdateUserRequest struct {
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
}

// UserData holds the information extracted from the context.
type UserData struct {
	UserID   int
	Username string
}

type UserMarkers struct {
	MarkersWithPhotos []MarkerSimpleWithDescrption `json:"markers"`
	CurrentPage       int                          `json:"currentPage"`
	TotalPages        int                          `json:"totalPages"`
	TotalMarkers      int                          `json:"totalMarkers"`
}
