package util

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/slack-go/slack"
)

const DELAY_THRESHOLD = 10.0

var (
	SLACK_BOT_TOKEN  = os.Getenv("SLACK_BOT_TOKEN")
	SLACK_CHANNEL_ID = os.Getenv("SLACK_CHANNEL_ID")
)

func SendSlackNotification(duration time.Duration, statusCode int, clientIP, method, path, userAgent, queryParams, referer string) {
	client := slack.New(SLACK_BOT_TOKEN)

	currentTime := time.Now().Format("2006-01-02 (Mon) 15:04:05")

	// Header section with bold text and warning emoji
	headerText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Performance Alert* (threshold: %.2f s) :warning:", DELAY_THRESHOLD), false, false)
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

func SendDeploymentSuccessNotification(serverName, environment string) {
	client := slack.New(SLACK_BOT_TOKEN)

	currentTime := time.Now().Format("2006-01-02 (Mon) 15:04:05")

	// Header section with bold text and celebration emoji
	headerText := slack.NewTextBlockObject("mrkdwn", "*Deployment Success* :tada:", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Divider to separate header from body
	dividerSection := slack.NewDividerBlock()

	// Fields for detailed information about the deployment
	fields := make([]*slack.TextBlockObject, 0)
	fields = append(fields, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Server:*\n`%s`", serverName), false, false))
	fields = append(fields, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Environment:*\n`%s`", environment), false, false))
	fields = append(fields, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Deployed At:*\n`%s`", currentTime), false, false))
	bodySection := slack.NewSectionBlock(nil, fields, nil)

	// Combine the message blocks
	msg := slack.MsgOptionBlocks(headerSection, dividerSection, bodySection)

	// Post the message to Slack
	_, _, err := client.PostMessage(SLACK_CHANNEL_ID, msg)
	if err != nil {
		log.Printf("Error sending deployment success message to Slack: %v\n", err)
	}
}
