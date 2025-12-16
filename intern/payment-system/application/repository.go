package application

import (
	"context"
	"final-project/intern/payment-system/domain"
)

type PaymentInfoRepository interface {
	CreateRecord(ctx context.Context, key string, info domain.PaymentInfo) error
	GetRecord(ctx context.Context, key string) (domain.PaymentInfo, error)
}
