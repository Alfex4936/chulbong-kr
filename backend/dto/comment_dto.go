package dto

type CommentRequest struct {
	MarkerID    int    `json:"markerId"`
	CommentText string `json:"commentText"`
}

type CommentLoadParams struct {
	N    int `query:"n"`
	Page int `query:"page"`
}
