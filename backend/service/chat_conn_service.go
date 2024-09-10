package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/Alfex4936/chulbong-kr/util"

	sonic "github.com/bytedance/sonic"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/redis/rueidis"
)

// SaveConnection stores a WebSocket connection associated with a markerID in app memory
func (s *ChatService) SaveConnection(markerID, clientID string, conn *websocket.Conn) {
	// s.WebSocketManager.mu.Lock()         // Lock at the start of the method
	// defer s.WebSocketManager.mu.Unlock() // Unlock when the method returns

	key := fmt.Sprintf("marker_%s", markerID)

	newConn := &ChulbongConn{
		Socket:       conn,
		UserID:       clientID,
		Send:         make(chan []byte, 10), // Buffered channel
		InActiveChan: make(chan struct{}),
	}
	newConn.UpdateLastSeen()

	// newConn.Socket.SetPongHandler(func(string) error {
	// 	newConn.LastSeen = time.Now()
	// 	return nil
	// })

	go newConn.writePump()

	// Check if there's already a connection list for this markerID
	// s.WebSocketManager.connections.SetIfAbsent(key, []*ChulbongConn{newConn}) // csmap

	conns, loaded := s.WebSocketManager.connections.LoadOrStore(key, []*ChulbongConn{newConn})
	if !loaded {
		return
	}

	// conns, ok := s.WebSocketManager.connections.Load(key) // Doesn't have GetOrSet
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
		if item.UserID == clientID {
			// Connection for clientID already exists, avoid adding a duplicate
			return
		}
	}

	updatedConns := append(conns, newConn)

	// Update the map with the new or modified slice
	s.WebSocketManager.connections.Store(key, updatedConns)
}

// REDIS
// To see which users are in a room easily (in case distributed servers)
func (s *ChatService) AddConnectionRoomToRedis(markerID, userID, username string) error {
	ctx := context.Background()
	key := fmt.Sprintf("room:%s:connections", markerID)
	connID := uuid.New().String() // unique identifier for the connection

	connInfo := dto.ConnectionInfo{
		UserID:   userID,
		RoomID:   markerID,
		Username: username,
		ConnID:   connID,
	}

	jsonConnInfo, err := sonic.Marshal(connInfo)
	if err != nil {
		return err
	}

	// Use HSET to set the field value
	setCmd := s.Redis.Core.Client.B().Hset().
		Key(key).
		FieldValue().
		FieldValue(userID, rueidis.BinaryString(jsonConnInfo)).
		Build()

	// Execute the HSET command
	if err := s.Redis.Core.Client.Do(ctx, setCmd).Error(); err != nil {
		return err
	}

	// Set expiration on the whole hash key
	expireCmd := s.Redis.Core.Client.B().Expire().Key(key).Seconds(int64(time.Hour / time.Second)).Build()
	if err := s.Redis.Core.Client.Do(ctx, expireCmd).Error(); err != nil {
		return err
	}
	return nil
}

func (s *ChatService) CheckDuplicateConnection(markerID, userID string) (bool, error) {
	ctx := context.Background()
	key := fmt.Sprintf("room:%s:connections", markerID)

	// Build the HEXISTS command
	cmd := s.Redis.Core.Client.B().Hexists().Key(key).Field(userID).Build()

	// Execute the HEXISTS command
	exists, err := s.Redis.Core.Client.Do(ctx, cmd).AsBool()
	if err != nil {
		return false, err // Proper error handling
	}

	return exists, nil
}

func (s *ChatService) CheckDuplicateConnectionByLocal(markerID, userID string) (bool, error) {
	key := fmt.Sprintf("marker_%s", markerID)

	// Retrieve the current list of connections for the room.
	if conns, ok := s.WebSocketManager.connections.Load(key); ok {
		// Find and remove the specified connection from the slice.
		for _, c := range conns {
			if c.UserID == userID {
				return true, nil
			}
		}
	}

	return false, nil
}

