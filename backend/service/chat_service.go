package service

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
	csmap "github.com/mhmtszr/concurrent-swiss-map"
	"github.com/puzpuzpuz/xsync/v3"
	"github.com/redis/rueidis"
	"github.com/zeebo/xxh3"

	"github.com/google/uuid"
)

type RemovalTask struct {
	MarkerID  string
	RequestID string
}

// const initialSize = 2 << 10

var (
	// Queue for holding removal tasks that need to be retried
	retryQueue = make(chan RemovalTask, 100)

	// Context for managing the lifecycle of the background retry goroutine
	retryCtx, _ = context.WithCancel(context.Background())
)

type ChatService struct {
	DB               *sqlx.DB
	Redis            *RedisService
	WebSocketManager *RoomConnectionManager
	Ticker           *time.Ticker
}

func NewChatService(lifecycle fx.Lifecycle, db *sqlx.DB, redis *RedisService, manager *RoomConnectionManager) *ChatService {
	service := &ChatService{
		DB:               db,
		Redis:            redis,
		WebSocketManager: manager,
		Ticker:           nil,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go service.processRetryQueue(retryCtx)
			return nil
		},
		OnStop: func(context.Context) error {
			return nil
		},
	})

	return service
}

type ChulbongConn struct {
	LastSeen     int64
	UserID       string
	Socket       *websocket.Conn
	Send         chan []byte
	InActiveChan chan struct{}
	// mu           sync.Mutex
}
type RoomConnectionManager struct {
	// connections       *haxmap.Map[string, []*ChulbongConn] // roomid and users
	connections       *xsync.MapOf[string, []*ChulbongConn] // roomid and users
	processedMessages *csmap.CsMap[string, struct{}]        // uid (struct{} does not occupy any space)
	// mu                sync.Mutex
}

func CustomXXH3Hasher(s string) uintptr {
	return uintptr(xxh3.HashString(s))
}

// NewRoomConnectionManager initializes a ConnectionManager with a new haxmap instance
func NewRoomConnectionManager(lifecycle fx.Lifecycle) *RoomConnectionManager {
	hasher := func(key string) uint64 {
		return xxh3.HashString(key)
	}

	manager := &RoomConnectionManager{
		connections: xsync.NewMapOf[string, []*ChulbongConn](),
		processedMessages: csmap.Create(
			csmap.WithShardCount[string, struct{}](64),
			csmap.WithCustomHasher[string, struct{}](hasher),
		),
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				ticker := time.NewTicker(30 * time.Minute)
				defer ticker.Stop()

				for range ticker.C {
					manager.StartConnectionChecker()
				}
			}()
			return nil
		},
		OnStop: func(context.Context) error {
			return nil
		},
	})

	return manager
}

func (s *ChatService) BroadcastUserCountToRoom(roomID string) {
	userCount, err := s.GetUserCountInRoom(context.Background(), roomID)
	if err != nil {
		log.Printf("Error getting user count: %v", err)
		return
	}

	// LAVINMQ:
	if userCount > 0 {
		message := fmt.Sprintf("%s (%d명 접속 중)", roomID, userCount)
		// PublishMessageToAMQP(context.Background(), roomID, message, "chulbong-kr", "")
		s.BroadcastMessageToRoom(roomID, message, "chulbong-kr", "")
	}
}

func (s *ChatService) BroadcastUserCountToRoomByLocal(roomID string) {
	userCount, err := s.GetUserCountInRoomByLocal(roomID)
	if err != nil {
		log.Printf("Error getting user count: %v", err)
		return
	}

	// LAVINMQ:
	if userCount > 0 {
		message := fmt.Sprintf("%s (%d명 접속 중)", roomID, userCount)
		// PublishMessageToAMQP(context.Background(), roomID, message, "chulbong-kr", "")
		s.BroadcastMessageToRoom(roomID, message, "chulbong-kr", "")
	}
}

