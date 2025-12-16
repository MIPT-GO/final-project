package email

import (
	"testing"

	"final-project/intern/notification/config"
	"final-project/intern/notification/domain"

	"github.com/stretchr/testify/assert"
)

func TestSMTPSender_Send_NoCredentials(t *testing.T) {
	cfg := config.Config{
		SMTPHost:     "smtp.gmail.com",
		SMTPPort:     587,
		SMTPUsername: "",
		SMTPPassword: "",
	}

	sender := NewSMTPSender(&cfg)
	err := sender.Send(t.Context(), domain.EmailNotification{
		To:      "user@example.com",
		Subject: "subject",
		Body:    "<b>hello</b>",
	})

	assert.ErrorIs(t, err, errSMTPNotConfigured)
}

func TestSMTPSender_Send_EmptyRecipient(t *testing.T) {
	cfg := config.Config{
		SMTPHost:     "smtp.gmail.com",
		SMTPPort:     587,
		SMTPUsername: "sender@example.com",
		SMTPPassword: "app-password",
	}

	sender := NewSMTPSender(&cfg)
	err := sender.Send(t.Context(), domain.EmailNotification{
		To:      "   ",
		Subject: "subject",
		Body:    "<b>hello</b>",
	})

	assert.ErrorIs(t, err, errEmptyRecipientEmail)
}

func TestBuildEmailMessage_HTMLHeaders(t *testing.T) {
	raw := buildEmailMessage("from@example.com", "to@example.com", "Booking created", "<html>ok</html>")

	assert.Contains(t, raw, "From: from@example.com\r\n")
	assert.Contains(t, raw, "To: to@example.com\r\n")
	assert.Contains(t, raw, "Subject: ")
	assert.Contains(t, raw, "utf-8")
	assert.Contains(t, raw, "MIME-Version: 1.0\r\n")
	assert.Contains(t, raw, "Content-Type: text/html; charset=\"utf-8\"\r\n")
	assert.Contains(t, raw, "\r\n<html>ok</html>")
}

func TestBuildEmailMessage_EmptySubject(t *testing.T) {
	raw := buildEmailMessage("from@example.com", "to@example.com", "", "<html>ok</html>")

	assert.Contains(t, raw, "From: from@example.com\r\n")
	assert.Contains(t, raw, "To: to@example.com\r\n")
	assert.NotContains(t, raw, "Subject:")
	assert.Contains(t, raw, "Content-Type: text/html; charset=\"utf-8\"\r\n")
}


