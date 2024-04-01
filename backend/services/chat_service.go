package services

import (
	"chulbong-kr/dto"
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
	csmap "github.com/mhmtszr/concurrent-swiss-map"
	"github.com/redis/go-redis/v9"
	"github.com/zeebo/xxh3"

	"github.com/google/uuid"
)

var adjectives = []string{
	"귀여운",     // Cute
	"멋진",      // Cool
	"착한",      // Kind
	"용감한",     // Brave
	"영리한",     // Clever
	"재미있는",    // Fun
	"행복한",     // Happy
	"사랑스러운",   // Lovely
	"기운찬",     // Energetic
	"빛나는",     // Shining
	"평화로운",    // Peaceful
	"신비로운",    // Mysterious
	"자유로운",    // Free
	"매력적인",    // Charming
	"섬세한",     // Delicate
	"우아한",     // Elegant
	"활발한",     // Lively
	"강인한",     // Strong
	"독특한",     // Unique
	"무서운",     // Scary
	"꿈꾸는",     // Dreamy
	"느긋한",     // Relaxed
	"열정적인",    // Passionate
	"소중한",     // Precious
	"신선한",     // Fresh
	"창의적인",    // Creative
	"우수한",     // Excellent
	"재치있는",    // Witty
	"감각적인",    // Sensual
	"흥미로운",    // Interesting
	"유명한",     // Famous
	"현명한",     // Wise
	"대담한",     // Bold
	"침착한",     // Calm
	"신속한",     // Swift
	"화려한",     // Gorgeous
	"정열적인",    // Passionate (Alternate)
	"끈기있는",    // Persistent
	"애정이 깊은",  // Affectionate
	"민첩한",     // Agile
	"빠른",      // Quick
	"조용한",     // Quiet
	"명랑한",     // Cheerful
	"정직한",     // Honest
	"용서하는",    // Forgiving
	"용기있는",    // Courageous
	"성실한",     // Sincere
	"호기심이 많은", // Curious
	"겸손한",     // Humble
	"관대한",     // Generous
}

// 9 names
var names = []string{
	"라이언", // Ryan
	"어피치", // Apeach
	"콘",   // Con
	"무지",  // Muzi
	"네오",  // Neo
	"프로도", // Frodo
	"제이지", // Jay-G
	"튜브",  // Tube
	"철봉",  // chulbong
}

var WsRoomManager *RoomConnectionManager = NewRoomConnectionManager()

type ChulbongConn struct {
	Socket *websocket.Conn
	UserID string
	Mutex  *sync.Mutex
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
	return manager
}

func (manager *RoomConnectionManager) StartConnectionChecker() {
	ticker := time.NewTicker(3 * time.Minute)
	go func() {
		for range ticker.C {
			manager.CheckConnections()
		}
	}()
}

// Clean every hour
func (manager *RoomConnectionManager) StartCleanUpProcessedMsg() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			manager.processedMessages.Clear()
		}
	}()
}

// CheckConnections iterates over all connections and sends a ping message.
// Connections that fail to respond can be considered inactive and removed.
func (manager *RoomConnectionManager) CheckConnections() {
	manager.connections.Range(func(id string, conns []*ChulbongConn) bool {
		var activeConns []*ChulbongConn // Store active connections

		for _, conn := range conns {
			conn.Mutex.Lock()
			conn.Socket.SetWriteDeadline(time.Now().Add(30 * time.Second))
			if err := conn.Socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				// If ping fails, log and skip adding this connection to activeConns
				// log.Printf("Ping failed for connection: %v. Marking as inactive.", err)
				continue
			}
			conn.Mutex.Unlock()
			// If ping succeeds, add connection to activeConns
			activeConns = append(activeConns, conn)
		}

		// Update the map only if there are active connections left
		if len(activeConns) == 0 && id != "suwon" {
			manager.connections.Delete(id)
		} else {
			manager.connections.Store(id, activeConns)
		}
		return true
	})
}

// SaveConnection stores a WebSocket connection associated with a markerID in app memory
func (manager *RoomConnectionManager) SaveConnection(markerID, clientId string, conn *websocket.Conn) {
	manager.mu.Lock()         // Lock at the start of the method
	defer manager.mu.Unlock() // Unlock when the method returns

	key := fmt.Sprintf("marker_%s", markerID)

	newConn := &ChulbongConn{
		Socket: conn,
		UserID: clientId,
		Mutex:  &sync.Mutex{},
	}

	// Check if there's already a connection list for this markerID
	manager.connections.SetIfAbsent(key, []*ChulbongConn{newConn})

	conns, ok := manager.connections.Load(key) // Doesn't have GetOrSet
	if !ok {
		return
	}

	// If we reach here, it means the list existed, so we must check for duplicates and append if necessary
	for _, item := range conns {
		if item.UserID == clientId {
			// Connection for clientId already exists, avoid adding a duplicate
			return
		}
	}

	updatedConns := append(conns, newConn)

	// Update the map with the new or modified slice
	manager.connections.Store(key, updatedConns)
}

