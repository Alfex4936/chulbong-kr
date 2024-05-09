package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Alfex4936/chulbong-kr/dto"

	"github.com/goccy/go-json"
)

// ProcessMessageFromSubscription processes a message from a Redis subscription
func (s *ChatService) ProcessMessageFromSubscription(msg []byte) {
	var broadcastMsg dto.BroadcastMessage
	err := json.Unmarshal(msg, &broadcastMsg)
	if err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		return
	}

	if s.WebSocketManager.hasProcessed(broadcastMsg.UID) {
		return // Skip processing if we've already handled this message
	}

	s.WebSocketManager.markAsProcessed(broadcastMsg.UID) // Mark the message as processed locally

	// then broadcast
	if err := s.BroadcastMessageToRoom(broadcastMsg.RoomID, broadcastMsg.Message, broadcastMsg.UserNickname, broadcastMsg.UserID); err != nil {
		log.Printf("Error broadcasting message: %v", err)
		return
	}
}

// Clean every hour
func (manager *RoomConnectionManager) StartCleanUpProcessedMsg() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			manager.processedMessages.Clear()
		}
	}()
}

// Enqueue a task for retry
func enqueueRemovalTask(markerID, requestID string) {
	retryQueue <- RemovalTask{MarkerID: markerID, RequestID: requestID}
}

// Background task for processing the retry queue
func (s *ChatService) processRetryQueue(ctx context.Context) {
	for {
		select {
		case task := <-retryQueue:
			// Attempt removal again
			ctx := context.Background()
			key := fmt.Sprintf("room:%s:connections", task.MarkerID)

			hdelCmd := s.Redis.Core.Client.B().Hdel().Key(key).Field(task.RequestID).Build()
			if err := s.Redis.Core.Client.Do(ctx, hdelCmd).Error(); err != nil {
				log.Printf("Retry failed for removal task, consider further action: %v", err)
			}
		case <-ctx.Done():
			// cancelRetryCtx()
			return // Exit the goroutine if the context is cancelled
		}
	}
}

func (manager *RoomConnectionManager) hasProcessed(uid string) bool {
	_, exists := manager.processedMessages.Load(uid)
	return exists
}

func (manager *RoomConnectionManager) markAsProcessed(uid string) {
	manager.processedMessages.Store(uid, struct{}{})
}
