package service

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/Alfex4936/chulbong-kr/config"
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
    <title>Email verification for k-pullup.com</title>
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
    <title>Password Reset request for k-pullup.com</title>
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
                    </td>
                </tr>
                <tr>
                    <td style="padding-bottom: 20px;" align="center">
                        <p>You have requested to reset your password. Please click the link below to proceed:</p>
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

// SendVerificationEmail sends a verification email to the user
func (s *SmtpService) SendVerificationEmail(to, token string) error {

	// Define email headers including content type for HTML
	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: k-pullup Email Verification\r\nMIME-Version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n", s.Config.SmtpUsername, to)

	// Replace the {{TOKEN}} placeholder in the template with the actual token
	htmlBody := strings.Replace(emailTemplate, "{{TOKEN}}", token, -1)

	// Combine headers and HTML body into a single raw email message
	message := []byte(headers + htmlBody)

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
	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: Password Reset for k-pullup\r\nMIME-Version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n", s.Config.SmtpUsername, to)

	// Replace the {{RESET_LINK}} placeholder with the actual reset link
	clientURL := fmt.Sprintf("%s?token=%s&email=%s", s.Config.FrontendResetRouter, token, to)
	htmlBody := strings.Replace(emailTemplateForReset, "{{RESET_LINK}}", clientURL, -1)

	// Combine headers and HTML body into a single raw email message
	message := []byte(headers + htmlBody)

	// Connect to the SMTP server and send the email
	auth := smtp.PlainAuth("", s.Config.SmtpUsername, s.Config.SmtpPassword, s.Config.SmtpServer)
	err := smtp.SendMail(s.Config.SmtpServer+":"+s.Config.SmtpPort, auth, s.Config.SmtpUsername, []string{to}, message)
	if err != nil {
		return err
	}
	return nil
}