func (manager *RoomConnectionManager) BroadcastUserCountToRoom(roomID string) {
	userCount, err := GetUserCountInRoom(context.Background(), roomID)
	if err != nil {
		log.Printf("Error getting user count: %v", err)
		return
	}

	message := fmt.Sprintf("%s (%d명 접속 중)", roomID, userCount)
	// Your existing method to broadcast a message to all users in the room
	manager.BroadcastMessageToRoom(roomID, message, "chulbong-kr", "")
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

	markAsProcessed(broadcastMsg.UID)

	// Retrieve the slice of connections for the given roomID from the manager's connections map
	if conns, ok := manager.connections.Load(key); ok {
		// Iterate over the connections and send the message
		for _, conn := range conns {

			conn.Mutex.Lock()
			if err := conn.Socket.WriteMessage(websocket.TextMessage, msgJSON); err != nil {
				// log.Printf("Failed to send message to connection in room %s: %v", roomID, err)
				continue
			}
			conn.Mutex.Unlock()
		}
	}
	//  else {
	// 	log.Printf("No connections found for room %s", roomID)
	// }
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

// RemoveConnection removes a WebSocket connection associated with a id
func (manager *RoomConnectionManager) RemoveWsFromRoom(markerID string, conn *websocket.Conn) {
	markerIDStr := fmt.Sprintf("marker_%s", markerID)

	// Retrieve the current list of connections for the room.
	if conns, ok := manager.connections.Load(markerIDStr); ok {
		// Find and remove the specified connection from the slice.
		for i, c := range conns {
			if c.Socket == conn {
				conns = append(conns[:i], conns[i+1:]...)
				break
			}
		}

		// If the slice is empty after removal, delete the entry from the map.
		if len(conns) == 0 {
			manager.connections.Delete(markerIDStr)
		} else {
			// Otherwise, update the map with the modified slice.
			manager.connections.Store(markerIDStr, conns)
		}
	}
}

// REDIS
// To see which users are in a room easily (in case distributed servers)
func AddConnectionRoomToRedis(markerID, userID, username string) error {
	ctx := context.Background()
	key := fmt.Sprintf("room:%s:connections", markerID)
	connID := uuid.New().String() // unique identifier for the connection

	connInfo := dto.ConnectionInfo{
		UserID:   userID,
		RoomID:   markerID,
		Username: username,
		ConnID:   connID,
	}

	jsonConnInfo, err := json.Marshal(connInfo)
	if err != nil {
		return err
	}

	// Use HSET to store the connection information, indexed by the userID
	err = RedisStore.Conn().HSet(ctx, key, userID, jsonConnInfo).Err()
	if err != nil {
		return err
	}

	RedisStore.Conn().Expire(ctx, key, 1*time.Hour)
	return nil
}

func CheckDuplicateConnection(markerID, userID string) (bool, error) {
	ctx := context.Background()
	key := fmt.Sprintf("room:%s:connections", markerID)

	// Check if there's an entry for this userID
	exists, err := RedisStore.Conn().HExists(ctx, key, userID).Result()
	if err != nil {
		return false, err
	}

	return exists, nil
}

func RemoveConnectionFromRedis(markerID, xRequestID string) error {
	ctx := context.Background()
	key := fmt.Sprintf("room:%s:connections", markerID)

	// Use HDEL to remove the connection information, indexed by the UserID
	return RedisStore.Conn().HDel(ctx, key, xRequestID).Err()
}

func GetUserCountInRoom(ctx context.Context, markerID string) (int64, error) {
	key := fmt.Sprintf("room:%s:connections", markerID)
	return RedisStore.Conn().HLen(ctx, key).Result()
}

// GetAllRedisConnectionsFromRoom retrieves all connection information for a given room.
func GetAllRedisConnectionsFromRoom(markerID string) ([]dto.ConnectionInfo, error) {
	ctx := context.Background()
	key := fmt.Sprintf("room:%s:connections", markerID)

	// Retrieve all field-value pairs from the hash
	results, err := RedisStore.Conn().HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	// Deserialize the connection information
	var connections []dto.ConnectionInfo
	for _, jsonConnInfo := range results {
		var connInfo dto.ConnectionInfo
		err := json.Unmarshal([]byte(jsonConnInfo), &connInfo)
		if err != nil {
			// Handle the error as appropriate: log it, skip it, etc.
			log.Printf("Error unmarshaling connection info: %v", err)
			continue
		}
		connections = append(connections, connInfo)
	}

	return connections, nil
}

// SubscribeToChatUpdate to subscribe to messages
func SubscribeToChatUpdate(markerID string) *redis.PubSub {
	key := fmt.Sprintf("room:%s:messages", markerID)

	pubsub := RedisStore.Conn().Subscribe(context.Background(), key)
	return pubsub
}

func ProcessMessageFromSubscription(msg []byte) {
	var broadcastMsg dto.BroadcastMessage
	err := json.Unmarshal(msg, &broadcastMsg)
	if err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		return
	}

	if hasProcessed(broadcastMsg.UID) {
		return // Skip processing if we've already handled this message
	}

	markAsProcessed(broadcastMsg.UID) // Mark the message as processed locally

	// then broadcast
	WsRoomManager.BroadcastMessageToRoom(broadcastMsg.RoomID, broadcastMsg.Message, broadcastMsg.UserNickname, broadcastMsg.UserID)
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

func GenerateKoreanNickname() string {

	// Select a random
	adjective := adjectives[rand.Intn(len(adjectives))]

	name := names[rand.Intn(len(names))]

	// Generate a unique identifier
	uid := uuid.New().String()

	// Use the first 8 characters of the UUID to keep it short
	shortUID := uid[:8]

	// possibilities for conflict
	// highly unlikely.
	// 25 * 9 * 16^8 (UUID first 8 characters)
	// UUID can conflict by root(16*8) = 65,536
	return fmt.Sprintf("%s %s [%s]", adjective, name, shortUID)
}

func hasProcessed(uid string) bool {
	_, exists := WsRoomManager.processedMessages.Load(uid)
	return exists
}

func markAsProcessed(uid string) {
	WsRoomManager.processedMessages.Store(uid, struct{}{})
}
