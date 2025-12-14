package storage

import (
	"final-project/intern/payment-system/domain"

	"github.com/google/uuid"
)

type RandomKeyGenerator struct{}

func (generator *RandomKeyGenerator) Generate(info *domain.PaymentInfo) string {
	id := uuid.New()
	return id.String()
}