// BroadcastMessageToRoom sends a WebSocket message to all users in a specific room LAVINMQ:
func (s *ChatService) BroadcastMessageToRoom2(roomID string, msgJSON []byte) {
	key := fmt.Sprintf("marker_%s", roomID)

	// markAsProcessed(broadcastMsg.UID)

	// Retrieve the slice of connections for the given roomID from the manager's connections map
	if conns, ok := s.WebSocketManager.connections.Load(key); ok {
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
func (s *ChatService) BroadcastMessageToRoom(roomID, message, userNickname, userID string) error {
	key := fmt.Sprintf("marker_%s", roomID)
	broadcastMsg := dto.BroadcastMessage{
		UID:          uuid.New().String(),
		Message:      message,
		UserID:       userID,
		UserNickname: userNickname,
		RoomID:       roomID,
		Timestamp:    time.Now().UnixMilli(),
	}

	// Serialize the message struct to JSON
	msgJSON, err := json.Marshal(broadcastMsg)
	if err != nil {
		log.Printf("Error marshalling message to JSON: %v", err)
		return err
	}

	// markAsProcessed(broadcastMsg.UID)

	// Retrieve the slice of connections for the given roomID from the manager's connections map
	if conns, ok := s.WebSocketManager.connections.Load(key); ok {
		// Iterate over the connections and send the message
		for _, conn := range conns {
			select {
			case conn.Send <- msgJSON:
				// Message enqueued to be sent by writePump goroutine
			default:
				// full send channel
				return fmt.Errorf("send channel for connection in room %s is full", roomID)
			}
		}
	}

	return nil
}

// BroadcastMessage sends a WebSocket message to all users
func (s *ChatService) BroadcastMessage(message []byte, userID, roomId, userNickname string) {
	broadcastMsg := dto.BroadcastMessage{
		UID:          uuid.New().String(),
		Message:      string(message),
		UserID:       userID,
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

	s.WebSocketManager.connections.Range(func(_ string, conns []*ChulbongConn) bool {
		// Iterate over the connections and send the message
		for _, conn := range conns {
			select {
			case conn.Send <- msgJSON:
				// Message enqueued to be sent by writePump goroutine
			default:
				// Handle full send channel, e.g., by logging or closing the connection
			}
		}
		return true // Continue iteration
	})
}

// Atomic update of the LastSeen timestamp.
func (c *ChulbongConn) UpdateLastSeen() {
	now := time.Now().UnixNano() // Get current time as Unix nano timestamp
	atomic.StoreInt64(&c.LastSeen, now)
}

// Atomic read of the LastSeen timestamp.
func (c *ChulbongConn) GetLastSeen() time.Time {
	unixNano := atomic.LoadInt64(&c.LastSeen)
	return time.Unix(0, unixNano)
}

// SubscribeToChatUpdate to subscribe to messages
func (s *ChatService) SubscribeToChatUpdate(markerID string) {
	key := fmt.Sprintf("room:%s:messages", markerID)

	// Using a dedicated connection for subscription
	dedicatedClient, cancel := s.Redis.Core.Client.Dedicate()
	defer cancel() // Ensure resources are cleaned up properly

	// Set up pub/sub hooks
	wait := dedicatedClient.SetPubSubHooks(rueidis.PubSubHooks{
		OnMessage: func(m rueidis.PubSubMessage) {
			// Handle incoming messages
			fmt.Printf("Received message: %s\n", m.Message)
		},
	})

	// Subscribe to the channel
	dedicatedClient.Do(context.Background(), dedicatedClient.B().Subscribe().Channel(key).Build())

	// Use the wait channel to handle disconnection
	err := <-wait // will receive a value if the client disconnects
	if err != nil {
		fmt.Printf("Disconnected due to: %v\n", err)
	}
}

// PublishChatToRoom publishes a chat message to a specific room
func (s *ChatService) PublishChatToRoom(markerID string, message []byte) error {
	key := fmt.Sprintf("room:%s:messages", markerID)

	// Publish the serialized message to the Redis pub/sub system
	err := s.Redis.Core.Client.Do(context.Background(), s.Redis.Core.Client.B().Publish().Channel(key).Message(rueidis.BinaryString(message)).Build()).Error()
	return err
}
