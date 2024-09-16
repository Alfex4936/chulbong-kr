package middleware

import (
	"io"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a new Fiber app with the SlidingWindow middleware
func setupApp(cfg limiter.Config) *fiber.App {
	app := fiber.New()
	app.Use(limiter.New(cfg))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	return app
}

func TestSlidingWindowRateLimiting(t *testing.T) {
	app := setupApp(limiter.Config{
		Max:               5,
		Expiration:        time.Second * 10,
		KeyGenerator:      keyGenerator,
		LimiterMiddleware: &SlidingWindow{},
		LimitReached: func(c *fiber.Ctx) error {
			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			c.Status(429).SendString("Too Many Requests")
			return nil
		},
	})

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		resp, _ := app.Test(req, -1)
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "OK", string(body))
	}

	// The next request should be rate limited
	req := httptest.NewRequest("GET", "/", nil)
	resp, _ := app.Test(req, -1)
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, 429, resp.StatusCode)
	assert.Equal(t, "Too Many Requests", string(body))

	// Check headers
	// log.Printf("Checking headers: %+v", resp.Header)
	// limit := resp.Header.Get("X-RateLimit-Limit")
	// remaining := resp.Header.Get("X-RateLimit-Remaining")
	reset := resp.Header.Get("Retry-After")

	// assert.Equal(t, strconv.Itoa(5), limit)
	// assert.Equal(t, strconv.Itoa(0), remaining)
	assert.NotEmpty(t, reset)
}

func TestSlidingWindowConcurrency(t *testing.T) {
	app := setupApp(limiter.Config{
		Max:               10,
		Expiration:        time.Second * 10,
		KeyGenerator:      func(c *fiber.Ctx) string { return c.IP() },
		LimiterMiddleware: &SlidingWindow{},
		LimitReached: func(c *fiber.Ctx) error {
			// Custom response when rate limit is exceeded
			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			c.Status(429).SendString("Too many requests, please try again later.")
			return nil
		},
	})

	wg := sync.WaitGroup{}
	concurrentRequests := 100

	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest("GET", "/", nil)
			resp, _ := app.Test(req, -1)
			assert.Contains(t, []int{200, 429}, resp.StatusCode)
		}()
	}

	wg.Wait()
}

func TestSlidingWindowSkipSuccessfulRequests(t *testing.T) {
	app := setupApp(limiter.Config{
		Max:                    5,
		Expiration:             time.Second * 10,
		SkipSuccessfulRequests: true,
		KeyGenerator:           func(c *fiber.Ctx) string { return c.IP() },
		LimiterMiddleware:      &SlidingWindow{},
		LimitReached: func(c *fiber.Ctx) error {
			// Custom response when rate limit is exceeded
			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			c.Status(429).SendString("Too many requests, please try again later.")
			return nil
		},
	})

	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		resp, _ := app.Test(req, -1)
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "OK", string(body))
	}

	// Even after 10 requests, since successful requests are skipped, it should not limit
	req := httptest.NewRequest("GET", "/", nil)
	resp, _ := app.Test(req, -1)
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "OK", string(body))
}

func benchmarkMiddleware(b *testing.B, cfg limiter.Config) {
	app := setupApp(cfg)

	// Reset the timer to ignore setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		resp, _ := app.Test(req, -1)
		resp.Body.Close()
	}
}

func BenchmarkLimiterSlidingWindow(b *testing.B) {
	cfg := limiter.Config{
		Max:               5,
		Expiration:        time.Second * 10,
		KeyGenerator:      keyGenerator,
		LimiterMiddleware: &SlidingWindow{},
		LimitReached: func(c *fiber.Ctx) error {
			c.Status(429).SendString("Too Many Requests")
			return nil
		},
	}
	benchmarkMiddleware(b, cfg)
}

func BenchmarkLimiterOriginal(b *testing.B) {
	cfg := limiter.Config{
		Max:               5,
		Expiration:        time.Second * 10,
		KeyGenerator:      keyGenerator,
		LimiterMiddleware: limiter.SlidingWindow{},
		LimitReached: func(c *fiber.Ctx) error {
			c.Status(429).SendString("Too Many Requests")
			return nil
		},
	}
	benchmarkMiddleware(b, cfg)
}
