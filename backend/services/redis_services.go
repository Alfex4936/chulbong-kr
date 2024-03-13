package services

import (
	"encoding/json"
	"time"

	"github.com/gofiber/storage/redis/v3"
)

var (
	RedisStore *redis.Storage
)

const (
	ALL_MARKERS_KEY  string = "all_markers"
	USER_PROFILE_KEY string = "user_profile"
	USER_FAV_KEY     string = "user_fav"
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
