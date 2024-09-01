package middleware

import (
	"strconv"
	"sync/atomic"
	"time"

	"github.com/Alfex4936/chulbong-kr/middleware/ratelimit"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

const (
	// X-RateLimit-* headers
	xRateLimitLimit     = "X-RateLimit-Limit"
	xRateLimitRemaining = "X-RateLimit-Remaining"
	xRateLimitReset     = "X-RateLimit-Reset"
	retry               = "Retry-After"
)

// SlidingWindow implements LimiterHandler with an optimized approach
type SlidingWindow struct{}

// New creates a new sliding window middleware handler
func (SlidingWindow) New(cfg limiter.Config) fiber.Handler {
	var (
		max        = strconv.Itoa(cfg.Max)
		expiration = uint64(cfg.Expiration.Seconds())
	)

	// Create manager to simplify storage operations
	manager := ratelimit.NewManager(cfg.Storage)

	// Update timestamp every second
	ratelimit.StartTimeStampUpdater()

	return func(c *fiber.Ctx) error {
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		key := cfg.KeyGenerator(c)
		ts := uint64(atomic.LoadUint32(&ratelimit.Timestamp))

		// Get or initialize the item from the manager
		e := manager.Get(key)
		if e.Exp == 0 {
			e.Exp = ts + expiration
		} else if ts >= e.Exp {
			// Handle expiration
			e.PrevHits = e.CurrHits
			e.CurrHits = 0
			e.Exp = ts + expiration - (ts - e.Exp)
		}

		// Increment current hits
		e.CurrHits++

		// Calculate remaining requests
		resetInSec := e.Exp - ts
		weight := float64(resetInSec) / float64(expiration)
		rate := int(float64(e.PrevHits)*weight) + e.CurrHits
		remaining := cfg.Max - rate

		// Update the manager with the new item state
		manager.Set(key, e, time.Duration(resetInSec+expiration)*time.Second)

		if remaining < 0 {
			c.Set(xRateLimitLimit, max)
			c.Set(xRateLimitRemaining, strconv.Itoa(remaining))
			c.Set(xRateLimitReset, strconv.FormatUint(resetInSec, 10))
			c.Set(fiber.HeaderRetryAfter, strconv.FormatUint(resetInSec, 10))
			return cfg.LimitReached(c)
		}

		err := c.Next()

		if (cfg.SkipSuccessfulRequests && c.Response().StatusCode() < fiber.StatusBadRequest) ||
			(cfg.SkipFailedRequests && c.Response().StatusCode() >= fiber.StatusBadRequest) {
			e.CurrHits--
			remaining++
			manager.Set(key, e, cfg.Expiration)
		}

		c.Set(xRateLimitLimit, max)
		c.Set(xRateLimitRemaining, strconv.Itoa(remaining))
		c.Set(xRateLimitReset, strconv.FormatUint(resetInSec, 10))

		return err
	}
}
