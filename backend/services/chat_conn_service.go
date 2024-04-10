package services

import (
	"chulbong-kr/dto"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

// SaveConnection stores a WebSocket connection associated with a markerID in app memory
func (manager *RoomConnectionManager) SaveConnection(markerID, clientId string, conn *websocket.Conn) {
	// manager.mu.Lock()         // Lock at the start of the method
	// defer manager.mu.Unlock() // Unlock when the method returns

	key := fmt.Sprintf("marker_%s", markerID)

	newConn := &ChulbongConn{
		Socket:       conn,
		UserID:       clientId,
		Send:         make(chan []byte, 256), // Buffered channel
		InActiveChan: make(chan struct{}, 10),
		mu:           &sync.Mutex{},
	}
	newConn.UpdateLastSeen()

	// newConn.Socket.SetPongHandler(func(string) error {
	// 	newConn.LastSeen = time.Now()
	// 	return nil
	// })

	go newConn.writePump()

	// Check if there's already a connection list for this markerID
	// manager.connections.SetIfAbsent(key, []*ChulbongConn{newConn}) // csmap

	conns, loaded := manager.connections.LoadOrStore(key, []*ChulbongConn{newConn})
	if !loaded {
		return
	}

	// conns, ok := manager.connections.Load(key) // Doesn't have GetOrSet
	// if !ok {
	// 	return
	// }

	// LAVINMQ: Ensure only one subscription per room
	// if _, exists := ActiveSubscriptions.Load(key); !exists {
	// 	go SubscribeAndBroadcastFromAMQP(markerID)
	// 	ActiveSubscriptions.Store(key, struct{}{}) // marker_%s
	// }

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
	err = RedisStore.HSet(ctx, key, userID, jsonConnInfo).Err()
	if err != nil {
		return err
	}

	RedisStore.Expire(ctx, key, 1*time.Hour)
	return nil
}

func CheckDuplicateConnection(markerID, userID string) (bool, error) {
	ctx := context.Background()
	key := fmt.Sprintf("room:%s:connections", markerID)

	// Check if there's an entry for this userID
	exists, err := RedisStore.HExists(ctx, key, userID).Result()
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (manager *RoomConnectionManager) CheckDuplicateConnectionByLocal(markerID, userID string) (bool, error) {
	key := fmt.Sprintf("marker_%s", markerID)

	// Retrieve the current list of connections for the room.
	if conns, ok := manager.connections.Load(key); ok {
		// Find and remove the specified connection from the slice.
		for _, c := range conns {
			if c.UserID == userID {
				return true, nil
			}
		}
	}

	return false, nil
}

func (manager *RoomConnectionManager) UpdateLastPing(markerID string, conn *websocket.Conn) {

	key := fmt.Sprintf("marker_%s", markerID)

	// Retrieve the current list of connections for the room.
	if conns, ok := manager.connections.Load(key); ok {
		// Find and remove the specified connection from the slice.
		for _, c := range conns {
			if c.Socket == conn {
				c.UpdateLastSeen()
				return
			}
		}
	}
}

func (manager *RoomConnectionManager) StartConnectionChecker() {
	ticker := time.NewTicker(60 * time.Minute)
	go func() {
		for range ticker.C {
			manager.CheckConnections()
		}
	}()
}

// CheckConnections iterates over all connections and sends a ping message.
// Connections that fail to respond can be considered inactive and removed.
func (manager *RoomConnectionManager) CheckConnections() {
	// manager.mu.Lock()         // Lock at the start of the method
	// defer manager.mu.Unlock() // Unlock when the method returns

	gracePeriod := 10 * time.Minute
	now := time.Now()

	// Temporary storage for IDs to delete and update
	// toDelete := []string{}
	toDelete := map[string][]*ChulbongConn{}
	toUpdate := map[string][]*ChulbongConn{}

	manager.connections.Range(func(id string, conns []*ChulbongConn) bool {
		var activeConns []*ChulbongConn // Store active connections
		var deadConns []*ChulbongConn   // Store dead connections

		for _, conn := range conns {

			if now.Sub(conn.GetLastSeen()) <= gracePeriod {
				activeConns = append(activeConns, conn)
			} else {
				deadConns = append(deadConns, conn)
			}

		}

		// Decide whether to update or delete the marker based on active connections
		if len(activeConns) > 0 {
			toUpdate[id] = activeConns
		} else {
			toDelete[id] = deadConns
		}

		return true
	})

	// Apply updates and deletions
	for id, conns := range toUpdate {
		manager.connections.Store(id, conns)
	}

	for id, conns := range toDelete {
		manager.connections.Delete(id)
		for _, conn := range conns {
			select {
			case conn.InActiveChan <- struct{}{}:
				// Message enqueued to be sent by writePump goroutine
			default:
				// Handle full send channel, e.g., by logging or closing the connection
			}
		}
	}
}

func RemoveConnectionFromRedis(markerID, xRequestID string) {
	ctx := context.Background()
	key := fmt.Sprintf("room:%s:connections", markerID)

	// Use HDEL to remove the connection information, indexed by the UserID
	if err := RedisStore.HDel(ctx, key, xRequestID).Err(); err != nil {
		enqueueRemovalTask(markerID, xRequestID)
	}
}

// RemoveConnection removes a WebSocket connection associated with a id
func (manager *RoomConnectionManager) RemoveWsFromRoom(markerID, clientId string, conn *websocket.Conn) {
	// manager.mu.Lock()         // Lock at the start of the method
	// defer manager.mu.Unlock() // Unlock when the method returns

	markerIDStr := fmt.Sprintf("marker_%s", markerID)

	// Retrieve the current list of connections for the room.
	if conns, ok := manager.connections.Load(markerIDStr); ok {
		// Find and remove the specified connection from the slice.
		for i, c := range conns {
			if c.UserID == clientId {
				conns = append(conns[:i], conns[i+1:]...)
				break
			}
		}

		// If the slice is empty after removal, delete the entry from the map.
		if len(conns) == 0 {
			manager.connections.Delete(markerIDStr)
			// LAVINMQ:
			// ActiveSubscriptions.Delete(markerIDStr)
			// StopSubscriptionForRoom(markerIDStr)
		} else {
			// Otherwise, update the map with the modified slice.
			manager.connections.Store(markerIDStr, conns)
		}
	}
}

// KickUserFromRoom closes the connection for a user in a specified room.
func (manager *RoomConnectionManager) KickUserFromRoom(markerID, userID string) error {
	key := fmt.Sprintf("marker_%s", markerID)
	conns, ok := manager.connections.Load(key)
	if !ok {
		return errors.New("room not found")
	}

	// Iterate over the connections to find the user
	for i, conn := range conns {
		if conn.UserID == userID {
			// Close the connection
			close(conn.Send) // Signal the writePump goroutine to exit
			conn.Socket.Close()

			// Remove the connection from the slice
			if len(conns) == 0 {
				manager.connections.Delete(key)
			} else {
				manager.connections.Store(key, append(conns[:i], conns[i+1:]...))
			}
			return nil
		}
	}
	return errors.New("user not found in room")
}

func GetUserCountInRoom(ctx context.Context, markerID string) (int64, error) {
	key := fmt.Sprintf("room:%s:connections", markerID)
	return RedisStore.HLen(ctx, key).Result()
}

func (manager *RoomConnectionManager) GetUserCountInRoomByLocal(markerID string) (int, error) {
	markerIDStr := fmt.Sprintf("marker_%s", markerID)
	conns, ok := manager.connections.Load(markerIDStr)
	if !ok {
		return 0, nil
	}

	return len(conns), nil
}

// GetAllRedisConnectionsFromRoom retrieves all connection information for a given room.
func GetAllRedisConnectionsFromRoom(markerID string) ([]dto.ConnectionInfo, error) {
	ctx := context.Background()
	key := fmt.Sprintf("room:%s:connections", markerID)

	// Retrieve all field-value pairs from the hash
	results, err := RedisStore.HGetAll(ctx, key).Result()
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

func (conn *ChulbongConn) writePump() {
	for {
		select {
		case <-conn.InActiveChan:
			// conn.Socket.WriteJSON(fiber.Map{"error": "inactive connection"})
			// time.Sleep(500 * time.Millisecond)
			conn.Socket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "inactive"))
			// conn.Socket.Close()
			return
		case message, ok := <-conn.Send:
			if !ok {
				conn.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// conn.Socket.SetWriteDeadline(time.Now().Add(30 * time.Second))

			if err := conn.Socket.WriteMessage(websocket.TextMessage, message); err != nil {
				// log.Printf("📆 Error sending text message: %v", err)
				continue
			}
		}
	}
}

func (manager *RoomConnectionManager) BanUser(markerID, userID string, duration time.Duration) error {
	// First, kick the user from the room.
	err := manager.KickUserFromRoom(markerID, userID)
	if err != nil {
		log.Printf("Error kicking user %s from room %s: %v", userID, markerID, err)
	}

	// Then, set the ban in Redis.
	banKey := fmt.Sprintf("ban_%s_%s", markerID, userID)
	return RedisStore.Set(context.Background(), banKey, "banned", duration).Err()
}

func (manager *RoomConnectionManager) GetBanDetails(markerID, userID string) (banned bool, remainingTime time.Duration, err error) {
	banKey := fmt.Sprintf("ban_%s_%s", markerID, userID)
	exists, err := RedisStore.Exists(context.Background(), banKey).Result()
	if err != nil {
		return false, 0, err
	}
	if exists == 0 {
		return false, 0, nil
	}

	// Get the remaining TTL of the ban
	ttl, err := RedisStore.TTL(context.Background(), banKey).Result()
	if err != nil {
		return true, 0, err
	}

	return true, ttl, nil
}