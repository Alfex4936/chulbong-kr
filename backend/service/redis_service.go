package service

import (
	"context"
	"sync"
	"time"

	"github.com/Alfex4936/chulbong-kr/config"

	"github.com/goccy/go-json"
	"github.com/redis/rueidis"
)

type RedisClient struct {
	Mu     sync.RWMutex
	Client rueidis.Client
}

func (src *RedisClient) Reconnect(newClient rueidis.Client) {
	src.Mu.Lock()
	defer src.Mu.Unlock()
	src.Client.Close()
	src.Client = newClient
}

type RedisService struct {
	RedisConfig *config.RedisConfig
	Core        *RedisClient
}

// NewRedisService creates a new instance of RedisService with the provided configuration and Redis client.
func NewRedisService(redisConfig *config.RedisConfig, redis *RedisClient) *RedisService {
	return &RedisService{
		RedisConfig: redisConfig,
		Core:        redis,
	}
}

// TODO: cannot use Generic as Fx doesn't support it directly maybe
// SetCacheEntry sets a cache entry with the given key and value, with an expiration time.
func (s *RedisService) SetCacheEntry(key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	ctx := context.Background()
	setCmd := s.Core.Client.B().Set().Key(key).Value(rueidis.BinaryString(jsonValue)).Nx().Ex(expiration).Build()
	return s.Core.Client.Do(ctx, setCmd).Error()
}

// TODO: cannot use Generic as Fx doesn't support it directly maybe
// GetCacheEntry retrieves a cache entry by its key and unmarshals it into the provided type.
func (s *RedisService) GetCacheEntry(key string, target interface{}) error {
	ctx := context.Background()
	getCmd := s.Core.Client.B().Get().Key(key).Build()
	resp, err := s.Core.Client.Do(ctx, getCmd).ToString()

	if err != nil {
		return err
	}
	if resp != "" {
		err = json.Unmarshal([]byte(resp), target)
	}
	return err
}

// ResetCache invalidates cache entries by deleting the specified key
func (s *RedisService) ResetCache(key string) error {
	ctx := context.Background()

	// Build and execute the DEL command using the client
	delCmd := s.Core.Client.B().Del().Key(key).Build()

	// Execute the DELETE command
	if err := s.Core.Client.Do(ctx, delCmd).Error(); err != nil {
		return err
	}

	return nil
}

// ResetAllCache invalidates cache entries by deleting all keys matching a given pattern.
func (s *RedisService) ResetAllCache(pattern string) error {
	ctx := context.Background()

	var cursor uint64
	for {
		// Build the SCAN command with the current cursor
		scanCmd := s.Core.Client.B().Scan().Cursor(cursor).Match(pattern).Count(10).Build()

		// Execute the SCAN command to find keys matching the pattern
		scanEntry, err := s.Core.Client.Do(ctx, scanCmd).AsScanEntry()
		if err != nil {
			return err
		}

		// Use the ScanEntry for cursor and keys directly
		cursor = scanEntry.Cursor
		keys := scanEntry.Elements

		// Delete keys using individual DEL commands
		for _, key := range keys {
			delCmd := s.Core.Client.B().Del().Key(key).Build()
			if err := s.Core.Client.Do(ctx, delCmd).Error(); err != nil {
				return err
			}
		}

		// If the cursor returned by SCAN is 0, iterated through all the keys
		if cursor == 0 {
			break
		}
	}

	return nil
}
