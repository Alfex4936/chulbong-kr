package services

import (
	"encoding/json"
	"time"

	"github.com/gofiber/storage/redis/v3"
)

var RedisStore *redis.Storage

const (
	ALL_MARKERS_KEY  string = "all_markers"
	USER_PROFILE_KEY string = "user_profile"
)

// SetCacheEntry sets a cache entry with the given key and value, with an expiration time.
func SetCacheEntry[T any](key string, value T, expiration time.Duration) error {
	// Marshal the fiber.Map to JSON
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Set the cache entry with expiration
	err = RedisStore.Set(key, jsonValue, expiration)
	if err != nil {
		return err
	}

	return nil
}

// GetCacheEntry retrieves a cache entry by its key and unmarshals it into
func GetCacheEntry[T any](key string) (T, error) {
	var result T

	// Get the cache entry
	jsonValue, err := RedisStore.Get(key)
	if err != nil {
		return result, err
	}

	// Unmarshal JSON into
	if len(jsonValue) > 0 {
		err = json.Unmarshal(jsonValue, &result)
		if err != nil {
			return result, err
		}
		return result, nil
	}

	// Return an empty result if no data found
	return result, nil
}

// ResetCache invalidates cache entries for both the metadata and body of a specific endpoint.
func ResetCache(key string) error {
	// Construct the keys for the metadata and body
	metadataKey := key
	bodyKey := key + "_body"

	// Delete the metadata key
	err := RedisStore.Delete(metadataKey)
	if err != nil {
		return err
	}

	// Delete the body key
	RedisStore.Delete(bodyKey)
	// if err != nil {
	//     return err
	// }

	return nil
}
