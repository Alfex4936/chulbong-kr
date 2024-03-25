package middlewares

import (
	"chulbong-kr/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ZapLogMiddleware(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Start timer
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Gather request/response details
		statusCode := c.Response().StatusCode()
		method := c.Method()
		path := c.OriginalURL()

		clientIP := c.Get("Fly-Client-IP")
		// If Fly-Client-IP is not found, fall back to X-Forwarded-For
		if clientIP == "" {
			clientIP = c.Get("X-Forwarded-For")
		}
		// If X-Forwarded-For is also empty, use c.IP() as the last resort
		if clientIP == "" {
			clientIP = c.IP()
		}

		userAgent := c.Get(fiber.HeaderUserAgent)
		referer := c.Get(fiber.HeaderReferer)
		queryParams := c.OriginalURL()[len(c.Path()):]

		if duration.Seconds() > utils.DELAY_THRESHOLD {
			go utils.SendSlackNotification(duration, statusCode, clientIP, method, path, userAgent, queryParams, referer) // Send notification in a non-blocking way
		}

		// Choose the log level and construct the log message
		level := zap.InfoLevel
		if statusCode >= 500 {
			level = zap.ErrorLevel
		} else if statusCode >= 400 {
			level = zap.WarnLevel
		}

		// Construct the structured log
		logger.Check(level, "HTTP request processed").
			Write(
				zap.Int("status", statusCode),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("client_ip", clientIP),
				zap.String("user_agent", userAgent),
				zap.String("referer", referer),
				zap.Duration("duration", duration),
				zap.Error(err), // Include the error if present
			)

		return err
	}
}
