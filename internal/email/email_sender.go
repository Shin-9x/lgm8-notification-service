package email

import (
	"fmt"
	"log"
	"net/smtp"
)

// EmailSender handles sending emails using an SMTP server
type EmailSender struct {
	Enabled    bool
	SMTPServer string
	Port       int
	Auth       smtp.Auth
	From       string
}

// NewEmailSender initializes an email sender with the given SMTP credentials
func NewEmailSender(enabled bool, smtpServer string, port int, username, password, from string) *EmailSender {
	auth := smtp.PlainAuth("", username, password, smtpServer)
	return &EmailSender{
		Enabled:    enabled,
		SMTPServer: smtpServer,
		Port:       port,
		Auth:       auth,
		From:       from,
	}
}

// SendEmail sends an email with the given subject and body
func (e *EmailSender) SendEmail(to, subject, body string) error {
	if e.Enabled {
		log.Println("Email notification enabled. Proceed sending the mail...")

		msg := []byte(
			"From: " + e.From + "\r\n" +
				"To: " + to + "\r\n" +
				"Subject: " + subject + "\r\n\r\n" +
				body + "\r\n",
		)

		addr := fmt.Sprintf("%s:%d", e.SMTPServer, e.Port)

		err := smtp.SendMail(addr, e.Auth, e.From, []string{to}, msg)
		if err != nil {
			return fmt.Errorf("failed to send email: [%w]", err)
		}
	} else {
		log.Println("Email notification disabled.")
	}

	return nil
}
