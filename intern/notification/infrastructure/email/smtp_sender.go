package email

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mime"
	"net"
	"net/smtp"
	"strconv"
	"strings"

	"final-project/intern/notification/config"
	"final-project/intern/notification/domain"
)

var (
	errSMTPNotConfigured   = errors.New("smtp credentials are not configured")
	errEmptyRecipientEmail = errors.New("empty recipient email")
)

type SMTPSender struct {
	config *config.Config
}

func NewSMTPSender(cfg *config.Config) *SMTPSender {
	return &SMTPSender{
		config: cfg,
	}
}

func (sender *SMTPSender) Send(_ context.Context, notification domain.EmailNotification) error {
	if sender.config.SMTPUsername == "" || sender.config.SMTPPassword == "" {
		return errSMTPNotConfigured
	}

	recipient := strings.TrimSpace(notification.To)
	if recipient == "" {
		return errEmptyRecipientEmail
	}

	from := sender.config.SMTPFrom
	if from == "" {
		from = sender.config.SMTPUsername
	}

	auth := smtp.PlainAuth(
		"",
		sender.config.SMTPUsername,
		sender.config.SMTPPassword,
		sender.config.SMTPHost,
	)

	smtpAddress := net.JoinHostPort(sender.config.SMTPHost, strconv.Itoa(sender.config.SMTPPort))

	message := buildEmailMessage(from, recipient, notification.Subject, notification.Body)

	if err := smtp.SendMail(smtpAddress, auth, from, []string{recipient}, []byte(message)); err != nil {
		return fmt.Errorf("failed to send email via smtp: %w", err)
	}

	slog.Info("email sent", "to", recipient, "subject", notification.Subject)
	return nil
}

func buildEmailMessage(from, to, subject, body string) string {
	var messageBuilder strings.Builder

	messageBuilder.WriteString("From: " + from + "\r\n")
	messageBuilder.WriteString("To: " + to + "\r\n")

	if subject != "" {
		messageBuilder.WriteString("Subject: " + mime.QEncoding.Encode("utf-8", subject) + "\r\n")
	}

	messageBuilder.WriteString("MIME-Version: 1.0\r\n")
	messageBuilder.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
	messageBuilder.WriteString("\r\n")
	messageBuilder.WriteString(body)

	return messageBuilder.String()
}
