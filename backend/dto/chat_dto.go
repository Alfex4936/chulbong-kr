package dto

// ConnectionInfo structure to hold connection metadata
type ConnectionInfo struct {
	UserID   string `json:"userID"`
	RoomID   string `json:"roomID"`
	Username string `json:"username"`
	ConnID   string `json:"connID"`
}

type BroadcastMessage struct {
	Timestamp    int64  `json:"timestamp"` // Unix timestamp
	UID          string `json:"uid"`
	Message      string `json:"message"`
	UserID       string `json:"userId"`
	UserNickname string `json:"userNickname"`
	RoomID       string `json:"roomID"`
	// IsOwner      bool   `json:"isOwner,omitempty"`
}

type UserCountMessage struct {
	RoomID       string `json:"roomID"`
	UserNickname string `json:"userNickname"`
	UserCount    int    `json:"userCount"`
}
