package notification

import "github.com/goccy/go-json"

type NotificationRedis struct {
	NotificationId   int64           `json:"notificationId" db:"NotificationId"`
	UserId           string          `json:"userId" db:"UserId"`
	NotificationType string          `json:"type" db:"NotificationType"`
	Title            string          `json:"title" db:"Title"`
	Message          string          `json:"message" db:"Message"`
	Metadata         json.RawMessage `json:"metadata" db:"Metadata"`
}

type NotificationMarkerMetadata struct {
	MarkerID  int64   `json:"markerID"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
}

type NotificationLikeMetadata struct {
	MarkerID int `json:"markerID"`
	UserId   int `json:"userId"`
	LikerId  int `json:"likerId"`
}
