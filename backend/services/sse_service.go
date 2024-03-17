package services

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// PublishMarkerUpdate to publish messages
func PublishMarkerUpdate(message string) {
	err := RedisStore.Conn().Publish(context.Background(), "markerUpdates", message).Err()
	if err != nil {
		panic(err)
	}
}

// SubscribeToMarkerUpdates to subscribe to messages
func SubscribeToMarkerUpdates() *redis.PubSub {
	pubsub := RedisStore.Conn().Subscribe(context.Background(), "markerUpdates")
	return pubsub
}
