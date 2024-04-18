package services

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/goccy/go-json"
	"github.com/redis/rueidis"
	// redis "github.com/gofiber/storage/redis/v3"
)

var (
	RedisStore rueidis.Client
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

	ctx := context.Background()

	// Build the SET command
	setCmd := RedisStore.B().Set().
		Key(key).
		Value(rueidis.BinaryString(jsonValue)).
		Nx().
		Ex(expiration).
		Build()

	// Execute the SET command
	if err := RedisStore.Do(ctx, setCmd).Error(); err != nil {
		return err
	}

	// Since we are not reusing the command, no need to pin or unpin
	return nil
}

// GetCacheEntry retrieves a cache entry by its key and unmarshals it into the provided type.
func GetCacheEntry[T any](key string) (T, error) {
	var result T

	ctx := context.Background()

	// Build the GET command using the client's command builder
	getCmd := RedisStore.B().Get().Key(key).Build()

	// Execute the GET command
	resp, err := RedisStore.Do(ctx, getCmd).ToString()

	// Handle potential errors
	if err != nil {
		if err == rueidis.ErrNoSlot {
			// Key does not exist
			return result, errors.New("key does not exist")
		}
		return result, err
	}

	// Unmarshal JSON into the type provided if data was found
	if resp != "" {
		err = json.Unmarshal([]byte(resp), &result)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

// ResetCache invalidates cache entries by deleting the specified key
func ResetCache(key string) error {
	ctx := context.Background()

	// Build and execute the DEL command using the client
	delCmd := RedisStore.B().Del().Key(key).Build()

	// Execute the DELETE command
	if err := RedisStore.Do(ctx, delCmd).Error(); err != nil {
		return err
	}

	return nil
}

// ResetAllCache invalidates cache entries by deleting all keys matching a given pattern.
func ResetAllCache(pattern string) error {
	ctx := context.Background()

	var cursor uint64 = 0
	for {
		// Build the SCAN command with the current cursor
		scanCmd := RedisStore.B().Scan().Cursor(cursor).Match(pattern).Count(10).Build()

		// Execute the SCAN command to find keys matching the pattern
		resp, err := RedisStore.Do(ctx, scanCmd).ToArray()
		if err != nil {
			return err
		}

		// First element is the new cursor, subsequent elements are the keys
		if len(resp) > 0 {
			newCursorStr := resp[0].String()
			newCursor, err := strconv.ParseUint(newCursorStr, 10, 64)
			if err != nil {
				return err // handle parsing error
			}
			cursor = newCursor

			keys := make([]string, 0, len(resp)-1)
			for _, msg := range resp[1:] {
				keys = append(keys, msg.String())
			}

			// Delete keys using individual DEL commands
			for _, key := range keys {
				delCmd := RedisStore.B().Del().Key(key).Build()
				if err := RedisStore.Do(ctx, delCmd).Error(); err != nil {
					return err
				}
			}
		}

		// If the cursor returned by SCAN is 0, we've iterated through all the keys
		if cursor == 0 {
			break
		}
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
