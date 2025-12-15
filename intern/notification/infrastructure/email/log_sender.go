package email

import (
	"context"
	"log/slog"

	"final-project/intern/notification/domain"
)

type LogEmailSender struct{}

func NewLogEmailSender() *LogEmailSender {
	return &LogEmailSender{}
}

func (sender *LogEmailSender) Send(_ context.Context, notification domain.EmailNotification) error {
	slog.Info(
		"send email notification (log only)",
		"type", string(notification.Type),
		"to", notification.To,
		"subject", notification.Subject,
		"body", notification.Body,
	)

	return nil
}