func (s *ChatService) UpdateLastPing(markerID string, conn *websocket.Conn) {

	key := fmt.Sprintf("marker_%s", markerID)

	// Retrieve the current list of connections for the room.
	if conns, ok := s.WebSocketManager.connections.Load(key); ok {
		// Find and remove the specified connection from the slice.
		for _, c := range conns {
			if c.Socket == conn {
				c.UpdateLastSeen()
				return
			}
		}
	}
}

func (s *RoomConnectionManager) StartConnectionChecker() {
	ticker := time.NewTicker(30 * time.Minute)
	go func() {
		for range ticker.C {
			s.CheckConnections()
		}
	}()
}

// CheckConnections iterates over all connections and sends a ping message.
// Connections that fail to respond can be considered inactive and removed.
func (s *RoomConnectionManager) CheckConnections() {
	// s.WebSocketManager.mu.Lock()         // Lock at the start of the method
	// defer s.WebSocketManager.mu.Unlock() // Unlock when the method returns

	gracePeriod := 10 * time.Minute
	now := time.Now()

	// It is safe to modify the map while iterating it
	s.connections.Range(func(id string, conns []*ChulbongConn) bool {
		var activeConns []*ChulbongConn // Store active connections

		for _, conn := range conns {

			if now.Sub(conn.GetLastSeen()) > gracePeriod {
				select {
				case conn.InActiveChan <- struct{}{}:
					// Message enqueued to be sent by writePump goroutine
				default:
					// Handle full send channel, e.g., by logging or closing the connection
				}
			} else {
				activeConns = append(activeConns, conn)
			}

		}

		// Decide whether to update or delete the marker based on active connections
		if len(activeConns) > 0 {
			s.connections.Store(id, activeConns)
		} else {
			s.connections.Delete(id)
		}
		return true
	})
}

func (s *ChatService) RemoveConnectionFromRedis(markerID, xRequestID string) {
	ctx := context.Background()
	key := fmt.Sprintf("room:%s:connections", markerID)

	// Use HDEL to remove the connection information, indexed by the UserID
	cmd := s.Redis.Core.Client.B().Hdel().Key(key).Field(xRequestID).Build()
	if err := s.Redis.Core.Client.Do(ctx, cmd).Error(); err != nil {
		enqueueRemovalTask(markerID, xRequestID)
	}
}

// RemoveConnection removes a WebSocket connection associated with a id
func (s *ChatService) RemoveWsFromRoom(markerID, clientID string) {
	// s.WebSocketManager.mu.Lock()         // Lock at the start of the method
	// defer s.WebSocketManager.mu.Unlock() // Unlock when the method returns

	markerIDStr := fmt.Sprintf("marker_%s", markerID)

	// Retrieve the current list of connections for the room.
	if conns, ok := s.WebSocketManager.connections.Load(markerIDStr); ok {
		// Find and remove the specified connection from the slice.
		for i, c := range conns {
			if c.UserID == clientID {
				close(c.Send)
				close(c.InActiveChan)

				conns = append(conns[:i], conns[i+1:]...)
				break
			}
		}

		// If the slice is empty after removal, delete the entry from the map.
		if len(conns) == 0 {
			s.WebSocketManager.connections.Delete(markerIDStr)
			// LAVINMQ:
			// ActiveSubscriptions.Delete(markerIDStr)
			// StopSubscriptionForRoom(markerIDStr)
		} else {
			// Otherwise, update the map with the modified slice.
			s.WebSocketManager.connections.Store(markerIDStr, conns)
		}
	}
}

// KickUserFromRoom closes the connection for a user in a specified room.
func (s *ChatService) KickUserFromRoom(markerID, userID string) error {
	key := fmt.Sprintf("marker_%s", markerID)
	conns, ok := s.WebSocketManager.connections.Load(key)
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
				s.WebSocketManager.connections.Delete(key)
			} else {
				s.WebSocketManager.connections.Store(key, append(conns[:i], conns[i+1:]...))
			}
			return nil
		}
	}
	return errors.New("user not found in room")
}

