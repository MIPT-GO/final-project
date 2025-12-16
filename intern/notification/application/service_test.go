package application_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"final-project/intern/notification/application"
	"final-project/intern/notification/domain"

	"github.com/stretchr/testify/assert"
)

type captureSender struct {
	mu    sync.Mutex
	last  domain.EmailNotification
	calls int
	err   error
}

func (sender *captureSender) Send(_ context.Context, notification domain.EmailNotification) error {
	sender.mu.Lock()
	defer sender.mu.Unlock()
	sender.calls++
	sender.last = notification
	return sender.err
}

func TestNotificationService_SendsThroughSender(t *testing.T) {
	sender := &captureSender{}
	service := application.NewNotificationService(sender)

	notification := domain.EmailNotification{
		Type:    domain.NotificationTypeBookingCreated,
		To:      "user@example.com",
		Subject: "subject",
		Body:    "<html/>",
	}

	err := service.SendEmailNotification(t.Context(), notification)
	assert.NoError(t, err)
	assert.Equal(t, 1, sender.calls)
	assert.Equal(t, notification, sender.last)
}

func TestNotificationService_ReturnsSenderError(t *testing.T) {
	expectedErr := errors.New("send failed")
	sender := &captureSender{err: expectedErr}
	service := application.NewNotificationService(sender)

	err := service.SendEmailNotification(t.Context(), domain.EmailNotification{To: "user@example.com"})
	assert.ErrorIs(t, err, expectedErr)
}


