package middleware

import (
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/Alfex4936/chulbong-kr/middleware/ratelimit"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/utils"
)

const (
	// X-RateLimit-* headers
	xRateLimitLimit     = "X-RateLimit-Limit"
	xRateLimitRemaining = "X-RateLimit-Remaining"
	xRateLimitReset     = "X-RateLimit-Reset"
	retry               = "Retry-After"
)

type SlidingWindow struct {
	entries *ratelimit.ConcurrentMap
}

var bufPool = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 64) // Preallocate buffer
		return &buf             // Return pointer to the buffer
	},
}

// New creates a new sliding window middleware handler
func (sw *SlidingWindow) New(cfg limiter.Config) fiber.Handler {
	maxStr := strconv.Itoa(cfg.Max)
	expiration := int64(cfg.Expiration.Seconds())

	sw.entries = ratelimit.NewConcurrentMap(256)

	// Use the optimized key generator
	if cfg.KeyGenerator == nil {
		cfg.KeyGenerator = keyGenerator
	}

	// Update timestamp every second
	utils.StartTimeStampUpdater()

	return func(c *fiber.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Get key from request
		key := cfg.KeyGenerator(c)

		// Get timestamp
		ts := int64(atomic.LoadUint32(&utils.Timestamp))

		// Get or create the entry
		var e *ratelimit.Item
		e, ok := sw.entries.Get(key)
		if !ok {
			e = &ratelimit.Item{}
			sw.entries.Set(key, e)
		}

		// Use per-entry mutex for thread safety
		var rate, remaining int
		var resetInSec int64

		// Lock the entry
		e.Mu.Lock()

		// Check if the entry has expired
		if e.Exp == 0 || ts >= int64(e.Exp) {
			e.PrevHits = e.CurrHits
			e.CurrHits = 0
			e.Exp = uint64(ts + expiration)
		}

		// Increment current hits
		e.CurrHits++

		// Calculate reset time and rate
		resetInSec = int64(e.Exp) - ts
		rate = (e.PrevHits*int(resetInSec))/int(expiration) + e.CurrHits
		remaining = cfg.Max - rate

		// Unlock the entry
		e.Mu.Unlock()

		// Check if the rate limit has been exceeded
		if remaining < 0 {
			c.Set(retry, strconv.FormatInt(resetInSec, 10))
			return cfg.LimitReached(c)
		}

		// Continue to the next handler
		err := c.Next()

		// Handle requests that should be skipped based on status code
		if (cfg.SkipSuccessfulRequests && c.Response().StatusCode() < fiber.StatusBadRequest) ||
			(cfg.SkipFailedRequests && c.Response().StatusCode() >= fiber.StatusBadRequest) {
			// Lock the entry
			e.Mu.Lock()
			e.CurrHits--
			remaining++
			e.Mu.Unlock()
		}

		// Set rate limit headers
		c.Set(xRateLimitLimit, maxStr)
		c.Set(xRateLimitRemaining, strconv.Itoa(remaining))
		c.Set(xRateLimitReset, strconv.FormatInt(resetInSec, 10))

		return err
	}
}

// Key generator with caching
func keyGenerator(c *fiber.Ctx) string {
	if key, ok := c.Locals("ratelimit_key").(string); ok {
		return key
	}
	key := c.IP()
	c.Locals("ratelimit_key", key)
	return key
}
