package handlers

import (
	"chulbong-kr/dto/notification"
	"chulbong-kr/middlewares"
	"chulbong-kr/services"
	"log"
	"strconv"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/rueidis"
)

// RegisterNotificationRoutes sets up the routes for Notification handling within the application.
func RegisterNotificationRoutes(api fiber.Router, websocketConfig websocket.Config) {
	api.Get("/ws/notification", middlewares.AuthSoftMiddleware, func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}, websocket.New(func(c *websocket.Conn) {
		var userId string

		if id, ok := c.Locals("userID").(int); ok {
			// Convert integer userID to string if it exists and is valid
			userId = strconv.Itoa(id)
		} else {
			// anonymous users
			userId = c.Query("request-id")
		}

		if userId == "" {
			c.WriteJSON(fiber.Map{"error": "wrong user id"})
			c.Close()
			return
		}

		WsNotificationHandler(c, userId)
	}, websocketConfig))

	api.Post("/notification", middlewares.AdminOnly, PostNotificationHandler)
}

// PostNotificationHandler handles POST requests to send notifications
func PostNotificationHandler(c *fiber.Ctx) error {
	var req notification.NotificationRedis
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	// Call the PostNotification function to insert the notification into the database and publish to Redis
	err := services.PostNotification(req.UserId, req.NotificationType, req.Title, req.Message, req.Metadata)
	if err != nil {
		log.Printf("Error posting notification: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to post notification"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Notification posted successfully"})
}

// user id can be anonymous too. it should check auth and if authenticated, use cookie value as userId
func WsNotificationHandler(c *websocket.Conn, userId string) {
	// Fetch unviewed notifications at the start of the WebSocket connection
	unviewedNotifications, err := services.GetNotifications(userId)
	if err != nil {
		log.Printf("Error fetching unviewed notifications: %v", err)
		return
	}
	for _, notification := range unviewedNotifications {
		jsonData, err := json.Marshal(notification)
		if err != nil {
			log.Printf("Error marshaling notification: %v", err)
			continue
		}
		if err := c.WriteMessage(websocket.TextMessage, jsonData); err != nil {
			log.Printf("Error sending unviewed notification to WebSocket: %v", err)
			continue
		}

		// Mark notifications as viewed based on their type
		services.MarkNotificationAsViewed(notification.NotificationId, notification.NotificationType, userId)
	}

	// Subscribe to Redis on a per-connection basis
	cancelSubscription, err := services.SubscribeNotification(userId, func(msg rueidis.PubSubMessage) {
		// This function will be called for each message received
		if err := c.WriteMessage(websocket.TextMessage, []byte(msg.Message)); err != nil {
			log.Printf("Error sending message to WebSocket: %v", err)
		} else {
			var notification notification.NotificationRedis
			json.Unmarshal([]byte(msg.Message), &notification)
			services.MarkNotificationAsViewed(notification.NotificationId, notification.NotificationType, userId)
		}
	})

	if err != nil {
		log.Printf("Error subscribing to notifications: %v", err)
		return
	}
	defer cancelSubscription() // unsubscribe when the WebSocket closes

	// Simple loop to keep the connection open and log any incoming messages which could include pings
	for {
		messageType, p, err := c.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}
		log.Printf("Received message of type %d: %s", messageType, string(p))
	}

	// Keep-alive loop: handle ping/pong and ensure the connection stays open
	// for {
	// 	// if err := c.PingHandler()("ping"); err != nil {
	// 	// 	log.Printf("Ping failed: %v", err)
	// 	// 	break // Exit loop if ping fails
	// 	// }
	// 	if err := c.SetReadDeadline(time.Now().Add(time.Second * 300)); err != nil {
	// 		break
	// 	}
	// 	// Wait for a pong response to keep the connection alive
	// 	if _, _, err := c.ReadMessage(); err != nil {
	// 		log.Printf("Error reading pong: %v", err)
	// 		break
	// 	}
	// }
}
