package notification

// Notification is the base struct for all notifications.
// It includes a Type field to indicate the specific type of notification.
type Notification struct {
	Type string `json:"type"`
}

// BroadcastNotification represents a notification for all
type BroadcastNotification struct {
	Notification        // Embedding Notification struct to inherit the Type field.
	Notice       string `json:"notice"`
}

// LikeNotification represents a notification for a like event.
type LikeNotification struct {
	Notification     // Embedding Notification struct to inherit the Type field.
	UserID       int `json:"userId"`
	MarkerID     int `json:"markerId"`
}

// CommentNotification represents a notification for a comment event.
type CommentNotification struct {
	Notification        // Embedding Notification struct to inherit the Type field.
	UserID       int    `json:"userId"`
	MarkerID     int    `json:"markerId"`
	Comment      string `json:"comment"`
}
