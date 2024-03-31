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
	"github.com/google/uuid"
)

// SaveConnection stores a WebSocket connection associated with a markerID in app memory
func (manager *RoomConnectionManager) SaveConnection(markerID, clientId string, conn *websocket.Conn) {
	manager.mu.Lock()         // Lock at the start of the method
	defer manager.mu.Unlock() // Unlock when the method returns

	key := fmt.Sprintf("marker_%s", markerID)

	newConn := &ChulbongConn{
		Socket: conn,
		UserID: clientId,
		Mutex:  &sync.Mutex{},
		Send:   make(chan []byte, 256), // Buffered channel
	}

	go newConn.writePump()

	// Check if there's already a connection list for this markerID
	manager.connections.SetIfAbsent(key, []*ChulbongConn{newConn})

	conns, ok := manager.connections.Load(key) // Doesn't have GetOrSet
	if !ok {
		return
	}

	// Ensure only one subscription per room
	if _, exists := ActiveSubscriptions.Load(key); !exists {
		go SubscribeAndBroadcastFromAMQP(markerID)
		ActiveSubscriptions.Store(key, struct{}{}) // marker_%s
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

func (manager *RoomConnectionManager) StartConnectionChecker() {
	ticker := time.NewTicker(3 * time.Minute)
	go func() {
		for range ticker.C {
			manager.CheckConnections()
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
			ActiveSubscriptions.Delete(id)
			StopSubscriptionForRoom(id) // id is "marker_%s"
		} else {
			manager.connections.Store(id, activeConns)
		}
		return true
	})
}

func RemoveConnectionFromRedis(markerID, xRequestID string) {
	ctx := context.Background()
	key := fmt.Sprintf("room:%s:connections", markerID)

	// Use HDEL to remove the connection information, indexed by the UserID
	if err := RedisStore.Conn().HDel(ctx, key, xRequestID).Err(); err != nil {
		enqueueRemovalTask(markerID, xRequestID)
	}
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
			ActiveSubscriptions.Delete(markerIDStr)
			StopSubscriptionForRoom(markerIDStr)
		} else {
			// Otherwise, update the map with the modified slice.
			manager.connections.Store(markerIDStr, conns)
		}
	}
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

func (conn *ChulbongConn) writePump() {
	for {
		message, ok := <-conn.Send
		if !ok {
			// The channel was closed, indicating the connection should be closed
			conn.Socket.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		if err := conn.Socket.WriteMessage(websocket.TextMessage, message); err != nil {
			// Handle the error, possibly breaking out of the loop
			break
		}
	}
}
