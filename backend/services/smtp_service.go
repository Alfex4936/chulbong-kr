package services

import (
	"chulbong-kr/database"
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"net/smtp"
	"os"
	"strings"
	"time"
)

var (
	smtpServer          = os.Getenv("SMTP_SERVER")
	smtpPort            = os.Getenv("SMTP_PORT")
	smtpUsername        = os.Getenv("SMTP_USERNAME")
	smtpPassword        = os.Getenv("SMTP_PASSWORD")
	frontendResetRouter = os.Getenv("FRONTEND_RESET_ROUTER")
)

var emailTemplate = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html lang="ko" xmlns="http://www.w3.org/1999/xhtml">
<head>
    <title>chulbong-kr</title>
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
    <title>Password Reset for chulbong-kr</title>
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

// GenerateToken generates a secure random token that is 6 digits long
func GenerateSixDigitToken() (string, error) {
	// Define the maximum value (999999) for a 6-digit number
	max := big.NewInt(999999)

	// Generate a random number between 0 and max
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	// Format the number as a 6-digit string with leading zeros if necessary
	token := fmt.Sprintf("%06d", n.Int64())

	return token, nil
}

func GenerateAndSaveSignUpToken(email string) (string, error) {
	token, err := GenerateSixDigitToken()
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(5 * time.Minute)

	// Attempt to insert or update the token for the user
	_, err = database.DB.Exec(`
        INSERT INTO PasswordTokens (Email, Token, ExpiresAt, Verified)
        VALUES (?, ?, ?, FALSE)
        ON DUPLICATE KEY UPDATE Token=VALUES(Token), ExpiresAt=VALUES(ExpiresAt), Verified=FALSE`,
		email, token, expiresAt)
	if err != nil {
		return "", fmt.Errorf("error saving or updating token: %w", err)
	}

	return token, nil
}

func ValidateToken(token string, email string) (bool, error) {
	// Start transaction
	tx, err := database.DB.Beginx()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	var expiresAt time.Time
	err = tx.QueryRow("SELECT ExpiresAt FROM PasswordTokens WHERE Token = ? AND Email = ? AND ExpiresAt > NOW() LIMIT 1", token, email).Scan(&expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // Token not found or expired
		}
		return false, err // Database or other error
	}

	// Update the Verified status
	_, err = tx.Exec("UPDATE PasswordTokens SET Verified = TRUE WHERE Token = ? AND Email = ?", token, email)
	if err != nil {
		return false, err
	}

	tx.Commit()
	return true, nil // Token is valid, not expired, and now marked as verified
}

func IsTokenVerified(email string) (bool, error) {
	var verified bool
	err := database.DB.Get(&verified, "SELECT Verified FROM PasswordTokens WHERE Email = ? AND ExpiresAt > NOW() AND Verified = TRUE LIMIT 1", email)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // No verified token found
		}
		return false, err // An error occurred
	}
	return verified, nil // A verified token exists
}

// SendVerificationEmail sends a verification email to the user
func SendVerificationEmail(to, token string) error {

	// Define email headers including content type for HTML
	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: chulbong-kr Email Verification\r\nMIME-Version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n", smtpUsername, to)

	// Replace the {{TOKEN}} placeholder in the template with the actual token
	htmlBody := strings.Replace(emailTemplate, "{{TOKEN}}", token, -1)

	// Combine headers and HTML body into a single raw email message
	message := []byte(headers + htmlBody)

	// Connect to the SMTP server and send the email
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)
	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, smtpUsername, []string{to}, message)
	if err != nil {
		return err
	}
	return nil
}

func SendPasswordResetEmail(to, token string) error {
	// Define email headers
	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: Password Reset for chulbong-kr\r\nMIME-Version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n", smtpUsername, to)

	// Replace the {{RESET_LINK}} placeholder with the actual reset link
	clientUrl := fmt.Sprintf("%s?token=%s&email=%s", frontendResetRouter, token, to)
	htmlBody := strings.Replace(emailTemplateForReset, "{{RESET_LINK}}", clientUrl, -1)

	// Combine headers and HTML body into a single raw email message
	message := []byte(headers + htmlBody)

	// Connect to the SMTP server and send the email
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)
	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, smtpUsername, []string{to}, message)
	if err != nil {
		return err
	}
	return nil
}
