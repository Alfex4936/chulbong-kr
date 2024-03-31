package services

import (
	"chulbong-kr/dto"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zeebo/xxh3"

	csmap "github.com/mhmtszr/concurrent-swiss-map"
)

var (
	LavinMQClient       *amqp.Connection
	ActiveSubscriptions = csmap.Create(
		csmap.WithShardCount[string, struct{}](64),
		csmap.WithCustomHasher[string, struct{}](func(key string) uint64 {
			return xxh3.HashString(key)
		}),
	)

	// Map to store cancellation functions for each room subscription
	cancellationFunctions = csmap.Create(
		csmap.WithShardCount[string, context.CancelFunc](64),
		csmap.WithCustomHasher[string, context.CancelFunc](func(key string) uint64 {
			return xxh3.HashString(key)
		}),
	)
)

func PublishMessageToAMQP(ctx context.Context, roomID, message, userNickname, userId string) {
	ch, err := LavinMQClient.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close() // Ensure the channel is closed when function returns

	queueName := fmt.Sprintf("chat_room_%s", roomID)

	broadcastMsg := dto.BroadcastMessage{
		UID:          uuid.New().String(),
		Message:      message,
		UserID:       userId,
		UserNickname: userNickname,
		RoomID:       roomID,
		Timestamp:    time.Now().UnixMilli(),
	}

	// Serialize the message struct to JSON
	msgJSON, err := json.Marshal(broadcastMsg)
	if err != nil {
		log.Printf("Error marshalling message to JSON: %v", err)
		return
	}

	// Publish a message to the queue
	ch.PublishWithContext(
		ctx,
		"",        // exchange - Using the default exchange which routes based on queue name
		queueName, // routing key (queue name)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msgJSON,
		},
	)
	// failOnError(err, "Failed to declare a queue")
}

func ListenFromAMQP(ctx context.Context, queueName, roomID string, callback func(string, []byte)) {
	ch, err := LavinMQClient.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close() // Ensure the channel is closed when function returns

	// _, err = ch.QueueDelete(
	// 	queueName, // queue name
	// 	true,      // ifUnused - only delete if unused
	// 	true,      // ifEmpty - only delete if empty
	// 	false,     // noWait - don't wait for server to confirm the deletion
	// )
	// if err != nil {
	// 	log.Printf("Failed to delete queue: %s", err)
	// 	return
	// }

	_, err = ch.QueuePurge(queueName, false)
	failOnError(err, "Failed to purge the queue")

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		amqp.Table{ // arguments - Use amqp.Table for passing queue arguments
			"x-expires":          6000000, // Queue expires after 6000 seconds of not being used
			"x-message-ttl":      30000,   // Messages expire after 30 seconds
			"x-max-length":       1000,    // Maximum length of the queue (number of messages)
			"x-max-length-bytes": 1000000, // Maximum size of the queue (in bytes)
			// "x-dead-letter-exchange": "myDLX", // Example for setting a dead-letter exchange
		}, // arguments,
	)
	failOnError(err, "Failed to declare a queue")

	consumerId := fmt.Sprintf("chat-%s", roomID)
	// log.Printf("[✅] Consumer! %s", consumerId)

	msgs, err := ch.Consume(
		q.Name,     // queue
		consumerId, // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)

	failOnError(err, "Failed to register a consumer")

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case d := <-msgs:
				if len(d.Body) == 0 {
					continue
				}
				// log.Printf("[✅] Message Received: %s", d.Body)
				callback(roomID, d.Body)
			case <-ctx.Done():
				// log.Printf("[✅] Done sub\n")
				return
			}
		}
	}()

	wg.Wait()
}

func SubscribeAndBroadcastFromAMQP(roomID string) {
	queueName := fmt.Sprintf("chat_room_%s", roomID)
	key := fmt.Sprintf("marker_%s", roomID)

	// Check if we already have an active subscription for this room
	if _, exists := cancellationFunctions.Load(key); !exists {
		ctx, cancel := context.WithCancel(context.Background())

		// Store the cancel function immediately to mark this room as having an active listener
		cancellationFunctions.Store(key, cancel)

		ListenFromAMQP(ctx, queueName, roomID, func(roomID string, messageJson []byte) {
			if len(messageJson) == 0 {
				return
			}

			WsRoomManager.BroadcastMessageToRoom2(roomID, messageJson)
		})
	}
	//  else {
	// 	// Log or handle the case where we're attempting to subscribe to a room that already has an active listener
	// 	log.Printf("Attempted to subscribe to room %s, but a subscription already exists.", roomID)
	// }
}

// StopSubscriptionForRoom to delete a room and stop its subscription
func StopSubscriptionForRoom(roomID string) { // marker_%s
	if cancel, exists := cancellationFunctions.Load(roomID); exists {
		log.Printf("[✅] Stopping subscription for %s", roomID)

		cancel() // This stops the ListenFromAMQP goroutine for the room
		cancellationFunctions.Delete(roomID)
	}
}

func StopConsuming(ch *amqp.Channel, consumerTag string) error {
	if err := ch.Cancel(consumerTag, false); err != nil {
		log.Printf("Failed to cancel consumer: %s", err)
		return err
	}
	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
