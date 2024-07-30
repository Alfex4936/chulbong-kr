package handler

import (
	"log"
	"strconv"

	"github.com/Alfex4936/chulbong-kr/dto/notification"
	"github.com/Alfex4936/chulbong-kr/middleware"
	"github.com/Alfex4936/chulbong-kr/service"

	sonic "github.com/bytedance/sonic"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/rueidis"
)

type NotificationHandler struct {
	NotiService *service.NotificationService
}

// NewNotificationHandler creates a new NotificationHandler with dependencies injected
func NewNotificationHandler(notification *service.NotificationService,
) *NotificationHandler {
	return &NotificationHandler{
		NotiService: notification,
	}
}

// RegisterNotificationRoutes sets up the routes for Notification handling within the application.
func RegisterNotificationRoutes(api fiber.Router, websocketConfig websocket.Config, handler *NotificationHandler, authMiddleware *middleware.AuthMiddleware) {
	api.Get("/ws/notification", authMiddleware.VerifySoft, func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}, websocket.New(func(c *websocket.Conn) {
		var userID string

		if id, ok := c.Locals("userID").(int); ok {
			// Convert integer userID to string if it exists and is valid
			userID = strconv.Itoa(id)
		} else {
			// anonymous users
			userID = c.Query("request-id")
		}

		if userID == "" {
			c.WriteJSON(fiber.Map{"error": "wrong user id"})
			c.Close()
			return
		}

		handler.WsNotificationHandler(c, userID)
	}, websocketConfig))

	api.Post("/notification", authMiddleware.CheckAdmin, handler.PostNotificationHandler)
}

// PostNotificationHandler handles POST requests to send notifications
func (h *NotificationHandler) PostNotificationHandler(c *fiber.Ctx) error {
	var req notification.NotificationRedis
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	// Call the PostNotification function to insert the notification into the database and publish to Redis
	err := h.NotiService.PostNotification(req.UserId, req.NotificationType, req.Title, req.Message, req.Metadata)
	if err != nil {
		log.Printf("Error posting notification: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to post notification"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Notification posted successfully"})
}

// user id can be anonymous too. it should check auth and if authenticated, use cookie value as userID
func (h *NotificationHandler) WsNotificationHandler(c *websocket.Conn, userID string) {
	// Fetch unviewed notifications at the start of the WebSocket connection
	unviewedNotifications, err := h.NotiService.GetNotifications(userID)
	if err != nil {
		log.Printf("Error fetching unviewed notifications: %v", err)
		return
	}
	for _, notification := range unviewedNotifications {
		jsonData, err := sonic.Marshal(notification)
		if err != nil {
			log.Printf("Error marshaling notification: %v", err)
			continue
		}
		if err := c.WriteMessage(websocket.TextMessage, jsonData); err != nil {
			log.Printf("Error sending unviewed notification to WebSocket: %v", err)
			continue
		}

		// Mark notifications as viewed based on their type
		h.NotiService.MarkNotificationAsViewed(notification.NotificationId, notification.NotificationType, userID)
	}

	// Subscribe to Redis on a per-connection basis
	cancelSubscription, err := h.NotiService.SubscribeNotification(userID, func(msg rueidis.PubSubMessage) {
		// This function will be called for each message received
		if err := c.WriteMessage(websocket.TextMessage, []byte(msg.Message)); err != nil {
			log.Printf("Error sending message to WebSocket: %v", err)
		} else {
			var notification notification.NotificationRedis
			sonic.Unmarshal([]byte(msg.Message), &notification)
			h.NotiService.MarkNotificationAsViewed(notification.NotificationId, notification.NotificationType, userID)
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
