package handlers

import (
	"bytes"
	"chulbong-kr/services"
	"chulbong-kr/util"
	"log"
	"strings"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// HandleChatRoomHandler manages chat rooms and messaging
func HandleChatRoomHandler(c *websocket.Conn, markerID, reqID string) {
	// clientId := c.Locals("userID").(int)
	// clientNickname := c.Locals("username").(string)
	if markerID == "" || strings.Contains(markerID, "&") {
		c.WriteJSON(fiber.Map{"error": "wrong marker id"})
		c.Close()
		return
	}
	clientId := reqID

	exists, _ := services.WsRoomManager.CheckDuplicateConnectionByLocal(markerID, clientId)
	if exists {
		c.WriteJSON(fiber.Map{"error": "duplicate connection"})
		c.Close()
		return
	}

	// clientId := rand.Int()

	// clientNickname := "user-" + uuid.New().String()
	clientNickname := util.GenerateKoreanNickname()

	// WsRoomManager = connections *haxmap.Map[string, []*websocket.Conn] // concurrent map
	services.WsRoomManager.SaveConnection(markerID, clientId, c) // saves to local websocket conncetions
	// services.AddConnectionRoomToRedis(markerID, clientId, clientNickname) // saves to redis, "room:%s:connections"

	// Broadcast join message
	// broadcasts directly by app memory objects
	// services.PublishMessageToAMQP(context.Background(), markerID, clientNickname+" 님이 입장하셨습니다.", clientNickname, clientId)
	services.WsRoomManager.BroadcastMessageToRoom(markerID, clientNickname+" 님이 입장하셨습니다.", clientNickname, clientId)
	services.WsRoomManager.BroadcastUserCountToRoomByLocal(markerID) // sends how many users in the room

	defer func() {
		// On disconnect, remove the client from the room
		services.WsRoomManager.RemoveWsFromRoom(markerID, clientId)
		// services.RemoveConnectionFromRedis(markerID, reqID)

		// Broadcast leave message
		// services.PublishMessageToAMQP(context.Background(), markerID, clientNickname+" 님이 퇴장하셨습니다.", clientNickname, clientId)
		services.WsRoomManager.BroadcastMessageToRoom(markerID, clientNickname+" 님이 퇴장하셨습니다.", clientNickname, clientId)
		services.WsRoomManager.BroadcastUserCountToRoomByLocal(markerID) // sends how many users in the room
	}()

	// c.SetPingHandler(func(appData string) error {
	// 	// Respond with a pong
	// 	return c.WriteMessage(websocket.PongMessage, []byte(appData))
	// })

	for {
		if err := c.SetReadDeadline(time.Now().Add(time.Second * 60)); err != nil {
			break
		}
		_, message, err := c.ReadMessage()
		if err != nil {
			// log.Printf("Error reading message: %v", err)
			break
		}

		if bytes.Equal(message, []byte(`{"type":"ping"}`)) {
			// if mType == 9 || mType == 10 {
			services.WsRoomManager.UpdateLastPing(markerID, c)
			// if err := c.WriteMessage(websocket.TextMessage, []byte("pong")); err != nil {
			// 	log.Printf("Error sending 'pong': %v", err)
			// }
			continue // Skip further processing for this message
		}

		message = bytes.TrimSpace(message)
		if len(message) == 0 {
			continue
		}

		messageString := string(message) // Convert to string only when necessary

		// First, remove URLs from the message
		messageWithoutURLs := util.RemoveURLs(messageString)

		// Then, replace bad words with asterisks in the message string
		cleanMessage, err := util.ReplaceBadWords(messageWithoutURLs)
		if err != nil {
			log.Printf("Error replacing bad words: %v", err)
			continue
		}

		if cleanMessage == "" {
			continue
		}

		// Publish the valid message to the RabbitMQ queue for this chat room
		// services.PublishMessageToAMQP(context.Background(), markerID, cleanMessage, clientNickname, clientId)

		// Broadcast received message
		services.WsRoomManager.UpdateLastPing(markerID, c)
		if err := services.WsRoomManager.BroadcastMessageToRoom(markerID, cleanMessage, clientNickname, clientId); err != nil {
			break
		}
	}
}

func GetRoomUsersHandler(c *fiber.Ctx) error {
	markerID := c.Params("markerID")

	// Call your function to get connection infos
	connections, err := services.GetAllRedisConnectionsFromRoom(markerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to get connections"})
	}

	return c.JSON(fiber.Map{"connections": connections, "total_users": len(connections)})
}
