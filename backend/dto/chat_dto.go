package dto

// ConnectionInfo structure to hold connection metadata
type ConnectionInfo struct {
	UserID   int    `json:"userID"`
	RoomID   string `json:"roomID"`
	Username string `json:"username"`
	ConnID   string `json:"connID"` // Unique connection identifier
}

type BroadcastMessage struct {
	UID          string `json:"uid"`
	Message      string `json:"message"`
	UserID       int    `json:"userId"`
	UserNickname string `json:"userNickname"`
	RoomID       string `json:"roomID"`
	Timestamp    int64  `json:"timestamp"` // Unix timestamp
}
