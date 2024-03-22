package services

import (
	"chulbong-kr/dto/notification"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"

	"github.com/gofiber/contrib/websocket"
	"github.com/redis/go-redis/v9"

	"github.com/alphadose/haxmap"
)

var WsManager *ConnectionManager = NewConnectionManager()

// PublishMarkerUpdate to publish messages
func PublishMarkerUpdate(message string) {
	err := RedisStore.Conn().Publish(context.Background(), "markerUpdates", message).Err()
	if err != nil {
		panic(err)
	}
}

// SubscribeToMarkerUpdates to subscribe to messages
func SubscribeToMarkerUpdates() *redis.PubSub {
	pubsub := RedisStore.Conn().Subscribe(context.Background(), "markerUpdates")
	return pubsub
}

// PublishLikeEvent
func PublishLikeEvent(userID string) {
	err := RedisStore.Conn().Publish(context.Background(), "markerLikes", userID).Err()
	if err != nil {
		panic(err)
	}
}

// SubscribeToMarkerUpdates
func SubscribeToLikeEvent() *redis.PubSub {
	pubsub := RedisStore.Conn().Subscribe(context.Background(), "markerLikes")
	return pubsub
}

func sendBroadcastNotification(c *websocket.Conn, message string) {
	notif := notification.BroadcastNotification{
		Notification: notification.Notification{Type: "all"},
		Notice:       message,
	}

	jsonNotif, err := json.Marshal(notif)
	if err != nil {
		log.Println("Error marshalling notice notification:", err)
		return
	}

	if err := c.WriteMessage(websocket.TextMessage, jsonNotif); err != nil {
		log.Println("Error sending notice notification:", err)
	}
}

func MakeLikeNotification(userID, markerID int) []byte {
	notif := notification.LikeNotification{
		Notification: notification.Notification{Type: "like"},
		UserID:       userID,
		MarkerID:     markerID,
	}

	jsonNotif, err := json.Marshal(notif)
	if err != nil {
		log.Println("Error marshalling like notification:", err)
		return nil
	}

	return jsonNotif
}

func MakeCommentNotification(userID, markerID int, comment string) []byte {
	notif := notification.CommentNotification{
		Notification: notification.Notification{Type: "comment"},
		UserID:       userID,
		MarkerID:     markerID,
		Comment:      comment,
	}

	jsonNotif, err := json.Marshal(notif)
	if err != nil {
		log.Println("Error marshalling comment notification:", err)
		return nil
	}

	return jsonNotif
}

type ConnectionManager struct {
	connections *haxmap.Map[string, *websocket.Conn]
}

// NewConnectionManager initializes a ConnectionManager with a new haxmap instance
func NewConnectionManager() *ConnectionManager {
	manager := &ConnectionManager{
		connections: haxmap.New[string, *websocket.Conn](),
	}
	// Start the connection checker
	manager.StartConnectionChecker()
	return manager
}

// AddConnection stores a WebSocket connection associated with a userID
func (manager *ConnectionManager) AddConnection(userID int, conn *websocket.Conn) {
	var key string
	if userID > 0 {
		key = fmt.Sprintf("user_%d", userID)
	} else {
		// Generate a unique identifier for a non-logged-in user.
		uniqueID := uuid.New().String()
		key = fmt.Sprintf("guest_%s", uniqueID)
	}
	manager.connections.Set(key, conn)
}

// RemoveConnection removes a WebSocket connection associated with a userID
func (manager *ConnectionManager) RemoveConnection(userID int) {
	userIDStr := fmt.Sprint(userID)
	manager.connections.Del(userIDStr)
}

// SendMessageToUser sends a WebSocket message to a specific user
func (manager *ConnectionManager) SendMessageToUser(userID int, message string) error {
	userIDStr := fmt.Sprintf("user_%d", userID)

	conn, ok := manager.connections.Get(userIDStr)
	if !ok {
		return fmt.Errorf("connection not found for user %d", userID)
	}
	// Send the message via the WebSocket connection
	return conn.WriteMessage(websocket.TextMessage, []byte(message))
}

// BroadcastMessage sends a WebSocket message to all users
func (manager *ConnectionManager) BroadcastMessage(message []byte) {
	manager.connections.ForEach(func(key string, conn *websocket.Conn) bool {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Failed to send broadcast message to user %s: %v", key, err)
		}
		return true // Continue iteration
	})
}

// CheckConnections iterates over all connections and sends a ping message.
// Connections that fail to respond can be considered inactive and removed.
func (manager *ConnectionManager) CheckConnections() {
	manager.connections.ForEach(func(userIDStr string, conn *websocket.Conn) bool {
		// Set a write deadline before sending the ping.
		conn.SetWriteDeadline(time.Now().Add(60 * time.Second))
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			log.Printf("Ping failed for user %s: %v. Removing connection.", userIDStr, err)
			manager.RemoveConnectionByStr(userIDStr)
		}
		return true
	})
}

func (manager *ConnectionManager) StartConnectionChecker() {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for range ticker.C {
			manager.CheckConnections()
		}
	}()
}

func (manager *ConnectionManager) RemoveConnectionByStr(userIDStr string) {
	manager.connections.Del(userIDStr)
}