func (s *ChatService) GetUserCountInRoom(ctx context.Context, markerID string) (int64, error) {
	key := fmt.Sprintf("room:%s:connections", markerID)

	// Use HLEN to get the number of connections in the room
	cmd := s.Redis.Core.Client.B().Hlen().Key(key).Build()
	return s.Redis.Core.Client.Do(ctx, cmd).AsInt64()
}

func (s *ChatService) GetUserCountInRoomByLocal(markerID string) (int, error) {
	markerIDStr := fmt.Sprintf("marker_%s", markerID)
	conns, ok := s.WebSocketManager.connections.Load(markerIDStr)
	if !ok {
		return 0, nil
	}

	return len(conns), nil
}

// GetAllRedisConnectionsFromRoom retrieves all connection information for a given room.
func (s *ChatService) GetAllRedisConnectionsFromRoom(markerID string) ([]dto.ConnectionInfo, error) {
	ctx := context.Background()
	key := fmt.Sprintf("room:%s:connections", markerID)

	// Use HGETALL to retrieve all field-value pairs from the hash
	cmd := s.Redis.Core.Client.B().Hgetall().Key(key).Build()

	// Execute the command and get the result as a string map
	result, err := s.Redis.Core.Client.Do(ctx, cmd).AsStrMap()
	if err != nil {
		return nil, err
	}

	// Deserialize the connection information from JSON stored in each field
	connections := make([]dto.ConnectionInfo, 0, len(result))
	for _, jsonConnInfo := range result {
		var connInfo dto.ConnectionInfo
		// Use StringToBytes to avoid unnecessary memory allocation
		jsonBytes := util.StringToBytes(jsonConnInfo)

		if err := sonic.Unmarshal(jsonBytes, &connInfo); err != nil {
			log.Printf("Error unmarshaling connection info: %v", err)
			continue
		}
		connections = append(connections, connInfo)
	}

	return connections, nil
}

func (s *ChatService) BanUser(markerID, userID string, duration time.Duration) error {
	// First, kick the user from the room.
	err := s.KickUserFromRoom(markerID, userID)
	if err != nil {
		log.Printf("Error kicking user %s from room %s: %v", userID, markerID, err)
	}

	// Then, set the ban in Redis.
	banKey := fmt.Sprintf("ban_%s_%s", markerID, userID)
	ctx := context.Background()

	// Use Set to set the ban in Redis
	cmd := s.Redis.Core.Client.B().Set().Key(banKey).Value("banned").Nx().Ex(duration).Build()
	return s.Redis.Core.Client.Do(ctx, cmd).Error()
}

func (s *ChatService) GetBanDetails(markerID, userID string) (banned bool, remainingTime time.Duration, err error) {
	banKey := fmt.Sprintf("ban_%s_%s", markerID, userID)
	ctx := context.Background()

	// Use HGETALL to retrieve all field-value pairs from the hash
	cmd := s.Redis.Core.Client.B().Exists().Key(banKey).Build()

	// Execute the command and get the result as a string map
	exists, err := s.Redis.Core.Client.Do(ctx, cmd).AsBool()
	if err != nil {
		return false, 0, err
	}
	if !exists {
		return false, 0, nil
	}

	// Use TTL to get the remaining TTL of the ban
	cmd = s.Redis.Core.Client.B().Ttl().Key(banKey).Build()

	// Execute the command and get the result as a string map
	ttl, err := s.Redis.Core.Client.Do(ctx, cmd).AsInt64()
	if err != nil {
		return true, 0, err
	}

	// Check if the TTL indicates the key does not exist or has no expiration
	if ttl == -2 {
		return false, 0, nil // Key does not exist
	} else if ttl == -1 {
		return true, 0, nil // Key exists but no expiration is set
	}

	// Convert the TTL from seconds to time.Duration
	duration := time.Duration(ttl) * time.Second

	return true, duration, nil
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
