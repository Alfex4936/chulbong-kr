package services

import (
	"chulbong-kr/dto"
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/alphadose/haxmap"
	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
	"github.com/redis/go-redis/v9"

	"github.com/google/uuid"
)

var WsRoomManager *RoomConnectionManager = NewRoomConnectionManager()

type ChulbongConn struct {
	Socket *websocket.Conn
	UserID string
}
type RoomConnectionManager struct {
	connections       *haxmap.Map[string, []*ChulbongConn] // roomid and users
	processedMessages *haxmap.Map[string, struct{}]        // uid (struct{} does not occupy any space)
}

// NewRoomConnectionManager initializes a ConnectionManager with a new haxmap instance
func NewRoomConnectionManager() *RoomConnectionManager {
	manager := &RoomConnectionManager{
		connections:       haxmap.New[string, []*ChulbongConn](),
		processedMessages: haxmap.New[string, struct{}](),
	}
	// Start the connection checker
	manager.StartConnectionChecker()
	manager.StartCleanUpProcessedMsg()
	return manager
}

func (manager *RoomConnectionManager) StartConnectionChecker() {
	ticker := time.NewTicker(30 * time.Second)
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
			manager.processedMessages = haxmap.New[string, struct{}]()
		}
	}()
}

// CheckConnections iterates over all connections and sends a ping message.
// Connections that fail to respond can be considered inactive and removed.
func (manager *RoomConnectionManager) CheckConnections() {
	manager.connections.ForEach(func(id string, conns []*ChulbongConn) bool {
		var inactiveConns []int // Store indexes of inactive connections
		for i, conn := range conns {
			conn.Socket.SetWriteDeadline(time.Now().Add(30 * time.Second)) // Adjust deadline as needed
			if err := conn.Socket.WriteMessage(websocket.PingMessage, nil); err != nil {
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
			manager.connections.Del(id)
		} else {
			manager.connections.Set(id, conns) // Update the slice in the map
		}
		return true
	})

}

// SaveConnection stores a WebSocket connection associated with a markerID in app memory
func (manager *RoomConnectionManager) SaveConnection(markerID, clientId string, conn *websocket.Conn) {
	key := fmt.Sprintf("marker_%s", markerID)

	// First, try to get the current list of connections for the room.
	conns, ok := manager.connections.Get(key)

	if !ok {
		// If it doesn't exist, initialize it with the new connection.
		conns = []*ChulbongConn{
			{
				Socket: conn,
				UserID: clientId,
			},
		}
	} else {
		// If it exists, append the new connection, ensuring not to duplicate.
		if !contains(conns, conn) {
			conns = append(conns, &ChulbongConn{
				Socket: conn,
				UserID: clientId,
			})
		}
	}

	// Update the map with the new or modified slice.
	manager.connections.Set(key, conns)

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

	MarkAsProcessed(broadcastMsg.UID)
	// go PublishChatToRoom(roomID, msgJSON) // TODO: in distributed server

	// Retrieve the slice of connections for the given roomID from the manager's connections map
	if conns, ok := manager.connections.Get(key); ok {
		// Iterate over the connections and send the message
		for _, conn := range conns {
			if conn.UserID == userId {
				broadcastMsg.IsOwner = true
			}
			// Serialize the message struct to JSON
			msgJSON, err := json.Marshal(broadcastMsg)
			if err != nil {
				log.Printf("Error marshalling message to JSON: %v", err)
				continue
			}

			if err := conn.Socket.WriteMessage(websocket.TextMessage, msgJSON); err != nil {
				log.Printf("Failed to send message to connection in room %s: %v", roomID, err)
			}
		}
	} else {
		log.Printf("No connections found for room %s", roomID)
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

	manager.connections.ForEach(func(key string, conns []*ChulbongConn) bool {
		for _, conn := range conns {
			if err := conn.Socket.WriteMessage(websocket.TextMessage, msgJSON); err != nil {
				log.Printf("Failed to send broadcast message to connection for %s: %v", key, err)
			}
		}
		return true // Continue iteration
	})
}

// RemoveConnection removes a WebSocket connection associated with a id
func (manager *RoomConnectionManager) RemoveWsFromRoom(markerID string, conn *websocket.Conn) {
	markerIDStr := fmt.Sprintf("marker_%s", markerID)

	// Retrieve the current list of connections for the room.
	if conns, ok := manager.connections.Get(markerIDStr); ok {
		// Find and remove the specified connection from the slice.
		for i, c := range conns {
			if c.Socket == conn {
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

	MarkAsProcessed(broadcastMsg.UID) // Mark the message as processed locally

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
	// 25 adj
	adjectives := []string{
		"귀여운",   // Cute
		"멋진",    // Cool
		"착한",    // Kind
		"용감한",   // Brave
		"영리한",   // Clever
		"재미있는",  // Fun
		"행복한",   // Happy
		"사랑스러운", // Lovely
		"기운찬",   // Energetic
		"빛나는",   // Shining
		"평화로운",  // Peaceful
		"신비로운",  // Mysterious
		"자유로운",  // Free
		"매력적인",  // Charming
		"섬세한",   // Delicate
		"우아한",   // Elegant
		"활발한",   // Lively
		"강인한",   // Strong
		"독특한",   // Unique
		"무서운",   // Scary
		"꿈꾸는",   // Dreamy
		"느긋한",   // Relaxed
		"열정적인",  // Passionate
		"소중한",   // Precious
		"신선한",   // Fresh
	}

	// 9 names
	names := []string{
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
	return fmt.Sprintf("%s%s-%s", adjective, name, shortUID)
}

// contains function to check if the slice already contains the given connection.
func contains(slice []*ChulbongConn, conn *websocket.Conn) bool {
	for _, item := range slice {
		if item.Socket == conn {
			return true
		}
	}
	return false
}

func retrieveConnID(markerID string, userID int) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("room:%s:connections", markerID)

	// Get the connection information stored under the UserID key
	jsonConnInfo, err := RedisStore.Conn().HGet(ctx, key, strconv.Itoa(userID)).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("no connection info found for user %d in room %s", userID, markerID)
		}
		return "", err
	}

	var connInfo dto.ConnectionInfo
	err = json.Unmarshal([]byte(jsonConnInfo), &connInfo)
	if err != nil {
		return "", err
	}

	return connInfo.ConnID, nil
}

func hasProcessed(uid string) bool {
	_, exists := WsRoomManager.processedMessages.Get(uid)
	return exists
}

func MarkAsProcessed(uid string) {
	WsRoomManager.processedMessages.Set(uid, struct{}{})
}
