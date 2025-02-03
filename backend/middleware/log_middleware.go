package middleware

import (
	"time"

	"github.com/Alfex4936/chulbong-kr/util"
	"github.com/jmoiron/sqlx"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

// TODO: Decode unicode?
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
		path := c.Path()

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
		queryParams := string(c.Context().URI().QueryString())

		if duration.Seconds() > util.DELAY_THRESHOLD {
			go util.SendSlackNotification(duration, statusCode, clientIP, method, path, userAgent, queryParams, referer) // Send notification in a non-blocking way
		}

		// Choose the log level and construct the log message
		var level zapcore.Level
		switch {
		case statusCode >= 500:
			level = zapcore.ErrorLevel
		case statusCode >= 400:
			level = zapcore.WarnLevel
		default:
			level = zapcore.InfoLevel
		}

		// Log the request
		logger.Check(level, processMsg).
			Write(
				zap.Int("status", statusCode),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("client_ip", clientIP),
				zap.String("user_agent", userAgent),
				zap.String("referer", referer),
				zap.String("query_params", queryParams),
				zap.Duration("duration", duration),
				zap.Error(err),
			)

		return err
	}
}
