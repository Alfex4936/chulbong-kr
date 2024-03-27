package handlers

import (
	"chulbong-kr/services"
	"chulbong-kr/utils"
	"log"
	"math/rand"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// HandleChatRoomHandler manages chat rooms and messaging
func HandleChatRoomHandler(c *websocket.Conn, markerID string) {
	// clientId := c.Locals("userID").(int)
	// clientNickname := c.Locals("username").(string)

	clientId := rand.Int()
	// clientNickname := "user-" + uuid.New().String()
	clientNickname := services.GenerateKoreanNickname()

	// WsRoomManager = connections *haxmap.Map[string, []*websocket.Conn] // concurrent map
	services.WsRoomManager.SaveConnection(markerID, c)                    // saves to local websocket conncetions
	services.AddConnectionRoomToRedis(markerID, clientId, clientNickname) // saves to redis, "room:%s:connections"

	// Broadcast join message
	// broadcasts directly by app memory objects
	services.WsRoomManager.BroadcastMessageToRoom(markerID, clientNickname+" has joined the chat.", clientNickname, clientId)

	defer func() {
		// On disconnect, remove the client from the room
		services.WsRoomManager.RemoveWsFromRoom(markerID, c)
		services.RemoveConnectionFromRedis(markerID, clientId)

		// Broadcast leave message
		services.WsRoomManager.BroadcastMessageToRoom(markerID, clientNickname+" has left the chat.", clientNickname, clientId)
	}()

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

	return c.JSON(connections)
}
