package domain

type NotificationType string

const (
	NotificationTypeBookingCreated NotificationType = "booking_created"
)

type EmailNotification struct {
	Type    NotificationType `json:"type"`
	To      string           `json:"to"`
	Subject string           `json:"subject"`
	Body    string           `json:"body"`
}
