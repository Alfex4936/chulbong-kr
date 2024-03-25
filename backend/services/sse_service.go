package services

import (
	"chulbong-kr/dto/notification"
	"context"
	"fmt"
	"log"
	"time"
	"unsafe"

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
	connections *haxmap.Map[string, []*websocket.Conn]
	// connectionsToRoom *haxmap.Map[uintptr, int] // *websocket.Conn cannot be key (not hashable)
}

// NewConnectionManager initializes a ConnectionManager with a new haxmap instance
func NewConnectionManager() *ConnectionManager {
	manager := &ConnectionManager{
		connections: haxmap.New[string, []*websocket.Conn](),
		// connectionsToRoom: haxmap.New[uintptr, int](),
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

	// First, try to get the current list of connections for the room.
	conns, ok := manager.connections.Get(key)
	if !ok {
		// If it doesn't exist, initialize it with the new connection.
		conns = []*websocket.Conn{conn}
	} else {
		// If it exists, append the new connection, ensuring not to duplicate.
		if !contains(conns, conn) {
			conns = append(conns, conn)
		}
	}

	// Update the map with the new or modified slice.
	manager.connections.Set(key, conns)

}

// AddConnectionRoom stores a WebSocket connection associated with a markerID
func (manager *ConnectionManager) AddConnectionRoom(markerID int, conn *websocket.Conn) {
	markerIDStr := fmt.Sprintf("%d", markerID)

	// First, try to get the current list of connections for the room.
	conns, ok := manager.connections.Get(markerIDStr)
	if !ok {
		// If it doesn't exist, initialize it with the new connection.
		conns = []*websocket.Conn{conn}
	} else {
		// If it exists, append the new connection, ensuring not to duplicate.
		if !contains(conns, conn) {
			conns = append(conns, conn)
		}
	}

	// Update the map with the new or modified slice.
	manager.connections.Set(markerIDStr, conns)
}

// func (cm *ConnectionManager) AddToRoom(conn *websocket.Conn, roomID int) {
// 	connKey := connPtrToUintptr(conn)
// 	cm.connectionsToRoom.Set(connKey, roomID)
// }

// RemoveConnection removes a WebSocket connection associated with a id
func (manager *ConnectionManager) RemoveConnectionFromRoom(markerID int, conn *websocket.Conn) {
	markerIDStr := fmt.Sprintf("%d", markerID)

	// Retrieve the current list of connections for the room.
	if conns, ok := manager.connections.Get(markerIDStr); ok {
		// Find and remove the specified connection from the slice.
		for i, c := range conns {
			if c == conn {
				conns = append(conns[:i], conns[i+1:]...)
				break
			}
		}

		// If the slice is empty after removal, delete the entry from the map.
		if len(conns) == 0 {
			manager.connections.Del(markerIDStr)
		} else {
			// Otherwise, update the map with the modified slice.
			manager.connections.Set(markerIDStr, conns)
		}
	}
}

// func (manager *ConnectionManager) RemoveRoomConnection(conn *websocket.Conn) {
// 	connKey := connPtrToUintptr(conn)
// 	// Use the connection-to-room mapping to find the corresponding markerID.
// 	if markerID, ok := manager.connectionsToRoom.Get(connKey); ok {
// 		// Convert markerID back to int if needed, then call RemoveConnectionFromRoom
// 		manager.RemoveConnectionFromRoom(markerID, conn)
// 	}

// 	// Finally, remove the connection-to-room mapping.
// 	manager.connectionsToRoom.Del(connKey)
// }

// SendMessageToUser sends a WebSocket message to a specific user
func (manager *ConnectionManager) SendMessageToUser(userID int, message string) error {
	userIDStr := fmt.Sprintf("user_%d", userID)

	conns, ok := manager.connections.Get(userIDStr)
	if !ok || len(conns) == 0 {
		return fmt.Errorf("no connections found for userID %d", userID)
	}

	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			// Log the error but continue sending to other connections
			log.Printf("Error sending message to user %d: %v", userID, err)
		}
	}

	return nil
}

// SendMessageToRoom sends a WebSocket message to a specific user
func (manager *ConnectionManager) SendMessageToRoom(markerID int, message string) error {
	markerIDStr := fmt.Sprintf("%d", markerID)

	conns, ok := manager.connections.Get(markerIDStr)
	if !ok || len(conns) == 0 {
		return fmt.Errorf("no connections found for markerID %d", markerID)
	}

	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			// Log the error but continue sending to other connections
			log.Printf("Error sending message to connection in room %d: %v", markerID, err)
		}
	}

	return nil
}

// BroadcastMessage sends a WebSocket message to all users
func (manager *ConnectionManager) BroadcastMessage(message []byte) {
	manager.connections.ForEach(func(key string, conns []*websocket.Conn) bool {
		for _, conn := range conns {
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Failed to send broadcast message to connection for %s: %v", key, err)
			}
		}
		return true // Continue iteration
	})
}

// CheckConnections iterates over all connections and sends a ping message.
// Connections that fail to respond can be considered inactive and removed.
func (manager *ConnectionManager) CheckConnections() {
	manager.connections.ForEach(func(userIDStr string, conns []*websocket.Conn) bool {
		var inactiveConns []int // Store indexes of inactive connections
		for i, conn := range conns {
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second)) // Adjust deadline as needed
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Ping failed for connection: %v. Marking as inactive.", err)
				inactiveConns = append(inactiveConns, i)
			}
		}
		// Remove inactive connections
		for _, i := range inactiveConns {
			conns = append(conns[:i], conns[i+1:]...)
		}
		// If all connections for a user are inactive, delete the user's entry
		if len(conns) == 0 {
			manager.connections.Del(userIDStr)
		} else {
			manager.connections.Set(userIDStr, conns) // Update the slice in the map
		}
		return true
	})

	// manager.connectionsToRoom.ForEach(func(connPtr uintptr, roomId int) bool {
	// 	conn := uintptrToConnPtr(connPtr)
	// 	// Set a write deadline before sending the ping.
	// 	conn.SetWriteDeadline(time.Now().Add(60 * time.Second))
	// 	if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
	// 		log.Printf("Ping failed for conn %s: %v. Removing connection.", conn.IP(), err)
	// 		manager.RemoveRoomConnection(conn)
	// 	}
	// 	return true
	// })
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

// Utility functions to convert between *websocket.Conn and uintptr
func connPtrToUintptr(conn *websocket.Conn) uintptr {
	return uintptr(unsafe.Pointer(conn))
}

func uintptrToConnPtr(u uintptr) *websocket.Conn {
	return (*websocket.Conn)(unsafe.Pointer(u))
}

// contains function to check if the slice already contains the given connection.
func contains(slice []*websocket.Conn, conn *websocket.Conn) bool {
	for _, item := range slice {
		if item == conn {
			return true
		}
	}
	return false
}
