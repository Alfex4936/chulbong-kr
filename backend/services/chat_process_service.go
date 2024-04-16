package services

import (
	"chulbong-kr/dto"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/goccy/go-json"
)

func ProcessMessageFromSubscription(msg []byte) {
	var broadcastMsg dto.BroadcastMessage
	err := json.Unmarshal(msg, &broadcastMsg)
	if err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		return
	}

	if hasProcessed(broadcastMsg.UID) {
		return // Skip processing if we've already handled this message
	}

	markAsProcessed(broadcastMsg.UID) // Mark the message as processed locally

	// then broadcast
	WsRoomManager.BroadcastMessageToRoom(broadcastMsg.RoomID, broadcastMsg.Message, broadcastMsg.UserNickname, broadcastMsg.UserID)
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
func processRetryQueue(ctx context.Context) {
	for {
		select {
		case task := <-retryQueue:
			// Attempt removal again
			ctx := context.Background()
			key := fmt.Sprintf("room:%s:connections", task.MarkerID)

			hdelCmd := RedisStore.B().Hdel().Key(key).Field(task.RequestID).Build()
			if err := RedisStore.Do(ctx, hdelCmd).Error(); err != nil {
				log.Printf("Retry failed for removal task, consider further action: %v", err)
			}
		case <-ctx.Done():
			// cancelRetryCtx()
			return // Exit the goroutine if the context is cancelled
		}
	}
}

func hasProcessed(uid string) bool {
	_, exists := WsRoomManager.processedMessages.Load(uid)
	return exists
}

func markAsProcessed(uid string) {
	WsRoomManager.processedMessages.Store(uid, struct{}{})
}
