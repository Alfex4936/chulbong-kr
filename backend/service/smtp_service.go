package service

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/Alfex4936/chulbong-kr/config"
	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/Alfex4936/chulbong-kr/util"
)

type SmtpService struct {
	Config *config.SmtpConfig
}

func NewSmtpService(config *config.SmtpConfig) *SmtpService {
	return &SmtpService{
		Config: config,
	}
}

var emailTemplate = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html lang="ko" xmlns="http://www.w3.org/1999/xhtml">
<head>
    <title>k-pullup.com 이메일 인증</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body>
<table style="background-color: #f3f3f3; color: #333; font-size: 16px; text-align: center; margin: 0; padding: 0;" width="100%%" cellspacing="0" cellpadding="0">
    <tr>
        <td align="center">
            <table style="margin: 40px auto; border-collapse: separate;" width="600" cellspacing="0" cellpadding="0">
                <tr>
                    <td style="padding-bottom: 20px;" align="center">
                        <h1 style="color: #e5b000;">이메일 인증 토큰</h1>
                    </td>
                </tr>
                <tr>
                    <td style="padding-bottom: 20px;" align="center">
                        <p>아래의 토큰을 사용하여 이메일을 인증해주세요:</p>
                    </td>
                </tr>
                <tr>
                    <td align="center">
                        <div style="background-color: #fff; border: 2px dashed #e5b000; padding: 10px 20px; margin-top: 20px; font-size: 20px; font-weight: bold; letter-spacing: 2px;">{{TOKEN}}</div>
                    </td>
                </tr>
            </table>
        </td>
    </tr>
</table>
</body>
</html>`

var emailTemplateForReset = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html lang="ko" xmlns="http://www.w3.org/1999/xhtml">
<head>
    <title>k-pullup.com 비밀번호 초기화 요청</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body>
<table style="background-color: #f3f3f3; color: #333; font-size: 16px; text-align: center; margin: 0; padding: 0;" width="100%%" cellspacing="0" cellpadding="0">
    <tr>
        <td align="center">
            <table style="margin: 40px auto; border-collapse: separate;" width="600" cellspacing="0" cellpadding="0">
                <tr>
                    <td style="padding-bottom: 20px;" align="center">
                        <h1 style="color: #e5b000;">Password Reset Request</h1>
                        <h1 style="color: #e5b000;">비밀번호 초기화 요청</h1>
                    </td>
                </tr>
                <tr>
                    <td style="padding-bottom: 20px;" align="center">
                        <p>You have requested to reset your password. Please click the link below to proceed:</p>
                        <p>비밀번호 초기화를 요청하셨습니다. 아래 링크를 들어가 계속하세요!:</p>
                    </td>
                </tr>
                <tr>
                    <td align="center">
                        <a href="{{RESET_LINK}}" style="display: inline-block; background-color: #e5b000; color: #fff; padding: 10px 20px; margin-top: 20px; font-size: 20px; font-weight: bold; text-decoration: none; letter-spacing: 2px;">Reset Password</a>
                    </td>
                </tr>
            </table>
        </td>
    </tr>
</table>
</body>
</html>`

var emailTemplateForPendingReports = `<!DOCTYPE html>
<html lang="ko">
<head>
    <title>Daily Pending Reports</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body style="margin:0; padding:0; font-family: Arial, sans-serif; background-color: #f3f3f3; color: #333; text-align: center;">
<table width="100%" cellspacing="0" cellpadding="0" style="background-color: #f3f3f3; margin: 0; padding: 0;">
    <tr>
        <td align="center">
            <table width="600" cellspacing="0" cellpadding="0" style="margin: 40px auto; border-collapse: collapse; background-color: #fff; border: 1px solid #ddd;">
                <thead>
                    <tr style="background-color: #e5b000;">
                        <th style="padding: 10px; border: 1px solid #ddd; color: #fff;">Report ID</th>
                        <th style="padding: 10px; border: 1px solid #ddd; color: #fff;">Description</th>
                        <th style="padding: 10px; border: 1px solid #ddd; color: #fff;">Link</th>
                    </tr>
                </thead>
                <tbody>
                    {{REPORTS}}
                </tbody>
            </table>
        </td>
    </tr>
</table>
</body>
</html>`

