package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"html"
	"strconv"
	"strings"

	"final-project/intern/notification/application"
	"final-project/intern/notification/domain"
)

var errUnexpectedEvent = errors.New("unexpected event type")

type bookingEventMessage struct {
	Event      string `json:"event"`
	Email      string `json:"email"`
	OwnerEmail string `json:"owner_email"`
	Hotel      string `json:"hotel"`
	Number     int    `json:"number"`
	Start      string `json:"start"`
	End        string `json:"end"`
}

type EmailNotificationHandler struct {
	notificationService *application.NotificationService
}

func NewEmailNotificationHandler(notificationService *application.NotificationService) *EmailNotificationHandler {
	return &EmailNotificationHandler{
		notificationService: notificationService,
	}
}

func (handler *EmailNotificationHandler) Handle(value []byte) error {
	var message bookingEventMessage
	if err := json.Unmarshal(value, &message); err != nil {
		return err
	}

	if domain.NotificationType(message.Event) != domain.NotificationTypeBookingCreated {
		return errUnexpectedEvent
	}

	subject, clientBody, ownerBody := buildEmailTemplates(message)

	if message.Email != "" {
		clientNotification := domain.EmailNotification{
			Type:    domain.NotificationType(message.Event),
			To:      message.Email,
			Subject: subject,
			Body:    clientBody,
		}
		if err := handler.notificationService.SendEmailNotification(context.Background(), clientNotification); err != nil {
			return err
		}
	}

	if message.OwnerEmail != "" {
		ownerNotification := domain.EmailNotification{
			Type:    domain.NotificationType(message.Event),
			To:      message.OwnerEmail,
			Subject: subject,
			Body:    ownerBody,
		}
		if err := handler.notificationService.SendEmailNotification(context.Background(), ownerNotification); err != nil {
			return err
		}
	}

	return nil
}

func buildEmailTemplates(message bookingEventMessage) (string, string, string) {
	subject := "Booking created"

	hotel := html.EscapeString(message.Hotel)
	room := html.EscapeString(strconv.Itoa(message.Number))
	start := html.EscapeString(message.Start)
	end := html.EscapeString(message.End)
	clientEmail := html.EscapeString(message.Email)

	clientBody := buildHTMLBody(subject, hotel, room, start, end, "")
	ownerBody := buildHTMLBody(subject, hotel, room, start, end, clientEmail)

	return subject, clientBody, ownerBody
}

func buildHTMLBody(title, hotel, room, start, end, clientEmail string) string {
	bodyBuilder := strings.Builder{}

	bodyBuilder.WriteString("<!doctype html><html><head><meta charset=\"utf-8\"/>")
	bodyBuilder.WriteString("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"/>")
	bodyBuilder.WriteString("<title>" + title + "</title>")
	bodyBuilder.WriteString("</head><body style=\"margin:0;padding:0;background:#f5f7fb;color:#111827;\">")

	bodyBuilder.WriteString("<div style=\"max-width:640px;margin:0 auto;padding:24px;\">")
	bodyBuilder.WriteString("<div style=\"background:#ffffff;border-radius:16px;padding:24px;border:1px solid #e5e7eb;\">")

	bodyBuilder.WriteString("<div style=\"font-size:20px;font-weight:700;line-height:28px;margin:0 0 12px 0;\">")
	bodyBuilder.WriteString(title)
	bodyBuilder.WriteString("</div>")

	bodyBuilder.WriteString("<div style=\"font-size:14px;line-height:20px;color:#374151;margin:0 0 20px 0;\">")
	bodyBuilder.WriteString("Your booking has been created. Details are below.")
	bodyBuilder.WriteString("</div>")

	bodyBuilder.WriteString("<table role=\"presentation\" cellspacing=\"0\" cellpadding=\"0\" style=\"width:100%;border-collapse:separate;border-spacing:0;\">")
	bodyBuilder.WriteString(row("Hotel", hotel))
	bodyBuilder.WriteString(row("Room", room))
	bodyBuilder.WriteString(row("Start", start))
	bodyBuilder.WriteString(row("End", end))
	if clientEmail != "" {
		bodyBuilder.WriteString(row("Client", clientEmail))
	}
	bodyBuilder.WriteString("</table>")

	bodyBuilder.WriteString("<div style=\"margin-top:20px;font-size:12px;line-height:18px;color:#6b7280;\">")
	bodyBuilder.WriteString("This email was sent automatically.")
	bodyBuilder.WriteString("</div>")

	bodyBuilder.WriteString("</div></div></body></html>")

	return bodyBuilder.String()
}

func row(label, value string) string {
	rowBuilder := strings.Builder{}
	rowBuilder.WriteString("<tr>")
	rowBuilder.WriteString("<td style=\"padding:10px 12px;border-top:1px solid #e5e7eb;font-size:13px;color:#6b7280;width:140px;\">")
	rowBuilder.WriteString(label)
	rowBuilder.WriteString("</td>")
	rowBuilder.WriteString("<td style=\"padding:10px 12px;border-top:1px solid #e5e7eb;font-size:14px;color:#111827;font-weight:600;\">")
	rowBuilder.WriteString(value)
	rowBuilder.WriteString("</td>")
	rowBuilder.WriteString("</tr>")
	return rowBuilder.String()
}
