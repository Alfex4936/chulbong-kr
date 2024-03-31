package handlers

import (
	"bytes"
	"chulbong-kr/services"
	"chulbong-kr/utils"
	"context"
	"log"
	"strings"

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

	exists, _ := services.CheckDuplicateConnection(markerID, clientId)
	if exists {
		c.WriteJSON(fiber.Map{"error": "duplicate connection"})
		c.Close()
		return
	}

	// clientId := rand.Int()

	// clientNickname := "user-" + uuid.New().String()
	clientNickname := utils.GenerateKoreanNickname()

	// WsRoomManager = connections *haxmap.Map[string, []*websocket.Conn] // concurrent map
	services.WsRoomManager.SaveConnection(markerID, clientId, c)          // saves to local websocket conncetions
	services.AddConnectionRoomToRedis(markerID, clientId, clientNickname) // saves to redis, "room:%s:connections"

	// Broadcast join message
	// broadcasts directly by app memory objects
	services.PublishMessageToAMQP(context.Background(), markerID, clientNickname+" 님이 입장하셨습니다.", clientNickname, clientId)
	// services.WsRoomManager.BroadcastMessageToRoom(markerID, clientNickname+" 님이 입장하셨습니다.", clientNickname, clientId)

	services.WsRoomManager.BroadcastUserCountToRoom(markerID) // sends how many users in the room

	defer func() {
		// On disconnect, remove the client from the room
		services.WsRoomManager.RemoveWsFromRoom(markerID, c)
		services.RemoveConnectionFromRedis(markerID, reqID)

		// Broadcast leave message
		services.PublishMessageToAMQP(context.Background(), markerID, clientNickname+" 님이 퇴장하셨습니다.", clientNickname, clientId)
		// services.WsRoomManager.BroadcastMessageToRoom(markerID, clientNickname+" 님이 퇴장하셨습니다.", clientNickname, clientId)
		services.WsRoomManager.BroadcastUserCountToRoom(markerID) // sends how many users in the room
	}()

	// c.SetPingHandler(func(appData string) error {
	// 	// Respond with a pong
	// 	return c.WriteMessage(websocket.PongMessage, []byte(appData))
	// })

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			// log.Printf("Error reading message: %v", err)
			break
		}

		if bytes.Equal(message, []byte("ping")) {
			if err := c.WriteMessage(websocket.TextMessage, []byte("pong")); err != nil {
				log.Printf("Error sending 'pong': %v", err)
			}
			continue // Skip further processing for this message
		}

		message = bytes.TrimSpace(message)
		if len(message) == 0 {
			continue
		}

		messageString := string(message) // Convert to string only when necessary
		bad, _ := utils.CheckForBadWords(messageString)
		if bad {
			continue
		}

		// Publish the valid message to the RabbitMQ queue for this chat room
		services.PublishMessageToAMQP(context.Background(), markerID, messageString, clientNickname, clientId)

		// // Broadcast received message
		// services.WsRoomManager.BroadcastMessageToRoom(markerID, messageString, clientNickname, clientId)
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
