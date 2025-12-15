package application

import (
	"context"

	"final-project/intern/notification/domain"
)

type EmailSender interface {
	Send(ctx context.Context, notification domain.EmailNotification) error
}

type NotificationService struct {
	sender EmailSender
}

func NewNotificationService(emailSender EmailSender) *NotificationService {
	return &NotificationService{
		sender: emailSender,
	}
}

func (service *NotificationService) SendEmailNotification(
	ctx context.Context,
	notification domain.EmailNotification,
) error {
	return service.sender.Send(ctx, notification)
}
