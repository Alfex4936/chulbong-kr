package middlewares

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

const DELAY_THRESHOLD = 10.0

var (
	SLACK_BOT_TOKEN  = os.Getenv("SLACK_BOT_TOKEN")
	SLACK_CHANNEL_ID = os.Getenv("SLACK_CHANNEL_ID")
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
		clientIP := c.IP()
		userAgent := c.Get(fiber.HeaderUserAgent)
		referer := c.Get(fiber.HeaderReferer)
		queryParams := c.OriginalURL()[len(c.Path()):]

		if duration.Seconds() > DELAY_THRESHOLD {
			go sendSlackNotification(duration, statusCode, clientIP, method, path, userAgent, queryParams, referer) // Send notification in a non-blocking way
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

		return nil
	}
}

func sendSlackNotification(duration time.Duration, statusCode int, clientIP, method, path, userAgent, queryParams, referer string) {
	client := slack.New(SLACK_BOT_TOKEN)

	currentTime := time.Now().Format("2006-01-02 (Mon) 15:04:05")

	// Header section with bold text and warning emoji
	headerText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Performance Alert* (threshold: %.2f ms) :warning:", DELAY_THRESHOLD), false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Divider to separate header from body
	dividerSection := slack.NewDividerBlock()

	// Fields for detailed information
	fields := make([]*slack.TextBlockObject, 0)
	fields = append(fields, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Method:*\n`%s`", method), false, false))
	fields = append(fields, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Path:*\n`%s`", path), false, false))
	fields = append(fields, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Query:*\n`%s`", queryParams), false, false))
	fields = append(fields, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Status:*\n`%d`", statusCode), false, false))
	fields = append(fields, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Duration:*\n`%s`", duration), false, false))
	fields = append(fields, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Client IP:*\n`%s`", clientIP), false, false))
	fields = append(fields, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*User-Agent:*\n`%s`", userAgent), false, false))
	fields = append(fields, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Referer:*\n`%s`", referer), false, false))
	bodySection := slack.NewSectionBlock(nil, fields, nil)

	// Time section
	timeText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*At Server Time:*\n`%s`", currentTime), false, false)
	timeSection := slack.NewSectionBlock(timeText, nil, nil)

	// Combine the message blocks
	msg := slack.MsgOptionBlocks(headerSection, dividerSection, bodySection, timeSection)

	// Post the message to Slack
	_, _, err := client.PostMessage(SLACK_CHANNEL_ID, msg)
	if err != nil {
		log.Printf("Error sending message to Slack: %v\n", err)
	}
}