// SendVerificationEmail sends a verification email to the user
func (s *SmtpService) SendVerificationEmail(to, token string) error {

	// Define email headers including content type for HTML
	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: k-pullup 이메일 인증\r\nMIME-Version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n", s.Config.SmtpUsername, to)

	// Replace the {{TOKEN}} placeholder in the template with the actual token
	htmlBody := strings.Replace(emailTemplate, "{{TOKEN}}", token, -1)

	// Preallocate the byte slice with the combined length of headers and htmlBody
	totalLength := len(headers) + len(htmlBody)
	message := make([]byte, totalLength)

	// Use copy to avoid creating intermediate allocations
	copy(message, headers)                 // Copy headers into message
	copy(message[len(headers):], htmlBody) // Append htmlBody after headers

	// Connect to the SMTP server and send the email
	auth := smtp.PlainAuth("", s.Config.SmtpUsername, s.Config.SmtpPassword, s.Config.SmtpServer)
	err := smtp.SendMail(s.Config.SmtpServer+":"+s.Config.SmtpPort, auth, s.Config.SmtpUsername, []string{to}, message)
	if err != nil {
		return err
	}
	return nil
}

func (s *SmtpService) SendPasswordResetEmail(to, token string) error {
	// Define email headers
	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: k-pullup 비밀번호 초기화\r\nMIME-Version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n", s.Config.SmtpUsername, to)

	// Replace the {{RESET_LINK}} placeholder with the actual reset link
	clientURL := fmt.Sprintf("%s?token=%s&email=%s", s.Config.FrontendResetRouter, token, to)
	htmlBody := strings.Replace(emailTemplateForReset, "{{RESET_LINK}}", clientURL, -1)

	// Combine headers and HTML body into a single raw email message
	// Preallocate the byte slice with the combined length of headers and htmlBody
	totalLength := len(headers) + len(htmlBody)
	message := make([]byte, totalLength)

	// Use copy to avoid creating intermediate allocations
	copy(message, headers)                 // Copy headers into message
	copy(message[len(headers):], htmlBody) // Append htmlBody after headers

	// Connect to the SMTP server and send the email
	auth := smtp.PlainAuth("", s.Config.SmtpUsername, s.Config.SmtpPassword, s.Config.SmtpServer)
	err := smtp.SendMail(s.Config.SmtpServer+":"+s.Config.SmtpPort, auth, s.Config.SmtpUsername, []string{to}, message)
	if err != nil {
		return err
	}
	return nil
}

func (s *SmtpService) SendPendingReportsEmail(to string, reports []dto.MarkerReportResponse) error {
	// Define email headers
	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: Daily Pending Reports\r\nMIME-Version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n", s.Config.SmtpUsername, to)

	// Build the reports table rows for email
	var reportRows string
	var slackReportRows string
	for _, report := range reports {
		link := fmt.Sprintf("https://k-pullup.com/pullup/%d", report.MarkerID)
		reportRows += fmt.Sprintf("<tr><td style=\"padding: 8px; border: 1px solid #ddd;\">%d</td><td style=\"padding: 8px; border: 1px solid #ddd;\">%s</td><td style=\"padding: 8px; border: 1px solid #ddd;\"><a href=\"%s\" style=\"color: #e5b000; text-decoration: none;\">View Report</a></td></tr>", report.ReportID, report.Description, link)
		slackReportRows += fmt.Sprintf("*Report ID:* %d\n*Description:* %s\n*Link:* <%s|View Report>\n\n", report.ReportID, report.Description, link)
	}

	// Replace the {{REPORTS}} placeholder in the template with the actual report rows
	htmlBody := strings.Replace(emailTemplateForPendingReports, "{{REPORTS}}", reportRows, -1)

	// Combine headers and HTML body into a single raw email message
	// Preallocate the byte slice with the combined length of headers and htmlBody
	totalLength := len(headers) + len(htmlBody)
	message := make([]byte, totalLength)

	// Use copy to avoid creating intermediate allocations
	copy(message, headers)                 // Copy headers into message
	copy(message[len(headers):], htmlBody) // Append htmlBody after headers

	// Connect to the SMTP server and send the email
	auth := smtp.PlainAuth("", s.Config.SmtpUsername, s.Config.SmtpPassword, s.Config.SmtpServer)
	err := smtp.SendMail(s.Config.SmtpServer+":"+s.Config.SmtpPort, auth, s.Config.SmtpUsername, []string{to}, message)
	if err != nil {
		return err
	}

	// Send notification to Slack
	util.SendSlackReportNotification(slackReportRows)

	return nil
}
