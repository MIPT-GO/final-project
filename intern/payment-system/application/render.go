package application

import (
	"final-project/intern/payment-system/config"
	"io"
)

type MessagePageData struct {
	Message string
}

type PaymentPageData struct {
	Amount  string
	Message string
	Key     string
	Prefix  string
}

type Render interface {
	New(config *config.Config)
	RenderPaymentPage(writer io.Writer, data PaymentPageData) error
	RenderMessagePage(writer io.Writer, data MessagePageData) error
}
