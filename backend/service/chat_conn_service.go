package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/Alfex4936/chulbong-kr/util"
	"github.com/puzpuzpuz/xsync/v3"

	sonic "github.com/bytedance/sonic"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/redis/rueidis"
)

// LAVINMQ: Ensure only one subscription per room
// if _, exists := ActiveSubscriptions.Load(key); !exists {
// 	go SubscribeAndBroadcastFromAMQP(markerID)
// 	ActiveSubscriptions.Store(key, struct{}{}) // marker_%s
// }

// SaveConnection stores a WebSocket connection associated with a markerID in app memory
func (s *ChatService) SaveConnection(markerID, clientID, clientNickname string, conn *websocket.Conn) (bool, error) {
	// s.Logger.Info("Saving connection", zap.String("markerID", markerID), zap.String("clientID", clientID))

	// Get or create the inner map for the room
	roomConns, _ := s.WebSocketManager.rooms.LoadOrCompute(markerID, func() *xsync.MapOf[string, *ChulbongConn] {
		return xsync.NewMapOf[string, *ChulbongConn]()
	})

	// Check if the clientID already exists in the inner map
	_, loaded := roomConns.LoadOrStore(clientID, func() *ChulbongConn {
		// Create the new connection object
		newConn := &ChulbongConn{
			Socket:       conn,
			UserID:       clientID,
			Nickname:     clientNickname,
			Send:         make(chan []byte, 10), // Buffered channel
			InActiveChan: make(chan struct{}),
		}
		newConn.UpdateLastSeen()

		// Start the writePump in a separate goroutine
		go newConn.writePump()

		return newConn
	}())

	if loaded {
		// Connection for clientID already exists, avoid adding a duplicate
		return false, fmt.Errorf("duplicate connection")
	}

	// No duplicate found, connection has been added successfully
	return true, nil
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

func (s *ChatService) CheckDuplicateConnectionByLocal(markerID, clientID string) (bool, error) {
	roomConns, ok := s.WebSocketManager.rooms.Load(markerID)
	if !ok {
		return false, nil // No room found, so no duplicate
	}

	_, exists := roomConns.Load(clientID)
	return exists, nil
}

func (s *ChatService) UpdateLastPing(markerID, clientID string) {
	if roomConns, ok := s.WebSocketManager.rooms.Load(markerID); ok {
		if conn, ok := roomConns.Load(clientID); ok {
			conn.UpdateLastSeen()
		}
	}
}

func (s *RoomConnectionManager) StartConnectionChecker() {
	ticker := time.NewTicker(30 * time.Minute)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			s.CheckConnections()
		}
	}()
}

// CheckConnections iterates over all connections and sends a ping message.
// Connections that fail to respond can be considered inactive and removed.
func (s *RoomConnectionManager) CheckConnections() {
	gracePeriod := 10 * time.Minute
	now := time.Now()

	s.rooms.Range(func(markerID string, roomConns *xsync.MapOf[string, *ChulbongConn]) bool {
		roomConns.Range(func(clientID string, conn *ChulbongConn) bool {
			if now.Sub(conn.GetLastSeen()) > gracePeriod {
				select {
				case conn.InActiveChan <- struct{}{}:
					// Inactive signal sent
				default:
					// Channel is full; handle accordingly
				}
				roomConns.Delete(clientID)
			}
			return true
		})
		if roomConns.Size() == 0 {
			s.rooms.Delete(markerID)
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
func (s *ChatService) RemoveWsFromRoom(markerID, clientID string) (string, error) {
	if roomConns, ok := s.WebSocketManager.rooms.Load(markerID); ok {
		if conn, ok := roomConns.Load(clientID); ok {
			clientNickname := conn.Nickname
			close(conn.Send)
			close(conn.InActiveChan)
			conn.Socket.Close()
			roomConns.Delete(clientID)

			// s.Logger.Info("Connection closed", zap.String("markerID", markerID), zap.String("clientID", clientID))

			if roomConns.Size() == 0 {
				s.WebSocketManager.rooms.Delete(markerID)
			}
			return clientNickname, nil
		}
	}
	return "", fmt.Errorf("connection not found")
}

// KickUserFromRoom closes the connection for a user in a specified room.
func (s *ChatService) KickUserFromRoom(markerID, clientID string) error {
	if roomConns, ok := s.WebSocketManager.rooms.Load(markerID); ok {
		if conn, ok := roomConns.Load(clientID); ok {
			close(conn.Send)
			close(conn.InActiveChan)
			conn.Socket.Close()
			roomConns.Delete(clientID)

			if roomConns.Size() == 0 {
				s.WebSocketManager.rooms.Delete(markerID)
			}
			return nil
		}
		return errors.New("user not found in room")
	}
	return errors.New("room not found")
}

func (s *ChatService) GetUserCountInRoom(ctx context.Context, markerID string) (int64, error) {
	key := fmt.Sprintf("room:%s:connections", markerID)

	// Use HLEN to get the number of connections in the room
	cmd := s.Redis.Core.Client.B().Hlen().Key(key).Build()
	return s.Redis.Core.Client.Do(ctx, cmd).AsInt64()
}

func (s *ChatService) GetUserCountInRoomByLocal(markerID string) (int, error) {
	if roomConns, ok := s.WebSocketManager.rooms.Load(markerID); ok {
		return int(roomConns.Size()), nil
	}
	return 0, nil
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

		if err := sonic.ConfigFastest.Unmarshal(jsonBytes, &connInfo); err != nil {
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
