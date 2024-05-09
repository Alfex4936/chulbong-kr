package handler

import (
	"bytes"

	"github.com/Alfex4936/chulbong-kr/service"
	"github.com/Alfex4936/chulbong-kr/util"

	"log"
	"strings"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type ChatHandler struct {
	ChatService *service.ChatService

	ChatUtil    *util.ChatUtil
	BadWordUtil *util.BadWordUtil
}

// NewChatHandler creates a new ChatHandler with dependencies injected
func NewChatHandler(chat *service.ChatService, cutil *util.ChatUtil, butil *util.BadWordUtil,
) *ChatHandler {
	return &ChatHandler{
		ChatService: chat,
		ChatUtil:    cutil,
		BadWordUtil: butil,
	}
}

// RegisterChatRoutes sets up the routes for chat handling within the application.
func RegisterChatRoutes(api fiber.Router, websocketConfig websocket.Config, handler *ChatHandler) {
	api.Get("/ws/:markerID", func(c *fiber.Ctx) error {
		// Extract markerID from the parameter
		markerID := c.Params("markerID")
		reqID := c.Query("request-id")

		// Use GetBanDetails to check if the user is banned and get the remaining ban time
		banned, remainingTime, err := handler.ChatService.GetBanDetails(markerID, reqID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}
		if banned {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":         "User is banned",
				"remainingTime": remainingTime.Seconds(), // Respond with remaining time in seconds
			})
		}

		// Proceed with WebSocket upgrade if not banned
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}, websocket.New(func(c *websocket.Conn) {
		// Extract markerID from the parameter again if necessary
		markerID := c.Params("markerID")
		reqID := c.Query("request-id")

		// Now, the connection is already upgraded to WebSocket, and passed the ban check.
		handler.HandleChatRoom(c, markerID, reqID)
	}, websocketConfig))
}

// HandleChatRoomHandler manages chat rooms and messaging
func (h *ChatHandler) HandleChatRoom(c *websocket.Conn, markerID, reqID string) {
	// clientID := c.Locals("userID").(int)
	// clientNickname := c.Locals("username").(string)
	if markerID == "" || strings.Contains(markerID, "&") {
		c.WriteJSON(fiber.Map{"error": "wrong marker id"})
		c.Close()
		return
	}
	clientID := reqID

	exists, _ := h.ChatService.CheckDuplicateConnectionByLocal(markerID, clientID)
	if exists {
		c.WriteJSON(fiber.Map{"error": "duplicate connection"})
		c.Close()
		return
	}

	// clientID := rand.Int()

	// clientNickname := "user-" + uuid.New().String()
	clientNickname := h.ChatUtil.GenerateKoreanNickname()

	// WsRoomManager = connections *haxmap.Map[string, []*websocket.Conn] // concurrent map
	h.ChatService.SaveConnection(markerID, clientID, c) // saves to local websocket conncetions
	// services.AddConnectionRoomToRedis(markerID, clientID, clientNickname) // saves to redis, "room:%s:connections"

	// Broadcast join message
	// broadcasts directly by app memory objects
	// services.PublishMessageToAMQP(context.Background(), markerID, clientNickname+" 님이 입장하셨습니다.", clientNickname, clientID)
	h.ChatService.BroadcastMessageToRoom(markerID, clientNickname+" 님이 입장하셨습니다.", clientNickname, clientID)
	h.ChatService.BroadcastUserCountToRoomByLocal(markerID) // sends how many users in the room

	defer func() {
		// On disconnect, remove the client from the room
		h.ChatService.RemoveWsFromRoom(markerID, clientID)
		// services.RemoveConnectionFromRedis(markerID, reqID)

		// Broadcast leave message
		// services.PublishMessageToAMQP(context.Background(), markerID, clientNickname+" 님이 퇴장하셨습니다.", clientNickname, clientID)
		h.ChatService.BroadcastMessageToRoom(markerID, clientNickname+" 님이 퇴장하셨습니다.", clientNickname, clientID)
		h.ChatService.BroadcastUserCountToRoomByLocal(markerID) // sends how many users in the room
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
			h.ChatService.UpdateLastPing(markerID, c)
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
		cleanMessage, err := h.BadWordUtil.ReplaceBadWords(messageWithoutURLs)
		if err != nil {
			log.Printf("Error replacing bad words: %v", err)
			continue
		}

		if cleanMessage == "" {
			continue
		}

		// Publish the valid message to the RabbitMQ queue for this chat room
		// services.PublishMessageToAMQP(context.Background(), markerID, cleanMessage, clientNickname, clientID)

		// Broadcast received message
		h.ChatService.UpdateLastPing(markerID, c)
		if err := h.ChatService.BroadcastMessageToRoom(markerID, cleanMessage, clientNickname, clientID); err != nil {
			break
		}
	}
}
