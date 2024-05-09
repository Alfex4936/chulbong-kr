package model

import (
	"time"

	"github.com/goccy/go-json"
)

// Notification represents the schema of the Notifications table
type Notification struct {
	NotificationId   int64           `json:"notificationId" db:"NotificationId"`
	UserId           string          `json:"userId" db:"UserId"`
	NotificationType string          `json:"notifi_type" db:"NotificationType"`
	Title            string          `json:"title" db:"Title"`
	Message          string          `json:"message" db:"Message"`
	Metadata         json.RawMessage `json:"metadata" db:"Metadata"`
	Viewed           bool            `json:"-" db:"Viewed"`
	CreatedAt        time.Time       `json:"-" db:"CreatedAt"`
	UpdatedAt        time.Time       `json:"-" db:"UpdatedAt"`
}
