package services

import (
	"chulbong-kr/dto"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
	csmap "github.com/mhmtszr/concurrent-swiss-map"
	"github.com/redis/go-redis/v9"
	"github.com/zeebo/xxh3"

	"github.com/google/uuid"
)

type RemovalTask struct {
	MarkerID  string
	RequestID string
}

var (
	// Queue for holding removal tasks that need to be retried
	retryQueue = make(chan RemovalTask, 100)

	// Context for managing the lifecycle of the background retry goroutine
	retryCtx, cancelRetryCtx = context.WithCancel(context.Background())
)

var WsRoomManager *RoomConnectionManager = NewRoomConnectionManager()

type ChulbongConn struct {
	Socket *websocket.Conn
	UserID string
	Mutex  *sync.Mutex
	Send   chan []byte
}
type RoomConnectionManager struct {
	connections       *csmap.CsMap[string, []*ChulbongConn] // roomid and users
	processedMessages *csmap.CsMap[string, struct{}]        // uid (struct{} does not occupy any space)
	mu                *sync.Mutex
}

// NewRoomConnectionManager initializes a ConnectionManager with a new haxmap instance
func NewRoomConnectionManager() *RoomConnectionManager {
	hasher := func(key string) uint64 {
		return xxh3.HashString(key)
	}

	manager := &RoomConnectionManager{
		connections: csmap.Create(
			csmap.WithShardCount[string, []*ChulbongConn](64),
			csmap.WithCustomHasher[string, []*ChulbongConn](hasher),
		),
		processedMessages: csmap.Create(
			csmap.WithShardCount[string, struct{}](64),
			csmap.WithCustomHasher[string, struct{}](hasher),
		),
		mu: &sync.Mutex{},
	}
	// Start the connection checker
	manager.StartConnectionChecker()
	manager.StartCleanUpProcessedMsg()

	go processRetryQueue(retryCtx)
	return manager
}

func (manager *RoomConnectionManager) BroadcastUserCountToRoom(roomID string) {
	userCount, err := GetUserCountInRoom(context.Background(), roomID)
	if err != nil {
		log.Printf("Error getting user count: %v", err)
		return
	}

	if userCount > 0 {
		message := fmt.Sprintf("%s (%d명 접속 중)", roomID, userCount)
		PublishMessageToAMQP(context.Background(), roomID, message, "chulbong-kr", "")
	}
	// manager.BroadcastMessageToRoom(roomID, message, "chulbong-kr", "")
}

// BroadcastMessageToRoom sends a WebSocket message to all users in a specific room
func (manager *RoomConnectionManager) BroadcastMessageToRoom2(roomID string, msgJSON []byte) {
	key := fmt.Sprintf("marker_%s", roomID)

	// markAsProcessed(broadcastMsg.UID)

	// Retrieve the slice of connections for the given roomID from the manager's connections map
	if conns, ok := manager.connections.Load(key); ok {
		// Iterate over the connections and send the message
		for _, conn := range conns {
			select {
			case conn.Send <- msgJSON:
				// Message enqueued to be sent by writePump goroutine
			default:
				// Handle full send channel, e.g., by logging or closing the connection
			}
		}
	}
}

// BroadcastMessageToRoom sends a WebSocket message to all users in a specific room
func (manager *RoomConnectionManager) BroadcastMessageToRoom(roomID, message, userNickname, userId string) {
	key := fmt.Sprintf("marker_%s", roomID)
	broadcastMsg := dto.BroadcastMessage{
		UID:          uuid.New().String(),
		Message:      message,
		UserID:       userId,
		UserNickname: userNickname,
		RoomID:       roomID,
		Timestamp:    time.Now().UnixMilli(),
	}

	// go PublishChatToRoom(roomID, msgJSON) // TODO: in distributed server

	// Serialize the message struct to JSON
	msgJSON, err := json.Marshal(broadcastMsg)
	if err != nil {
		log.Printf("Error marshalling message to JSON: %v", err)
		return
	}

	// markAsProcessed(broadcastMsg.UID)

	// Retrieve the slice of connections for the given roomID from the manager's connections map
	if conns, ok := manager.connections.Load(key); ok {
		// Iterate over the connections and send the message
		for _, conn := range conns {
			select {
			case conn.Send <- msgJSON:
				// Message enqueued to be sent by writePump goroutine
			default:
				// Handle full send channel, e.g., by logging or closing the connection
			}
		}
	}
}

// BroadcastMessage sends a WebSocket message to all users
func (manager *RoomConnectionManager) BroadcastMessage(message []byte, userId, roomId, userNickname string) {
	broadcastMsg := dto.BroadcastMessage{
		UID:          uuid.New().String(),
		Message:      string(message),
		UserID:       userId,
		UserNickname: userNickname,
		RoomID:       roomId,
		Timestamp:    time.Now().UnixMilli(),
	}
	// Serialize the message struct to JSON
	msgJSON, err := json.Marshal(broadcastMsg)
	if err != nil {
		log.Printf("Error marshalling message to JSON: %v", err)
		return
	}

	manager.connections.Range(func(key string, conns []*ChulbongConn) bool {
		for _, conn := range conns {
			conn.Mutex.Lock()
			if err := conn.Socket.WriteMessage(websocket.TextMessage, msgJSON); err != nil {
				log.Printf("Failed to send broadcast message to connection for %s: %v", key, err)
			}
			conn.Mutex.Unlock()
		}
		return true // Continue iteration
	})
}

// SubscribeToChatUpdate to subscribe to messages
func SubscribeToChatUpdate(markerID string) *redis.PubSub {
	key := fmt.Sprintf("room:%s:messages", markerID)

	pubsub := RedisStore.Conn().Subscribe(context.Background(), key)
	return pubsub
}

// PublishChatToRoom publishes a chat message to a specific room
func PublishChatToRoom(markerID string, message []byte) error {
	key := fmt.Sprintf("room:%s:messages", markerID)

	// Publish the serialized message to the Redis pub/sub system
	err := RedisStore.Conn().Publish(context.Background(), key, message).Err()
	if err != nil {
		return err
	}

	return nil
}
