package handlers

import (
	"chulbong-kr/services"
	"log"

	"github.com/gofiber/contrib/websocket"
)

// HandleChatRoomHandler manages chat rooms and messaging
func HandleChatRoomHandler(c *websocket.Conn, markerID int) {
	clientNickname := c.Locals("username").(string)
	// clientNickname := "user" + uuid.New().String()

	services.WsManager.AddConnectionRoom(markerID, c)
	// Broadcast join message
	services.WsManager.SendMessageToRoom(markerID, clientNickname+" has joined the chat.")

	defer func() {
		// On disconnect, remove the client from the room
		services.WsManager.RemoveConnectionFromRoom(markerID, c)

		// Broadcast leave message
		services.WsManager.SendMessageToRoom(markerID, clientNickname+" has left the chat.")
	}()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		// Broadcast received message
		services.WsManager.SendMessageToRoom(markerID, string(message))
	}
}
