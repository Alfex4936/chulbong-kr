package util

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/slack-go/slack"
)

const (
	DELAY_THRESHOLD = 10.0
	TIME_FORMAT_STR = "2006-01-02 (Mon) 15:04:05"
)

var (
	SLACK_BOT_TOKEN          = os.Getenv("SLACK_BOT_TOKEN")
	SLACK_CHANNEL_ID         = os.Getenv("SLACK_CHANNEL_ID")
	SLACK_CHANNEL_PENDING_ID = os.Getenv("SLACK_CHANNEL_PENDING_ID")
	slackClient              = slack.New(SLACK_BOT_TOKEN)
)

func SendSlackNotification(duration time.Duration, statusCode int, clientIP, method, path, userAgent, queryParams, referer string) {
	currentTime := time.Now().Format(TIME_FORMAT_STR)

	// Header section with bold text and warning emoji
	headerSection := slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", "*Performance Alert* (threshold: "+strconv.FormatFloat(DELAY_THRESHOLD, 'f', 2, 64)+" s) :warning:", false, false),
		nil, nil)

	// Divider to separate header from body
	dividerSection := slack.NewDividerBlock()

	// Fields for detailed information
	fields := make([]*slack.TextBlockObject, 8)
	fields[0] = slack.NewTextBlockObject("mrkdwn", "*Method:*\n`"+method+"`", false, false)
	fields[1] = slack.NewTextBlockObject("mrkdwn", "*Path:*\n`"+path+"`", false, false)
	fields[2] = slack.NewTextBlockObject("mrkdwn", "*Query:*\n`"+queryParams+"`", false, false)
	fields[3] = slack.NewTextBlockObject("mrkdwn", "*Status:*\n`"+strconv.Itoa(statusCode)+"`", false, false)
	fields[4] = slack.NewTextBlockObject("mrkdwn", "*Duration:*\n`"+duration.String()+"`", false, false)
	fields[5] = slack.NewTextBlockObject("mrkdwn", "*Client IP:*\n`"+clientIP+"`", false, false)
	fields[6] = slack.NewTextBlockObject("mrkdwn", "*User-Agent:*\n`"+userAgent+"`", false, false)
	fields[7] = slack.NewTextBlockObject("mrkdwn", "*Referer:*\n`"+referer+"`", false, false)
	bodySection := slack.NewSectionBlock(nil, fields, nil)

	// Time section
	timeSection := slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", "*At Server Time:*\n`"+currentTime+"`", false, false),
		nil, nil)

	// Combine the message blocks
	msg := slack.MsgOptionBlocks(headerSection, dividerSection, bodySection, timeSection)

	// Post the message to Slack
	_, _, err := slackClient.PostMessage(SLACK_CHANNEL_ID, msg)
	if err != nil {
		log.Printf("Error sending message to Slack: %v\n", err)
	}
}

func SendDeploymentSuccessNotification(serverName, environment string) {
	currentTime := time.Now().Format(TIME_FORMAT_STR)

	// Header section with bold text and celebration emoji
	headerSection := slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", "*Deployment Success* :tada:", false, false),
		nil, nil)

	// Divider to separate header from body
	dividerSection := slack.NewDividerBlock()

	// Fields for detailed information about the deployment
	fields := make([]*slack.TextBlockObject, 3)
	fields[0] = slack.NewTextBlockObject("mrkdwn", "*Server:*\n`"+serverName+"`", false, false)
	fields[1] = slack.NewTextBlockObject("mrkdwn", "*Environment:*\n`"+environment+"`", false, false)
	fields[2] = slack.NewTextBlockObject("mrkdwn", "*Deployed At:*\n`"+currentTime+"` (UTC)", false, false)
	bodySection := slack.NewSectionBlock(nil, fields, nil)

	// Combine the message blocks
	msg := slack.MsgOptionBlocks(headerSection, dividerSection, bodySection)

	// Post the message to Slack
	_, _, err := slackClient.PostMessage(SLACK_CHANNEL_ID, msg)
	if err != nil {
		log.Printf("Error sending deployment success message to Slack: %v\n", err)
	}
}

func SendSlackReportNotification(reportDetails string) {
	currentTime := time.Now().Format(TIME_FORMAT_STR)

	// Header section with bold text and information emoji
	headerSection := slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", "(`"+currentTime+"`) *Daily Pending Reports* :information_source:", false, false),
		nil, nil)

	// Divider to separate header from body
	dividerSection := slack.NewDividerBlock()

	// Fields for detailed information
	bodySection := slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", reportDetails, false, false),
		nil, nil)

	// Combine the message blocks
	msg := slack.MsgOptionBlocks(headerSection, dividerSection, bodySection)

	// Post the message to Slack
	_, _, err := slackClient.PostMessage(SLACK_CHANNEL_PENDING_ID, msg)
	if err != nil {
		log.Printf("Error sending message to Slack: %v\n", err)
	}
}
