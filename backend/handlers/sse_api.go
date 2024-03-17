package handlers

import (
	"chulbong-kr/services"

	"github.com/gofiber/fiber/v2"
)

func MarkerUpdatesHandler(c *fiber.Ctx) error {

	pubsub := services.SubscribeToMarkerUpdates()
	defer pubsub.Close()

	c.Set(fiber.HeaderContentType, "text/event-stream")
	c.Set(fiber.HeaderCacheControl, "no-cache")

	// Listen for messages
	for msg := range pubsub.Channel() {
		// Send SSE to the client
		c.Write([]byte("data: " + msg.Payload + "\n\n"))
	}

	return nil
}
