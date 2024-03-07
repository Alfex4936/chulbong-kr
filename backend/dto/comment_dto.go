package dto

type CommentRequest struct {
	MarkerID    int    `json:"markerId"`
	CommentText string `json:"commentText"`
}
