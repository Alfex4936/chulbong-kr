package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Alfex4936/chulbong-kr/dto"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zeebo/xxh3"

	csmap "github.com/mhmtszr/concurrent-swiss-map"
)

type MqService struct {
	LavinMqClient         *amqp.Connection
	ActiveSubscriptions   *csmap.CsMap[string, struct{}]
	CancellationFunctions *csmap.CsMap[string, context.CancelFunc]
	ChatService           *ChatService
}

func NewMqService(
	lavinClient *amqp.Connection,
	chatService *ChatService,
) *MqService {
	return &MqService{
		LavinMqClient: lavinClient,
		ActiveSubscriptions: csmap.Create(
			csmap.WithShardCount[string, struct{}](64),
			csmap.WithCustomHasher[string, struct{}](func(key string) uint64 {
				return xxh3.HashString(key)
			}),
		),

		// Map to store cancellation functions for each room subscription
		CancellationFunctions: csmap.Create(
			csmap.WithShardCount[string, context.CancelFunc](64),
			csmap.WithCustomHasher[string, context.CancelFunc](func(key string) uint64 {
				return xxh3.HashString(key)
			}),
		),
		ChatService: chatService,
	}
}

// LAVINMQ:
func (s *MqService) PublishMessageToAMQP(ctx context.Context, roomID, message, userNickname, userID string) {
	ch, err := s.LavinMqClient.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close() // Ensure the channel is closed when function returns

	// queueName := fmt.Sprintf("chat_room_%s", roomID)

	broadcastMsg := dto.BroadcastMessage{
		UID:          uuid.New().String(),
		Message:      message,
		UserID:       userID,
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

	routingKey := fmt.Sprintf("chat.room.%s", roomID)

	// Publish a message to the queue
	ch.PublishWithContext(
		ctx,
		"topic_exchange", // exchange - Using the default exchange which routes based on queue name
		routingKey,       // routing key (queue name)
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msgJSON,
		},
	)
	// failOnError(err, "Failed to declare a queue")
}

func (s *MqService) ListenFromAMQP(ctx context.Context, queueName, roomID string, callback func(string, []byte)) {
	ch, err := s.LavinMqClient.Channel()
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

	err = ch.ExchangeDeclare(
		"topic_exchange", // name
		"topic",          // type
		true,             // durable
		false,            // auto-delete
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	failOnError(err, "Failed to declare exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		amqp.Table{ // arguments - Use amqp.Table for passing queue arguments
			"x-expires":          6000000, // Queue expires after 6000 seconds of not being used
			"x-message-ttl":      30000,   // Messages expire after 30 seconds
			"x-max-length":       1000,    // Maximum length of the queue (number of messages)
			"x-max-length-bytes": 1000000, // Maximum size of the queue (in bytes)
			// "x-dead-letter-exchange": "myDLX", // Example for setting a dead-letter exchange
		}, // arguments,
	)
	failOnError(err, "Failed to declare a queue")

	// ch.QueuePurge(queueName, false) // doesn't matter if it fails

	// queueName = "chat_messages"                          // A general queue for chat messages
	roomIDPattern := fmt.Sprintf("chat.room.%s", roomID) // Specific pattern for the room

	// msgs, err := ch.Consume(
	// 	q.Name,     // queue
	// 	consumerId, // consumer
	// 	true,       // auto-ack
	// 	false,      // exclusive
	// 	false,      // no-local
	// 	false,      // no-wait
	// 	nil,        // args
	// )
	err = ch.QueueBind(
		q.Name,           // queue name
		roomIDPattern,    // routing key pattern
		"topic_exchange", // exchange name
		false,
		nil,
	)
	failOnError(err, "Failed to bind queue")

	// Consume messages from the queue
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer - leave empty for auto-generated consumer tag
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
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

func (s *MqService) SubscribeAndBroadcastFromAMQP(roomID string) {
	queueName := fmt.Sprintf("chat_room_%s", roomID)
	key := fmt.Sprintf("marker_%s", roomID)

	// Check if we already have an active subscription for this room
	if _, exists := s.CancellationFunctions.Load(key); !exists {
		ctx, cancel := context.WithCancel(context.Background())

		// Store the cancel function immediately to mark this room as having an active listener
		s.CancellationFunctions.Store(key, cancel)

		s.ListenFromAMQP(ctx, queueName, roomID, func(roomID string, messageJson []byte) {
			if len(messageJson) == 0 {
				return
			}

			s.ChatService.BroadcastMessageToRoom2(roomID, messageJson)
		})
	}
	//  else {
	// 	// Log or handle the case where we're attempting to subscribe to a room that already has an active listener
	// 	log.Printf("Attempted to subscribe to room %s, but a subscription already exists.", roomID)
	// }
}

// StopSubscriptionForRoom to delete a room and stop its subscription
func (s *MqService) StopSubscriptionForRoom(roomID string) { // marker_%s
	if cancel, exists := s.CancellationFunctions.Load(roomID); exists {
		// log.Printf("[✅] Stopping subscription for %s", roomID)

		cancel() // This stops the ListenFromAMQP goroutine for the room
		s.CancellationFunctions.Delete(roomID)
	}
}

func (s *MqService) StopConsuming(ch *amqp.Channel, consumerTag string) error {
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
