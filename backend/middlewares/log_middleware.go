package middlewares

import (
	"chulbong-kr/util"
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

		clientIP := util.GetUserIP(c)

		userAgent := c.Get(fiber.HeaderUserAgent)
		referer := c.Get(fiber.HeaderReferer)
		queryParams := c.OriginalURL()[len(c.Path()):]

		if duration.Seconds() > util.DELAY_THRESHOLD {
			go util.SendSlackNotification(duration, statusCode, clientIP, method, path, userAgent, queryParams, referer) // Send notification in a non-blocking way
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
				// zap.String("error", err.Error()),
				zap.Error(err), // Include the error if present
			)

		return err
	}
}
