package application

import "final-project/intern/payment-system/domain"

type KeyGenerator interface {
	Generate(info *domain.PaymentInfo) string
}
