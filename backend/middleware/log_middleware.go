package middleware

import (
	"net/url"
	"time"

	"github.com/Alfex4936/chulbong-kr/util"
	"github.com/jmoiron/sqlx"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	processMsg         = "HTTP request processed"
	insertVisitorQuery = "INSERT IGNORE INTO Visitors (IPAddress, VisitDate) VALUES (?, ?)"
)

type LogMiddleware struct {
	ChatUtil *util.ChatUtil
	DB       *sqlx.DB
}

func NewLogMiddleware(chatUtil *util.ChatUtil, db *sqlx.DB, logger *zap.Logger) *LogMiddleware {

	return &LogMiddleware{
		ChatUtil: chatUtil,
		DB:       db,
	}
}

func (l *LogMiddleware) ZapLogMiddleware(logger *zap.Logger) fiber.Handler {
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

		// Get the original URL (including query params)
		encodedPath := c.OriginalURL()

		// Decode URL to handle any percent-encoded characters (like Korean)
		path, encodeErr := url.QueryUnescape(encodedPath)
		if encodeErr != nil {
			// If decoding fails, fallback to the original encoded path
			path = encodedPath
		}

		clientIP := l.ChatUtil.GetUserIP(c)
		//visitDate := time.Now().Format("2006-01-02")

		// Perform the database logging in a goroutine
		// go func() {
		// 	if clientIP != "" {
		// 		// Check if the clientIP is a well-formed IP address
		// 		if ip := net.ParseIP(clientIP); ip != nil {
		// 			check, err := l.ChatUtil.IsIPFromSouthKorea(clientIP)
		// 			if check || err != nil {
		// 				l.DB.Exec(insertVisitorQuery, clientIP, visitDate)
		// 			}
		// 		}
		// 	}
		// }()

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

		// Include error details if an error occurred
		var errMsg string
		if err != nil {
			errMsg = err.Error()
		}

		// Construct the structured log
		logger.Check(level, processMsg).
			Write(
				zap.Int("status", statusCode),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("client_ip", clientIP),
				zap.String("user_agent", userAgent),
				zap.String("referer", referer),
				zap.Duration("duration", duration),
				zap.String("error", errMsg), // Log the error message
				// zap.Stack("stacktrace"),
				// zap.Error(err), // Include the error if present
			)

		return err
	}
}
