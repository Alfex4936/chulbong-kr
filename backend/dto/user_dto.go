package dto

import "chulbong-kr/models"

type SignUpRequest struct {
	Username *string `json:"username,omitempty"` // optional
	Email    string  `json:"email"`
	Password string  `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}
