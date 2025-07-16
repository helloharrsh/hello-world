package mail

import (
	"fmt"
	"net/smtp"
	"strings"

	"mailer_application/internal/config"
)

type Mailer struct {
	auth      smtp.Auth
	smtpHost  string
	smtpPort  string
	fromEmail string
}

// NewMailer creates a new Mailer instance
func NewMailer(cfg *config.Config) *Mailer {
	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPHost)

	return &Mailer{
		auth:      auth,
		smtpHost:  cfg.SMTPHost,
		smtpPort:  cfg.SMTPPort,
		fromEmail: cfg.FromEmail,
	}
}

// Send sends an email to the recipient
func (m *Mailer) Send(to string, subject string, body string) error {
	headers := []string{
		fmt.Sprintf("From: %s", m.fromEmail),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=\"utf-8\"",
		"",
		body,
	}

	message := strings.Join(headers, "\r\n")

	addr := fmt.Sprintf("%s:%s", m.smtpHost, m.smtpPort)
	return smtp.SendMail(addr, m.auth, m.fromEmail, []string{to}, []byte(message))
}
