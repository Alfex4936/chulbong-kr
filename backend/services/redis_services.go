package services

import (
	"context"
	"time"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	// redis "github.com/gofiber/storage/redis/v3"
)

var (
	RedisStore *redis.Client
)

const (
	ALL_MARKERS_KEY  string = "all_markers"
	USER_PROFILE_KEY string = "user_profile"
	USER_FAV_KEY     string = "user_fav"
)

// SetCacheEntry sets a cache entry with the given key and value, with an expiration time.
func SetCacheEntry[T any](key string, value T, expiration time.Duration) error {
	// Marshal the value to JSON
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	pipe := RedisStore.TxPipeline()
	ctx := context.Background()

	pipe.Set(ctx, key, jsonValue, expiration)
	pipe.Expire(ctx, key, expiration)

	_, err = pipe.Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}

// GetCacheEntry retrieves a cache entry by its key and unmarshals it into
func GetCacheEntry[T any](key string) (T, error) {
	var result T

	// Execute the Get command
	cmd := RedisStore.Get(context.Background(), key)
	jsonValue, err := cmd.Result()

	// Handle potential errors
	if err == redis.Nil {
		// Key does not exist
		return result, nil // or errors.New("key does not exist")
	} else if err != nil {
		return result, err
	}

	// Unmarshal JSON into
	if len(jsonValue) > 0 {
		err = json.Unmarshal([]byte(jsonValue), &result)
		if err != nil {
			return result, err
		}
		return result, nil
	}

	// Return an empty result if no data found
	return result, nil
}

// ResetCache invalidates cache entries
func ResetCache(key string) error {
	// Delete the metadata key
	pipe := RedisStore.TxPipeline()
	ctx := context.Background()

	pipe.Del(ctx, key)

	_, err := pipe.Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}

// ResetAllCache invalidates cache entries
func ResetAllCache(pattern string) error {
	// Delete the metadata key
	pipe := RedisStore.TxPipeline()
	ctx := context.Background()

	// Use SCAN to find all keys matching the pattern
	var cursor uint64
	var err error
	for {
		var keys []string
		keys, cursor, err = RedisStore.Scan(ctx, cursor, pattern, 10).Result()
		if err != nil {
			return err
		}

		// Delete the keys found in this scan iteration
		for _, key := range keys {
			pipe.Del(ctx, key)
		}

		// If the cursor returned by SCAN is 0, we've iterated through all the keys
		if cursor == 0 {
			break
		}
	}

	// Execute the delete commands in the pipeline
	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// func CompileBadWordsPattern(badWords []string) error {
// 	var pattern strings.Builder
// 	// Start of the group
// 	pattern.WriteString("(")
// 	for i, word := range badWords {
// 		// QuoteMeta escapes all regex meta characters in the bad word
// 		pattern.WriteString(regexp.QuoteMeta(word))
// 		// Separate words with a pipe, except the last word
// 		if i < len(badWords)-1 {
// 			pattern.WriteString("|")
// 		}
// 	}
// 	// End of the group
// 	pattern.WriteString(")")

// 	var err error
// 	badWordRegex, err = regexp.Compile(pattern.String())
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func LoadBadWordsIntoRedis(filePath string) {
// 	const batchSize = 500

// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)
// 	batch := make([]string, 0, batchSize)

// 	for scanner.Scan() {
// 		word := scanner.Text()
// 		batch = append(batch, word)

// 		// Once we've collected enough words, insert them in a batch.
// 		if len(batch) >= batchSize {
// 			err := AddBadWords(batch)
// 			if err != nil {
// 				fmt.Printf("Failed to insert batch: %v\n", err)
// 			}
// 			// Reset the batch slice for the next group of words
// 			batch = batch[:0]
// 		}
// 	}

// 	// Don't forget to insert any words left in the batch after finishing the loop
// 	if len(batch) > 0 {
// 		err := AddBadWords(batch)
// 		if err != nil {
// 			fmt.Printf("Failed to insert final batch: %v\n", err)
// 		}
// 	}

// 	if err := scanner.Err(); err != nil {
// 		panic(err)
// 	}
// }
