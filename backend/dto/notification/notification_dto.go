package notification

import "github.com/goccy/go-json"

type NotificationRedis struct {
	NotificationId   int64           `json:"notificationId" db:"NotificationId"`
	Metadata         json.RawMessage `json:"metadata" db:"Metadata"`
	UserId           string          `json:"userId" db:"UserId"`
	NotificationType string          `json:"type" db:"NotificationType"`
	Title            string          `json:"title" db:"Title"`
	Message          string          `json:"message" db:"Message"`
}

type NotificationMarkerMetadata struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	MarkerID  int64   `json:"markerID"`
	Address   string  `json:"address"`
}

type NotificationLikeMetadata struct {
	MarkerID int `json:"markerID"`
	UserId   int `json:"userId"`
	LikerId  int `json:"likerId"`
}
