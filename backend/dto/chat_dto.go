package dto

// ConnectionInfo structure to hold connection metadata
type ConnectionInfo struct {
	UserID   string `json:"userID"`
	RoomID   string `json:"roomID"`
	Username string `json:"username"`
	ConnID   string `json:"connID"`
}

type BroadcastMessage struct {
	UID          string `json:"uid"`
	Message      string `json:"message"`
	UserID       string `json:"userId"`
	UserNickname string `json:"userNickname"`
	RoomID       string `json:"roomID"`
	Timestamp    int64  `json:"timestamp"` // Unix timestamp
	// IsOwner      bool   `json:"isOwner,omitempty"`
}

type BroadcastRoomInfoMessage struct {
	RoomID     string `json:"roomID"`
	TotalUsers string `json:"totalUsers"`
	Timestamp  int64  `json:"timestamp"`
}
