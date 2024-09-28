package service

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gofiber/contrib/websocket"
	csmap "github.com/mhmtszr/concurrent-swiss-map"
	"github.com/puzpuzpuz/xsync/v3"
	"github.com/redis/rueidis"
	"github.com/zeebo/xxh3"

	sonic "github.com/bytedance/sonic"
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

	jsonBufferPool = sync.Pool{
		New: func() any {
			return &bytes.Buffer{}
		},
	}
)

type ChatService struct {
	DB               *sqlx.DB
	Redis            *RedisService
	WebSocketManager *RoomConnectionManager

	Logger *zap.Logger
}

func NewChatService(lifecycle fx.Lifecycle, db *sqlx.DB, redis *RedisService, manager *RoomConnectionManager, l *zap.Logger) *ChatService {
	service := &ChatService{
		DB:               db,
		Redis:            redis,
		WebSocketManager: manager,
		Logger:           l,
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
	Nickname     string
	Socket       *websocket.Conn
	Send         chan []byte
	InActiveChan chan struct{}
}

type RoomConnectionManager struct {
	// connections       *haxmap.Map[string, []*ChulbongConn] // roomid and users
	rooms             *xsync.MapOf[string, *xsync.MapOf[string, *ChulbongConn]]
	processedMessages *csmap.CsMap[string, struct{}] // uid (struct{} does not occupy any space)
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
		rooms: xsync.NewMapOf[string, *xsync.MapOf[string, *ChulbongConn]](),
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
		message := roomID + " (" + strconv.Itoa(userCount) + "명 접속 중)"
		// PublishMessageToAMQP(context.Background(), roomID, message, "chulbong-kr", "")

		// Broadcast the user count message
		s.BroadcastRawMessageToRoom(roomID, message)
	}
}

// TODO: LAVINMQ
// BroadcastMessageToRoom sends a WebSocket message to all users in a specific room LAVINMQ:
func (s *ChatService) BroadcastMessageToRoom2(roomID string, msgJSON []byte) {
	// key := fmt.Sprintf("marker_%s", roomID)

	// // markAsProcessed(broadcastMsg.UID)

	// // Retrieve the slice of connections for the given roomID from the manager's connections map
	// if conns, ok := s.WebSocketManager.connections.Load(key); ok {
	// 	// Iterate over the connections and send the message
	// 	for _, conn := range conns {
	// 		select {
	// 		case conn.Send <- msgJSON:
	// 			// Message enqueued to be sent by writePump goroutine
	// 		default:
	// 			// Handle full send channel, e.g., by logging or closing the connection
	// 		}
	// 	}
	// }
}

// BroadcastMessageToRoom sends a WebSocket message to all users in a specific room
func (s *ChatService) BroadcastMessageToRoom(markerID, message, senderNickname, senderUserID string) error {
	roomConns, ok := s.WebSocketManager.rooms.Load(markerID)
	if !ok {
		return nil // No connections in room
	}

	payload := createMessagePayload(markerID, message, senderNickname, senderUserID)
	roomConns.Range(func(clientID string, conn *ChulbongConn) bool {
		select {
		case conn.Send <- payload:
			// Message enqueued to be sent by writePump goroutine
			// s.Logger.Info("Broadcast message to room", zap.String("roomID", markerID), zap.String("clientID", clientID))
		default:
			// Handle full send channel if necessary
		}
		return true
	})
	return nil
}

func (s *ChatService) BroadcastRawMessageToRoom(markerID, message string) {
	roomConns, ok := s.WebSocketManager.rooms.Load(markerID)
	if !ok {
		return // No connections in room
	}

	payload := createMessagePayload(markerID, message, "chulbong-kr", "")
	roomConns.Range(func(clientID string, conn *ChulbongConn) bool {
		select {
		case conn.Send <- payload:
			// Message enqueued to be sent by writePump goroutine
			// s.Logger.Info("Broadcast user count to room", zap.String("roomID", markerID), zap.String("clientID", clientID))
		default:
			// Handle full send channel if necessary
		}
		return true
	})
}

// BroadcastMessage sends a WebSocket message to all users in all rooms
func (s *ChatService) BroadcastMessage(message []byte, userID, roomID, userNickname string) {
	broadcastMsg := dto.BroadcastMessage{
		UID:          uuid.New().String(),
		Message:      string(message),
		UserID:       userID,
		UserNickname: userNickname,
		RoomID:       roomID,
		Timestamp:    time.Now().UnixMilli(),
	}
	// Serialize the message struct to JSON
	msgJSON, err := sonic.ConfigFastest.Marshal(broadcastMsg)
	if err != nil {
		// s.Logger.Error("Error marshalling message to JSON", zap.Error(err))
		return
	}

	// Iterate over all rooms
	s.WebSocketManager.rooms.Range(func(_roomID string, roomConns *xsync.MapOf[string, *ChulbongConn]) bool {
		// if shouldBroadcastToRoom(_roomID) {
		// 	// Proceed with sending messages
		// }

		// Iterate over all connections in the room
		roomConns.Range(func(_clientID string, conn *ChulbongConn) bool {
			select {
			case conn.Send <- msgJSON:
				// Message enqueued to be sent by writePump goroutine
			default:
				// Handle full send channel, e.g., by logging or closing the connection
				// s.Logger.Warn("Send channel is full, message dropped", zap.String("clientID", _clientID))
			}
			return true // Continue iteration over connections
		})
		return true // Continue iteration over rooms
	})
}

// Atomic update of the LastSeen timestamp.
func (c *ChulbongConn) UpdateLastSeen() {
	atomic.StoreInt64(&c.LastSeen, time.Now().UnixNano())
}

// Atomic read of the LastSeen timestamp.
func (c *ChulbongConn) GetLastSeen() time.Time {
	return time.Unix(0, atomic.LoadInt64(&c.LastSeen))
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

func (s *ChatService) GetNickname(markerID, clientID string) (string, error) {
	if roomConns, ok := s.WebSocketManager.rooms.Load(markerID); ok {
		if conn, ok := roomConns.Load(clientID); ok {
			return conn.Nickname, nil
		}
	}
	return "", fmt.Errorf("connection not found")
}

func createMessagePayload(markerID, message, senderNickname, senderUserID string) []byte {
	buf := jsonBufferPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		jsonBufferPool.Put(buf)
	}()

	broadcastMsg := dto.BroadcastMessage{
		UID:          uuid.New().String(),
		Message:      message,
		UserID:       senderUserID,
		UserNickname: senderNickname,
		RoomID:       markerID,
		Timestamp:    time.Now().UnixMilli(),
	}

	encoder := sonic.ConfigFastest.NewEncoder(buf)
	if err := encoder.Encode(broadcastMsg); err != nil {
		// Handle error appropriately
		log.Printf("Error encoding broadcast message: %v", err)
		return nil
	}

	// Make a copy of the bytes before resetting the buffer
	payload := make([]byte, buf.Len())
	copy(payload, buf.Bytes())

	return payload
}
