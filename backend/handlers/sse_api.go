package handlers

import (
	"chulbong-kr/services"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/contrib/websocket"
)

func MarkerUpdateEventHandler(c *websocket.Conn) {
	pubsub := services.SubscribeToMarkerUpdates()
	defer pubsub.Close()

	// Channel to receive messages from Redis subscription
	ch := pubsub.Channel()

	userID, ok := c.Locals("userID").(int)
	if ok {
		services.WsManager.AddConnection(userID, c)
		defer services.WsManager.RemoveConnectionFromRoom(userID, c)
	}

	for msg := range ch {
		if err := services.WsManager.SendMessageToUser(userID, msg.Payload); err != nil {
			log.Println("websocket write error:", err)
			break
		}
	}
}

func MarkerLikeEventHandler(c *websocket.Conn) {
	pubsub := services.SubscribeToLikeEvent()
	defer pubsub.Close()

	ch := pubsub.Channel()

	userID, ok := c.Locals("userID").(int)
	if !ok {
		c.Close()
		c.WriteJSON("No Authorization for this session.")
		return
	}

	services.WsManager.AddConnection(userID, c)
	defer services.WsManager.RemoveConnectionFromRoom(userID, c)

	for msg := range ch {
		// Expecting msg.Payload to be like "ownerUserID-markerID"
		parts := strings.Split(msg.Payload, "-")
		if len(parts) != 2 {
			log.Println("Unexpected message format")
			continue
		}

		ownerUserID, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Printf("Failed to parse owner user ID: %v", err)
			continue
		}

		markerID, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Printf("Failed to parse marker ID: %v", err)
			continue
		}

		// Create the structured message
		message := string(services.MakeLikeNotification(ownerUserID, markerID))

		// Send the message to the marker owner
		if err := services.WsManager.SendMessageToUser(ownerUserID, message); err != nil {
			log.Printf("Failed to send like notification to user %d: %v", ownerUserID, err)
		}
	}
}
