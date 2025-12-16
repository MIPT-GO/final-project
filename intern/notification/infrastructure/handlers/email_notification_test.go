package handlers_test

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"testing"

	"final-project/intern/notification/application"
	"final-project/intern/notification/domain"
	"final-project/intern/notification/infrastructure/handlers"

	"github.com/stretchr/testify/assert"
)

type captureSender struct {
	mu            sync.Mutex
	notifications []domain.EmailNotification
}

func (sender *captureSender) Send(_ context.Context, notification domain.EmailNotification) error {
	sender.mu.Lock()
	defer sender.mu.Unlock()
	sender.notifications = append(sender.notifications, notification)
	return nil
}

func (sender *captureSender) all() []domain.EmailNotification {
	sender.mu.Lock()
	defer sender.mu.Unlock()
	out := make([]domain.EmailNotification, len(sender.notifications))
	copy(out, sender.notifications)
	return out
}

func TestEmailNotificationHandler_Handle_SendsToClientAndOwner(t *testing.T) {
	sender := &captureSender{}
	service := application.NewNotificationService(sender)
	handler := handlers.NewEmailNotificationHandler(service)

	message := map[string]any{
		"event":       "booking_created",
		"email":       "client@example.com",
		"owner_email": "owner@example.com",
		"hotel":       "Hilton",
		"number":      101,
		"start":       "2025-12-15T12:00:00Z",
		"end":         "2025-12-20T12:00:00Z",
	}
	data, _ := json.Marshal(message)

	err := handler.Handle(data)
	assert.NoError(t, err)

	sent := sender.all()
	assert.Len(t, sent, 2)

	assert.Equal(t, "Booking created", sent[0].Subject)
	assert.Equal(t, domain.NotificationTypeBookingCreated, sent[0].Type)
	assert.Equal(t, "client@example.com", sent[0].To)

	assert.Equal(t, "Booking created", sent[1].Subject)
	assert.Equal(t, domain.NotificationTypeBookingCreated, sent[1].Type)
	assert.Equal(t, "owner@example.com", sent[1].To)

	assert.True(t, strings.Contains(sent[0].Body, "Hilton"))
	assert.True(t, strings.Contains(sent[1].Body, "Hilton"))

	assert.False(t, strings.Contains(sent[0].Body, "Client"))
	assert.True(t, strings.Contains(sent[1].Body, "Client"))
	assert.True(t, strings.Contains(sent[1].Body, "client@example.com"))
}

func TestEmailNotificationHandler_Handle_EscapesHTML(t *testing.T) {
	sender := &captureSender{}
	service := application.NewNotificationService(sender)
	handler := handlers.NewEmailNotificationHandler(service)

	message := map[string]any{
		"event":       "booking_created",
		"email":       "client@example.com",
		"owner_email": "owner@example.com",
		"hotel":       "<b>Hilton</b>",
		"number":      101,
		"start":       "2025-12-15T12:00:00Z",
		"end":         "2025-12-20T12:00:00Z",
	}
	data, _ := json.Marshal(message)

	err := handler.Handle(data)
	assert.NoError(t, err)

	sent := sender.all()
	assert.Len(t, sent, 2)
	assert.Contains(t, sent[0].Body, "&lt;b&gt;Hilton&lt;/b&gt;")
}

func TestEmailNotificationHandler_Handle_UnexpectedEvent(t *testing.T) {
	sender := &captureSender{}
	service := application.NewNotificationService(sender)
	handler := handlers.NewEmailNotificationHandler(service)

	message := map[string]any{
		"event": "booking_cancelled",
	}
	data, _ := json.Marshal(message)

	err := handler.Handle(data)
	assert.EqualError(t, err, "unexpected event type")
	assert.Len(t, sender.all(), 0)
}

func TestEmailNotificationHandler_Handle_InvalidJSON(t *testing.T) {
	sender := &captureSender{}
	service := application.NewNotificationService(sender)
	handler := handlers.NewEmailNotificationHandler(service)

	err := handler.Handle([]byte("{not-json}"))
	assert.Error(t, err)
	assert.Len(t, sender.all(), 0)
}
