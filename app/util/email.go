package util

import (
	"fmt"
	"net/smtp"

	"riz.it/domped/app/config"
	"riz.it/domped/app/domain"
)

type EmailUtil struct {
	Config *config.Config
}

func NewEmailUtil(config *config.Config) domain.Email {
	return &EmailUtil{
		Config: config,
	}
}

// Send implements domain.Email.
func (e *EmailUtil) Send(to string, subject string, body string) error {
	// Set up authentication information.
	auth := smtp.PlainAuth("", e.Config.SMTP.User, e.Config.SMTP.Password, e.Config.SMTP.Host)

	// Construct the email message with proper headers for HTML content.
	msg := []byte(
		"From: " + e.Config.SMTP.User + "\n" +
			"To: " + to + "\n" +
			"Subject: " + subject + "\n" +
			"Content-Type: text/html; charset=UTF-8\n\n" + // Add Content-Type header for HTML
			body,
	)

	// Send the email.
	serverAddr := fmt.Sprintf("%s:%s", e.Config.SMTP.Host, e.Config.SMTP.Port)
	if err := smtp.SendMail(serverAddr, auth, e.Config.SMTP.User, []string{to}, msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
