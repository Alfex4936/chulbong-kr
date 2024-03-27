package handlers

import (
	"chulbong-kr/services"
	"chulbong-kr/utils"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// HandleChatRoomHandler manages chat rooms and messaging
func HandleChatRoomHandler(c *websocket.Conn, markerID, reqID string) {
	// clientId := c.Locals("userID").(int)
	// clientNickname := c.Locals("username").(string)
	clientId := reqID

	exists, _ := services.CheckDuplicateConnection(markerID, clientId)
	if exists {
		c.WriteJSON(fiber.Map{"error": "duplicate connection"})
		c.Close()
		return
	}

	// clientId := rand.Int()

	// clientNickname := "user-" + uuid.New().String()
	clientNickname := services.GenerateKoreanNickname()

	// WsRoomManager = connections *haxmap.Map[string, []*websocket.Conn] // concurrent map
	services.WsRoomManager.SaveConnection(markerID, clientId, c)          // saves to local websocket conncetions
	services.AddConnectionRoomToRedis(markerID, clientId, clientNickname) // saves to redis, "room:%s:connections"

	// Broadcast join message
	// broadcasts directly by app memory objects
	services.WsRoomManager.BroadcastMessageToRoom(markerID, clientNickname+" 님이 입장하셨습니다.", clientNickname, clientId)

	defer func() {
		// On disconnect, remove the client from the room
		services.WsRoomManager.RemoveWsFromRoom(markerID, c)
		services.RemoveConnectionFromRedis(markerID, reqID)

		// Broadcast leave message
		services.WsRoomManager.BroadcastMessageToRoom(markerID, clientNickname+" 님이 퇴장하셨습니다.", clientNickname, clientId)
	}()

	c.SetPingHandler(func(appData string) error {
		// Respond with a pong
		return c.WriteMessage(websocket.PongMessage, []byte(appData))
	})

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		messageString := string(message)
		bad, _ := utils.CheckForBadWords(messageString)
		if bad {
			continue
		}

		if messageString == "ping" {
			err = c.WriteMessage(websocket.TextMessage, []byte("pong"))
			if err != nil {
				log.Printf("Error sending 'pong': %v", err)
			}
			continue // Skip further processing for this message

		}

		// Broadcast received message
		services.WsRoomManager.BroadcastMessageToRoom(markerID, messageString, clientNickname, clientId)
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
