package application

import (
	"final-project/intern/payment-system/domain"
)

type WebHookSender interface {
	Send(webhook string, status domain.PaymentStatus) error
}
