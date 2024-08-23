package middleware

import (
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/utils"
)

const (
	xRateLimitLimit     = "X-RateLimit-Limit"
	xRateLimitRemaining = "X-RateLimit-Remaining"
	xRateLimitReset     = "X-RateLimit-Reset"
)

type SlidingWindow struct{}

// New creates a new sliding window middleware handler
func (SlidingWindow) New(cfg limiter.Config) fiber.Handler {
	var (
		max        = strconv.Itoa(cfg.Max)
		expiration = uint64(cfg.Expiration.Seconds())
	)

	manager := newManager(cfg.Storage)

	StartTimeStampUpdater()

	// Return new handler
	return func(c *fiber.Ctx) error {
		// Skip middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		key := cfg.KeyGenerator(c)
		ts := uint64(atomic.LoadUint32(&utils.Timestamp))
		e := manager.getOrCreateEntry(key, ts, expiration)

		// Increment hits using atomic operation
		currHits := atomic.AddUint32(&e.currHits, 1)

		// Calculate when it resets in seconds
		resetInSec := e.exp - ts

		// Calculate weighted rate
		weight := float64(resetInSec) / float64(expiration)
		rate := int(float64(e.prevHits)*weight) + int(currHits)

		// Calculate remaining hits
		remaining := cfg.Max - rate

		// Check if hits exceed the cfg.Max
		if remaining < 0 {
			c.Set(xRateLimitLimit, max)
			c.Set(xRateLimitRemaining, "0")
			c.Set(xRateLimitReset, strconv.FormatUint(resetInSec, 10))
			c.Set(fiber.HeaderRetryAfter, strconv.FormatUint(resetInSec, 10))
			return cfg.LimitReached(c)
		}

		manager.updateEntry(key, e, ts, expiration, resetInSec)

		// Continue stack for reaching c.Response().StatusCode()
		err := c.Next()

		// Handle request skipping logic
		if shouldSkipRequest(cfg, c) {
			atomic.AddUint32(&e.currHits, ^uint32(0)) // Decrement hits
			remaining++
		}

		// Set RateLimit headers
		c.Set(xRateLimitLimit, max)
		c.Set(xRateLimitRemaining, strconv.Itoa(remaining))
		c.Set(xRateLimitReset, strconv.FormatUint(resetInSec, 10))

		return err
	}
}

func shouldSkipRequest(cfg limiter.Config, c *fiber.Ctx) bool {
	return (cfg.SkipSuccessfulRequests && c.Response().StatusCode() < fiber.StatusBadRequest) ||
		(cfg.SkipFailedRequests && c.Response().StatusCode() >= fiber.StatusBadRequest)
}

type item struct {
	currHits uint32
	prevHits uint32
	exp      uint64
}

type manager struct {
	pool    sync.Pool
	entries sync.Map
	storage fiber.Storage
}

func newManager(storage fiber.Storage) *manager {
	return &manager{
		pool: sync.Pool{
			New: func() interface{} {
				return new(item)
			},
		},
		storage: storage,
	}
}

func (m *manager) getOrCreateEntry(key string, ts uint64, expiration uint64) *item {
	// Get or initialize entry
	entry, _ := m.entries.LoadOrStore(key, m.pool.Get().(*item))
	e := entry.(*item)

	// Handle expiration logic
	if ts >= e.exp {
		e.prevHits = e.currHits
		e.currHits = 0
		e.exp = ts + expiration
	}
	return e
}

func (m *manager) updateEntry(key string, e *item, ts uint64, expiration uint64, resetInSec uint64) {
	m.entries.Store(key, e)
	go m.releaseExpiredEntries(key, ts+resetInSec+expiration)
}

func (m *manager) releaseExpiredEntries(key string, expirationTime uint64) {
	time.Sleep(time.Until(time.Unix(int64(expirationTime), 0)))
	if entry, ok := m.entries.LoadAndDelete(key); ok {
		m.pool.Put(entry)
	}
}

var (
	timestampTimer sync.Once
	Timestamp      uint32
)

func StartTimeStampUpdater() {
	timestampTimer.Do(func() {
		atomic.StoreUint32(&Timestamp, uint32(time.Now().Unix()))
		go func(sleep time.Duration) {
			ticker := time.NewTicker(sleep)
			defer ticker.Stop()

			for t := range ticker.C {
				atomic.StoreUint32(&Timestamp, uint32(t.Unix()))
			}
		}(1 * time.Second)
	})
}
